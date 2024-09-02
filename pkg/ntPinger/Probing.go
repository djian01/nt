package ntPinger

import (
	"context"
	"fmt"
	"net"
	"time"
)

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
