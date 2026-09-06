package ntPinger

import (
	"fmt"
	"net"
	"time"
)

// New returns a new Pinger struct pointer.
func NewPinger(InputVars InputVars) (*Pinger, error) {
	if err := validateInputVars(InputVars); err != nil {
		return nil, err
	}

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

func validateInputVars(input InputVars) error {
	if input.Count < 0 {
		return fmt.Errorf("count cannot be negative")
	}
	if input.PayLoadSize < 0 {
		return fmt.Errorf("payload size cannot be negative")
	}
	if input.Type == "icmp" && input.PayLoadSize > 65507 {
		return fmt.Errorf("ICMP payload size cannot exceed 65507 bytes")
	}
	if input.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}
	if input.Interval < 0 {
		return fmt.Errorf("interval cannot be negative")
	}
	if (input.Type == "tcp" || input.Type == "http") && (input.DestPort < 1 || input.DestPort > 65535) {
		return fmt.Errorf("destination port must be between 1 and 65535")
	}
	if input.Type == "dns" && input.Dns_Protocol != "udp" && input.Dns_Protocol != "tcp" {
		return fmt.Errorf("DNS protocol must be udp or tcp")
	}
	return nil
}

// Func - Get Sleep Time
func GetSleepTime(_ bool, interval int, rtt time.Duration) time.Duration {
	delay := time.Duration(interval)*time.Second - rtt
	if delay < 0 {
		return 0
	}
	return delay
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
			url = fmt.Sprintf("%s://%s%s", Http_scheme, DestHost, Http_path)
		} else if Http_scheme == "https" && DestPort == 443 {
			url = fmt.Sprintf("%s://%s%s", Http_scheme, DestHost, Http_path)
		} else {
			url = fmt.Sprintf("%s://%s:%d%s", Http_scheme, DestHost, DestPort, Http_path)
		}
	}
	return url
}

// Func - resolveDestHost
func ResolveDestHost(DestHost string) (DestAddr []net.IP, err error) {
	// Check Name Resolution
	DestAddr, err = net.LookupIP(DestHost)
	return
}
