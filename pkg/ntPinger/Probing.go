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
func tcpProbingRun (p *Pinger) {

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// idx
	idx := 0

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
				close(p.probeChan)
				forLoopEnds = true
			default:
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.InputVars.Timeout)*time.Second)
				defer cancel()

				pkt, err := tcpProbing(ctx, idx, p.destAddr, p.InputVars.DestPort, 0)

				if err != nil {
					if strings.Contains(err.Error(), "timeout") {
						// Probe Timeout
					} else {
						panic(err)
					}
				} else {
					// Probe Success
				}
				p.UpdateStatistics(&pkt)
				pkt.UpdateStatistics(p.Stat)
				p.probeChan <- &pkt
				idx++

				// sleep for interval
				time.Sleep(GetSleepTime(pkt.Status, p.InputVars.Interval))
			}

		}

	} else {
		for i := 0; i < p.InputVars.Count; i++ {
			if forLoopEnds {
				break
			}
			select {
			case <-interruptChan:
				close(p.probeChan) // case interruptChan, close the channel & break the loop
				forLoopEnds = true
			default:
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.InputVars.Timeout)*time.Second)
				defer cancel()

				pkt, err := tcpProbing(ctx, idx, p.destAddr, p.InputVars.DestPort, 0)

				if err != nil {
					if strings.Contains(err.Error(), "timeout") {
						// Probe Timeout
					} else {
						panic(err)
					}
				} else {
					// Probe Success
				}
				p.UpdateStatistics(&pkt)
				pkt.UpdateStatistics(p.Stat)
				p.probeChan <- &pkt
				idx++

				// check the last loop of the probing, close probeChan
				if i == (p.InputVars.Count - 1) {
					close(p.probeChan)
				}

				// sleep for interval
				time.Sleep(GetSleepTime(pkt.Status, p.InputVars.Interval))
			}
		}
	}
}


// func: tcpProbing
func tcpProbing(ctx context.Context, idx int, destAddr string, destPort int, nbytes int) (PacketTCP, error) {

	// initial packet
	pkt := PacketTCP{
		Type:     "tcp",
		Idx:      idx,
		DestAddr: destAddr,
		DestPort: destPort,
		NBytes:   nbytes,
	}

	// setup Dialer
	d := net.Dialer{}

	// Record the start time
	pkt.SendTime = time.Now()

	// Ping Target
	pingTarget := fmt.Sprintf("%s:%d", destAddr, destPort)

	// Establish a connection with a context timeout
	conn, err := d.DialContext(ctx, pkt.Type, pingTarget)
	if err != nil {
		pkt.Status = false
		return pkt, err
	}
	defer conn.Close()

	// Create a packet of the desired size
	if nbytes != 0 {
		packetPayload := make([]byte, nbytes)

		// Send the packet
		_, err = conn.Write(packetPayload)
		if err != nil {
			return pkt, err
		}
	}

	// Record the Status
	pkt.Status = true

	// Calculate the RTT
	pkt.RTT = time.Since(pkt.SendTime)

	return pkt, nil
}
