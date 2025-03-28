// *************************
// go test -run ^Test_RecordingFunc$
// *************************

package record_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/ntTEST"
	"github.com/djian01/nt/pkg/record"
)

func Test_RecordingFunc(t *testing.T) {

	count := 20
	Type := "dns"
	dest := "8.8.8.8"

	var wg sync.WaitGroup

	// channel - NtResultChan: receiving results from probing
	probeChan := make(chan ntPinger.Packet, 1)

	// channel - recordingChan, close by the following code, no defer required
	recordingChan := make(chan ntPinger.Packet, 1)

	// record file name

	// recordingFile Path
	exeFileFolder, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// recordingFile Name
	timeStamp := time.Now().Format("20060102150405")
	recordingFileName := fmt.Sprintf("Record_%v_%v_%v.csv", Type, dest, timeStamp)

	// go routine, Recording Func
	go record.RecordingFunc(exeFileFolder, recordingFileName, 10, recordingChan, &wg)

	// go routine, Result Generator
	go ntTEST.ResultGenerate(count, Type, &probeChan)

	// start Generating Test result

	for pkg := range probeChan {
		// record the probe result
		fmt.Println(pkg)
		recordingChan <- pkg
	}

	wg.Add(1)
	close(recordingChan)
	// waiting the recording function to save the last records
	wg.Wait()
	fmt.Println("\n--- testing completed ---")
}
