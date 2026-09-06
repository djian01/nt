package ntScaner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/djian01/nt/pkg/ntPinger"

	"github.com/fatih/color"
)

// func - ScanMTURun
func ScanMTURun(ceilingSize int, DestAddr string, DestHost string) (err error) {
	if ceilingSize < 29 {
		return fmt.Errorf("ceiling MTU must be at least 29 bytes")
	}
	if ceilingSize > 65535 {
		return fmt.Errorf("ceiling MTU cannot exceed 65535 bytes")
	}

	// clear screen
	print("\033[H\033[2J")

	// print test title
	fmt.Printf("Max MTU Size Check for %s:\n", color.GreenString(DestHost))
	fmt.Println(strings.Repeat("-", 80))

	maxPayload := ceilingSize - 28
	payLoadSize, err := findLargestPayload(maxPayload, func(size int) (bool, error) {
		payload := ntPinger.GeneratePayloadData(size)
		pkt, probeErr := ntPinger.IcmpProbing(0, DestAddr, DestHost, size, true, 1, payload)
		if probeErr != nil {
			return false, probeErr
		}
		testTerminalOutput(DestAddr, pkt.Status, size)
		return pkt.Status, nil
	})
	if err != nil {
		return err
	}
	largestMTU := payLoadSize + 28
	// print result
	fmt.Printf("\nThe Max MTU Size to destination %s is %s bytes\n", color.CyanString(DestAddr), color.CyanString(strconv.Itoa(largestMTU)))
	fmt.Println("In this test:")
	fmt.Printf("Max MTU (%s) = IP Header (%s bytes) + ICMP Header (%s bytes) + ICMP Payload (%s bytes)\n\n", color.CyanString(strconv.Itoa(largestMTU)), color.CyanString("20"), color.CyanString("8"), color.CyanString(strconv.Itoa(payLoadSize)))
	return nil
}

func findLargestPayload(maxPayload int, probe func(int) (bool, error)) (int, error) {
	reachable, err := probe(0)
	if err != nil {
		return 0, err
	}
	if !reachable {
		return 0, fmt.Errorf("destination did not answer the minimum-size DF probe")
	}

	low, high := 1, maxPayload
	largest := 0
	for low <= high {
		mid := low + (high-low)/2
		success, err := probe(mid)
		if err != nil {
			return 0, err
		}
		if success {
			largest = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return largest, nil
}

// test Terminal output
func testTerminalOutput(DestAddr string, testStatus bool, testMTU int) {
	if testStatus {
		fmt.Printf("MTU Test - Destination: %s, TestMTU Size: %s, TestResult: %s\n", color.GreenString(DestAddr), color.GreenString(strconv.Itoa(testMTU+28)), color.GreenString(strconv.FormatBool(testStatus)))
	} else {
		fmt.Printf("MTU Test - Destination: %s, TestMTU Size: %s, TestResult: %s\n", color.GreenString(DestAddr), color.GreenString(strconv.Itoa(testMTU+28)), color.RedString(strconv.FormatBool(testStatus)))
	}
}
