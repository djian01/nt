package terminalOutput

import (
	"os"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/ntScaner"
	"golang.org/x/term"
)

// Main func for Output
func OutputFunc(outputChan <-chan ntPinger.Packet, len int, recording bool) {

	// initial displayRows
	displayTable := []ntPinger.Packet{}

	for i := 0; i < len; i++ {
		displayTable = append(displayTable, nil)
	}

	// clear the screen
	ClearScreen()

	// initial displayIdx
	displayIdx := 0

	// get current terminal window size
	windowW, windowH, err := getTerminaSize()
	if err != nil {
		return
	}

	// process Display Table from Channel NtResultChan
	for PacketPointer := range outputChan {
		w, h, _ := getTerminaSize()

		// if the terminal window size changed, clear Screen
		if w != windowW || h != windowH {
			windowW = w
			windowH = h
			ClearScreen()
		}

		idx := GetAvailableSliceItem(&displayTable)
		displayTable[idx] = PacketPointer
		TablePrint(&displayTable, len, recording, displayIdx)
		displayIdx++
	}
}

// Get the Slice Idx to inject the Result Object
func GetAvailableSliceItem(displayTable *[]ntPinger.Packet) int {

	// check if the slice been filled up, if not provide the empty slice id
	for i := 0; i < len(*displayTable); i++ {

		if (*displayTable)[i] == nil {
			return i
		}
	}

	// if all the slots have been filled, re-range the seq and provide the last index
	for i := 1; i < len(*displayTable); i++ {
		(*displayTable)[i-1] = (*displayTable)[i]
	}
	return (len(*displayTable) - 1)
}

// Main func for TCPScanOutput
func TcpScanOutputFunc(outputChan <-chan *[]ntScaner.TcpScanPort, recording bool, destHost string) {

	// clear the screen
	ClearScreen()

	// initial displayIdx
	displayIdx := 0

	// process Display Table from Channel NtResultChan
	for Ports := range outputChan {
		ScanTablePrint(Ports, recording, displayIdx, destHost)
		displayIdx++
	}
}

// Func - cleanScreen
func ClearScreen() {
	print("\033[H\033[2J")
}

// Func - get the terminal window size
func getTerminaSize() (width int, height int, err error) {
	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	return
}
