package output

import (
	"fmt"
	"nt/pkg/sharedStruct"
	"strings"

	"github.com/fatih/color"
)

// Table Print Func
func TablePrint(displayTable *[]sharedStruct.NtResult, len int, recording bool, displayIdx int) {

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
	switch (*displayTable)[0].Type {

	case "icmp":
		// Print the table header
		moveToRow(tableHeadRowIdx + 1)
		fmt.Printf("%-5s %-15s %-15s %-15s %-10s %-15s %-30s \n", "Seq", "Status", "HostName", "IP", "Size", "RTT", "Timestamp")
		fmt.Println(strings.Repeat("-", 106))

		// Print the table & statistics data
		for idx, t := range *displayTable {

			// ANSI escape code to move the cursor to a specific row (1-based index)
			moveToRow(idx + tableHeadRowIdx + 3)

			if t.Timestamp == "" {
				fmt.Printf("%-5s %-15s %-15s %-15s %-10s %-15s %-30s\n", "", "", "", "", "", "", "")
			} else {
				if t.Status == "ICMP_OK" {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-15s", t.Status)
					fmt.Printf("%-5d %-s %-15s %-15s %-10d %-15v %-30s       \n", t.Seq, color.GreenString(Status), t.HostName, t.IP, t.Size, t.RTT, t.Timestamp)
				} else if t.Status == "ICMP_Failed" {
					Status := fmt.Sprintf("%-15s", t.Status)
					fmt.Printf("%-5d %-s %-15s %-15s %-10d %-15v %-30s       \n", t.Seq, color.RedString(Status), t.HostName, t.IP, t.Size, t.RTT, t.Timestamp)
				}
			}

			// print the statistics
			if t.Timestamp != "" {
				moveToRow(len + tableHeadRowIdx + 3)
				fmt.Printf("\n--- %s ICMP ping statistics ---\n", t.IP)
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", t.PacketsSent, t.PacketsRecv, float64(t.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", t.MinRtt, t.AvgRtt, t.MaxRtt)
			}

		}

		// move the cursor to row
		moveToRow(len + tableHeadRowIdx + 8)

	case "tcp":

	case "http":

	case "dns":

	}

}

// Func - move cursor to x row
func moveToRow(row int) {
	// ANSI escape code to move the cursor to a specific row (1-based index)
	fmt.Printf("\033[%d;1H", row)
}
