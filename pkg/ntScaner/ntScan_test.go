// *************************
// sudo go test -run ^Test_ScanMTURun$
// sudo go test -run ^Test_ScanTcpRun$

// *************************

package ntScaner_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/djian01/nt/pkg/ntScaner"
)

func Test_ScanMTURun(t *testing.T) {
	if os.Getenv("NT_RUN_LIVE_TESTS") != "1" {
		t.Skip("set NT_RUN_LIVE_TESTS=1 to run live network tests")
	}

	highInput := 1500

	DestAddr := "192.168.1.1"
	DestHost := "google.com"

	err := ntScaner.ScanMTURun(highInput, DestAddr, DestHost)
	if err != nil {
		fmt.Println(err)
	}
}
