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
	"nt/pkg/output"
	"nt/pkg/record"
)

// Iniital tcpCmd
var tcpCmd = &cobra.Command{
	Use:   "tcp [flags] <Destination Host> <Destination Port>", // Sub-command, shown in the -h, Usage field
	Short: "tcp Ping Test Module",
	Long:  "tcp Ping test Module for tcp testing",
	Args:  cobra.ExactArgs(2), // 2 Args, <Destination Host> <Destination Port> are required
	Run:   TcpCommandLink,
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
		panic(err)
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

	// Channel - outputChan
	outputChan := make(chan ntPinger.Packet, 1)
	defer close(outputChan)

	// Channel - recordingChan, closed in the end of the testing, no need to defer close
	recordingChan := make(chan ntPinger.Packet, 1)

	// build the InputVar
	InputVar := ntPinger.InputVars{
		Type:     "tcp",
		Count:    count,
		NBypes:   size,
		Timeout:  timeout,
		Interval: interval,
		DestHost: destHost,
		DestPort: destPort,
	}

	// Start Ping Main Command, manually input display Len
	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		panic(err)
	}

	go p.Run()

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
	for pkt := range p.ProbeChan {
		// outputChan
		outputChan <- pkt

		// recordingChan
		if recording {
			recordingChan <- pkt
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

	// forLoopClose := false

	// for {
	// 	// check forLoopFlag
	// 	if forLoopClose {
	// 		break
	// 	}

	// 	// select chans
	// 	select {
	// 	case probingResult := <-probingChan:

	// 		// outputChan
	// 		outputChan <- probingResult

	// 		// recordingChan
	// 		if recording {
	// 			recordingChan <- probingResult
	// 		}

	// 	case <-doneChan:
	// 		// if recording is enabled, close the recordingchain and save the rest of the records to CSV
	// 		if recording {
	// 			wg.Add(1)
	// 			close(recordingChan)
	// 			// waiting the recording function to save the last records
	// 			wg.Wait()
	// 		}

	// 		fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
	// 		fmt.Println("\n--- testing completed ---")

	// 		forLoopClose = true

	// 	case <-interruptChan:
	// 		// if recording is enabled, close the recordingchain and save the rest of the records to CSV
	// 		if recording {
	// 			wg.Add(1)
	// 			close(recordingChan)
	// 			// waiting the recording function to save the last records
	// 			wg.Wait()
	// 		}

	// 		fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
	// 		fmt.Println("\n--- Interrupt received, stopping testing ---")

	// 		forLoopClose = true
	// 	}
	// }
	// if recording {
	// 	fmt.Printf("Recording CSV file is saved at: %s\n", color.GreenString(recordingFilePath))
	// }
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
	tcpCmd.Flags().IntVarP(&count, "count", "c", 0, "TCP Ping Test Count (default 0, Ping continuous till interruption)")

	// Flag - Ping size
	var size int
	tcpCmd.Flags().IntVarP(&size, "size", "s", 0, "TCP Ping Test Payload Size (Default value is 0 byte, no payload)")

	// Flag - Ping timeout
	var timeout int
	tcpCmd.Flags().IntVarP(&timeout, "timeout", "t", 4, "TCP Ping Test Count (default 4 seconds)")

	// Flag - Ping interval
	var interval int
	tcpCmd.Flags().IntVarP(&interval, "interval", "i", 1, "TCP Ping Test Interval (default 1 second)")
}
