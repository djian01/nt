package ntscan

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func ScanMTUMain() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <destination>")
		os.Exit(1)
	}

	dest := os.Args[1]
	ipAddr, err := net.ResolveIPAddr("ip", dest)
	if err != nil {
		fmt.Printf("Error resolving IP address: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Checking max MTU to %s\n", ipAddr.String())

	var maxMTU int
	if runtime.GOOS == "windows" {
		maxMTU = findMaxMTUWindows(dest)
	} else {
		maxMTU = findMaxMTULinux(ipAddr)
	}

	fmt.Printf("Maximum MTU to %s: %d\n", ipAddr.String(), maxMTU)
}

func findMaxMTULinux(dest *net.IPAddr) int {
	low := 576   // Minimum safe MTU
	high := 1500 // Common Ethernet MTU
	var lastWorkingMTU int

	for low <= high {
		mid := (low + high) / 2
		if sendPingLinux(dest, mid) {
			lastWorkingMTU = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return lastWorkingMTU
}

func sendPingLinux(dest *net.IPAddr, packetSize int) bool {
	// Use the Linux method for sending ICMP with raw sockets
	// (Similar to the previous code, for Linux)
	return true // Placeholder
}

func findMaxMTUWindows(dest string) int {
	low := 576
	high := 1500
	var lastWorkingMTU int

	for low <= high {
		mid := (low + high) / 2
		if sendPingWindows(dest, mid) {
			lastWorkingMTU = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return lastWorkingMTU
}

func sendPingWindows(dest string, packetSize int) bool {
	cmd := exec.Command("ping", dest, "-f", "-l", strconv.Itoa(packetSize), "-n", "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running ping: %v\n", err)
		return false
	}

	// Check for the "Packet needs to be fragmented but DF set" string in the output
	return !strings.Contains(string(output), "Packet needs to be fragmented")
}
