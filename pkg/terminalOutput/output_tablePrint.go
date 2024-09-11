package terminalOutput

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

		fmt.Printf("%-5s %-10s %-20s %-20s %-10s %-15s %-25s %-20s \n", "Seq", "Status", "HostName", "IP", "Payload", "RTT", "Timestamp", "AddInfo")
		fmt.Println(strings.Repeat("-", 125))

		// Print the table & statistics data
		for idx, t := range *displayTable {

			pkt := &ntPinger.PacketICMP{}

			if t != nil {
				pkt = t.(*ntPinger.PacketICMP)
			}

			// ANSI escape code to move the cursor to a specific row (1-based index)
			moveToRow(idx + tableHeadRowIdx + 3)

			if pkt.SendTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
				fmt.Printf("%-5s %-10s %-20s %-20s %-10s %-15s %-25s %-20s\n", "", "", "", "", "", "", "", "")
			} else {
				// AddInfo
				AddInfo := fmt.Sprintf("%-20s", pkt.AdditionalInfo)

				// check Status
				if pkt.Status {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-20s %-20s %-10d %-15v %-25s %-s       \n", pkt.Seq, color.GreenString(Status), pkt.DestHost, pkt.DestAddr, pkt.PayLoadSize, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				} else {
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-20s %-20s %-10d %-15v %-25s %-s      \n", pkt.Seq, color.RedString(Status), pkt.DestHost, pkt.DestAddr, pkt.PayLoadSize, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				}
			}

			// print the statistics
			if pkt.SendTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
				moveToRow(len + tableHeadRowIdx + 3)
				fmt.Printf("\n--- %s %s statistics ---\n", pkt.DestAddr, color.CyanString(fmt.Sprintf("%v Ping", pkt.Type)))
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", pkt.PacketsSent, pkt.PacketsRecv, float64(pkt.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", pkt.MinRtt, pkt.AvgRtt, pkt.MaxRtt)
			}

		}

		// move the cursor to row
		moveToRow(len + tableHeadRowIdx + 8)

	case "tcp":
		// Print the table header
		moveToRow(tableHeadRowIdx + 1)

		fmt.Printf("%-5s %-10s %-20s %-20s %-10s %-10s %-15s %-25s %-20s  \n", "Seq", "Status", "HostName", "IP", "Port", "Payload", "RTT", "Timestamp", "AddInfo")
		fmt.Println(strings.Repeat("-", 135))

		// Print the table & statistics data
		for idx, t := range *displayTable {

			pkt := &ntPinger.PacketTCP{}

			if t != nil {
				pkt = t.(*ntPinger.PacketTCP)
			}

			// ANSI escape code to move the cursor to a specific row (1-based index)
			moveToRow(idx + tableHeadRowIdx + 3)

			if pkt.SendTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
				fmt.Printf("%-5s %-10s %-20s %-20s %-10s %-10s %-15s %-25s %-20s \n", "", "", "", "", "", "", "", "", "")
			} else {

				// AddInfo
				AddInfo := fmt.Sprintf("%-20s", pkt.AdditionalInfo)

				// check Status
				if pkt.Status {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-20s %-20s %-10d %-10d %-15v %-25s %-s      \n", pkt.Seq, color.GreenString(Status), pkt.DestHost, pkt.DestAddr, pkt.DestPort, pkt.PayLoadSize, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				} else {
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-20s %-20s %-10d %-10d %-15v %-25s %-s       \n", pkt.Seq, color.RedString(Status), pkt.DestHost, pkt.DestAddr, pkt.DestPort, pkt.PayLoadSize, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				}
			}

			// print the statistics
			if pkt.SendTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
				moveToRow(len + tableHeadRowIdx + 3)
				fmt.Printf("\n--- %s %s statistics ---\n", pkt.DestAddr, color.CyanString(fmt.Sprintf("%v Ping", pkt.Type)))
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", pkt.PacketsSent, pkt.PacketsRecv, float64(pkt.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", pkt.MinRtt, pkt.AvgRtt, pkt.MaxRtt)
			}

		}

		// move the cursor to row
		moveToRow(len + tableHeadRowIdx + 8)
	case "http":
		// Print the table header
		moveToRow(tableHeadRowIdx + 1)

		fmt.Printf("%-5s %-10s %-10s %-33s %-15s %-15s %-25s %-20s  \n", "Seq", "Status", "Method", "URL", "Response_Code", "Response_Time", "Timestamp", "AddInfo")
		fmt.Println(strings.Repeat("-", 130))

		// Print the table & statistics data
		for idx, t := range *displayTable {

			pkt := &ntPinger.PacketHTTP{}

			if t != nil {
				pkt = t.(*ntPinger.PacketHTTP)
			}

			// ANSI escape code to move the cursor to a specific row (1-based index)
			moveToRow(idx + tableHeadRowIdx + 3)

			if pkt.SendTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
				fmt.Printf("%-5s %-10s %-10s %-33s %-15s %-15s %-25s %-20s \n", "", "", "", "", "", "", "", "")
			} else {

				// AddInfo
				AddInfo := fmt.Sprintf("%-20s", pkt.AdditionalInfo)

				// url
				url := ntPinger.ConstructURL(pkt.Http_scheme, pkt.DestHost, pkt.Http_path, pkt.DestPort)
				url = TruncateString(url, 30)

				// check Status
				if pkt.Status {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-10s %-33s %-15d %-15v %-25s %-s      \n", pkt.Seq, color.GreenString(Status), pkt.Http_method, url, pkt.Http_response_code, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				} else {
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-10s %-33s %-15d %-15v %-25s %-s      \n", pkt.Seq, color.RedString(Status), pkt.Http_method, url, pkt.Http_response_code, pkt.RTT, pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				}
			}

			// print the statistics
			if pkt.SendTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
				moveToRow(len + tableHeadRowIdx + 3)
				fmt.Printf("\n--- %s %s statistics ---\n", pkt.DestAddr, color.CyanString(fmt.Sprintf("%v Ping", pkt.Type)))
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", pkt.PacketsSent, pkt.PacketsRecv, float64(pkt.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", pkt.MinRtt, pkt.AvgRtt, pkt.MaxRtt)
			}

		}

		// move the cursor to row
		moveToRow(len + tableHeadRowIdx + 8)
	case "dns":
		// Print the table header
		moveToRow(tableHeadRowIdx + 1)

		fmt.Printf("%-5s %-10s %-15s %-25s %-25s %-12s %-10s %-15s %-20s %-20s  \n", "Seq", "Status", "Resolver", "Query", "Response", "Query_Type", "Protocol", "Response_Time","Send_Time","AddInfo")
		fmt.Println(strings.Repeat("-", 158))

		// Print the table & statistics data
		for idx, t := range *displayTable {

			pkt := &ntPinger.PacketDNS{}

			if t != nil {
				pkt = t.(*ntPinger.PacketDNS)
			}

			// ANSI escape code to move the cursor to a specific row (1-based index)
			moveToRow(idx + tableHeadRowIdx + 3)

			if pkt.SendTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
				fmt.Printf("%-5s %-10s %-15s %-25s %-25s %-12s %-10s %-15s %-20s %20s \n", "", "", "", "", "", "", "", "", "", "")
			} else {

				// AddInfo
				AddInfo := fmt.Sprintf("%-20s", pkt.AdditionalInfo)

				// Query
				Query := TruncateString(pkt.Dns_query, 22)

				// Response
				Response := TruncateString(pkt.Dns_response, 22)


				// check Status
				if pkt.Status {
					// When using the /fatih/color package, the colored string produced by color.GreenString(t.Status) is already
					// wrapped with escape sequences that apply the color in the terminal. This wrapping adds extra characters to the string,
					// which affects how the width specifier (like %-20s) is interpreted
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-15s %-25s %-25s %-12v %-10s %-15s %-20s %-s      \n", pkt.Seq, color.GreenString(Status), pkt.DestHost, Query, Response, pkt.Dns_queryType, pkt.Dns_protocol, pkt.RTT,  pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				} else {
					Status := fmt.Sprintf("%-10v", pkt.Status)
					fmt.Printf("%-5d %-s %-15s %-25s %-25s %-12v %-10s %-15s %-20s %-s      \n", pkt.Seq, color.RedString(Status), pkt.DestHost, Query, Response, pkt.Dns_queryType, pkt.Dns_protocol, pkt.RTT,  pkt.SendTime.Format("2006-01-02 15:04:05"), color.YellowString(AddInfo))
				}
			}

			// print the statistics
			if pkt.SendTime.String() != "0001-01-01 00:00:00 +0000 UTC" {
				moveToRow(len + tableHeadRowIdx + 3)
				fmt.Printf("\n--- %s %s statistics ---\n", pkt.DestAddr, color.CyanString(fmt.Sprintf("%v Ping", pkt.Type)))
				fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", pkt.PacketsSent, pkt.PacketsRecv, float64(pkt.PacketLoss*100))
				fmt.Printf("round-trip min/avg/max = %v/%v/%v       \n", pkt.MinRtt, pkt.AvgRtt, pkt.MaxRtt)
			}

		}

		// move the cursor to row
		moveToRow(len + tableHeadRowIdx + 8)
	}

}

// Func - move cursor to x row
func moveToRow(row int) {
	// ANSI escape code to move the cursor to a specific row (1-based index)
	fmt.Printf("\033[%d;1H", row)
}

// TruncateString truncates a string to a maximum length and appends "..." if it exceeds the max length
func TruncateString(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength-3] + "..." // Subtract 3 to account for "..."
	}
	return s
}
