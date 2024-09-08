package ntPinger

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// func: tcpProbingRun
func tcpProbingRun(p *Pinger, errChan chan<- error) {

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Sequence
	Seq := 0

	// forLoopEnds Flag
	forLoopEnds := false

	// count
	if p.InputVars.Count == 0 {
		for {
			if forLoopEnds {
				break
			}

			select {
			case <-interruptChan: // case interruptChan, close the channel & break the loop
				close(p.ProbeChan)
				forLoopEnds = true
			default:
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.InputVars.Timeout)*time.Second)
				defer cancel()

				pkt, err := tcpProbing(&ctx, Seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.DestPort, p.InputVars.PayLoadSize)
				if err != nil {
					errChan <- err
				}

				p.UpdateStatistics(&pkt)
				pkt.UpdateStatistics(p.Stat)
				p.ProbeChan <- &pkt
				Seq++

				// sleep for interval
				time.Sleep(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT))
			}

		}

	} else {
		for i := 0; i < p.InputVars.Count; i++ {
			if forLoopEnds {
				break
			}
			select {
			case <-interruptChan:
				close(p.ProbeChan) // case interruptChan, close the channel & break the loop
				forLoopEnds = true
			default:
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.InputVars.Timeout)*time.Second)
				defer cancel()

				pkt, err := tcpProbing(&ctx, Seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.DestPort, p.InputVars.PayLoadSize)
				if err != nil {
					errChan <- err
				}

				p.UpdateStatistics(&pkt)
				pkt.UpdateStatistics(p.Stat)
				p.ProbeChan <- &pkt
				Seq++

				// check the last loop of the probing, close probeChan
				if i == (p.InputVars.Count - 1) {
					close(p.ProbeChan)
				}

				// sleep for interval
				time.Sleep(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT))
			}
		}
	}
}

// func: tcpProbing
func tcpProbing(ctx *context.Context, Seq int, destAddr string, desetHost string, destPort int, PayLoadSize int) (PacketTCP, error) {

	// initial packet
	pkt := PacketTCP{
		Type:        "tcp",
		Seq:         Seq,
		DestHost:    desetHost,
		DestAddr:    destAddr,
		DestPort:    destPort,
		PayLoadSize: PayLoadSize,
	}

	// setup Dialer
	d := net.Dialer{}

	// Record the start time
	pkt.SendTime = time.Now()

	// Ping Target
	pingTarget := fmt.Sprintf("%s:%d", destAddr, destPort)

	// Establish a connection with a context timeout - 3-way handshake
	conn, err := d.DialContext(*ctx, pkt.Type, pingTarget)
	if err != nil {
		pkt.Status = false

		switch {
			// Error: "connection refused"
		case strings.Contains(err.Error(), "refused"):
			// Calculate the RTT
			pkt.RTT = time.Since(pkt.SendTime)
			// Add Info
			pkt.AdditionalInfo = "Conn_Refused"
			return pkt, nil
	
			// Error: "no route to host"
		case strings.Contains(err.Error(), "no route"):
			// Add Info
			pkt.AdditionalInfo = "No_Route"
			return pkt, nil
	
			// Error: "timeout"
		case strings.Contains(err.Error(), "timeout"):
			// Add Info
			pkt.AdditionalInfo = "Conn_Timeout"
			return pkt, nil
	
			// Error: Else
		default:
			return pkt, fmt.Errorf("conn error: %w", err)
		}			
		// // Error: "connection refused"
		// if strings.Contains(err.Error(), "refused") {
		// 	// Calculate the RTT
		// 	pkt.RTT = time.Since(pkt.SendTime)
		// 	// Add Info
		// 	pkt.AdditionalInfo = "Conn_Refused"
		// 	return pkt, nil

		// 	// Error: "no route to host"
		// } else if strings.Contains(err.Error(), "no route") {
		// 	// Add Info
		// 	pkt.AdditionalInfo = "No_Route"
		// 	return pkt, nil

		// 	// Error: "timeout"
		// } else if strings.Contains(err.Error(), "timeout") {
		// 	// Add Info
		// 	pkt.AdditionalInfo = "Conn_Timeout"
		// 	return pkt, nil

		// 	// Error: Else
		// } else {
		// 	return pkt, fmt.Errorf("conn error: %w", err)
		// }
	}
	defer conn.Close()

	// create and send payload if required
	if PayLoadSize != 0 {
		packetPayload := make([]byte, PayLoadSize)

		// Send the payload
		_, err = conn.Write(packetPayload)
		if err != nil {
			return pkt, fmt.Errorf("conn write error: %w", err)
		}
	}

	// Record the Status
	pkt.Status = true

	// Calculate the RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// Check Latency
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil
}
