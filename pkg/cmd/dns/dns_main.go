package dns

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/record"
	output "github.com/djian01/nt/pkg/terminalOutput"
)

// Initial dnsCmd
var dnsCmd = &cobra.Command{
	Use:   "dns [flags] <DNS Resolver IP> <DNS Query Name>", // Sub-command, shown in the -h, Usage field
	Short: "DNS Ping Test Module",
	Long:  "DNS Ping test Module for dns testing",
	Args:  cobra.ExactArgs(2), // 2 Args, <Destination Host> <Destination Port> are required
	Run:   DnsCommandLink,
	Example: `
# Example: DNS ping to "8.8.8.8" with query "google.com" and have recording enabled
nt -r dns 8.8.8.8 google.com

# Example: DNS ping to "4.2.2.2" with query "abc.com" with count: 10 and interval: 2 sec
nt dns -c 10 -i 2 4.2.2.2 abc.com
`,
}

// Initial the bucket
var bucket = 10

// Func - DnsCommandLink: obtain Flags and call DnsCommandMain()
func DnsCommandLink(cmd *cobra.Command, args []string) {

	// GFlag -r
	recording, _ := cmd.Flags().GetBool("recording")

	// GFlag -d
	displayRow, _ := cmd.Flags().GetInt("displayrow")

	// Arg - destHost
	destHost := args[0]

	// dns_query
	Dns_query := args[1]

	// Flag -c
	count, _ := cmd.Flags().GetInt("count")

	// Flag -o
	Dns_protocol, _ := cmd.Flags().GetString("protocol")

	// Flag -t
	timeout, _ := cmd.Flags().GetInt("timeout")

	// Flag -i
	interval, _ := cmd.Flags().GetInt("interval")

	// call func TcpCommandMain
	err := DnsCommandMain(recording, displayRow, destHost, Dns_query, Dns_protocol, count, timeout, interval)
	if err != nil {
		panic(err) // panic all error from under
	}
}

// Func - TcpCommandMain
func DnsCommandMain(recording bool, displayRow int, destHost string, Dns_query string, Dns_protocol string, count int, timeout int, interval int) error {

	// Wait Group
	var wgRecord sync.WaitGroup

	// recording row
	recordingRow := 0
	if recording {
		recordingRow = 1
	}

	// recordingFilePath
	recordingFilePath := ""

	// Channel - outputChan (if there are N go routine, the channel deep is N)
	outputChan := make(chan ntPinger.Packet, 1)
	defer close(outputChan)

	// Channel - error (for Go Routines)
	errChan := make(chan error, 1)
	defer close(errChan)

	// Channel - recordingChan, closed in the end of the testing, no need to defer close
	recordingChan := make(chan ntPinger.Packet, 1)

	// build the InputVar
	InputVar := ntPinger.InputVars{
		Type:         "dns",
		Count:        count,
		Timeout:      timeout,
		Interval:     interval,
		DestHost:     destHost,
		Dns_query:    Dns_query,
		Dns_Protocol: Dns_protocol,
	}

	// Start Ping Main Command, manually input display Len
	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		return err // return err from NewPinger including resolve error
	}

	go p.Run(errChan)

	// Output
	//// Go Routine: OutputFunc
	go output.OutputFunc(outputChan, displayRow, recording)

	// Recording
	if recording {

		// recordingFile Path
		exeFileFolder, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// recordingFile Name
		timeStamp := time.Now().Format("20060102150405")
		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", InputVar.Type, destHost, timeStamp)
		recordingFilePath = filepath.Join(exeFileFolder, recordingFileName)

		// Go Routine: RecordingFunc
		go record.RecordingFunc(exeFileFolder, recordingFileName, bucket, recordingChan, &wgRecord)
	}

	// harvest the result
	loopClose := false
	for {
		// check loopClose Flag
		if loopClose {
			break
		}

		// select option
		select {
		case pkt, ok := <-p.ProbeChan:
			if !ok {
				loopClose = true
				break // break select, bypass "outputChan <- pkt"
			}

			// outputChan
			outputChan <- pkt

			// recordingChan
			if recording {
				recordingChan <- pkt
			}
		case err := <-errChan:
			return err
		}
	}

	// wait for the last interval (1 sec)
	time.Sleep(time.Duration(1) * time.Second)

	// close recordingChan
	if recording {
		wgRecord.Add(1)
		close(recordingChan)
		// waiting the recording function to save the last records
		wgRecord.Wait()
	} else {
		close(recordingChan)
	}

	// display testing completed
	fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
	fmt.Println("\n--- testing completed ---")

	// if recording is enabled, display the recording file path
	if recording {
		fmt.Printf("Recording CSV file is saved at: %s\n", color.GreenString(recordingFilePath))
	}

	return nil
}

// Func - DnsCommand
func DnsCommand() *cobra.Command {
	return dnsCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	dnsCmd.Flags().IntVarP(&count, "count", "c", 0, "DNS Ping Count (default 0 - Non Stop till Ctrl+C)")

	// Flag - DNS Protocol Type
	var protocol string
	dnsCmd.Flags().StringVarP(&protocol, "protocol", "o", "udp", "DNS Ping Protocol Type (default: udp)")

	// Flag - Ping timeout
	var timeout int
	dnsCmd.Flags().IntVarP(&timeout, "timeout", "t", 4, "DNS Ping Timeout (default: 4 sec)")

	// Flag - Ping interval
	var interval int
	dnsCmd.Flags().IntVarP(&interval, "interval", "i", 1, "DNS Ping Interval (default: 1 sec)")
}
