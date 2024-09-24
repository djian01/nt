package ntScaner

import (
	"context"
	"sync"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
)

func ScanTcpWorker(TestPortChan <-chan *TcpScanPort, errChan chan<- error) {

	for TestPort := range TestPortChan {

		// successFlag
		successFlag := false

		// WaitGroup to wait for all goroutines
		var wg sync.WaitGroup

		// Mutex to safely update the status of the TcpScanPort object
		var mutex sync.Mutex

		// Channel to signal when a testResult is found
		testResultChan := make(chan bool, 1)

		// Run 3 tests concurrently. Any one comes back with success, return sucess immediately.
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(TestPort.Timeout)*time.Second)
				defer cancel()

				var pkt ntPinger.PacketTCP

				pkt, err := ntPinger.TcpProbing(&ctx, 1, TestPort.DestAddr, TestPort.DestAddr, TestPort.Port, 0)
				if err != nil {
					errChan <- err
					return
				}

				// update the AdditionalInfo
				TestPort.AdditionalInfo = pkt.AdditionalInfo

				if pkt.Status {
					successFlag = true
					testResultChan <- true
				} else {
					testResultChan <- false
				}
			}()
		}

		// Wait for either all tests to complete or for one to succeed
		go func() {
			wg.Wait()
			close(testResultChan) // Close the success channel when all tests complete
		}()

		//
		for result := range testResultChan {
			if result {
				// Received success
				mutex.Lock()
				TestPort.Status = 2 // set the status as 2, success
				mutex.Unlock()
			} else {
				// Received failed
				if !successFlag {
					mutex.Lock()
					TestPort.Status = 3 // set the status as 3, failed
					mutex.Unlock()
				}
			}
		}
	}
}

// func tcpScanStat
func TcpScanStat(Ports *[]TcpScanPort) (countTested int, countSuccess int, countFail int) {

	for _, port := range *Ports {
		if port.Status == 2 {
			countSuccess++
			countTested++
		} else if port.Status == 3 {
			countFail++
			countTested++
		}
	}
	return
}
