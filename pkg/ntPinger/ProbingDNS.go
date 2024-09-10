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

// func: dnsProbingRun
func dnsProbingRun(p *Pinger, errChan chan<- error) {

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
				// Perform DNS probing
				pkt, err := DnsProbing(Seq, p.InputVars.DestHost, p.InputVars.Dns_query, p.InputVars.Dns_Protocol, p.InputVars.Timeout)
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
				// Perform DNS probing
				pkt, err := DnsProbing(Seq, p.InputVars.DestHost, p.InputVars.Dns_query, p.InputVars.Dns_Protocol, p.InputVars.Timeout)
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

func DnsProbing(Seq int, destHost string, Dns_query string, Dns_Protocol string, timeout int) (PacketDNS, error) {

	// Initial PacketDNS
	pkt := PacketDNS{
		Type:      "dns",
		Status:    false,
		Seq:       Seq,
		DestHost:  destHost,
		Dns_query: Dns_query,
	}

	// Set up the DNS resolver with custom protocol (TCP or UDP)
	resolver := &net.Resolver{}

	if destHost != "" {
		// Custom DNS resolver
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				// Use the protocol ("tcp" or "udp") to establish the connection
				return net.Dial(Dns_Protocol, fmt.Sprintf("%s:53", destHost))
			},
		}
	}

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Record the start time
	pkt.SendTime = time.Now()

	// Attempt to resolve "A" record first within the context
	var err error
	Dns_response_slice, err := resolver.LookupHost(ctx, Dns_query)
	pkt.Dns_response = IPSlideToString(Dns_response_slice)

	if err == nil {
		pkt.Status = true
		pkt.Dns_queryType = "A"

	} else {
		// If "A" record is not found, attempt "CNAME" within the context
		pkt.Dns_response, err = resolver.LookupCNAME(ctx, Dns_query)
		if err == nil {
			pkt.Status = true
			pkt.Dns_queryType = "CNAME"
			//pkt.AdditionalInfo = fmt.Sprintf("CNAME resolves to: %s", cname)
		} else {
			// If neither "A" nor "CNAME" is found, handle the error
			// Capture specific error message in AdditionalInfo
			switch {
			case strings.Contains(err.Error(), "no such host"):
				pkt.AdditionalInfo = "No_Such_Host"
			case strings.Contains(err.Error(), "timeout"):
				pkt.AdditionalInfo = "Timeout"
			case strings.Contains(err.Error(), "temporary failure"):
				pkt.AdditionalInfo = "Temporary_failure"
			case strings.Contains(err.Error(), "query refused"):
				pkt.AdditionalInfo = "DNS_query_refused"
			case strings.Contains(err.Error(), "invalid domain name"):
				pkt.AdditionalInfo = "Invalid_domain_name"
			case strings.Contains(err.Error(), "permission denied"):
				pkt.AdditionalInfo = "Permission_denied"
			default:
				return pkt, fmt.Errorf("dns query failed: %w", err) // return error
			}
		}
	}

	// Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// Check Latency
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil
}

// resonved IPs []string -> string
func IPSlideToString(IPSlide []string) string {

	resultSlide := []string{}

	for _, ip := range IPSlide {
		if isIPv4(ip) {
			resultSlide = append(resultSlide, ip)
		}
	}

	return strings.Join(resultSlide, ",")

}

// isIPv4 checks if a given string is a valid IPv4 address
func isIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false // Not a valid IP address
	}
	return parsedIP.To4() != nil // To4() returns non-nil only for IPv4 addresses
}
