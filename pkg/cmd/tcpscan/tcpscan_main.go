package tcpscan

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/ntScaner"
	"github.com/djian01/nt/pkg/record"
	"github.com/djian01/nt/pkg/terminalOutput"
	"github.com/spf13/cobra"
)

// Initial tcpCmd
var tcpScanCmd = &cobra.Command{
	Use:   "tcpscan [flags] <Destination Host> <Destination Port> ...", // Sub-command, shown in the -h, Usage field
	Short: "TCP SCAN Test Module",
	Long:  "TCP SCAN Test Module for target host listening port scanning",
	Args:  cobra.MinimumNArgs(2), // Minimum 2 Arg, <Destination Host> <Destination Port> ...
	Run:   TcpScanCommandLink,
	Example: `
# Example: TCP Scan to "10.123.1.10" for port "80, 443, 8080 & 1500-1505" with recording enabled
nt -r tcpscan 10.123.1.10 80 443 8080 1500-1505

# Example: TCP SCAN to "10.2.3.10" for port "22, 1522-1525 & 8433" with custom timeout: 5 sec and custom interval: 2 sec
nt tcpscan -t 5 -i 2 10.2.3.10 22 1522-1525 8433
`,
}

// Func - TcpScanCommandLink: obtain Flags and call TcpScanCommandMain()
func TcpScanCommandLink(cmd *cobra.Command, args []string) {

	// GFlag -r
	recording, _ := cmd.Flags().GetBool("recording")

	// Arg - destHost
	destHost := args[0]

	// Arg - Ports
	var Ports []int

	for _, arg := range args[1:] {
		valid, ports := IsValidInput(arg)
		if !valid {
			panic("Errors with the input Port(s)")
		}
		Ports = append(Ports, ports...)
	}

	// Flag -t
	timeout, _ := cmd.Flags().GetInt("timeout")

	// call func TcpScanCommandMain
	err := TcpScanCommandMain(recording, destHost, Ports, timeout)
	if err != nil {
		panic(err) // panic all error from under
	}
}

// Func TcpScanCommandMain
func TcpScanCommandMain(recording bool, destHost string, Ports []int, timeout int) error {

	// get the total number of the ports
	PortCount := len(Ports)
	countTested := 0

	// resolve destHost
	destAddr := ""

	resolvedIPs, err := ntPinger.ResolveDestHost(destHost)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to resolve domain: %v", destHost))
	}

	// Get the 1st IPv4 IP from resolved IPs
	for _, ip := range resolvedIPs {
		// To4() returns nil if it's not an IPv4 address
		if ip.To4() != nil {
			destAddr = ip.String()
			break
		}
	}

	// if the total number of ports is over 50, error and exit
	if PortCount > 50 {
		return fmt.Errorf("error: total number of input ports is over 50")
	}

	// Create a list of 50 empty TcpScanPort items
	PortsTable := make([]ntScaner.TcpScanPort, 50)

	for i := 0; i < PortCount; i++ {
		PortsTable[i].ID = i
		PortsTable[i].Port = Ports[i]
		PortsTable[i].Timeout = timeout
		PortsTable[i].DestAddr = destAddr
		PortsTable[i].DestHost = destHost
		PortsTable[i].Status = 1
	}

	// Go Routine - Terminal Output
	go func() {

		destplayIdx := 0

		for {
			// Display the items in a 10×5 table
			terminalOutput.ScanTablePrint(&PortsTable, recording, destplayIdx, destHost)

			time.Sleep(time.Duration(1) * time.Second)

			if countTested == PortCount {
				// display the final output
				terminalOutput.ScanTablePrint(&PortsTable, recording, destplayIdx, destHost)
				break
			}

			destplayIdx++
		}
	}()

	// Create TcpScanPort & errChan
	TcpScanPort := make(chan *ntScaner.TcpScanPort, 1)
	errChan := make(chan error, 1)

	// create 5 workers
	for i := 0; i < 5; i++ {
		go ntScaner.ScanTcpWorker(TcpScanPort, errChan)
	}

	// go routine - TcpScanPort <- &PortsTable[i]
	go func() {
		for i := 0; i < PortCount; i++ {
			TcpScanPort <- &PortsTable[i]
		}
	}()

	// loop
	loopBreak := false

	for {
		if loopBreak {
			break
		}

		select {
		case err := <-errChan:
			loopBreak = true
			return err
		default:
			countTested, _, _ = ntScaner.TcpScanStat(&PortsTable)
			if countTested == PortCount {
				loopBreak = true
			}
		}
	}

	// sleep and wait for 1 sec
	time.Sleep(time.Duration(1) * time.Second)

	// recording
	if recording {

		// recordingFile Path
		exeFileFolder, err := os.Getwd()
		if err != nil {
			return err
		}
		// recordingFile Name
		timeStamp := time.Now().Format("20060102150405")
		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", "TcpScan", destHost, timeStamp)
		recordingFilePath := filepath.Join(exeFileFolder, recordingFileName)

		record.TcpScan_Recording(recordingFilePath, PortsTable)
	}

	return nil
}

// Func - TcpScanCommand
func TcpScanCommand() *cobra.Command {
	return tcpScanCmd
}

// Func - init()
func init() {

	// Flag - Ping timeout
	var timeout int
	tcpScanCmd.Flags().IntVarP(&timeout, "timeout", "t", 4, "TCP Ping Timeout (default: 4 sec)")
}

// isValidInput checks if the argument is either a valid integer or a valid range of integers.
func IsValidInput(arg string) (bool, []int) {

	// Check if it's a single integer
	port, err := strconv.Atoi(arg)
	if err == nil {
		if port > 65535 { // Check if the port exceeds 65535
			return false, nil
		}
		return true, []int{port}
	}

	// Check if it's a range (e.g., 100-105)
	rangeRegex := regexp.MustCompile(`^(\d+)-(\d+)$`)

	matches := rangeRegex.FindStringSubmatch(arg)
	if len(matches) == 3 {
		start, _ := strconv.Atoi(matches[1])
		end, _ := strconv.Atoi(matches[2])

		// Ensure the range is valid (start <= end)
		if start <= end && end <= 65535 {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return true, result
		}
	}

	// Invalid input or port exceeds 65535
	return false, nil
}
