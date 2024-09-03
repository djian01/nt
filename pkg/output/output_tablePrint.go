package output

import (
	"fmt"
	"nt/pkg/ntPinger"
	"strings"

	"github.com/fatih/color"
)

// Table Print Func
func TablePrint(displayTable *[]ntPinger.Packet, len int, recording bool, displayIdx int) {

	// Set the 1st row for table head
	var tableHeadRowIdx int
	if recording {
		tableHeadRowIdx = 1
	} else {
		tableHeadRowIdx = 0
	}

	// Display REC if recording enabled
	if recording {
		moveToRow(1)
		if displayIdx%2 == 0 {
			fmt.Printf("%s", color.RedString("REC ●"))
		} else {
			fmt.Printf("%s", color.RedString("REC    "))
		}

	}

	// Pring Table by "type"
	switch (*displayTable)[0].GetType() {

	case "icmp":
		// Print the table header
		moveToRow(tableHeadRowIdx + 1)

		fmt.Printf("%-5s %-15s %-15s %-15s %-10s %-15s %-30s \n", "Seq", "Status", "HostName", "IP", "Size", "RTT", "Timestamp")
		fmt.Println(strings.Repeat("-", 106))

		// Print the table & statistics data
		for idx, t := range *displayTable {

			pkt := &ntPinger.PacketICMP{}

			if t != nil {
				pkt = t.(*ntPinger.PacketICMP)
			}

			// ANSI escape code to move the cursor to a specific row (1-based index)
			moveToRow(idx + tableHeadRowIdx + 3)

			if pkt.SendTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
				fmt.Printf("%-5s %-15s %-15s %-15s %-10s %-15s %-30s\n", "", "", "", "", "", "", "")
			} else {
				if pkt.Status {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-15v", pkt.Status)
					fmt.Printf("%-5d %-s %-15s %-15s %-10d %-15v %-30s       \n", pkt.Seq, color.GreenString(Status), pkt.DestHost, pkt.DestAddr, pkt.NBytes, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"))
				} else {
					Status := fmt.Sprintf("%-15v", pkt.Status)
					fmt.Printf("%-5d %-s %-15s %-15s %-10d %-15v %-30s       \n", pkt.Seq, color.RedString(Status), pkt.DestHost, pkt.DestAddr, pkt.NBytes, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"))
				}
			}

			// print the statistics
			if pkt.SendTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
				moveToRow(len + tableHeadRowIdx + 3)
				fmt.Printf("\n--- %s %s Ping statistics ---\n", pkt.DestAddr, pkt.Type)
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", pkt.PacketsSent, pkt.PacketsRecv, float64(pkt.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", pkt.MinRtt, pkt.AvgRtt, pkt.MaxRtt)
			}

		}

		// move the cursor to row
		moveToRow(len + tableHeadRowIdx + 8)

	case "tcp":
		// Print the table header
		moveToRow(tableHeadRowIdx + 1)

		fmt.Printf("%-5s %-15s %-15s %-15s %-10s %-12s %-15s %-30s \n", "Seq", "Status", "HostName", "IP", "Port", "Size", "RTT", "Timestamp")
		fmt.Println(strings.Repeat("-", 114))

		// Print the table & statistics data
		for idx, t := range *displayTable {

			pkt := &ntPinger.PacketTCP{}

			if t != nil {
				pkt = t.(*ntPinger.PacketTCP)
			}

			// ANSI escape code to move the cursor to a specific row (1-based index)
			moveToRow(idx + tableHeadRowIdx + 3)

			if pkt.SendTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
				fmt.Printf("%-5s %-15s %-15s %-15s %-10s %-12s %-15s %-30s\n", "", "", "", "", "", "", "", "")
			} else {
				if pkt.Status {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-15v", pkt.Status)
					Size := fmt.Sprintf("%d bytes", pkt.NBytes)
					fmt.Printf("%-5d %-s %-15s %-15s %-10d %-12s %-15v %-30s       \n", pkt.Seq, color.GreenString(Status), pkt.DestHost, pkt.DestAddr, pkt.DestPort, Size, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"))
				} else {
					Status := fmt.Sprintf("%-15v", pkt.Status)
					Size := fmt.Sprintf("%d bytes", pkt.NBytes)
					fmt.Printf("%-5d %-s %-15s %-15s %-10d %-12s %-15v %-30s       \n", pkt.Seq, color.RedString(Status), pkt.DestHost, pkt.DestAddr, pkt.DestPort, Size, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"))
				}
			}

			// print the statistics
			if pkt.SendTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
				moveToRow(len + tableHeadRowIdx + 3)
				fmt.Printf("\n--- %s %s Ping statistics ---\n", pkt.DestAddr, pkt.Type)
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", pkt.PacketsSent, pkt.PacketsRecv, float64(pkt.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", pkt.MinRtt, pkt.AvgRtt, pkt.MaxRtt)
			}

		}

		// move the cursor to row
		moveToRow(len + tableHeadRowIdx + 8)
	case "http":

	case "dns":

	}

}

// Func - move cursor to x row
func moveToRow(row int) {
	// ANSI escape code to move the cursor to a specific row (1-based index)
	fmt.Printf("\033[%d;1H", row)
}
