package output

import (
	"nt/pkg/ntPinger"
)

// Main func for Output
func OutputFunc(outputChan <-chan ntPinger.Packet, len int, recording bool) {

	// initial displayRows
	displayTable := []ntPinger.Packet{}

	for i := 0; i < len; i++ {
		displayTable = append(displayTable, nil)
	}

	// clear the screen
	clearScreen()

	// initial displayIdx
	displayIdx := 0

	// process Display Table from Channel NtResultChan
	for PacketPointer := range outputChan {
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

// Func - cleanScreen
func clearScreen() {
	print("\033[H\033[2J")
}
