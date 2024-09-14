package ntScan

import (
	"fmt"
	"nt/pkg/ntPinger"
	"strconv"

	"github.com/fatih/color"
)

func ScanMTUMain(highInput int, DestAddr string) (largestMTU int, err error) {
	// initial vars
	low := 576
	high := highInput
	testMTU := 0
	var pkt ntPinger.PacketICMP

	// Larget MTU scan
	for {
		// ************ divide mode ************
		if (high - low) > 10 {

			// Get test MTU
			testMTU = getMidMtu(high,low)

			// generate payload
			payLoad := ntPinger.GeneratePayloadData(testMTU)			

			// IcmpProbing
			pkt, err = ntPinger.IcmpProbing(0, DestAddr, DestAddr, testMTU, true, 1, payLoad)
			if err != nil {
				return 
			}

			// display test result
			testTerminalOutput(DestAddr, pkt.Status, testMTU)

			// if the testMTU success
			if pkt.Status {
				low = testMTU
				// if the test MTU fail
			} else {
				high = testMTU
			}


		// ************** increase mode ****************
		} else {
			// update test MTU
			testMTU = low

			// generate payload
			payLoad := ntPinger.GeneratePayloadData(testMTU)			

			// IcmpProbing
			pkt, err = ntPinger.IcmpProbing(0, DestAddr, DestAddr, testMTU, true, 1, payLoad)
			if err != nil {
				return 
			}

			// display test result
			testTerminalOutput(DestAddr, pkt.Status, testMTU)

			// if the testMTU success
			if pkt.Status {
				low = testMTU + 1
				// if the test MTU fail
			} else {
				largestMTU = testMTU -1
				return
			}
		
		}
	}
}



// func - get the mid MTU
func getMidMtu(high, low int) int {
	return ((high - low)/2 + low)
}

// test Terminal output
func testTerminalOutput (DestAddr string, testStatus bool, testMTU int){
	if testStatus {
		fmt.Printf("MTU Test - Destination: %s, TestMTU Size: %s, TestResult: %s\n", color.GreenString(DestAddr), color.GreenString(strconv.Itoa(testMTU)), color.GreenString(strconv.FormatBool(testStatus)))
	} else {
		fmt.Printf("MTU Test - Destination: %s, TestMTU Size: %s, TestResult: %s\n", color.GreenString(DestAddr), color.GreenString(strconv.Itoa(testMTU)), color.RedString(strconv.FormatBool(testStatus)))
	}	
}