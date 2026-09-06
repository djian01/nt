package ntScaner

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
)

func TestScanTCPPortSuccessCannotBeOverwritten(t *testing.T) {
	originalProbe := tcpProbe
	defer func() { tcpProbe = originalProbe }()

	var calls atomic.Int32
	tcpProbe = func(ctx *context.Context, seq int, destAddr, destHost string, destPort, payloadSize int) (ntPinger.PacketTCP, error) {
		call := calls.Add(1)
		return ntPinger.PacketTCP{Status: call == 1, AdditionalInfo: "Conn_Refused"}, nil
	}

	result := scanTCPPort(TcpScanPort{Port: 443, Timeout: 1, DestAddr: "127.0.0.1", DestHost: "localhost"})
	if result.Err != nil {
		t.Fatal(result.Err)
	}
	if result.Port.Status != 2 {
		t.Fatalf("status=%d, want 2", result.Port.Status)
	}
	if calls.Load() != 3 {
		t.Fatalf("probe calls=%d, want 3", calls.Load())
	}
}

func TestScanTCPWorkerResultsExits(t *testing.T) {
	originalProbe := tcpProbe
	defer func() { tcpProbe = originalProbe }()
	tcpProbe = func(ctx *context.Context, seq int, destAddr, destHost string, destPort, payloadSize int) (ntPinger.PacketTCP, error) {
		return ntPinger.PacketTCP{Status: false, AdditionalInfo: "Conn_Refused"}, nil
	}

	jobs := make(chan TcpScanPort, 2)
	results := make(chan TcpScanResult, 2)
	done := make(chan struct{})
	go func() {
		defer close(done)
		ScanTcpWorkerResults(jobs, results)
	}()
	for id := 0; id < 2; id++ {
		jobs <- TcpScanPort{ID: id, Port: 443, Timeout: 1, DestAddr: "127.0.0.1", DestHost: "localhost"}
	}
	close(jobs)

	for i := 0; i < 2; i++ {
		result := <-results
		if result.Err != nil {
			t.Fatal(result.Err)
		}
		if result.Port.Status != 3 {
			t.Fatalf("status=%d, want 3", result.Port.Status)
		}
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("worker did not exit after jobs closed")
	}
}
