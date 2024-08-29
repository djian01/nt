package tcp

// import (
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"path/filepath"
// 	"sync"
// 	"syscall"
// 	"time"

// 	"github.com/spf13/cobra"

// 	"nt/pkg/output"
// 	"nt/pkg/record"
// 	"nt/pkg/sharedStruct"
// )

// // Iniital tcpCmd
// var tcpCmd = &cobra.Command{
// 	Use:   "tcp [flags] <host>", // Sub-command, shown in the -h, Usage field
// 	Short: "tcp Ping Test Module",
// 	Long:  "tcp Ping test Module for tcp testing",
// 	Args:  cobra.ExactArgs(1), // Only 1 Arg (dest) is required
// 	Run:   TcpCommandLink,
// }

// // Initial the bucket
// var bucket = 10

// // Func - TcpCommandLink: obtain Flags and call IcmpCommandMain()
// func TcpCommandLink(cmd *cobra.Command, args []string) {

// 	// GFlag -r
// 	recording, _ := cmd.Flags().GetBool("recording")

// 	// GFlag -d
// 	displayRow, _ := cmd.Flags().GetInt("displayrow")

// 	// Arg - dest
// 	dest := args[0]

// 	// Flag -c
// 	count, _ := cmd.Flags().GetInt("count")

// 	// Flag -s
// 	size, _ := cmd.Flags().GetInt("size")

// 	// Flag -p
// 	port, _ := cmd.Flags().GetInt("port")

// 	// Flag -i
// 	interval, _ := cmd.Flags().GetInt("interval")

// 	// call func TcpCommandMain
// 	err := TcpCommandMain(recording, displayRow, dest, port, count, size, interval)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// // Func - IcmpCommandMain
// func TcpCommandMain(recording bool, displayRow int, dest string, port int, count int, size int, interval int) error {

// 	// Wait Group
// 	var wg sync.WaitGroup

// 	// Channel - probingChan
// 	probingChan := make(chan sharedStruct.NtResult, 1)
// 	defer close(probingChan)

// 	// Channel - outputChan
// 	outputChan := make(chan sharedStruct.NtResult, 1)
// 	defer close(outputChan)

// 	// Channel - recordingChan, no need to defer close
// 	recordingChan := make(chan sharedStruct.NtResult, 1)

// 	// Channel - signal pinger.Run() is done
// 	doneChan := make(chan bool, 1)
// 	defer close(doneChan)

// 	// Create a channel to listen for SIGINT (Ctrl+C)
// 	interruptChan := make(chan os.Signal, 1)
// 	defer close(interruptChan)
// 	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	// Start Ping Main Command, manually input display Len
// 	err := IcmpProbingFunc(dest, count, size, interval, probingChan, doneChan)
// 	if err != nil {
// 		return err
// 	}

// 	// Output
// 	//// Go Routine: OutputFunc
// 	go output.OutputFunc(outputChan, displayRow)

// 	// Recording
// 	if recording {
// 		timeStamp := time.Now().Format("20060102150405")
// 		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", "icmp", dest, timeStamp)
// 		recordingFilePath := filepath.Join(recordingPath, recordingFileName)

// 		// Go Routine: RecordingFunc
// 		go record.RecordingFunc("icmp", recordingFilePath, bucket, recordingChan, &wg)
// 	}

// 	// for loop for getting the ntResult
// 	forLoopClose := false

// 	for {
// 		// check forLoopFlag
// 		if forLoopClose {
// 			break
// 		}

// 		// select chans
// 		select {
// 		case probingResult := <-probingChan:

// 			// outputChan
// 			outputChan <- probingResult

// 			// recordingChan
// 			if recording {
// 				recordingChan <- probingResult
// 			}

// 		case <-doneChan:
// 			// if recording is enabled, close the recordingchain and save the rest of the records to CSV
// 			if recording {
// 				wg.Add(1)
// 				close(recordingChan)
// 				// waiting the recording function to save the last records
// 				wg.Wait()
// 			}

// 			fmt.Printf("\033[%d;1H", (displayRow + 7))
// 			fmt.Println("\n--- testing completed ---")

// 			forLoopClose = true

// 		case <-interruptChan:
// 			// if recording is enabled, close the recordingchain and save the rest of the records to CSV
// 			if recording {
// 				wg.Add(1)
// 				close(recordingChan)
// 				// waiting the recording function to save the last records
// 				wg.Wait()
// 			}

// 			fmt.Printf("\033[%d;1H", (displayRow + 7))
// 			fmt.Println("\n--- Interrupt received, stopping testing ---")

// 			forLoopClose = true
// 		}
// 	}
// 	return nil
// }

// // Func - IcmpCommand
// func IcmpCommand() *cobra.Command {
// 	return icmpCmd
// }

// // Func - init()
// func init() {

// 	// Flag - Ping count
// 	var count int
// 	icmpCmd.Flags().IntVarP(&count, "count", "c", 0, "Ping Test Count (default 0, Ping continuous till interruption)")

// 	// Flag - Ping size
// 	var size int
// 	icmpCmd.Flags().IntVarP(&size, "size", "s", 24, "Ping Test Packet Size (must be larger than 24 Bytes)")

// 	// Flag - Ping interval
// 	var interval int
// 	icmpCmd.Flags().IntVarP(&interval, "interval", "i", 1, "Ping Test Interval")
// }
