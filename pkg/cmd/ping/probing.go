package ping

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nt/pkg/sharedStruct"

	probing "github.com/prometheus-community/pro-bing"
)

// func - probingFunc
func ProbingFunc(pingHost string, count, size, interval int, ProbingChain chan sharedStruct.NtResult) {

	// Create a new ping instance for the given host
	pinger, err := probing.NewPinger(pingHost)
	if err != nil {
		panic(err)
	}

	// Use the unprivileged mode
	pinger.SetPrivileged(true)

	// Set the number of packets to send
	if count != 0 {
		pinger.Count = count
	}

	// Set the size of the packet
	if size > 24 {
		pinger.Size = size
	}

	// Set the interval between each packet
	if interval != 1 {
		pinger.Interval = time.Duration(interval) * time.Second
	}

	// Create a channel to signal when pinger.Run() is done
	doneChannel := make(chan bool, 1)
	defer close(doneChannel)

	// missing receive packet
	processingPkgId := -1
	pinger.OnSend = func(pkt *probing.Packet) {

		// if pkt.Seg != (processingPkgId + 1), the processingPkgId packet is missing
		if pkt.Seq != (processingPkgId + 1) {
			//fmt.Printf("Ping error: Missing Ping reply Packet ID: %v from %v\n", (processingPkgId + 1), pingHost)

			// return result
			ntr := sharedStruct.NtResult{
				HostName:  pingHost,
				IP:        pkt.Addr,
				Status:    "PING_Failed",
				RTT:       pkt.Rtt,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				Seq:       pkt.Seq,
			}
			ProbingChain <- ntr

			// reset the processingPkgId
			processingPkgId = pkt.Seq - 1

			// last ping packet for the given count
			if pkt.Seq == (count - 1) {
				go func() {
					// wait for the last interval
					time.Sleep(time.Duration(interval) * time.Second)
					// if the processingPkgId != the last packet seq
					if processingPkgId != pkt.Seq {
						// fmt.Printf("Ping error: Missing Ping reply Packet ID: %v from %v\n", pkt.Seq, pingHost)

						// return result
						ntr := sharedStruct.NtResult{
							HostName:  pingHost,
							IP:        pkt.Addr,
							Status:    "PING_Failed",
							RTT:       pkt.Rtt,
							Timestamp: time.Now().Format("2006-01-02 15:04:05"),
							Seq:       pkt.Seq,
						}
						ProbingChain <- ntr

						// close the doneChannel
						doneChannel <- true
					}
				}()
			}
		}
	}

	// Run the ping
	// fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.OnRecv = func(pkt *probing.Packet) {
		// update processingPkgId
		processingPkgId = pkt.Seq

		// print receive packet info
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)

		// return result
		stat := pinger.Statistics()
		ntr := sharedStruct.NtResult{
			Seq:         pkt.Seq,
			Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
			HostName:    pingHost,
			IP:          pkt.Addr,
			Size:        pkt.Nbytes,
			Status:      "PING_OK",
			RTT:         pkt.Rtt,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
			PacketLoss:  stat.PacketLoss,
			MinRtt:      stat.MinRtt,
			MaxRtt:      stat.MaxRtt,
			AvgRtt:      stat.AvgRtt,
		}
		ProbingChain <- ntr
	}

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChannel := make(chan os.Signal, 1)
	defer close(interruptChannel)

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run the pinger in a separate goroutine
	go func() {
		err = pinger.Run()
		if err != nil {
			panic(err)
		}
		doneChannel <- true
	}()

	// Wait for either the completion of the ping or an interrupt signal
	select {
	case <-doneChannel:
		// Ping count completed
	case <-interruptChannel:
		// Ctrl+C pressed
		fmt.Printf("\033[%d;1H", 22)
		fmt.Println("\n--- Interrupt received, stopping ping ---")
	}

	// Get the statistics of the ping
	// stats := pinger.Statistics()
	// fmt.Printf("\n--- %s ping statistics ---\n", pinger.Addr())
	// fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
	// 	stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	// fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
	// 	stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

}
