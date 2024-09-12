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
	case "tcp", "icmp", "dns", "http":

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
func GetSleepTime(PacketStatus bool, Interval int, RTT time.Duration) time.Duration {

	// if RTT is not 0
	if RTT.String() != "0001-01-01 00:00:00 +0000 UTC" {
		return time.Duration(Interval) * time.Second
		// if Packet Status is "true"
	} else if PacketStatus {
		return time.Duration(Interval) * time.Second
		// else return 0 sleep time ( in the case of timeout)
	} else {
		return time.Duration(0) * time.Second
	}
}

// Func - Check Latency. If the AvgRTT is over 10ms and current RTT > 2 * AvgRTT -> High_Latency
func CheckLatency(avgRtt, currentRtt time.Duration) bool {
	// HighLatencyFlag
	HighLatencyFlag := false

	// Check if the average RTT is larger than 10ms and if current RTT is more than double the average RTT
	if avgRtt > 10*time.Millisecond && currentRtt > 2*avgRtt {
		HighLatencyFlag = true
	}

	return HighLatencyFlag
}

// Func - Generate Payload Date []byte
func GeneratePayloadData(payLoadSize int) []byte {

	bin := make([]byte, payLoadSize)
	for i := 0; i < payLoadSize; i++ {
		bin[i] = byte(i) // Example payload data
	}

	return bin
}

// Construct URL
func ConstructURL(Http_scheme, DestHost, Http_path string, DestPort int) string {

	url := ""

	if Http_path == "" {
		if Http_scheme == "http" && DestPort == 80 {
			url = fmt.Sprintf("%s://%s", Http_scheme, DestHost)
		} else if Http_scheme == "https" && DestPort == 443 {
			url = fmt.Sprintf("%s://%s", Http_scheme, DestHost)
		} else {
			url = fmt.Sprintf("%s://%s:%d", Http_scheme, DestHost, DestPort)
		}

	} else {
		if Http_scheme == "http" && DestPort == 80 {
			url = fmt.Sprintf("%s://%s/%s", Http_scheme, DestHost, Http_path)
		} else if Http_scheme == "https" && DestPort == 443 {
			url = fmt.Sprintf("%s://%s/%s", Http_scheme, DestHost, Http_path)
		} else {
			url = fmt.Sprintf("%s://%s:%d/%s", Http_scheme, DestHost, DestPort, Http_path)
		}
	}
	return url
}
