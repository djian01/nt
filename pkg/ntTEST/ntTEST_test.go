// *************************
// go test -run ^Test_ResultGenerate$
// *************************

package ntTEST_test

import (
	"fmt"
	"nt/pkg/ntPinger"
	"nt/pkg/ntTEST"
	"testing"
)

func Test_ResultGenerate(t *testing.T) {

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
