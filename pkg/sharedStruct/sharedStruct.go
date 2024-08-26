package sharedStruct

import "time"

// result struct
type NtResult struct {
	Seq       int
	Timestamp string
	HostName  string
	URL       string
	IP        string
	Size      int
	TCP_Port  string
	Status    string
	RTT       time.Duration
	Type      string
	// statistics
	PacketsSent int
	PacketsRecv int
	PacketLoss  float64
	MinRtt      time.Duration
	AvgRtt      time.Duration
	MaxRtt      time.Duration
}
