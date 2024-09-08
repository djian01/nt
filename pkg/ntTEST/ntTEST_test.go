// *************************
// go test -run ^Test_ResultGenerateICMP$
// go test -run ^Test_ResultGenerateHTTP$
// *************************

package ntTEST_test

import (
	"fmt"
	"nt/pkg/ntPinger"
	"nt/pkg/ntTEST"
	"testing"
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
