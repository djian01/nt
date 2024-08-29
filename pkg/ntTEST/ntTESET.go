package ntTEST

import (
	"fmt"
	"math/rand"
	"nt/pkg/sharedStruct"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ResultGenerate(count int, Type string, NtResultChain chan sharedStruct.NtResult, doneChan chan bool) {

	// statistics
	var PacketsSent = 0
	var PacketsRecv = 0
	var PacketLoss float64
	var MinRtt time.Duration
	var AvgRtt time.Duration
	var MaxRtt time.Duration
	var status string

	// initial loopCount
	loopCount := 0

	// random Source
	source := rand.NewSource(time.Now().UnixNano())

	// random Error Seed
	errorMax := 10
	errorMin := 1
	errorSeed := rand.New(source).Intn(errorMax-errorMin+1) + errorMin

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// endless loop
	if count == 0 {

		// forLoopFlag
		forLoopFlag := true

		// create the input NtResult
		for {

			// check forLoopFlag
			if !forLoopFlag {
				break
			}

			// check Interrpution
			select {
			case <-interruptChan:
				// if doneChan <- true if interrupted
				fmt.Println("\n--- Interrupt received, stopping testing ---")
				forLoopFlag = false
				doneChan <- true
			default:
				// pass
			}

			loopCount++
			PacketsSent++

			// error check
			if PacketsSent%errorSeed != 0 {
				PacketsRecv++
				status = "OK"
			} else {
				status = "NOT_OK"
			}

			min := 200
			max := 700
			ranRTT := time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

			// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
			if loopCount == 1 {
				ranRTT = time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond
				MinRtt = ranRTT
				MaxRtt = ranRTT
				AvgRtt = ranRTT
				// after Initialization
			} else if status == "OK" {
				// generate RTT
				ranRTT = time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

				// RTT Statistics
				if ranRTT > MaxRtt {
					MaxRtt = ranRTT
				} else if ranRTT < MinRtt {
					MinRtt = ranRTT
				}
				AvgRtt = time.Duration(((int64(AvgRtt)/1000000)*(int64(PacketsRecv-1))+(int64(ranRTT)/1000000))/int64(PacketsRecv)) * time.Millisecond
			}

			// generate statistic
			PacketLoss = (1 - float64(PacketsRecv)/float64(PacketsSent))

			NtResult := sharedStruct.NtResult{
				Seq:       loopCount - 1,
				Type:      Type,
				HostName:  "google.com",
				IP:        "1.2.3.4",
				Size:      56,
				Status:    status,
				RTT:       ranRTT,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),

				PacketsSent: PacketsSent,
				PacketsRecv: PacketsRecv,
				PacketLoss:  PacketLoss,
				MinRtt:      MinRtt,
				MaxRtt:      MaxRtt,
				AvgRtt:      AvgRtt,
			}

			NtResultChain <- NtResult
			time.Sleep(1 * time.Second)
		}

	} else {
		// create the input NtResult
		for i := 0; i < count; i++ {

			PacketsSent++

			// error check
			if PacketsSent%errorSeed != 0 {
				PacketsRecv++
				status = "OK"
			} else {
				status = "NOT_OK"
			}

			min := 200
			max := 700
			ranRTT := time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

			// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
			if i == 0 {
				ranRTT := time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond
				MinRtt = ranRTT
				MaxRtt = ranRTT
				AvgRtt = ranRTT
				// after Initialization
			} else {

				// generate RTT
				ranRTT = time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

				if ranRTT > MaxRtt {
					MaxRtt = ranRTT
				} else if ranRTT < MinRtt {
					MinRtt = ranRTT
				}
				AvgRtt = time.Duration(((int64(AvgRtt)/1000000)*(int64(PacketsRecv-1))+(int64(ranRTT)/1000000))/int64(PacketsRecv)) * time.Millisecond
			}

			// generate statistic
			PacketLoss = (1 - float64(PacketsRecv)/float64(PacketsSent))

			NtResult := sharedStruct.NtResult{
				Seq:       i,
				Type:      Type,
				HostName:  "google.com",
				IP:        "1.2.3.4",
				Size:      56,
				Status:    status,
				RTT:       ranRTT,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),

				PacketsSent: PacketsSent,
				PacketsRecv: PacketsRecv,
				PacketLoss:  PacketLoss,
				MinRtt:      MinRtt,
				MaxRtt:      MaxRtt,
				AvgRtt:      AvgRtt,
			}

			NtResultChain <- NtResult
			time.Sleep(1 * time.Second)
		}

		// doneChan = true
		doneChan <- true
	}
}
