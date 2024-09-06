package ntPinger

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// func: httpProbingRun
func httpProbingRun(p *Pinger, errChan chan<- error) {

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Sequence number
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
				// Perform HTTP probing
				pkt, err := HttpProbing(Seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.DestPort, p.InputVars.Timeout, p.InputVars.Http_tls)
				if err != nil {
					errChan <- err
				}

				// Update statistics and send packet
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
				close(p.ProbeChan)
				forLoopEnds = true
			default:
				// Perform HTTP probing
				pkt, err := HttpProbing(Seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.DestPort, p.InputVars.Timeout, p.InputVars.Http_tls)
				if err != nil {
					errChan <- err
				}

				// Update statistics and send packet
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

// func: HttpProbing
func HttpProbing(Seq int, destAddr string, destHost string, destPort int, timeout int, Http_tls bool) (PacketHTTP, error) {

	// Initial PacketHTTP
	pkt := PacketHTTP{
		Type:     "http",
		Status:   false,
		Seq:      Seq,
		DestAddr: destAddr,
		DestHost: destHost,
		DestPort: destPort,
		Http_tls: Http_tls,
	}

	// check http scheme
	scheme := ""
	if Http_tls {
		scheme = "https"
	} else {
		scheme = "http"
	}

	// Construct the URL for HTTP or HTTPS
	url := fmt.Sprintf("%s://%s:%d", scheme, destAddr, destPort)

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Record the start time
	pkt.SendTime = time.Now()

	// Perform the GET request
	resp, err := client.Get(url)
	if err != nil {
		pkt.RTT = time.Since(pkt.SendTime)
		if strings.Contains(err.Error(), "timeout") {
			// Timeout error
			pkt.AdditionalInfo = "Timeout"
			return pkt, nil
		} else if strings.Contains(err.Error(), "connection refused") {
			// Connection refused error
			pkt.AdditionalInfo = "Conn_Refused"
			return pkt, nil
		} else {
			return pkt, fmt.Errorf("failed to send request: %w", err)
		}
	}
	defer resp.Body.Close()

	// Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// Mark packet as successful
	pkt.Status = true

	// Check for latency
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil
}
