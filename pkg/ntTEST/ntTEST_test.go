// *************************
// go test -run ^Test_ResultGenerateICMP$
// go test -run ^Test_ResultGenerateHTTP$
// go test -run ^Test_ResultGenerateDNS$

// *************************

package ntTEST_test

import (
	"fmt"
	"testing"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/ntTEST"
)

func Test_ResultGenerateICMP(t *testing.T) {

	count := 20
	Type := "icmp"

	// channel - probeChan: receiving results from probing
	// probeChan will be closed by the ResultGenerate()
	probeChan := make(chan ntPinger.Packet, 1)

	go ntTEST.ResultGenerate(count, Type, &probeChan)

	// start Generating Test result
	for pkt := range probeChan {

		fmt.Println(pkt)
	}

	fmt.Println("\n--- ntTEST Testing Completed ---")
}

func Test_ResultGenerateHTTP(t *testing.T) {

	count := 0
	Type := "http"

	// channel - probeChan: receiving results from probing
	// probeChan will be closed by the ResultGenerate()
	probeChan := make(chan ntPinger.Packet, 1)

	go ntTEST.ResultGenerate(count, Type, &probeChan)

	// start Generating Test result
	for pkt := range probeChan {

		fmt.Println(pkt)
	}

	fmt.Println("\n--- ntTEST Testing Completed ---")
}

func Test_ResultGenerateDNS(t *testing.T) {

	count := 0
	Type := "dns"

	// channel - probeChan: receiving results from probing
	// probeChan will be closed by the ResultGenerate()
	probeChan := make(chan ntPinger.Packet, 1)

	go ntTEST.ResultGenerate(count, Type, &probeChan)

	// start Generating Test result
	for pkt := range probeChan {

		fmt.Println(pkt)
	}

	fmt.Println("\n--- ntTEST Testing Completed ---")
}
