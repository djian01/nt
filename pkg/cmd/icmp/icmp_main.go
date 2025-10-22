package icmp

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/record"

	output "github.com/djian01/nt/pkg/terminalOutput"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Initial icmpCmd
var icmpCmd = &cobra.Command{
	Use:   "icmp [flags] <host>", // Sub-command, shown in the -h, Usage field
	Short: "ICMP Ping Test Module",
	Long:  "ICMP Ping test Module for ICMP testing",
	Args:  cobra.ExactArgs(1), // Only 1 Arg (dest) is required
	Run:   IcmpCommandLink,
	Example: `
# Example: ICMP ping to "google.com" with recording enabled
nt -r icmp google.com

# Example: ICMP ping to "10.2.3.10" with count: 10, interval: 2 sec,  payload 48 bytes
nt icmp -c 10 -i 2 -s 48 10.2.3.10
`,
}

// Initial the bucket
var bucket = 10

// Func - IcmpCommandLink: obtain Flags and call IcmpCommandMain()
func IcmpCommandLink(cmd *cobra.Command, args []string) {

	// GFlag -r
	recording, _ := cmd.Flags().GetBool("recording")

	// GFlag -d
	displayRow, _ := cmd.Flags().GetInt("displayrow")

	// Arg - destHost
	destHost := args[0]

	// Flag -c
	count, _ := cmd.Flags().GetInt("count")

	// Flag -s
	size, _ := cmd.Flags().GetInt("size")

	// Flag -i
	interval, _ := cmd.Flags().GetInt("interval")

	// Flag -t
	timeout, _ := cmd.Flags().GetInt("timeout")

	// Flag -d
	df, _ := cmd.Flags().GetBool("df")

	// call func IcmpCommandMain
	err := IcmpCommandMain(recording, displayRow, destHost, count, size, df, timeout, interval)
	if err != nil {
		panic(err)
	}
}

// Func - IcmpCommandMain
func IcmpCommandMain(recording bool, displayRow int, destHost string, count int, size int, df bool, timeout int, interval int) error {

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
		Type:        "icmp",
		Count:       count,
		PayLoadSize: size,
		Timeout:     timeout,
		Interval:    interval,
		DestHost:    destHost,
		Icmp_DF:     df,
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

// Func - IcmpCommand
func IcmpCommand() *cobra.Command {
	return icmpCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	icmpCmd.Flags().IntVarP(&count, "count", "c", 0, "ICMP Ping Count (default 0 - Non Stop till Ctrl+C)")

	// Flag - Ping Payload size
	var size int
	// Total L2 Frame Size  = L2 Header (14) + L3 Header (20) + ICMP Header (8) + Payload
	// Total L3 Packet Size = L3 Header (20) + ICMP Header (8) + Payload
	icmpCmd.Flags().IntVarP(&size, "size", "s", 32, "ICMP Ping Payload Size (default: 32 bytes)")

	// Flag - Ping timeout
	var timeout int
	icmpCmd.Flags().IntVarP(&timeout, "timeout", "t", 4, "ICMP Ping Timeout (default: 4 sec)")

	// Flag - Ping interval
	var interval int
	icmpCmd.Flags().IntVarP(&interval, "interval", "i", 1, "ICMP Ping Interval (default: 1 sec)")

	// Flag - de-fregmentation bit
	var df bool
	icmpCmd.Flags().BoolVarP(&df, "df", "d", false, "ICMP Ping de-fregmentation")
}
