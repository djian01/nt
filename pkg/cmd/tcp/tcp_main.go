package tcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"nt/pkg/ntPinger"
	"nt/pkg/record"
	output "nt/pkg/terminalOutput"
)

// Iniital tcpCmd
var tcpCmd = &cobra.Command{
	Use:   "tcp [flags] <Destination Host> <Destination Port>", // Sub-command, shown in the -h, Usage field
	Short: "tcp Ping Test Module",
	Long:  "tcp Ping test Module for tcp testing",
	Args:  cobra.ExactArgs(2), // 2 Args, <Destination Host> <Destination Port> are required
	Run:   TcpCommandLink,
	Example: `
# Example: TCP ping to "google.com:443" with recording enabled
nt -r tcp google.com 443

# Example: TCP ping to "10.2.3.10:22" with count: 10 and interval: 2 sec
nt tcp -c 10 -i 2 10.2.3.10 22
`,
}

// Initial the bucket
var bucket = 10

// Func - IcmpCommandLink: obtain Flags and call IcmpCommandMain()
func TcpCommandLink(cmd *cobra.Command, args []string) {

	// GFlag -r
	recording, _ := cmd.Flags().GetBool("recording")

	// GFlag -d
	displayRow, _ := cmd.Flags().GetInt("displayrow")

	// Arg - destHost
	destHost := args[0]

	// Arg - destPort
	destPort, err := strconv.Atoi(args[1])
	if err != nil {
		panic("Input port number is NOT int!")
	}

	// Flag -c
	count, _ := cmd.Flags().GetInt("count")

	// Flag -s
	size, _ := cmd.Flags().GetInt("size")

	// Flag -t
	timeout, _ := cmd.Flags().GetInt("timeout")

	// Flag -i
	interval, _ := cmd.Flags().GetInt("interval")

	// call func TcpCommandMain
	err = TcpCommandMain(recording, displayRow, destHost, destPort, count, size, timeout, interval)
	if err != nil {
		panic(err) // panic all error from under
	}
}

// Func - TcpCommandMain
func TcpCommandMain(recording bool, displayRow int, destHost string, destPort int, count int, size int, timeout int, interval int) error {

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
		Type:        "tcp",
		Count:       count,
		PayLoadSize: size,
		Timeout:     timeout,
		Interval:    interval,
		DestHost:    destHost,
		DestPort:    destPort,
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
		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", "tcp", destHost, timeStamp)
		recordingFilePath = filepath.Join(exeFileFolder, recordingFileName)

		// Go Routine: RecordingFunc
		go record.RecordingFunc(recordingFilePath, bucket, recordingChan, &wgRecord)
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

	// wait for the last interval
	time.Sleep(time.Duration(interval) * time.Second)

	// if recording Enabled
	if recording {
		wgRecord.Add(1)
		close(recordingChan)
		// waiting the recording function to save the last records
		wgRecord.Wait()
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

// Func - IcmpCommand
func TcpCommand() *cobra.Command {
	return tcpCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	tcpCmd.Flags().IntVarP(&count, "count", "c", 0, "TCP Ping Count (default 0 - Non Stop till Ctrl+C)")

	// Flag - Ping size
	var size int
	tcpCmd.Flags().IntVarP(&size, "size", "s", 0, "TCP Ping Payload Size (default: 0 byte - no payload)")

	// Flag - Ping timeout
	var timeout int
	tcpCmd.Flags().IntVarP(&timeout, "timeout", "t", 4, "TCP Ping Timeout (default: 4 sec)")

	// Flag - Ping interval
	var interval int
	tcpCmd.Flags().IntVarP(&interval, "interval", "i", 1, "TCP Ping Interval (default: 1 sec)")
}
