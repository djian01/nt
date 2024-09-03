// *************************
// sudo go test -run ^Test_OutputMain$
// *************************

package output_test

import (
	"fmt"
	"nt/pkg/ntPinger"
	"nt/pkg/ntTEST"
	"nt/pkg/output"
	"testing"
)

// test - func Output
func Test_OutputMain(t *testing.T) {

	count := 5
	Type := "tcp"
	recording := true
	displayRow := 10

	// recording row
	recordingRow := 0
	if recording {
		recordingRow = 1
	}

	// channel - NtResultChan: receiving results from probing
	probeChan := make(chan ntPinger.Packet, 1)

	// channel - NtResultChan: receiving results from probing
	OutputChan := make(chan ntPinger.Packet, 1)
	defer close(OutputChan)

	// go routine, Result Generator
	go ntTEST.ResultGenerate(count, Type, &probeChan)

	// starts func SliceProcessing
	go output.OutputFunc(OutputChan, displayRow, recording)

	for pkt := range probeChan {
		OutputChan <- pkt
	}

	// start Generating Test result
	fmt.Printf("\033[%d;1H", (displayRow + recordingRow + 7))
	fmt.Println("\n--- Output testing completed ---")

}
func Test_mytest(t *testing.T) {
	pkt := &ntPinger.PacketTCP{}

	if pkt.SendTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
		fmt.Println("ok")
	}
}
