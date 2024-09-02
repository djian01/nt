package ntPinger

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// CLI Imput variables Struct
type InputVars struct {
	Type           string // Select one of these: tcp, icmp, http, dns
	Count          int    // Default is 0 which means nonstop till interruption.
	NBypes         int    // Specific the payload. ICMP default payload is 24 bytes. TCP/HTTP/ICMP have no payload by default.
	Timeout        int    // default timeout is 4 seconds
	Interval       int    // Interval is the wait time between each packet send. Default is 1s.
	SourceHost     string
	DestHost       string
	DestPort       int
	Http_path      string
	Http_tls       bool
	Icmp_dfragment bool // ipv4 only
	Dns_request    string
}

// Packet Interface
type Packet interface {
	GetStatus() bool
	GetRtt() time.Duration
	UpdateStatistics(s Statistics)
}

// PacketTCP Struct
type PacketTCP struct {
	Type     string
	Status   bool
	Idx      int
	DestAddr string
	DestHost string
	DestPort int
	NBytes   int
	SendTime time.Time
	RTT      time.Duration
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	StatusDetails string
}

func (pkg PacketTCP) GetStatus() bool {
	return pkg.Status
}
func (pkg PacketTCP) GetRtt() time.Duration {
	return pkg.RTT
}
func (pkg *PacketTCP) UpdateStatistics(s Statistics) {
	pkg.AvgRtt = s.AvgRtt
	pkg.MaxRtt = s.MaxRtt
	pkg.MinRtt = s.MinRtt
	pkg.PacketsSent = s.PacketsSent
	pkg.PacketsRecv = s.PacketsRecv
	pkg.PacketLoss = s.PacketLoss
}

// PacketHTTP Struct
type PacketHTTP struct {
	Type               string
	Status             bool
	Idx                int
	DestAddr           string
	DestHost           string
	DestPort           int
	NBytes             int
	SendTime           time.Time
	RTT                time.Duration
	Http_path          string
	Http_tls           bool
	Http_response_code int
	Http_response      string
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	StatusDetails string
}

func (pkg PacketHTTP) GetStatus() bool {
	return pkg.Status
}
func (pkg PacketHTTP) GetRtt() time.Duration {
	return pkg.RTT
}
func (pkg *PacketHTTP) UpdateStatistics(s Statistics) {
	pkg.AvgRtt = s.AvgRtt
	pkg.MaxRtt = s.MaxRtt
	pkg.MinRtt = s.MinRtt
	pkg.PacketsSent = s.PacketsSent
	pkg.PacketsRecv = s.PacketsRecv
	pkg.PacketLoss = s.PacketLoss
}

// PacketICMP Struct
type PacketICMP struct {
	Type           string
	Status         bool
	Idx            int
	DestAddr       string
	DestHost       string
	DestPort       int
	NBytes         int
	SendTime       time.Time
	RTT            time.Duration
	Icmp_dfragment bool // ipv4 only
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	StatusDetails string
}

func (pkg PacketICMP) GetStatus() bool {
	return pkg.Status
}
func (pkg PacketICMP) GetRtt() time.Duration {
	return pkg.RTT
}
func (pkg *PacketICMP) UpdateStatistics(s Statistics) {
	pkg.AvgRtt = s.AvgRtt
	pkg.MaxRtt = s.MaxRtt
	pkg.MinRtt = s.MinRtt
	pkg.PacketsSent = s.PacketsSent
	pkg.PacketsRecv = s.PacketsRecv
	pkg.PacketLoss = s.PacketLoss
}

// PacketDNS Struct
type PacketDNS struct {
	Type         string
	Status       bool
	Idx          int
	DestAddr     string
	DestHost     string
	DestPort     int
	NBytes       int
	SendTime     time.Time
	RTT          time.Duration
	Dns_request  string
	Dns_response string
	// statistics
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	// Status update
	StatusDetails string
}

func (pkg PacketDNS) GetStatus() bool {
	return pkg.Status
}
func (pkg PacketDNS) GetRtt() time.Duration {
	return pkg.RTT
}
func (pkg *PacketDNS) UpdateStatistics(s Statistics) {
	pkg.AvgRtt = s.AvgRtt
	pkg.MaxRtt = s.MaxRtt
	pkg.MinRtt = s.MinRtt
	pkg.PacketsSent = s.PacketsSent
	pkg.PacketsRecv = s.PacketsRecv
	pkg.PacketLoss = s.PacketLoss
}

// Pinger Struct
type Pinger struct {

	// ************** General Fields **********************
	// Input Vars
	InputVars InputVars

	// statistics
	Stat Statistics

	// statistics Mutex
	statsMu sync.RWMutex

	// source IP details
	SourceAddr   string
	SourceIpAddr net.IP

	// destination IP details
	destAddr   string
	destIpAddr net.IP

	// OnSend is called when Pinger sends a packet
	OnSend func(Packet)

	// OnRecv is called when Pinger receives and processes a packet
	OnRecv func(Packet)

	// probeChan
	probeChan chan Packet
}

// Method (Pinger) - Resolve does the DNS lookup for the Pinger address and sets IP protocol.
func (p *Pinger) Resolve() error {
	if len(p.InputVars.DestHost) == 0 {
		return errors.New("destination Host cannot be empty")
	}
	// Check Name Resolution
	resolvedIPs, err := net.LookupIP(p.InputVars.DestHost)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to resolve domain: %v", p.InputVars.DestHost))
	}

	// Get the 1st IPv4 IP from resolved IPs
	for _, ip := range resolvedIPs {
		// To4() returns nil if it's not an IPv4 address
		if ip.To4() != nil {
			p.destIpAddr = ip
			p.destAddr = ip.String()
			break
		}
	}
	return nil
}

// Method (Pinger) - Update pinger statistics
func (p *Pinger) UpdateStatistics(pkt Packet) {

	// lock the Statistic
	p.statsMu.Lock()
	defer p.statsMu.Unlock()

	// PacketsSent
	p.Stat.PacketsSent++

	// PacketsRecv
	if pkt.GetStatus() {
		p.Stat.PacketsRecv++

		// PacketLoss
		p.Stat.updatePacketLoss()

		// MinRtt
		if p.Stat.MinRtt > pkt.GetRtt() {
			p.Stat.MinRtt = pkt.GetRtt()
		}

		// MaxRtt
		if p.Stat.MaxRtt < pkt.GetRtt() {
			p.Stat.MaxRtt = pkt.GetRtt()
		}

		// AvgRtt
		p.Stat.AvgRtt = time.Duration(((int64(p.Stat.AvgRtt)/1000000)*(int64(p.Stat.PacketsRecv-1))+(int64(pkt.GetRtt())/1000000))/int64(p.Stat.PacketsRecv)) * time.Millisecond
	}
}

// Method (Pinger) - Update pinger statistics
func (p *Pinger) Run() {

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// idx
	idx := 0

	// forLoopEnds Flag
	forLoopEnds := false

	switch p.InputVars.Type {

	// Type: tcp
	case "tcp":
		// count
		if p.InputVars.Count == 0 {
			for {
				if forLoopEnds {
					break
				}

				select {
				case <-interruptChan: // case interruptChan, close the channel & break the loop
					close(p.probeChan)
					forLoopEnds = true
				default:
					ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.InputVars.Timeout)*time.Second)
					defer cancel()

					pkt, err := tcpProbing(ctx, idx, p.destAddr, p.InputVars.DestPort, 0)

					if err != nil {
						if strings.Contains(err.Error(), "timeout") {
							// Probe Timeout
						} else {
							panic(err)
						}
					} else {
						// Probe Success
					}
					p.UpdateStatistics(&pkt)
					pkt.UpdateStatistics(p.Stat)
					fmt.Println(pkt)
					idx++
				}

			}

		} else {
			for i := 0; i < p.InputVars.Count; i++ {
				if forLoopEnds {
					break
				}
				select {
				case <-interruptChan:
					close(p.probeChan) // case interruptChan, close the channel & break the loop
					forLoopEnds = true
				default:
					ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.InputVars.Timeout)*time.Second)
					defer cancel()

					pkt, err := tcpProbing(ctx, idx, p.destAddr, p.InputVars.DestPort, 0)

					if err != nil {
						if strings.Contains(err.Error(), "timeout") {
							// Probe Timeout
						} else {
							panic(err)
						}
					} else {
						// Probe Success
					}
					p.UpdateStatistics(&pkt)
					pkt.UpdateStatistics(p.Stat)
					fmt.Println(pkt)
					idx++

					// check the last loop of the probing, close probeChan
					if i == (p.InputVars.Count - 1) {
						close(p.probeChan)
					}

					// sleep for interval
					time.Sleep(GetSleepTime(pkt.Status, p.InputVars.Interval))
				}
			}
		}

	case "icmp":

	case "http":

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
func (s *Statistics) updatePacketLoss() {
	s.PacketLoss = (1 - float64(s.PacketsRecv)/float64(s.PacketsSent))
}
