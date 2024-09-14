// *************************
// sudo go test -run ^Test_ScanMTUMain$

// *************************

package ntScan_test

import (
	"fmt"
	"nt/pkg/ntScan"
	"testing"
)


func Test_ScanMTUMain(t *testing.T) {

	highInput := 1500

	DestAddr := "192.168.1.1"

	largestMTU, err := ntScan.ScanMTUMain(highInput,DestAddr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(largestMTU)
	}

}