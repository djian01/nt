package output

import (
	"nt/pkg/sharedStruct"
)

// Main func for Output
func OutputFunc(outputChan <-chan sharedStruct.NtResult, len int) {

	// initial displayRows
	displayTable := []sharedStruct.NtResult{}

	for i := 0; i < len; i++ {
		displayTable = append(displayTable, sharedStruct.NtResult{})
	}

	// clear the screen
	clearScreen()

	// process Display Table from Channel NtResultChan
	for NtResult := range outputChan {
		idx := GetAvailableSliceItem(&displayTable)
		displayTable[idx] = NtResult
		TablePrint(&displayTable, len)
	}
}

// Get the Slice Idx to inject the Result Object
func GetAvailableSliceItem(displayTable *[]sharedStruct.NtResult) int {

	// check if the slice been filled up, if not provide the empty slice id
	for i := 0; i < len(*displayTable); i++ {
		if (*displayTable)[i].Timestamp == "" {
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
