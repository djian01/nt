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
				pkt, err := dnsProbing(Seq, p.InputVars.DestHost, p.InputVars.Dns_queryType, p.InputVars.Timeout, p.InputVars.DestHost)
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
				pkt, err := dnsProbing(Seq, p.InputVars.DestHost, p.InputVars.Dns_queryType, p.InputVars.Timeout, p.InputVars.DestHost)
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

// func: dnsProbing
func dnsProbing(Seq int, destHost string, Dns_queryType string, timeout int, dnsResolver string) (PacketDNS, error) {

	// Initial PacketDNS
	pkt := PacketDNS{
		Type:          "dns",
		Status:        false,
		Seq:           Seq,
		DestHost:      destHost,
		Dns_queryType: Dns_queryType,
	}

	// Set up the DNS resolver
	resolver := &net.Resolver{}
	if dnsResolver != "" {
		// Custom DNS resolver
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial(network, fmt.Sprintf("%s:53", dnsResolver))
			},
		}
	}

	// Record the start time
	pkt.SendTime = time.Now()

	var err error

	// Perform DNS query based on the query type
	switch strings.ToUpper(Dns_queryType) {
	case "A":
		_, err = resolver.LookupHost(context.Background(), destHost)
	case "AAAA":
		_, err = resolver.LookupIPAddr(context.Background(), destHost)
	case "CNAME":
		_, err = resolver.LookupCNAME(context.Background(), destHost)
	default:
		err = fmt.Errorf("unsupported query type: %s", Dns_queryType)
	}

	// Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			pkt.AdditionalInfo = "No_Such_Host"
		} else if strings.Contains(err.Error(), "timeout") {
			pkt.AdditionalInfo = "Timeout"
		} else {
			return pkt, fmt.Errorf("dns query failed: %w", err)
		}
		pkt.Status = false
		return pkt, nil
	}

	// Mark packet as successful
	pkt.Status = true

	// Check Latency
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil
}
