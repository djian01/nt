package terminalOutput

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/djian01/nt/pkg/sharedstruct"
	"github.com/fatih/color"
)

// Set the rows & cols
const (
	rows = 5
	cols = 10
)

// ScanTablePrint function to display the values in a 10x5 table with color based on status (The input item lengh is 50)
func ScanTablePrint(Ports *[]sharedstruct.TcpScanPort, recording bool, displayIdx int, destHost string) {

	// clear the screen
	ClearScreen()

	// Set the 1st row for table head
	var tableHeadRowIdx int
	if recording {
		tableHeadRowIdx = 1
	} else {
		tableHeadRowIdx = 0
	}

	// Display REC if recording enabled
	if recording {
		moveToRow(1)
		if displayIdx%2 == 0 {
			fmt.Printf("%s", color.RedString("REC ●"))
		} else {
			fmt.Printf("%s", color.RedString("REC    "))
		}
	}

	// print table head
	moveToRow(tableHeadRowIdx + 1)
	fmt.Printf("TCP Scan for %s\n", color.CyanString(destHost))
	fmt.Println(strings.Repeat("-", 100))

	// Define grey color using FgHiBlack (highlighted black which appears as grey)
	grey := color.New(color.FgHiBlack).SprintFunc()

	// initial index
	index := 0

	// printing table
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			Port := (*Ports)[index]

			// Set color based on status
			switch Port.Status {

			case 1: // Not Been Tested
				cell := fmt.Sprintf("%-10d", Port.Port)
				fmt.Printf("%s", grey(cell))

			case 2: // Test Success
				cell := fmt.Sprintf("%-10d", Port.Port)
				fmt.Printf("%s", color.GreenString(cell))

			case 3: // Test Failed
				cell := fmt.Sprintf("%-10d", Port.Port)
				fmt.Printf("%s", color.RedString(cell))

			default: // invisible
				cell := fmt.Sprintf("%-10s", "")
				fmt.Printf("%s", cell) // Default (empty string)
			}

			index++
		}
		// new row
		fmt.Println()
	}

	// Get statistic
	countTested, countSuccess, countFail := TcpScanStat(Ports)

	// Print statistic
	fmt.Printf("\n")
	fmt.Printf("Tested Port(s): %s, Success Port(s): %s, Failed Port(s): %s \n", color.CyanString(strconv.Itoa(countTested)), color.CyanString(strconv.Itoa(countSuccess)), color.CyanString(strconv.Itoa(countFail)))
	fmt.Printf("\n")

}

// func tcpScanStat
func TcpScanStat(Ports *[]sharedstruct.TcpScanPort) (countTested int, countSuccess int, countFail int) {

	for _, port := range *Ports {
		if port.Status == 2 {
			countSuccess++
			countTested++
		} else if port.Status == 3 {
			countFail++
			countTested++
		}
	}
	return
}
