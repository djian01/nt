//go:build linux
// +build linux

package ntPinger

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// func: IcmpProbing
func IcmpProbing(Seq int, destAddr string, destHost string, PayLoadSize int, Icmp_DF bool, timeout int, payload []byte) (PacketICMP, error) {

	// Initial PacketICMP
	pkt := PacketICMP{
		Type:        "icmp",
		Status:      false,
		Seq:         Seq,
		DestAddr:    destAddr,
		DestHost:    destHost,
		PayLoadSize: PayLoadSize,
		Icmp_DF:     Icmp_DF,
	}

	// initial cmdOutput
	cmdOutput := []byte{}
	var err error

	// SEND - the ICMP Request
	pkt.SendTime = time.Now()

	// Construct the ping command for Linux
	args := []string{"-c", "1", "-W", strconv.Itoa(timeout), "-s", strconv.Itoa(PayLoadSize)}
	if Icmp_DF {
		args = append(args, "-M", "do")
	}
	args = append(args, destAddr)

	cmd := exec.Command("ping", args...)

	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 1") || strings.Contains(err.Error(), "exit status 2") {
			// PASS: error occurs if packets are dropped
		} else {
			fmt.Printf("Error running ping: %v\n", err)
			return pkt, err
		}
	}

	status, rtt, AdditionalInfo := parseLinuxPingOutput(string(cmdOutput))

	pkt.Status = status
	pkt.RTT = rtt
	pkt.AdditionalInfo = AdditionalInfo

	return pkt, nil
}

// func parseLinuxPingOutput
func parseLinuxPingOutput(output string) (status bool, rtt time.Duration, AdditionalInfo string) {
	// Regular expression to capture RTT from the "time=" part
	rttRegex := regexp.MustCompile(`time[=<](\d+\.?\d*)\s*ms`)

	// Check if there's a reply indicating success
	if strings.Contains(output, "bytes from") {
		rttMatch := rttRegex.FindStringSubmatch(output)
		if len(rttMatch) > 1 {
			rttStr := rttMatch[1]
			rttFloat, err := strconv.ParseFloat(rttStr, 64)
			if err == nil {
				rtt = time.Duration(rttFloat * float64(time.Millisecond))
			} else {
				fmt.Println("Error parsing RTT:", err)
				rtt = 0
			}
		}

		status = true

	} else if strings.Contains(output, "Message too long") || strings.Contains(output, "frag needed and DF set") {
		status = false
		rtt = 0
		AdditionalInfo = "MTU Exceed, DF set"
	} else if strings.Contains(output, "100% packet loss") {
		status = false
		rtt = 0
		AdditionalInfo = "Timeout"
	} else {
		status = false
		rtt = 0
		AdditionalInfo = "Unknown_Error"
	}

	return status, rtt, AdditionalInfo
}
