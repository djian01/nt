package icmp

import "github.com/spf13/cobra"

// Iniital icmpCmd
var icmpCmd = &cobra.Command{
	Use:   "icmp [flags] <host>", // Sub-command, shown in the -h, Usage field
	Short: "ICMP Ping Test Module",
	Long:  "ICMP Ping test Module for ICMP testing",
	Args:  cobra.ExactArgs(1), // Only 1 Arg (dest) is required
	Run:   IcmpCommandLink,
	Example: `
# Example: ICMP ping to "google.com" with recording enabled
nt -r icmp google.com

# Example: ICMP ping to "10.2.3.10" with count: 10, interval: 2 sec and payload 48 bytes
nt icmp -c 10 -i 2 -s 48 10.2.3.10 22
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

	// Arg - dest
	dest := args[0]

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
	err := IcmpCommandMain(recording, displayRow, dest, count, size, timeout, interval, df)
	if err != nil {
		panic(err)
	}
}

// Func - IcmpCommandMain
func IcmpCommandMain(recording bool, displayRow int, dest string, count int, size int, timeout int, interval int, df bool) error {

	return nil
}

// 	// Wait Group
// 	var wg sync.WaitGroup

// 	// recording row
// 	recordingRow := 0
// 	if recording {
// 		recordingRow = 1
// 	}

// 	// recordingFilePath
// 	recordingFilePath := ""

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
// 	go output.OutputFunc(outputChan, displayRow, recording)

// 	// Recording
// 	if recording {

// 		// recordingFile Path
// 		exeFileFolder, err := os.Getwd()
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			os.Exit(1)
// 		}

// 		// recordingFile Name
// 		timeStamp := time.Now().Format("20060102150405")
// 		recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", "icmp", dest, timeStamp)
// 		recordingFilePath = filepath.Join(exeFileFolder, recordingFileName)

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

// 			fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
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

// 			fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
// 			fmt.Println("\n--- Interrupt received, stopping testing ---")

// 			forLoopClose = true
// 		}
// 	}
// 	if recording {
// 		fmt.Printf("Recording CSV file is saved at: %s\n", color.GreenString(recordingFilePath))
// 	}
// 	return nil
// }

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
	icmpCmd.Flags().IntVarP(&size, "size", "s", 32, "ICMP Ping Payload Size (default: 32 bytes - Total Packet Size = 44 bytes header + Payload Size)")

	// Flag - Ping timeout
	var timeout int
	icmpCmd.Flags().IntVarP(&timeout, "timeout", "t", 4, "ICMP Ping Timeout (default: 4 sec)")

	// Flag - Ping interval
	var interval int
	icmpCmd.Flags().IntVarP(&interval, "interval", "i", 1, "ICMP Ping Interval (default: 1 sec)")

	// Flag - de-fregmentation bit
	var df bool
	icmpCmd.Flags().BoolVarP(&df, "df", "d", false, "ICMP Ping de-fregmentation (default: false)")
}
