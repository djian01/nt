package ntTEST

import (
	"math/rand"
	"nt/pkg/ntPinger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ResultGenerate(count int, Type string, probeChan *chan ntPinger.Packet) {

	// statistics
	var PacketsSent = 0
	var PacketsRecv = 0
	var PacketLoss float64
	var MinRtt time.Duration
	var AvgRtt time.Duration
	var MaxRtt time.Duration
	var status bool

	// random Source
	source := rand.NewSource(time.Now().UnixNano())

	// random Error Seed
	errorMax := 12
	errorMin := 2
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

			// PacketSent ++
			PacketsSent++

			// check Interrpution
			select {
			case <-interruptChan:
				// if doneChan <- true if interrupted
				forLoopFlag = false
			default:
				// pass
			}

			// error check
			if PacketsSent%errorSeed != 0 {
				PacketsRecv++
				status = true
			} else {
				status = false
			}

			min := 200
			max := 700
			ranRTT := time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

			// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
			if PacketsSent == 1 {
				ranRTT = time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond
				MinRtt = ranRTT
				MaxRtt = ranRTT
				AvgRtt = ranRTT
				// after Initialization
			} else if status {
				// generate RTT
				ranRTT = time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

				// RTT Statistics
				if ranRTT > MaxRtt {
					MaxRtt = ranRTT
				} else if ranRTT < MinRtt {
					MinRtt = ranRTT
				}
				AvgRtt = (AvgRtt*time.Duration(PacketsRecv-1) + ranRTT) / time.Duration(PacketsRecv)

			}

			// generate statistic
			PacketLoss = (1 - float64(PacketsRecv)/float64(PacketsSent))

			switch Type {
			case "icmp":
				probeResult := ntPinger.PacketICMP{
					Seq:      PacketsSent - 1,
					Type:     Type,
					DestHost: "google.com",
					DestAddr: "1.2.3.4",
					NBytes:   56,
					Status:   status,
					RTT:      ranRTT,
					SendTime: time.Now(),

					PacketsSent: PacketsSent,
					PacketsRecv: PacketsRecv,
					PacketLoss:  PacketLoss,
					MinRtt:      MinRtt,
					MaxRtt:      MaxRtt,
					AvgRtt:      AvgRtt,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "tcp":
				destPort := 443

				probeResult := ntPinger.PacketTCP{
					Seq:      PacketsSent - 1,
					Type:     Type,
					DestHost: "google.com",
					DestAddr: "1.2.3.4",
					NBytes:   56,
					Status:   status,
					RTT:      ranRTT,
					SendTime: time.Now(),
					DestPort: destPort,

					PacketsSent: PacketsSent,
					PacketsRecv: PacketsRecv,
					PacketLoss:  PacketLoss,
					MinRtt:      MinRtt,
					MaxRtt:      MaxRtt,
					AvgRtt:      AvgRtt,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "http":

			case "dns":

			}
		}

		// code finish
		close(*probeChan)

	} else {

		// forLoopFlag
		forLoopFlag := true

		// create the input NtResult
		for i := 0; i < count; i++ {

			// check forLoopFlag
			if !forLoopFlag {
				break
			}

			// PacketSent ++
			PacketsSent++

			// check Interrpution
			select {
			case <-interruptChan:
				// if doneChan <- true if interrupted
				forLoopFlag = false

			default:
				// pass
			}

			// error check
			if PacketsSent%errorSeed != 0 {
				PacketsRecv++
				status = true
			} else {
				status = false
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

			switch Type {
			case "icmp":
				probeResult := ntPinger.PacketICMP{
					Seq:      PacketsSent - 1,
					Type:     Type,
					DestHost: "google.com",
					DestAddr: "1.2.3.4",
					NBytes:   56,
					Status:   status,
					RTT:      ranRTT,
					SendTime: time.Now(),

					PacketsSent: PacketsSent,
					PacketsRecv: PacketsRecv,
					PacketLoss:  PacketLoss,
					MinRtt:      MinRtt,
					MaxRtt:      MaxRtt,
					AvgRtt:      AvgRtt,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "tcp":
				destPort := 443

				probeResult := ntPinger.PacketTCP{
					Seq:      PacketsSent - 1,
					Type:     Type,
					DestHost: "google.com",
					DestAddr: "1.2.3.4",
					NBytes:   56,
					Status:   status,
					RTT:      ranRTT,
					SendTime: time.Now(),
					DestPort: destPort,

					PacketsSent: PacketsSent,
					PacketsRecv: PacketsRecv,
					PacketLoss:  PacketLoss,
					MinRtt:      MinRtt,
					MaxRtt:      MaxRtt,
					AvgRtt:      AvgRtt,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "http":

			case "dns":
			}
		}

		// code finish
		close(*probeChan)
	}
}
