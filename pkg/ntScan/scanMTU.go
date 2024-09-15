package ntScan

import (
	"fmt"
	"nt/pkg/ntPinger"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// func - ScanMTURun
func ScanMTURun(ceilingSize int, DestAddr string, DestHost string) (err error) {
	// initial vars
	low := 576
	high := ceilingSize
	payLoadSize := 0
	var pkt ntPinger.PacketICMP
	largestMTU := 0

	// clear screen
	print("\033[H\033[2J")

	// print test title
	fmt.Printf("MAX MTU Size Check for %s:\n", color.GreenString(DestHost))
	fmt.Println(strings.Repeat("-", 80))

	// Larget MTU scan
	for {
		// ************ divide mode ************
		if (high - low) > 10 {

			// Get test MTU
			payLoadSize = getMidMtu(high,low)

			// generate payload
			payLoad := ntPinger.GeneratePayloadData(payLoadSize)			

			// IcmpProbing
			pkt, err = ntPinger.IcmpProbing(0, DestAddr, DestAddr, payLoadSize, true, 1, payLoad)
			if err != nil {
				return 
			}

			// display test result
			testTerminalOutput(DestAddr, pkt.Status, payLoadSize)

			// if the testMTU success
			if pkt.Status {
				low = payLoadSize
				// if the test MTU fail
			} else {
				high = payLoadSize
			}


		// ************** increase mode ****************
		} else {
			// update test MTU
			payLoadSize = low

			// generate payload
			payLoad := ntPinger.GeneratePayloadData(payLoadSize)			

			// IcmpProbing
			pkt, err = ntPinger.IcmpProbing(0, DestAddr, DestAddr, payLoadSize, true, 1, payLoad)
			if err != nil {
				return 
			}

			// display test result
			testTerminalOutput(DestAddr, pkt.Status, payLoadSize)

			// if the testMTU success
			if pkt.Status {
				low = payLoadSize + 1
				// if the test MTU fail
			} else {
				largestMTU = payLoadSize -1 + 28 // the larget MTU = 20 byptes (IP Header) + 8 bytes (ICMP Header) + testMTU (Payload)
				break				
			}	
		}
	}
	// print result
	fmt.Printf("\nThe MAX MTU Size to destination %s is %s bytes\n", color.CyanString(DestAddr), color.CyanString(strconv.Itoa(largestMTU)))
	fmt.Println("In this test:")
	fmt.Printf("Max MTU (%s) = IP Header (%s bytes) + ICMP Header (%s bytes) + ICMP Payload (%s bytes)\n\n", color.CyanString(strconv.Itoa(largestMTU)), color.CyanString("20"), color.CyanString("8"), color.CyanString(strconv.Itoa(payLoadSize-1)))
	return nil
}



// func - get the mid MTU
func getMidMtu(high, low int) int {
	return ((high - low)/2 + low)
}

// test Terminal output
func testTerminalOutput (DestAddr string, testStatus bool, testMTU int){
	if testStatus {
		fmt.Printf("MTU Test - Destination: %s, TestMTU Size: %s, TestResult: %s\n", color.GreenString(DestAddr), color.GreenString(strconv.Itoa(testMTU + 28)), color.GreenString(strconv.FormatBool(testStatus)))
	} else {
		fmt.Printf("MTU Test - Destination: %s, TestMTU Size: %s, TestResult: %s\n", color.GreenString(DestAddr), color.GreenString(strconv.Itoa(testMTU + 28)), color.RedString(strconv.FormatBool(testStatus)))
	}	
}