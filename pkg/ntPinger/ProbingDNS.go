package ntPinger

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// func: dnsProbingRun
func dnsProbingRun(p *Pinger, errChan chan<- error) {
	runProbes(p, errChan, func(seq int) (Packet, error) {
		pkt, err := DnsProbing(seq, p.InputVars.DestHost, p.InputVars.Dns_query, p.InputVars.Dns_Protocol, p.InputVars.Timeout)
		return &pkt, err
	})
}

// func DnsProbing
func DnsProbing(Seq int, destHost string, Dns_query string, Dns_Protocol string, timeout int) (PacketDNS, error) {

	// Initial PacketDNS
	pkt := PacketDNS{
		Type:         "dns",
		Status:       false,
		Seq:          Seq,
		DestHost:     destHost,
		Dns_query:    Dns_query,
		Dns_protocol: Dns_Protocol,
	}

	// Set up the DNS resolver with custom protocol (TCP or UDP)
	resolver := &net.Resolver{
		// Custom DNS resolver
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// Use the protocol ("tcp" or "udp") to establish the connection
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, Dns_Protocol, net.JoinHostPort(destHost, "53"))
		},
	}

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Record the start time
	pkt.SendTime = time.Now()

	// LookupHost within the context
	var err error
	Dns_response_slice, err := resolver.LookupHost(ctx, Dns_query)

	// Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// DNS response
	pkt.Dns_response = IPSlideToString(Dns_response_slice)

	// check timeout
	if pkt.RTT >= time.Duration(timeout)*time.Second {
		err = errors.New("timeout")
	}

	// err check
	var cname string
	var cnameErr error

	if err == nil {
		cname, cnameErr = resolver.LookupCNAME(ctx, Dns_query)
		pkt.Status = true

	} else {
		// If neither "A" nor "CNAME" is found, handle the error
		// Capture specific error message in AdditionalInfo
		switch {
		case strings.Contains(err.Error(), "no such host"):
			pkt.AdditionalInfo = "Non_Existent_Domain"
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

	// CNAME check
	if err == nil {
		if cnameErr == nil {
			cname = strings.TrimSuffix(cname, ".")
		}

		if cname != Dns_query {
			// If CNAME lookup succeeds and the CNAME is different from the query
			pkt.Dns_queryType = "CNAME"
		} else {
			pkt.Dns_queryType = "A"
		}
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
