package tcpscan

import (
	"regexp"
	"strconv"

	"github.com/djian01/nt/pkg/ntScan"
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

	err := ntScan.ScanTcpRun(recording, destHost, Ports, timeout)
	if err != nil {
		return err
	} else {
		return nil
	}
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
