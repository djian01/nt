package customPinger

import (
	"net"
	"sync"
	"time"
)

// ***************** type Structs ********************
type Packet struct {
	Type          string
	idx           int
	addr          net.Addr
	http_url      string
	http_tls      bool
	http_response string
	dns_request   string
	dns_response  string
	port          int
	nbytes        int
	sendTime      time.Time
	receiveTime   time.Time
	Rtt           time.Duration
}

type Pinger struct {

	// Pinger Type (TCP, HTTP, DNS, ICMP)
	Type string

	// Count tells pinger to stop after sending (and receiving) Count echo
	// packets. If this option is not specified, pinger will operate until
	// interrupted. Default is 0 which means nonstop till interruption.
	Count int

	// Interval is the wait time between each packet send. Default is 1s.
	Interval time.Duration

	// statistics
	Stat Statistics

	// statistics Mutex
	statsMu sync.RWMutex

	// Size of packet being sent
	Size int

	// Timeoiut value for ping test
	Timeout time.Duration

	// Source is the source IP address
	Source string

	// destination IP details
	addr   string
	ipaddr net.IP

	// destination host
	destHost string

	// OnSetup is called when Pinger has finished setting up the listening socket
	OnSetup func()

	// OnSend is called when Pinger sends a packet
	OnSend func(*Packet)

	// OnRecv is called when Pinger receives and processes a packet
	OnRecv func(*Packet)

	// OnFinish is called when Pinger exits
	OnFinish func(*Statistics)

	// OnSendError is called when an error occurs while Pinger attempts to send a packet
	OnSendError func(*Packet, error)

	// OnRecvError is called when an error occurs while Pinger attempts to receive a packet
	OnRecvError func(error)

	// Channel and mutex used to communicate when the Pinger should stop between goroutines.
	done chan bool
	lock sync.Mutex
}

// Statistics represent the stats of a currently running or finished
// pinger operation.
type Statistics struct {
	// PacketsRecv is the number of packets received.
	PacketsRecv int

	// PacketsSent is the number of packets sent.
	PacketsSent int

	// PacketLoss is the percentage of packets lost.
	PacketLoss float64

	// IPAddr is the address of the host being pinged.
	IPAddr net.IP

	// Addr is the string address of the host being pinged.
	Addr string

	// MinRtt is the minimum round-trip time sent via this pinger.
	MinRtt time.Duration

	// MaxRtt is the maximum round-trip time sent via this pinger.
	MaxRtt time.Duration

	// AvgRtt is the average round-trip time sent via this pinger.
	AvgRtt time.Duration
}
