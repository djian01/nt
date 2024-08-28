package icmp

import (
	"fmt"
	"net"
	"time"

	"nt/pkg/sharedStruct"

	probing "github.com/prometheus-community/pro-bing"
)

// func - probingFunc
func IcmpProbingFunc(pingHost string, count, size, interval int, probingChan chan<- sharedStruct.NtResult, doneChan chan<- bool) error {

	// Check Name Resolution
	targetIP, err := net.LookupIP(pingHost)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to resolve domain: %v", pingHost))
	}

	// *********  Setup pinger ************
	// Create a new ping instance for the given host
	pinger, err := probing.NewPinger(targetIP[0].String())
	if err != nil {
		return err
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

	// ********** func - pinger.OnSend ***********
	processingPkgId := -1
	pinger.OnSend = func(pkt *probing.Packet) {

		// if pkt.Seg != (processingPkgId + 1), the processingPkgId packet is missing
		if pkt.Seq != (processingPkgId + 1) {
			//fmt.Printf("Ping error: Missing Ping reply Packet ID: %v from %v\n", (processingPkgId + 1), pingHost)

			// return result
			ntr := sharedStruct.NtResult{
				HostName:  pingHost,
				IP:        pkt.Addr,
				Status:    "ICMP_Failed",
				RTT:       pkt.Rtt,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				Seq:       pkt.Seq,
				Type:      "icmp",
			}
			probingChan <- ntr

			// reset the processingPkgId
			processingPkgId = pkt.Seq - 1

			// last ping packet for the given count
			if pkt.Seq == (count - 1) {
				go func() {
					// wait for the last interval
					time.Sleep(time.Duration(interval) * time.Second)

					// if the processingPkgId != the last packet seq
					if processingPkgId != pkt.Seq {

						// return result
						ntr := sharedStruct.NtResult{
							HostName:  pingHost,
							IP:        pkt.Addr,
							Status:    "ICMP_Failed",
							RTT:       pkt.Rtt,
							Timestamp: time.Now().Format("2006-01-02 15:04:05"),
							Seq:       pkt.Seq,
							Type:      "icmp",
						}
						probingChan <- ntr
					}
				}()
			}
		}
	}

	// ********** func - pinger.OnRecv ***********
	pinger.OnRecv = func(pkt *probing.Packet) {

		// update processingPkgId
		processingPkgId = pkt.Seq

		// return result
		stat := pinger.Statistics()
		ntr := sharedStruct.NtResult{
			Seq:         pkt.Seq,
			Type:        "icmp",
			Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
			HostName:    pingHost,
			IP:          pkt.Addr,
			Size:        pkt.Nbytes,
			Status:      "ICMP_OK",
			RTT:         pkt.Rtt,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
			PacketLoss:  stat.PacketLoss,
			MinRtt:      stat.MinRtt,
			MaxRtt:      stat.MaxRtt,
			AvgRtt:      stat.AvgRtt,
		}
		probingChan <- ntr
	}

	// ********** Go Routine to run pinger ***********

	// Go Routine - Monitor count is completed
	go func() {
		err = pinger.Run()
		if err != nil {
			panic(err)
		}
		// doneChan: true when count is completed
		doneChan <- true
	}()

	return nil
}
