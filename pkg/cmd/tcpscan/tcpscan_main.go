package tcpscan

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
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

# Example: TCP SCAN to "10.2.3.10" for port "22, 1522-1525 & 8433" with custom timeout: 5 sec
nt tcpscan -t 5 10.2.3.10 22 1522-1525 8433
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
	if PortCount == 0 {
		return fmt.Errorf("at least one port is required")
	}
	if PortCount > 50 {
		return fmt.Errorf("total number of input ports is over 50")
	}
	if timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}
	for _, port := range Ports {
		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535: %d", port)
		}
	}

	// resolve destHost
	destAddr := ""

	resolvedIPs, err := ntPinger.ResolveDestHost(destHost)
	if err != nil {
		return fmt.Errorf("failed to resolve domain %s: %w", destHost, err)
	}

	// Get the 1st IPv4 IP from resolved IPs
	for _, ip := range resolvedIPs {
		// To4() returns nil if it's not an IPv4 address
		if ip.To4() != nil {
			destAddr = ip.String()
			break
		}
	}
	if destAddr == "" {
		return fmt.Errorf("destination has no IPv4 address: %s", destHost)
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

	/// clear screen (only need to clear screen once)
	terminalOutput.ClearScreen()

	terminalOutput.ScanTablePrint(&PortsTable, recording, 0, destHost)

	jobs := make(chan ntScaner.TcpScanPort, PortCount)
	results := make(chan ntScaner.TcpScanResult, PortCount)

	// create 5 workers
	workerCount := 5
	if PortCount < workerCount {
		workerCount = PortCount
	}
	var workers sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			ntScaner.ScanTcpWorkerResults(jobs, results)
		}()
	}

	// go routine - TcpScanPort <- &PortsTable[i]
	for i := 0; i < PortCount; i++ {
		jobs <- PortsTable[i]
	}
	close(jobs)

	var scanErr error
	for completed := 0; completed < PortCount; completed++ {
		result := <-results
		if result.Err != nil {
			if scanErr == nil {
				scanErr = result.Err
			}
			continue
		}
		PortsTable[result.Port.ID] = result.Port
		terminalOutput.ScanTablePrint(&PortsTable, recording, completed+1, destHost)
	}
	workers.Wait()
	if scanErr != nil {
		return scanErr
	}

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
		if port < 1 || port > 65535 {
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
		if start >= 1 && start <= end && end <= 65535 {
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
