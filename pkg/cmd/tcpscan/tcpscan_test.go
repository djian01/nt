// *************************
// go test -run ^Test_IsValidInput$
// go test -run ^Test_ScanTcpRun$
// *************************

package tcpscan_test

import (
	"fmt"
	"testing"

	"github.com/djian01/nt/pkg/cmd/tcpscan"
)

func Test_IsValidInput(t *testing.T) {

	fault, Ports := tcpscan.IsValidInput("100-150")
	if !fault {
		fmt.Println("Errors with input Port(s)")
	}

	fmt.Println("Input Port:")
	for _, Port := range Ports {
		fmt.Printf("%d\n", Port)
	}
}

func Test_ScanTcpRun(t *testing.T) {

	recording := false
	destHost := "google.com"
	PortsArgs := []string{"22", "80", "443", "100-120", "1500"}
	timeout := 4

	// turn PortsArgs into []int{}
	Ports := []int{}

	for _, PortString := range PortsArgs {
		Valid, Port := tcpscan.IsValidInput(PortString)
		if !Valid {
			fmt.Println("The input Port is not a valid port")
		} else {
			Ports = append(Ports, Port...)
		}
	}

	err := tcpscan.TcpScanCommandMain(recording, destHost, Ports, timeout)
	if err != nil {
		fmt.Println(err)
	}
}
