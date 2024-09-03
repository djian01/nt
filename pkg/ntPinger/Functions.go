package ntPinger

import (
	"fmt"
	"time"
)

// New returns a new Pinger struct pointer.
func NewPinger(InputVars InputVars) (*Pinger, error) {

	// Inital Pinger
	var p Pinger
	p.InputVars = InputVars

	// Fill in
	switch InputVars.Type {
	case "tcp", "http", "icmp", "dns":

		// resolve the destHost
		err := p.Resolve()
		if err != nil {
			return nil, err
		}

		// initial probeChan
		p.ProbeChan = make(chan Packet, 1)

		// statistic
		p.Stat = Statistics{
			PacketsSent: 0,
			PacketsRecv: 0,
			PacketLoss:  0,
			MinRtt:      time.Duration(1<<63 - 1), // // Set time.Duration to its maximum value
			MaxRtt:      time.Duration(-1 << 63),
			AvgRtt:      time.Duration(0),
		}

	default:
		return nil, fmt.Errorf("please select one of these types: tcp, http, icmp, dns")
	}
	return &p, nil
}

// Func - Get Sleep Time
func GetSleepTime(PacketStatus bool, Interval int) time.Duration {
	if PacketStatus {
		return time.Duration(Interval) * time.Second
	} else {
		return time.Duration(0) * time.Second
	}
}
