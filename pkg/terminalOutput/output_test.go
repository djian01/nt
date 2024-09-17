// *************************
// go test -run ^Test_OutputMain$
// *************************

package terminalOutput_test

import (
	"fmt"
	"testing"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/ntTEST"
	output "github.com/djian01/nt/pkg/terminalOutput"
)

// test - func Output
func Test_OutputMain(t *testing.T) {

	count := 20
	Type := "dns"
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
