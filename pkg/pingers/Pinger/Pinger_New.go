package tcppinger

import (
	"time"
)

// NewPinger returns a new Pinger and resolves the address.
func NewTcpPinger(addr string) (*TcpPinger, error) {
	p := New(addr)
	return p, p.Resolve()
}

// New returns a new Pinger struct pointer.
func New(addr string) *TcpPinger {

	return &TcpPinger{
		Count:    -1,
		Interval: time.Second,
		Stat: Statistics{
			PacketsSent: 0,
			PacketsRecv: 0,
			PacketLoss:  0,
			IPAddr:      nil,
			Addr:        "",
			MinRtt:      time.Duration(0),
			MaxRtt:      time.Duration(0),
			AvgRtt:      time.Duration(0),
		},
		Size:    24,
		Timeout: time.Duration(10), // default TCP timeout is 10s

		ipaddr: nil,
		addr:   "",

		done: make(chan bool, 1),
	}
}
