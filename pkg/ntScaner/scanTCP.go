package ntScaner

import (
	"context"
	"sync"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
)

var tcpProbe = ntPinger.TcpProbing

// TcpScanResult is the completed result for one scan job.
type TcpScanResult struct {
	Port TcpScanPort
	Err  error
}

// ScanTcpWorkerResults tests each port and publishes one immutable result.
func ScanTcpWorkerResults(testPortChan <-chan TcpScanPort, resultChan chan<- TcpScanResult) {
	for testPort := range testPortChan {
		resultChan <- scanTCPPort(testPort)
	}
}

// ScanTcpWorker preserves the original pointer-based API for existing callers.
// New code should use ScanTcpWorkerResults so one goroutine owns shared tables.
func ScanTcpWorker(testPortChan <-chan *TcpScanPort, errChan chan<- error) {
	for testPort := range testPortChan {
		result := scanTCPPort(*testPort)
		if result.Err != nil {
			errChan <- result.Err
			continue
		}
		*testPort = result.Port
	}
}

func scanTCPPort(testPort TcpScanPort) TcpScanResult {
	type attemptResult struct {
		packet ntPinger.PacketTCP
		err    error
	}

	attempts := make(chan attemptResult, 3)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(testPort.Timeout)*time.Second)
			defer cancel()
			pkt, err := tcpProbe(&ctx, 1, testPort.DestAddr, testPort.DestHost, testPort.Port, 0)
			attempts <- attemptResult{packet: pkt, err: err}
		}()
	}

	wg.Wait()
	close(attempts)

	var firstErr error
	for attempt := range attempts {
		if attempt.err != nil {
			if firstErr == nil {
				firstErr = attempt.err
			}
			continue
		}
		if attempt.packet.Status {
			testPort.Status = 2
			testPort.AdditionalInfo = attempt.packet.AdditionalInfo
			return TcpScanResult{Port: testPort}
		}
		if testPort.AdditionalInfo == "" {
			testPort.AdditionalInfo = attempt.packet.AdditionalInfo
		}
	}

	if firstErr != nil {
		return TcpScanResult{Port: testPort, Err: firstErr}
	}
	testPort.Status = 3
	return TcpScanResult{Port: testPort}
}

// TcpScanStat returns completion statistics for the current scan table.
func TcpScanStat(ports *[]TcpScanPort) (countTested int, countSuccess int, countFail int) {
	for _, port := range *ports {
		switch port.Status {
		case 2:
			countSuccess++
			countTested++
		case 3:
			countFail++
			countTested++
		}
	}
	return
}
