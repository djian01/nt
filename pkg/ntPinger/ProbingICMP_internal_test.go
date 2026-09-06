package ntPinger

import (
	"errors"
	"net"
	"runtime"
	"testing"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func TestParseNativeICMPReply(t *testing.T) {
	destIP := net.ParseIP("127.0.0.1").To4()
	peer := &net.UDPAddr{IP: destIP}
	reply := icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Code: 0,
		Body: &icmp.Echo{ID: 123, Seq: 7, Data: []byte("payload")},
	}
	wire, err := reply.Marshal(nil)
	if err != nil {
		t.Fatal(err)
	}

	done, status, info := parseNativeICMPReply(wire, peer, destIP, 7)
	if !done || !status || info != "" {
		t.Fatalf("unexpected reply result: done=%t status=%t info=%q", done, status, info)
	}

	done, _, _ = parseNativeICMPReply(wire, peer, destIP, 8)
	if done {
		t.Fatal("reply with a different sequence was accepted")
	}

	done, _, _ = parseNativeICMPReply([]byte{1, 2, 3}, peer, destIP, 7)
	if done {
		t.Fatal("malformed reply was accepted")
	}
}

func TestParseNativeICMPErrorReplies(t *testing.T) {
	destIP := net.ParseIP("127.0.0.1").To4()
	peer := &net.IPAddr{IP: destIP}
	tests := []struct {
		message icmp.Message
		want    string
	}{
		{
			message: icmp.Message{Type: ipv4.ICMPTypeDestinationUnreachable, Code: 4, Body: &icmp.DstUnreach{Data: make([]byte, 28)}},
			want:    "MTU Exceed, DF set",
		},
		{
			message: icmp.Message{Type: ipv4.ICMPTypeDestinationUnreachable, Code: 1, Body: &icmp.DstUnreach{Data: make([]byte, 28)}},
			want:    "NoRoute",
		},
		{
			message: icmp.Message{Type: ipv4.ICMPTypeTimeExceeded, Code: 0, Body: &icmp.TimeExceeded{Data: make([]byte, 28)}},
			want:    "TTL_Exceeded",
		},
	}

	for _, test := range tests {
		wire, err := test.message.Marshal(nil)
		if err != nil {
			t.Fatal(err)
		}
		done, status, info := parseNativeICMPReply(wire, peer, destIP, 1)
		if !done || status || info != test.want {
			t.Fatalf("unexpected error reply: done=%t status=%t info=%q", done, status, info)
		}
	}
}

func TestIcmpProbingTransportSelection(t *testing.T) {
	originalNative := nativeICMPProbe
	originalCommand := commandICMPProbe
	defer func() {
		nativeICMPProbe = originalNative
		commandICMPProbe = originalCommand
	}()

	nativeCalls, commandCalls := 0, 0
	nativeICMPProbe = func(seq int, destAddr, destHost string, payloadSize int, df bool, timeout int, payload []byte) (PacketICMP, error) {
		nativeCalls++
		pkt := newICMPPacket(seq, destAddr, destHost, payloadSize, df)
		pkt.Status = true
		return pkt, nil
	}
	commandICMPProbe = func(seq int, destAddr, destHost string, payloadSize int, df bool, timeout int, payload []byte) (PacketICMP, error) {
		commandCalls++
		return newICMPPacket(seq, destAddr, destHost, payloadSize, df), nil
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if _, err := IcmpProbing(1, "127.0.0.1", "localhost", 8, false, 1, make([]byte, 8)); err != nil {
			t.Fatal(err)
		}
		if nativeCalls != 1 || commandCalls != 0 {
			t.Fatalf("standard probe calls: native=%d command=%d", nativeCalls, commandCalls)
		}

		nativeICMPProbe = func(seq int, destAddr, destHost string, payloadSize int, df bool, timeout int, payload []byte) (PacketICMP, error) {
			nativeCalls++
			return PacketICMP{}, errNativeICMPUnavailable
		}
		if _, err := IcmpProbing(2, "127.0.0.1", "localhost", 8, false, 1, make([]byte, 8)); err != nil {
			t.Fatal(err)
		}
		if commandCalls != 1 {
			t.Fatalf("command fallback calls=%d, want 1", commandCalls)
		}
	}

	if _, err := IcmpProbing(3, "127.0.0.1", "localhost", 8, true, 1, make([]byte, 8)); err != nil {
		t.Fatal(err)
	}
	if commandCalls == 0 {
		t.Fatal("DF probe did not use the command compatibility path")
	}
}

func TestNativeICMPToLocalhostWithoutCommand(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("unprivileged ICMP datagram sockets are supported on macOS and Linux")
	}
	pkt, err := icmpNativeProbing(11, "127.0.0.1", "localhost", 16, false, 1, GeneratePayloadData(16))
	if errors.Is(err, errNativeICMPUnavailable) {
		t.Skipf("unprivileged ICMP is disabled: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if !pkt.Status || pkt.Seq != 11 || pkt.PayLoadSize != 16 || pkt.RTT <= 0 {
		t.Fatalf("unexpected localhost result: %+v", pkt)
	}
}

func TestRunProbesClosesOnceAndStopsAfterCount(t *testing.T) {
	p := &Pinger{
		InputVars: InputVars{Count: 2, Interval: 0},
		ProbeChan: make(chan Packet, 2),
		Stat:      Statistics{MinRtt: time.Duration(1<<63 - 1), MaxRtt: time.Duration(-1 << 63)},
	}
	errChan := make(chan error, 1)
	runProbes(p, errChan, func(seq int) (Packet, error) {
		return &PacketTCP{Type: "tcp", Seq: seq, Status: true, RTT: time.Millisecond}, nil
	})

	count := 0
	for range p.ProbeChan {
		count++
	}
	if count != 2 {
		t.Fatalf("received %d probes, want 2", count)
	}
	select {
	case err := <-errChan:
		t.Fatalf("unexpected error: %v", err)
	default:
	}
}

func TestRunProbesClosesAfterError(t *testing.T) {
	p := &Pinger{InputVars: InputVars{Count: 1}, ProbeChan: make(chan Packet, 1)}
	errChan := make(chan error, 1)
	wantErr := errors.New("probe failed")
	runProbes(p, errChan, func(int) (Packet, error) { return nil, wantErr })
	if _, ok := <-p.ProbeChan; ok {
		t.Fatal("probe channel remained open after error")
	}
	if err := <-errChan; !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
}

func TestRunProbesMarksLatencySpike(t *testing.T) {
	p := &Pinger{
		InputVars: InputVars{Count: 2, Interval: 0},
		ProbeChan: make(chan Packet, 2),
		Stat:      Statistics{MinRtt: time.Duration(1<<63 - 1), MaxRtt: time.Duration(-1 << 63)},
	}
	runProbes(p, make(chan error, 1), func(seq int) (Packet, error) {
		rtt := 20 * time.Millisecond
		if seq == 1 {
			rtt = 50 * time.Millisecond
		}
		return &PacketTCP{Type: "tcp", Seq: seq, Status: true, RTT: rtt}, nil
	})

	<-p.ProbeChan
	second := (<-p.ProbeChan).(*PacketTCP)
	if second.AdditionalInfo != "High_Latency" {
		t.Fatalf("additional info=%q, want High_Latency", second.AdditionalInfo)
	}
}

func TestGetSleepTime(t *testing.T) {
	if got := GetSleepTime(true, 1, 250*time.Millisecond); got != 750*time.Millisecond {
		t.Fatalf("sleep time=%v, want 750ms", got)
	}
	if got := GetSleepTime(false, 1, 0); got != time.Second {
		t.Fatalf("timeout sleep=%v, want 1s", got)
	}
	if got := GetSleepTime(true, 1, 2*time.Second); got != 0 {
		t.Fatalf("overrun sleep=%v, want 0", got)
	}
}

func TestValidateInputVars(t *testing.T) {
	tests := []InputVars{
		{Type: "icmp", Count: -1, Timeout: 1},
		{Type: "icmp", PayLoadSize: -1, Timeout: 1},
		{Type: "icmp", PayLoadSize: 65508, Timeout: 1},
		{Type: "icmp", Timeout: 0},
		{Type: "icmp", Timeout: 1, Interval: -1},
		{Type: "tcp", Timeout: 1, DestPort: 0},
		{Type: "tcp", Timeout: 1, DestPort: 65536},
		{Type: "dns", Timeout: 1, Dns_Protocol: "invalid"},
	}
	for _, input := range tests {
		if err := validateInputVars(input); err == nil {
			t.Fatalf("expected validation error for %+v", input)
		}
	}
}
