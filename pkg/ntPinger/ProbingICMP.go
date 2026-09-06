package ntPinger

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var errNativeICMPUnavailable = errors.New("native unprivileged ICMP is unavailable")

var (
	nativeICMPProbe  = icmpNativeProbing
	commandICMPProbe = icmpCommandProbing
)

// IcmpProbing sends one IPv4 echo request. Standard probes use an unprivileged
// ICMP datagram socket on supported systems. The OS ping command remains the
// compatibility path for DF probes and systems without unprivileged ICMP.
func IcmpProbing(seq int, destAddr string, destHost string, payloadSize int, df bool, timeout int, payload []byte) (PacketICMP, error) {
	if !df && (runtime.GOOS == "darwin" || runtime.GOOS == "linux") {
		pkt, err := nativeICMPProbe(seq, destAddr, destHost, payloadSize, df, timeout, payload)
		if err == nil {
			return pkt, nil
		}
		if !errors.Is(err, errNativeICMPUnavailable) {
			return pkt, err
		}
	}

	return commandICMPProbe(seq, destAddr, destHost, payloadSize, df, timeout, payload)
}

func newICMPPacket(seq int, destAddr string, destHost string, payloadSize int, df bool) PacketICMP {
	return PacketICMP{
		Type:        "icmp",
		Status:      false,
		Seq:         seq,
		DestAddr:    destAddr,
		DestHost:    destHost,
		PayLoadSize: payloadSize,
		Icmp_DF:     df,
	}
}

func icmpNativeProbing(seq int, destAddr string, destHost string, payloadSize int, df bool, timeout int, payload []byte) (PacketICMP, error) {
	pkt := newICMPPacket(seq, destAddr, destHost, payloadSize, df)
	destIP := net.ParseIP(destAddr).To4()
	if destIP == nil {
		return pkt, fmt.Errorf("invalid IPv4 destination: %s", destAddr)
	}

	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return pkt, fmt.Errorf("%w: %v", errNativeICMPUnavailable, err)
	}
	defer conn.Close()

	data := payload
	if len(data) != payloadSize {
		data = GeneratePayloadData(payloadSize)
	}
	request := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Seq: seq, Data: data},
	}
	wire, err := request.Marshal(nil)
	if err != nil {
		return pkt, fmt.Errorf("marshal ICMP echo request: %w", err)
	}

	pkt.SendTime = time.Now()
	deadline := pkt.SendTime.Add(time.Duration(timeout) * time.Second)
	if err := conn.SetDeadline(deadline); err != nil {
		return pkt, fmt.Errorf("set ICMP deadline: %w", err)
	}
	if _, err := conn.WriteTo(wire, &net.UDPAddr{IP: destIP}); err != nil {
		if errors.Is(err, os.ErrPermission) || errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
			return pkt, fmt.Errorf("%w: %v", errNativeICMPUnavailable, err)
		}
		return pkt, fmt.Errorf("send ICMP echo request: %w", err)
	}

	buffer := make([]byte, 65535)
	for {
		n, peer, err := conn.ReadFrom(buffer)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				pkt.AdditionalInfo = "Timeout"
				return pkt, nil
			}
			return pkt, fmt.Errorf("read ICMP reply: %w", err)
		}

		done, status, info := parseNativeICMPReply(buffer[:n], peer, destIP, seq)
		if !done {
			continue
		}
		pkt.Status = status
		pkt.AdditionalInfo = info
		if status {
			pkt.RTT = time.Since(pkt.SendTime)
		}
		return pkt, nil
	}
}

func parseNativeICMPReply(data []byte, peer net.Addr, destIP net.IP, seq int) (done bool, status bool, info string) {
	reply, err := icmp.ParseMessage(1, data)
	if err != nil {
		return false, false, ""
	}
	switch reply.Type {
	case ipv4.ICMPTypeEchoReply:
		echo, ok := reply.Body.(*icmp.Echo)
		if !ok || echo.Seq != seq || !icmpPeerMatches(peer, destIP) {
			return false, false, ""
		}
		return true, true, ""
	case ipv4.ICMPTypeDestinationUnreachable:
		if reply.Code == 4 {
			return true, false, "MTU Exceed, DF set"
		}
		return true, false, "NoRoute"
	case ipv4.ICMPTypeTimeExceeded:
		return true, false, "TTL_Exceeded"
	default:
		return false, false, ""
	}
}

func icmpPeerMatches(peer net.Addr, destIP net.IP) bool {
	switch addr := peer.(type) {
	case *net.UDPAddr:
		return addr.IP.Equal(destIP)
	case *net.IPAddr:
		return addr.IP.Equal(destIP)
	default:
		return false
	}
}
