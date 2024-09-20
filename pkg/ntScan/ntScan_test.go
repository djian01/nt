// *************************
// sudo go test -run ^Test_ScanMTURun$
// sudo go test -run ^Test_ScanTcpRun$

// *************************

package ntScan_test

import (
	"fmt"
	"testing"

	"github.com/djian01/nt/pkg/ntScan"
)

func Test_ScanMTURun(t *testing.T) {

	highInput := 1500

	DestAddr := "192.168.1.1"
	DestHost := "google.com"

	err := ntScan.ScanMTURun(highInput, DestAddr, DestHost)
	if err != nil {
		fmt.Println(err)
	}
}

// func Test_ScanTcpRun(t *testing.T) {

// 	recording := true
// 	destHost := "google.com"
// 	PortsArgs := []string{"22", "80", "443", "100-120", "1500"}
// 	timeout := 4

// 	// turn PortsArgs into []int{}
// 	Ports := []int{}

// 	for _, PortString := range PortsArgs {
// 		Valid, Port := tcpscan.IsValidInput(PortString)
// 		if !Valid {
// 			fmt.Println("The input Port is not a valid port")
// 		} else {
// 			Ports = append(Ports, Port...)
// 		}
// 	}

// 	err := ntScan.ScanTcpRun(recording, destHost, Ports, timeout)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
