package output

import (
	"fmt"
	"nt/pkg/sharedStruct"
	"strings"

	"github.com/fatih/color"
)

func TablePrint(displayTable *[]sharedStruct.NtResult, len int) {

	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Pring Table by "type"
	switch (*displayTable)[0].Type {

	case "ping":
		// Print the table header
		fmt.Printf("%-10s %-20s %-20s %-20s %-10s %-10s %-30s \n", "Seq", "Status", "HostName", "IP", "Size", "RTT", "Timestamp")
		fmt.Println(strings.Repeat("-", 116))

		// Print the table & statistics data
		for idx, t := range *displayTable {
			// ANSI escape code to move the cursor to a specific row (1-based index)
			fmt.Printf("\033[%d;1H", idx+3)

			if t.Timestamp == "" {
				fmt.Printf("%-10s %-20s %-20s %-20s %-10s %-10s %-30s\n", "", "", "", "", "", "", "")
			} else {
				if t.Status == "PING_OK" {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-20s", t.Status)
					fmt.Printf("%-10d %s %-20s %-20s %-10d %-10v %-30s       \n", t.Seq, color.GreenString(Status), t.HostName, t.IP, t.Size, t.RTT, t.Timestamp)
				} else if t.Status == "PING_Failed" {
					Status := fmt.Sprintf("%-20s", t.Status)
					fmt.Printf("%-10d %s %-20s %-20s %-10d %-10v %-30s       \n", t.Seq, color.RedString(Status), t.HostName, t.IP, t.Size, t.RTT, t.Timestamp)
				}
			}

			// print the statistics
			if t.Timestamp != "" {
				fmt.Printf("\033[%d;1H", 13)
				fmt.Printf("\n--- %s ping statistics ---\n", t.IP)
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", t.PacketsSent, t.PacketsRecv, float64(t.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", t.MinRtt, t.AvgRtt, t.MaxRtt)
			}

		}
		fmt.Printf("\033[%d;1H", 18)

	case "tcp":

	case "http":

	}

	// Print

}

func moveToRow(row int) {
	// ANSI escape code to move the cursor to a specific row (1-based index)
	fmt.Printf("\033[%d;1H", row)
}
