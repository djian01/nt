package ntPinger

import (
	"crypto/tls"
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
			// Loop End Signal
			if forLoopEnds {
				break
			}

			// Pinger end Singal
			if p.PingerEnd {
				interruptChan <- os.Interrupt //send interrupt to interruptChan
			}

			// Perform HTTP probing
			pkt, err := HttpProbing(Seq, p.InputVars.DestHost, p.InputVars.DestPort, p.InputVars.Http_path, p.InputVars.Http_scheme, p.InputVars.Http_method, p.InputVars.Http_statusCodes, p.InputVars.Timeout)
			if err != nil {
				errChan <- err
			}

			// Update statistics and send packet
			p.UpdateStatistics(&pkt)
			pkt.UpdateStatistics(p.Stat)
			p.ProbeChan <- &pkt
			Seq++

			// sleep for interval
			select {
			case <-time.After(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT)):
				// wait for SleepTime
			case <-interruptChan: // case interruptChan, close the channel & break the loop
				forLoopEnds = true
				close(p.ProbeChan)
			}
		}

	} else {
		for i := 0; i < p.InputVars.Count; i++ {

			if forLoopEnds {
				break
			}

			// Perform HTTP probing
			pkt, err := HttpProbing(Seq, p.InputVars.DestHost, p.InputVars.DestPort, p.InputVars.Http_path, p.InputVars.Http_scheme, p.InputVars.Http_method, p.InputVars.Http_statusCodes, p.InputVars.Timeout)
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
			select {
			case <-time.After(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT)):
				// wait for SleepTime
			case <-interruptChan: // case interruptChan, close the channel & break the loop
				forLoopEnds = true
				close(p.ProbeChan)
			}
		}
	}
}

// HttpProbing performs HTTP probing with the ability to choose HTTP methods and ignore certificate errors
func HttpProbing(Seq int, destHost string, destPort int, Http_path string, Http_scheme string, Http_method string, Http_statusCodes []HttpStatusCode, timeout int) (PacketHTTP, error) {

	// Initial PacketHTTP
	pkt := PacketHTTP{
		Type:        "http",
		Status:      false,
		Seq:         Seq,
		DestHost:    destHost,
		DestPort:    destPort,
		Http_scheme: Http_scheme,
		Http_method: Http_method,
		Http_path:   Http_path,
	}

	// Construct the URL for HTTP or HTTPS
	url := ConstructURL(Http_scheme, destHost, Http_path, destPort)

	// Create a custom Transport to ignore certificate errors
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a new HTTP client with the custom Transport and timeout
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
	}

	// Record the start time
	pkt.SendTime = time.Now()

	// Create a new request based on the specified HTTP method
	req, err := http.NewRequest(Http_method, url, nil)
	if err != nil {
		return pkt, fmt.Errorf("failed to create request: %w", err)
	}

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "context deadline exceeded"):
			// Timeout error
			pkt.AdditionalInfo = "Conn_Timeout"
			return pkt, nil

		case strings.Contains(err.Error(), "timeout"):
			// Timeout error
			pkt.AdditionalInfo = "Conn_Timeout"
			return pkt, nil

		case strings.Contains(err.Error(), "handshake timeout"):
			// Connection handshake timeout
			pkt.AdditionalInfo = "Handshake_Timeout"
			return pkt, nil

		case strings.Contains(err.Error(), "network is unreachable"):
			// Connection network is unreachable
			pkt.AdditionalInfo = "Network_Unreachable"
			return pkt, nil

		case strings.Contains(err.Error(), "connection refused"):
			// Connection refused error
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "Conn_Refused"
			return pkt, nil

		case strings.Contains(err.Error(), "connection reset by peer"):
			// Connection reset by peer
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "Conn_Reset_by_Peer"
			return pkt, nil

		case strings.Contains(err.Error(), "EOF"):
			// EOF
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "EOF"
			return pkt, nil

		default:
			return pkt, fmt.Errorf("failed to send request: %w", err)
		}
	}
	defer resp.Body.Close()

	// Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// Fill in the HTTP response code & response phase
	pkt.Http_response_code = resp.StatusCode
	if parts := strings.SplitN(resp.Status, " ", 2); len(parts) == 2 {
		pkt.Http_response = parts[1]
	} else {
		pkt.Http_response = resp.Status
	}

	// Status decision based on allowed ranges
	pkt.Status = statusAllowed(resp.StatusCode, Http_statusCodes)
	if !pkt.Status {
		if pkt.AdditionalInfo == "" {
			pkt.AdditionalInfo = "StatusNotAllowed"
		}
	}

	// Check for latency (assuming CheckLatency is a separate function)
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil
}

// statusAllowed returns true if code falls within any allowed [Lower,Upper] range.
// If no ranges are provided, it falls back to treating 2xx/3xx as success (to match old behavior).
func statusAllowed(code int, allowed []HttpStatusCode) bool {
	if len(allowed) == 0 {
		return code >= 200 && code < 400
	}
	for _, r := range allowed {
		low, high := r.LowerCode, r.UpperCode
		if low > high {
			low, high = high, low
		}
		if code >= low && code <= high {
			return true
		}
	}
	return false
}
