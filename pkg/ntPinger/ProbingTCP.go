package ntPinger

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// func: tcpProbingRun
func tcpProbingRun(p *Pinger, errChan chan<- error) {
	runProbes(p, errChan, func(seq int) (Packet, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.InputVars.Timeout)*time.Second)
		defer cancel()
		pkt, err := TcpProbing(&ctx, seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.DestPort, p.InputVars.PayLoadSize)
		return &pkt, err
	})
}

// func: tcpProbing
func TcpProbing(ctx *context.Context, Seq int, destAddr string, desetHost string, destPort int, PayLoadSize int) (PacketTCP, error) {

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

		// Error: "unreachable"
		case strings.Contains(err.Error(), "unreachable"):
			// Add Info
			pkt.AdditionalInfo = "Network_Unreachable"
			return pkt, nil

			// Error: Else
		default:
			return pkt, fmt.Errorf("conn error: %w", err)
		}
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

	return pkt, nil
}
