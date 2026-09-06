// *************************
// go test -run ^Test_IsValidInput$
// go test -run ^Test_ScanTcpRun$
// *************************

package tcpscan_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/djian01/nt/pkg/cmd/tcpscan"
)

func TestIsValidInputRejectsInvalidPorts(t *testing.T) {
	for _, input := range []string{"0", "-1", "65536", "0-5"} {
		valid, _ := tcpscan.IsValidInput(input)
		if valid {
			t.Fatalf("port input %q was accepted", input)
		}
	}
}

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
	if os.Getenv("NT_RUN_LIVE_TESTS") != "1" {
		t.Skip("set NT_RUN_LIVE_TESTS=1 to run live network tests")
	}

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
