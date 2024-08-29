// *************************
// go test -run ^Test_RecordingFunc$
// *************************

package record_test

import (
	"fmt"
	"nt/pkg/ntTEST"
	"nt/pkg/record"
	"nt/pkg/sharedStruct"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func Test_RecordingFunc(t *testing.T) {

	count := 7
	Type := "icmp"
	dest := "google.com"

	var wg sync.WaitGroup

	// channel - NtResultChan: receiving results from probing
	NtResultChan := make(chan sharedStruct.NtResult, 1)
	defer close(NtResultChan)

	// channel - recordingChan, close by the following code, no defer required
	recordingChan := make(chan sharedStruct.NtResult, 1)

	// Channel - signal pinger.Run() is done
	doneChan := make(chan bool, 1)
	defer close(doneChan)

	// record file name
	parentFilePath := "/home/uadmin/go/nt/"
	timeStamp := time.Now().Format("20060102150405")
	recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", Type, dest, timeStamp)
	recordingFilePath := filepath.Join(parentFilePath, recordingFileName)

	// go routine, Recording Func
	go record.RecordingFunc(Type, recordingFilePath, 10, recordingChan, &wg)

	// go routine, Result Generator
	go ntTEST.ResultGenerate(count, "icmp", NtResultChan, doneChan)

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
			wg.Add(1)
			close(recordingChan)
			// waiting the recording function to save the last records
			wg.Wait()
			fmt.Println("\n--- testing completed ---")
		case r := <-NtResultChan:
			recordingChan <- r
		}
	}
}
