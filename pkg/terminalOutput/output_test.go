// *************************
// go test -run ^Test_OutputMain$
// go test -run ^Test_ScanTablePrint$
// go test -run ^Test_TcpScanOutputFunc$
// *************************

package terminalOutput_test

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
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

	min := 1000
	max := 65535

	// Create a list of 50 empty TcpScanPort items
	Ports := make([]ntScan.TcpScanPort, 50)

	for i := 0; i < 50; i++ {
		status := (i % 3) + 1 // Cycle between statuses 1, 2, and 3
		Ports[i] = ntScan.TcpScanPort{
			ID:     i + 1,
			Port:   randomInt(min, max),
			Status: status,
		}
	}

	// Display the items in a 10×5 table
	terminalOutput.ScanTablePrint(&Ports, recording, destplayIdx, destHost)
}

func Test_TcpScanOutputFunc(t *testing.T) {

	recording := true
	destHost := "8.8.8.8"

	min := 1000
	max := 65535

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// output chan
	outputChan := make(chan *[]ntScan.TcpScanPort, 1)

	// Create a empty list of 50 items with different statuses
	Ports := make([]ntScan.TcpScanPort, 50)

	for i := 0; i < 50; i++ {
		Ports[i] = ntScan.TcpScanPort{
			ID:   i + 1,
			Port: randomInt(min, max),
		}
	}

	// go routine
	go terminalOutput.TcpScanOutputFunc(outputChan, recording, destHost)

	// loop
	loopBreak := false
	for {
		if loopBreak {
			break
		}

		select {
		case <-interruptChan: // case interruptChan, close the channel & break the loop
			close(outputChan)
			loopBreak = true
		default:
			for i := 0; i < 50; i++ {
				Ports[i].Status = randomInt(1, 3)
			}
			outputChan <- &Ports
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// randomPort generates a random integer between min and max
func randomInt(min, max int) int {

	source := rand.NewSource(time.Now().UnixNano())
	ranPort := rand.New(source).Intn(max-min+1) + min

	return ranPort
}
