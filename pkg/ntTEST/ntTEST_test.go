// *************************
// go test -run ^Test_ResultGenerate$
// *************************

package ntTEST_test

import (
	"fmt"
	"nt/pkg/ntTEST"
	"nt/pkg/sharedStruct"
	"testing"
)

func Test_ResultGenerate(t *testing.T) {

	count := 0

	// channel - NtResultChan: receiving results from probing
	NtResultChan := make(chan sharedStruct.NtResult, 1)
	defer close(NtResultChan)

	// Channel - signal pinger.Run() is done
	doneChan := make(chan bool, 1)
	defer close(doneChan)

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
			fmt.Println("\n--- testing completed ---")
		case r := <-NtResultChan:
			fmt.Println(r)
			// default:
			// 	// do something which will not pause the for loop
		}

	}

}
