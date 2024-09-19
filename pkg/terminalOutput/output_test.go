// *************************
// go test -run ^Test_OutputMain$
// go test -run ^Test_ScanTablePrint$
// *************************

package terminalOutput_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/ntScan"
	"github.com/djian01/nt/pkg/ntTEST"
	"github.com/djian01/nt/pkg/terminalOutput"
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

func Test_ScanTablePrint(t *testing.T) {

	terminalOutput.ClearScreen()

	destplayIdx := 11
	recording := true
	destHost := "8.8.8.8"

	// Create a list of 50 items with different statuses
	Ports := make([]ntScan.TcpScanPort, 50)

	for i := 0; i < 50; i++ {
		status := (i % 3) + 1 // Cycle between statuses 1, 2, and 3
		Ports[i] = ntScan.TcpScanPort{
			ID:     i + 1,
			Port:   randomPort(),
			Status: status,
		}
	}

	// Display the items in a 10×5 table
	terminalOutput.ScanTablePrint(Ports, recording, destplayIdx, destHost)
}

// randomPort generates a random integer between min and max
func randomPort() int {

	min := 1000
	max := 65535

	source := rand.NewSource(time.Now().UnixNano())
	ranPort := rand.New(source).Intn(max-min+1) + min

	return ranPort
}
