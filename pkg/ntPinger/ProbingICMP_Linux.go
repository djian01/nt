//go:build linux
// +build linux

package ntPinger

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

// func: IcmpProbing
func IcmpProbing(Seq int, destAddr string, desetHost string, PayLoadSize int, Icmp_DF bool, timeout int, payload []byte) (PacketICMP, error) {

	// Initial PacketICMP
	pkt := PacketICMP{
		Type:        "icmp",
		Status:      false,
		Seq:         Seq,
		DestAddr:    destAddr,
		DestHost:    desetHost,
		PayLoadSize: PayLoadSize,
		Icmp_DF:     Icmp_DF,
	}

	// Convert the string to *net.IPAddr
	ipAddr, err := net.ResolveIPAddr("ip4", destAddr)
	if err != nil {
		return pkt, fmt.Errorf("failed to resolve IP address: %v", err)
	}

	// Create a raw socket with IPPROTO_ICMP
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		return pkt, fmt.Errorf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Set the "Don't Fragment" (DF) flag
	err = SetDFBit(conn, Icmp_DF)
	if err != nil {
		return pkt, err
	}

	// Prepare ICMP message
	icmpType := icmpv4EchoRequest
	icmpId := os.Getpid() & 0xffff

	icmpB := icmpBody{
		ID:   icmpId,
		Seq:  Seq,
		Data: payload,
	}

	//// build ICMP Request Binary
	BinIcmpReq, err := (&icmpMessage{
		Type: icmpType,
		Code: 0,
		Body: &icmpB,
	}).Marshal()

	if err != nil {
		return pkt, err
	}

	// TIMEOUT - Set timeout for the connection
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	// SEND - the ICMP Request
	pkt.SendTime = time.Now()

	if _, err = conn.Write(BinIcmpReq); err != nil {
		if strings.Contains(err.Error(), "message too long") {
			// timeout
			pkt.AdditionalInfo = "MTU Exceed, DF set"
			pkt.Status = false
			return pkt, nil
		} else {
			return pkt, fmt.Errorf("failed to send request: %w", err)
		}

	}

	// RECEIVE -  the ICMP Response
	BinIcmpRep := make([]byte, PayLoadSize+42)

	if _, err = conn.Read(BinIcmpRep); err != nil {
		if strings.Contains(err.Error(), "timeout") {
			// timeout
			pkt.AdditionalInfo = "Timeout"
			pkt.Status = false
			return pkt, nil
		} else {
			return pkt, fmt.Errorf("failed to read response: %w", err)
		}

	}

	// RTT - Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// Process the ICMP Response
	BinIcmpRep = ipv4Payload(BinIcmpRep)

	var icmpMsgRep *icmpMessage
	if icmpMsgRep, err = parseICMPMessage(BinIcmpRep); err != nil {
		return pkt, fmt.Errorf("failed to parsing the recived reply: %w", err)
	}

	// Verify the ICMP response seq and the request seq
	if (*icmpMsgRep).Body.Seq == Seq {
		pkt.Status = true
	} else {
		pkt.AdditionalInfo = "Seq_Mismatch"
	}

	// Check Latency
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil

}

// ************* ICMP Type const ******************
const (
	icmpv4EchoRequest = 8
	icmpv4EchoReply   = 0
	icmpv6EchoRequest = 128
	icmpv6EchoReply   = 129
)

// Visualization of the ICMP header for an Echo Request/Reply:

// 0               1               2               3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Type      |     Code      |          Checksum             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           Identifier          |        Sequence Number        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                             Payload                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// ************* Struct, Method ******************

// Struct - icmpMessage: Icmp Whole Message
type icmpMessage struct {
	Type     int       // type
	Code     int       // code
	Checksum int       // checksum
	Body     *icmpBody // body
}

// icmpMessage Method - Len()
func (icmpMsg *icmpMessage) Marshal() ([]byte, error) {

	// 4 x bytes of ICMP Header, 1st Byte: Type, 2nd Byte: Code, 3rd & 4th Bytes: Checksum
	binIcmpMsg := []byte{byte(icmpMsg.Type), byte(icmpMsg.Code), 0, 0}

	// if icmpMsg Body is not nil and Len is not 0, append ICMP Header & Body
	if icmpMsg.Body != nil && icmpMsg.Body.Len() != 0 {
		binIcmpBody, err := icmpMsg.Body.Marshal()
		if err != nil {
			return nil, err
		}
		binIcmpMsg = append(binIcmpMsg, binIcmpBody...)
	}

	// if the type is icmpv6EchoRequest or icmpv6EchoReply, return binIcmpMsg
	// IPv6 ICMP checksum is handled differently
	switch icmpMsg.Type {
	case icmpv6EchoRequest, icmpv6EchoReply:
		return binIcmpMsg, nil
	}

	// if the type is icmpv4EchoRequest or icmpv4EchoReply, move forward

	// Create Chechsum for the ICMP Message Bin
	csum := Checksum(binIcmpMsg)
	csumByte := ChecksumToByte(csum)

	// Place checksum back in header
	binIcmpMsg[2] = csumByte[0]
	binIcmpMsg[3] = csumByte[1]

	return binIcmpMsg, nil
}

// Struct - icmpEchoBody: Icmp request body. Satisfy Interface - icmpMessageBody
type icmpBody struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

// icmpEchoBody Method - Len()
func (icmpB *icmpBody) Len() int {
	if icmpB == nil {
		return 0
	}

	// the total Body len = 2 bytes (ID) + 2 bytes (Seq) + len(b.Data)
	return 4 + len(icmpB.Data)
}

// icmpEchoBody Method - Marshal()
func (b *icmpBody) Marshal() ([]byte, error) {

	// len is len(b.Data) + 2 bytes (ID) + 2 bytes (Seq)
	binIcmpBody := make([]byte, 4+len(b.Data))

	// high byte -> ID, low byte -> bin[0]
	binIcmpBody[0], binIcmpBody[1] = byte(b.ID>>8), byte(b.ID&0xff)

	// high byte -> Seq, low byte -> bin[0]
	binIcmpBody[2], binIcmpBody[3] = byte(b.Seq>>8), byte(b.Seq&0xff)

	// copy the payload start from 5th byte
	copy(binIcmpBody[4:], b.Data)
	return binIcmpBody, nil
}

// ************* ICMP Response Functions ******************

// func - parseICMPMessage, parses bin as an ICMP message binary []byte
func parseICMPMessage(bin []byte) (*icmpMessage, error) {

	msglen := len(bin)
	if msglen < 4 {
		return nil, errors.New("message too short")
	}

	icmpMsg := &icmpMessage{Type: int(bin[0]), Code: int(bin[1]), Checksum: int(bin[2])<<8 | int(bin[3])}

	if msglen > 4 {
		var err error
		switch icmpMsg.Type {
		case icmpv4EchoRequest, icmpv4EchoReply, icmpv6EchoRequest, icmpv6EchoReply:
			icmpMsg.Body, err = parseICMPEcho(bin[4:])
			if err != nil {
				return nil, err
			}
		}
	}
	return icmpMsg, nil
}

// parseICMPEcho parses b as an ICMP echo request or reply message body.
func parseICMPEcho(bin []byte) (*icmpBody, error) {

	bodylen := len(bin)

	icmpB := &icmpBody{ID: int(bin[0])<<8 | int(bin[1]), Seq: int(bin[2])<<8 | int(bin[3])}

	if bodylen > 4 {
		icmpB.Data = make([]byte, bodylen-4)
		copy(icmpB.Data, bin[4:])
	}
	return icmpB, nil
}

func ipv4Payload(bin []byte) []byte {

	// if the packet is less than 20 bytes (the minimum size of an IPv4 header), simply return bin as is.
	if len(bin) < 20 {
		return bin
	}
	headerLen := int(bin[0]&0x0f) << 2 // headerLen shifting left by 2 bits gives the total length of the IPv4 header in bytes
	return bin[headerLen:]
}

// ******************** Checksum Funcs ***************************

// Checksum calculates the checksum for an ICMP message
func Checksum(data []byte) uint32 {

	csumCoverage := len(data) - 1 // checksum coverage

	csum := uint32(0)

	for i := 0; i < csumCoverage; i += 2 {
		csum += uint32(data[i+1])<<8 | uint32(data[i])
	}
	if csumCoverage&1 == 0 {
		csum += uint32(data[csumCoverage])
	}
	csum = csum>>16 + csum&0xffff
	csum = csum + csum>>16

	return csum
}

// convert Checksum to []byte
func ChecksumToByte(csum uint32) []byte {
	bin := []byte{}
	bin = append(bin, byte(^csum&0xff), byte(^csum>>8))
	return bin
}

// Func - Set DF bit
func SetDFBit(conn *net.IPConn, df bool) error {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	if df {
		rawConn.Control(func(fd uintptr) {
			syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MTU_DISCOVER, syscall.IP_PMTUDISC_DO)
		})
	} else {
		rawConn.Control(func(fd uintptr) {
			syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MTU_DISCOVER, syscall.IP_PMTUDISC_DONT)
		})
	}
	return nil
}
