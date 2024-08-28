package icmp

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"nt/pkg/output"
	"nt/pkg/record"
	"nt/pkg/sharedStruct"
)

// Iniital pingCmd
var icmpCmd = &cobra.Command{
	Use:   "icmp [flags] <host>", // Sub-command, shown in the -h, Usage field
	Short: "ICMP Ping Test Module",
	Long:  "ICMP Ping test Module for ICMP testing",
	Args:  cobra.ExactArgs(1), // Only 1 Arg (dest) is required
	Run:   IcmpCommandLink,
}

// Initial the bucket
var bucket = 10

// Func - IcmpCommandLink: obtain Flags and call IcmpCommandMain()
func IcmpCommandLink(cmd *cobra.Command, args []string) {

	// GFlag -r
	recording, _ := cmd.Flags().GetBool("recording")

	// GFlag -p
	recordingPath, _ := cmd.Flags().GetString("path")

	// GFlag -d
	displayRow, _ := cmd.Flags().GetInt("displayrow")

	// Arg - dest
	dest := args[0]

	// Flag -c
	count, _ := cmd.Flags().GetInt("count")

	// Flag -s
	size, _ := cmd.Flags().GetInt("size")

	// Flag -i
	interval, _ := cmd.Flags().GetInt("interval")

	// call func IcmpCommandMain
	err := IcmpCommandMain(recording, recordingPath, displayRow, dest, count, size, interval)
	if err != nil {
		panic(err)
	}
}

// Func - IcmpCommandMain
func IcmpCommandMain(recording bool, recordingPath string, displayRow int, dest string, count int, size int, interval int) error {

	// Channel - signal pinger.Run() is done
	doneChan := make(chan bool, 1)
	defer close(doneChan)

	// Channel - probingChan
	probingChan := make(chan sharedStruct.NtResult, 1)
	defer close(probingChan)

	// Channel - outputChan
	outputChan := make(chan sharedStruct.NtResult, 1)
	defer close(outputChan)

	// Channel - recordingChan
	recordingChan := make(chan sharedStruct.NtResult, 1)
	defer close(recordingChan)

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start Ping Main Command, manually input display Len
	err := IcmpProbingFunc(dest, count, size, interval, probingChan, doneChan)
	if err != nil {
		return err
	}

	// Output
	//// Go Routine: OutputFunc
	go output.OutputFunc(outputChan, displayRow)

	// Recording
	if recording {
		timeStamp := time.Now().Format("20060102150405")
		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", "icmp", dest, timeStamp)
		recordingFilePath := filepath.Join(recordingPath, recordingFileName)

		// Go Routine: RecordingFunc
		go record.RecordingFunc("icmp", recordingFilePath, bucket, recordingChan)
	}

	// for loop for getting the ntResult
	for {
		select {
		case probingResult := <-probingChan:

			// outputChan
			outputChan <- probingResult

			// recordingChan
			if recording {
				recordingChan <- probingResult
			}

		case <-doneChan:
			fmt.Printf("\033[%d;1H", (displayRow + 7))
			fmt.Println("\n--- testing completed ---")

			// if recording is enabled, close the recordingchain and save the rest of the records to CSV
			if recording {
				close(recordingChan)
				time.Sleep(1 * time.Second) // giving time for saving the records
			}

			return nil

		case <-interruptChan:
			fmt.Printf("\033[%d;1H", (displayRow + 7))
			fmt.Println("\n--- Interrupt received, stopping testing ---")

			// if recording is enabled, close the recordingchain and save the rest of the records to CSV
			if recording {
				close(recordingChan)
				time.Sleep(1 * time.Second) // giving time for saving the records
			}

			return nil
		}
	}
}

// Func - IcmpCommand
func IcmpCommand() *cobra.Command {
	return icmpCmd
}

// Func - init()
func init() {

	// Flag - Ping count
	var count int
	icmpCmd.Flags().IntVarP(&count, "count", "c", 0, "Ping Test Count")

	// Flag - Ping size
	var size int
	icmpCmd.Flags().IntVarP(&size, "size", "s", 24, "Ping Test Packet Size (must be larger than 24 Bytes)")

	// Flag - Ping interval
	var interval int
	icmpCmd.Flags().IntVarP(&interval, "interval", "i", 1, "Ping Test Interval")
}
