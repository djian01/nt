package ntPinger

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// CLI Imput variables Struct
type InputVars struct {
	Type        string // Select one of these: tcp, icmp, http, dns
	Count       int    // Default is 0 which means nonstop till interruption.
	PayLoadSize int    // Specific the payload. ICMP default payload is 24 bytes. TCP/HTTP/ICMP have no payload by default.
	Timeout     int    // default timeout is 4 seconds
	Interval    int    // Interval is the wait time between each packet send. Default is 1s.
	SourceHost  string
	DestHost    string
	DestPort    int
	Http_path   string
	Http_scheme string
	Http_method string
	//Icmp_df     bool // ipv4 only
	Dns_query   string
	Dns_queryType string
}

// Packet Interface
type Packet interface {
	GetStatus() bool
	GetRtt() time.Duration
	UpdateStatistics(s Statistics)
	GetType() string
	GetSendTime() time.Time
}

// PacketTCP Struct
type PacketTCP struct {
	Type        string
	Status      bool
	Seq         int
	DestAddr    string
	DestHost    string
	DestPort    int
	PayLoadSize int
	SendTime    time.Time
	RTT         time.Duration
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	AdditionalInfo string
}

func (pkt PacketTCP) GetType() string {
	return pkt.Type
}
func (pkt PacketTCP) GetStatus() bool {
	return pkt.Status
}
func (pkt PacketTCP) GetRtt() time.Duration {
	return pkt.RTT
}
func (pkt PacketTCP) GetSendTime() time.Time {
	return pkt.SendTime
}
func (pkt *PacketTCP) UpdateStatistics(s Statistics) {
	pkt.AvgRtt = s.AvgRtt
	if s.MaxRtt != time.Duration(-1<<63) {
		pkt.MaxRtt = s.MaxRtt
	}
	if s.MinRtt != time.Duration(1<<63-1) {
		pkt.MinRtt = s.MinRtt
	}
	pkt.PacketsSent = s.PacketsSent
	pkt.PacketsRecv = s.PacketsRecv
	pkt.PacketLoss = s.PacketLoss
}

// PacketHTTP Struct
type PacketHTTP struct {
	Type               string
	Status             bool
	Seq                int
	DestAddr           string
	DestHost           string
	DestPort           int
	PayLoadSize        int
	SendTime           time.Time
	RTT                time.Duration
	Http_path          string
	Http_scheme        string
	Http_response_code int
	Http_response      string
	Http_method        string
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	AdditionalInfo string
}

func (pkt PacketHTTP) GetType() string {
	return pkt.Type
}
func (pkt PacketHTTP) GetStatus() bool {
	return pkt.Status
}
func (pkt PacketHTTP) GetRtt() time.Duration {
	return pkt.RTT
}
func (pkt PacketHTTP) GetSendTime() time.Time {
	return pkt.SendTime
}
func (pkt *PacketHTTP) UpdateStatistics(s Statistics) {
	pkt.AvgRtt = s.AvgRtt
	if s.MaxRtt != time.Duration(-1<<63) {
		pkt.MaxRtt = s.MaxRtt
	}
	if s.MinRtt != time.Duration(1<<63-1) {
		pkt.MinRtt = s.MinRtt
	}
	pkt.PacketsSent = s.PacketsSent
	pkt.PacketsRecv = s.PacketsRecv
	pkt.PacketLoss = s.PacketLoss
}

// PacketICMP Struct
type PacketICMP struct {
	Type        string
	Status      bool
	Seq         int
	DestAddr    string
	DestHost    string
	PayLoadSize int
	SendTime    time.Time
	RTT         time.Duration
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	AdditionalInfo string
}

func (pkt PacketICMP) GetType() string {
	return pkt.Type
}
func (pkt PacketICMP) GetStatus() bool {
	return pkt.Status
}
func (pkt PacketICMP) GetRtt() time.Duration {
	return pkt.RTT
}
func (pkt PacketICMP) GetSendTime() time.Time {
	return pkt.SendTime
}
func (pkt *PacketICMP) UpdateStatistics(s Statistics) {
	pkt.AvgRtt = s.AvgRtt
	if s.MaxRtt != time.Duration(-1<<63) {
		pkt.MaxRtt = s.MaxRtt
	}
	if s.MinRtt != time.Duration(1<<63-1) {
		pkt.MinRtt = s.MinRtt
	}
	pkt.PacketsSent = s.PacketsSent
	pkt.PacketsRecv = s.PacketsRecv
	pkt.PacketLoss = s.PacketLoss
}

// PacketDNS Struct
type PacketDNS struct {
	Type          string
	Status        bool
	Seq           int
	DestAddr      string
	DestHost      string
	SendTime      time.Time
	RTT           time.Duration
	Dns_query     string
	Dns_queryType string
	Dns_response  string
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	AdditionalInfo string
}

func (pkt PacketDNS) GetType() string {
	return pkt.Type
}
func (pkt PacketDNS) GetStatus() bool {
	return pkt.Status
}
func (pkt PacketDNS) GetRtt() time.Duration {
	return pkt.RTT
}
func (pkt PacketDNS) GetSendTime() time.Time {
	return pkt.SendTime
}
func (pkt *PacketDNS) UpdateStatistics(s Statistics) {
	pkt.AvgRtt = s.AvgRtt
	if s.MaxRtt != time.Duration(-1<<63) {
		pkt.MaxRtt = s.MaxRtt
	}
	if s.MinRtt != time.Duration(1<<63-1) {
		pkt.MinRtt = s.MinRtt
	}
	pkt.PacketsSent = s.PacketsSent
	pkt.PacketsRecv = s.PacketsRecv
	pkt.PacketLoss = s.PacketLoss
}

// Pinger Struct
type Pinger struct {

	// ************** General Fields **********************
	// Input Vars
	InputVars InputVars

	// statistics
	Stat Statistics

	// statistics Mutex
	StatsMu sync.RWMutex

	// source IP details
	SourceAddr   string
	SourceIpAddr net.IP

	// destination IP details
	DestAddr   string
	DestIpAddr net.IP

	// OnSend is called when Pinger sends a packet (for future use only)
	OnSend func(Packet)

	// OnRecv is called when Pinger receives and processes a packet (for future use only)
	OnRecv func(Packet)

	// probeChan
	ProbeChan chan Packet
}

// Method (Pinger) - Resolve does the DNS lookup for the Pinger address and sets IP protocol.
func (p *Pinger) Resolve() error {
	if len(p.InputVars.DestHost) == 0 {
		return errors.New("destination Host cannot be empty")
	}
	// Check Name Resolution
	resolvedIPs, err := net.LookupIP(p.InputVars.DestHost)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to resolve domain: %v", p.InputVars.DestHost))
	}

	// Get the 1st IPv4 IP from resolved IPs
	for _, ip := range resolvedIPs {
		// To4() returns nil if it's not an IPv4 address
		if ip.To4() != nil {
			p.DestIpAddr = ip
			p.DestAddr = ip.String()
			break
		}
	}
	return nil
}

// Method (Pinger) - Update pinger statistics
func (p *Pinger) UpdateStatistics(pkt Packet) {

	// lock the Statistic
	p.StatsMu.Lock()
	defer p.StatsMu.Unlock()

	// PacketsSent
	p.Stat.PacketsSent++

	// initial PacketLoss
	p.Stat.PacketLoss = 1

	// PacketsRecv
	if pkt.GetStatus() {
		p.Stat.PacketsRecv++

		// PacketLoss
		p.Stat.UpdatePacketLoss()

		// MinRtt
		if p.Stat.MinRtt > pkt.GetRtt() {
			p.Stat.MinRtt = pkt.GetRtt()
		}

		// MaxRtt
		if p.Stat.MaxRtt < pkt.GetRtt() {
			p.Stat.MaxRtt = pkt.GetRtt()
		}

		// AvgRtt
		if p.Stat.AvgRtt == time.Duration(0) {
			p.Stat.AvgRtt = pkt.GetRtt()
		} else {
			p.Stat.AvgRtt = (p.Stat.AvgRtt*time.Duration(p.Stat.PacketsRecv-1) + pkt.GetRtt()) / time.Duration(p.Stat.PacketsRecv)
		}
	} else {
		// PacketLoss
		p.Stat.UpdatePacketLoss()
	}
}

// Method (Pinger) - Update pinger statistics
func (p *Pinger) Run(errChan chan<- error) {

	switch p.InputVars.Type {

	// Type: tcp
	case "tcp":
		// Go Routine - tcpProbingRun
		go tcpProbingRun(p, errChan)

	case "icmp":
		// Go Routine - icmpProbingRun
		go icmpProbingRun(p, errChan)

	case "http":
		// Go Routine - icmpProbingRun
		go httpProbingRun(p, errChan)

	case "dns":

	}

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

	// MinRtt is the minimum round-trip time sent via this pinger.
	MinRtt time.Duration

	// MaxRtt is the maximum round-trip time sent via this pinger.
	MaxRtt time.Duration

	// AvgRtt is the average round-trip time sent via this pinger.
	AvgRtt time.Duration
}

// Method (Statistics) - Update PacketLoss
func (s *Statistics) UpdatePacketLoss() {
	s.PacketLoss = (1 - float64(s.PacketsRecv)/float64(s.PacketsSent))
}
