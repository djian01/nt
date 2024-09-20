package ntScan

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/record"
	"github.com/djian01/nt/pkg/sharedstruct"
	"github.com/djian01/nt/pkg/terminalOutput"
)

func ScanTcpRun(recording bool, destHost string, Ports []int, timeout int) error {

	// get the total number of the ports
	PortCount := len(Ports)
	countTested := 0

	// resolve destHost
	destAddr := ""

	resolvedIPs, err := ntPinger.ResolveDestHost(destHost)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to resolve domain: %v", destHost))
	}

	// Get the 1st IPv4 IP from resolved IPs
	for _, ip := range resolvedIPs {
		// To4() returns nil if it's not an IPv4 address
		if ip.To4() != nil {
			destAddr = ip.String()
			break
		}
	}

	// if the total number of ports is over 50, error and exit
	if PortCount > 50 {
		return fmt.Errorf("error: total number of input ports is over 50")
	}

	// Create a list of 50 empty TcpScanPort items
	PortsTable := make([]sharedstruct.TcpScanPort, 50)

	for i := 0; i < PortCount; i++ {
		PortsTable[i].ID = i
		PortsTable[i].Port = Ports[i]
		PortsTable[i].Timeout = timeout
		PortsTable[i].DestAddr = destAddr
		PortsTable[i].DestHost = destHost
		PortsTable[i].Status = 1
	}

	// Go Routine - Terminal Output
	go func() {

		destplayIdx := 0

		for {
			// Display the items in a 10×5 table
			terminalOutput.ScanTablePrint(&PortsTable, recording, destplayIdx, destHost)

			time.Sleep(time.Duration(1) * time.Second)

			if countTested == PortCount {
				// display the final output
				terminalOutput.ScanTablePrint(&PortsTable, recording, destplayIdx, destHost)
				break
			}

			destplayIdx++
		}
	}()

	// Create TcpScanPort & errChan
	TcpScanPort := make(chan *sharedstruct.TcpScanPort, 1)
	errChan := make(chan error, 1)

	// create 5 workers
	for i := 0; i < 5; i++ {
		go ScanTcpWorker(TcpScanPort, errChan)
	}

	// go routine - TcpScanPort <- &PortsTable[i]
	go func() {
		for i := 0; i < PortCount; i++ {
			TcpScanPort <- &PortsTable[i]
		}
	}()

	// loop
	loopBreak := false

	for {
		if loopBreak {
			break
		}

		select {
		case err := <-errChan:
			loopBreak = true
			return err
		default:
			countTested, _, _ = TcpScanStat(&PortsTable)
			if countTested == PortCount {
				loopBreak = true
			}
		}
	}

	// sleep and wait for 1 sec
	time.Sleep(time.Duration(1) * time.Second)

	// recording
	if recording {

		// recordingFile Path
		exeFileFolder, err := os.Getwd()
		if err != nil {
			return err
		}
		// recordingFile Name
		timeStamp := time.Now().Format("20060102150405")
		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", "TcpScan", destHost, timeStamp)
		recordingFilePath := filepath.Join(exeFileFolder, recordingFileName)

		record.TcpScan_Recording(recordingFilePath, PortsTable)
	}

	return nil
}

func ScanTcpWorker(TestPortChan <-chan *sharedstruct.TcpScanPort, errChan chan<- error) {

	for TestPort := range TestPortChan {

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

				if pkt.Status {
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
				//fmt.Printf("Port %d succeeded!\n", TestPort.Port)
			} else {
				// Received failed
				mutex.Lock()
				TestPort.Status = 3 // set the status as 3, failed
				mutex.Unlock()
				//fmt.Printf("Port %d Failed!\n", TestPort.Port)
			}
		}
	}
}

// func tcpScanStat
func TcpScanStat(Ports *[]sharedstruct.TcpScanPort) (countTested int, countSuccess int, countFail int) {

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
