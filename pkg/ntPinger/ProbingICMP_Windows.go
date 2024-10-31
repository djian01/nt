//go:build windows
// +build windows

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
func IcmpProbing(Seq int, destAddr string, desetHost string, PayLoadSize int, Icmp_DF bool, timeout int, payload []byte) (PacketICMP, error) {

	// Initial PacketICMP
	pkt := PacketICMP{
		Type:        "icmp",
		Status:      false,
		Seq:         Seq,
		DestAddr:    destAddr,
		DestHost:    desetHost,
		PayLoadSize: PayLoadSize,
		Icmp_DF:     Icmp_DF,
	}

	// initial cmdOutput
	cmdOutput := []byte{}
	var err error

	// SEND - the ICMP Request
	pkt.SendTime = time.Now()

	// if DF bit is set
	if Icmp_DF {
		cmd := exec.Command("ping", destAddr, "-f", "-l", strconv.Itoa(PayLoadSize), "-n", "1", "-w", strconv.Itoa(timeout*1000))
		// output
		cmdOutput, err = cmd.CombinedOutput()
		if err != nil {
			if strings.Contains(err.Error(), "exit status 1") {
				// PASS
			} else {
				fmt.Printf("Error running ping: %v\n", err)
				return pkt, err
			}
		}

		// if DF bit is NOT set
	} else {
		cmd := exec.Command("ping", destAddr, "-l", strconv.Itoa(PayLoadSize), "-n", "1", "-w", strconv.Itoa(timeout*1000))
		// output
		cmdOutput, err = cmd.CombinedOutput()
		if err != nil {
			if strings.Contains(err.Error(), "exit status 1") {
				// PASS
			} else {
				fmt.Printf("Error running ping: %v\n", err)
				return pkt, err
			}
		}
	}

	status, rtt, AdditionalInfo := parseWinPingOutput(string(cmdOutput))

	pkt.Status = status
	pkt.RTT = rtt
	pkt.AdditionalInfo = AdditionalInfo

	return pkt, nil

}

// func parseWinPingOutput
func parseWinPingOutput(output string) (status bool, rtt time.Duration, AdditionalInfo string) {

	// Regular expression to capture RTT from the "time=" part
	rttRegex := regexp.MustCompile(`time(=|<)(\d+ms)`)

	// Check if there's a reply indicating success
	if strings.Contains(output, "Reply from") {
		// Find the RTT value
		rttMatch := rttRegex.FindStringSubmatch(output)
		if len(rttMatch) > 1 {
			rttStr := rttMatch[2]

			// Convert RTT string (e.g., "13ms") to time.Duration
			rttDuration, err := time.ParseDuration(rttStr) // string "13ms" -> time.Duration(13)*time.Millisecond
			if err == nil {
				rtt = rttDuration
			} else {
				fmt.Println("Error parsing RTT:", err)
				rtt = 0
			}
		}

		status = true
	} else if strings.Contains(output, "Request timed out") && strings.Contains(output, "Packets: Sent = 1, Received = 0") {
		// In case of a timeout or packet loss
		status = false
		rtt = 0
		AdditionalInfo = "Timeout"
	} else if strings.Contains(output, "Packet needs to be fragmented but DF set") && strings.Contains(output, "Packets: Sent = 1, Received = 0") {
		// MTU_Exceed
		status = false
		AdditionalInfo = "MTU Exceed, DF set"
		rtt = 0
	} else {
		status = false
		rtt = 0
		AdditionalInfo = "Unknown_Error"
	}

	return status, rtt, AdditionalInfo
}
