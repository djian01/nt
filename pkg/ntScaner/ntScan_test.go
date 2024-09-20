// *************************
// sudo go test -run ^Test_ScanMTURun$
// sudo go test -run ^Test_ScanTcpRun$

// *************************

package ntScaner_test

import (
	"fmt"
	"testing"

	"github.com/djian01/nt/pkg/ntScaner"
)

func Test_ScanMTURun(t *testing.T) {

	highInput := 1500

	DestAddr := "192.168.1.1"
	DestHost := "google.com"

	err := ntScaner.ScanMTURun(highInput, DestAddr, DestHost)
	if err != nil {
		fmt.Println(err)
	}
}
