// *************************
// sudo go test -run ^Test_Output$
// *************************

package output_test

import (
	"fmt"
	"nt/pkg/ntTEST"
	"nt/pkg/output"
	"nt/pkg/sharedStruct"
	"testing"
)

// test - func Output
func Test_OutputMain(t *testing.T) {

	count := 0
	Type := "icmp"
	recording := true
	displayRow := 10

	// recording row
	recordingRow := 0
	if recording {
		recordingRow = 1
	}

	// channel - NtResultChan: receiving results from probing
	probingChan := make(chan sharedStruct.NtResult, 1)
	defer close(probingChan)

	// channel - NtResultChan: receiving results from probing
	OutputChan := make(chan sharedStruct.NtResult, 1)
	defer close(OutputChan)

	// Channel - signal pinger.Run() is done
	doneChan := make(chan bool, 1)
	defer close(doneChan)

	// go routine, Result Generator
	go ntTEST.ResultGenerate(count, Type, probingChan, doneChan)

	// starts func SliceProcessing
	go output.OutputFunc(OutputChan, displayRow, recording)

	// start Generating Test result
	forLoopFlag := true

	for {
		// check forLoopFlag
		if !forLoopFlag {
			break
		}
		select {
		case <-doneChan:
			forLoopFlag = false
			fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
			fmt.Println("\n--- Output testing completed ---")
		case r := <-probingChan:
			OutputChan <- r
		}
	}

}
