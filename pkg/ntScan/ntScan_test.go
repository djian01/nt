// *************************
// sudo go test -run ^Test_ScanMTURun$

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
