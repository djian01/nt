// *************************
// sudo go test -run ^Test_IcmpCommandMain$
// *************************

package icmp_test

import (
	"os"
	"testing"

	"github.com/djian01/nt/pkg/cmd/icmp"
)

func TestIcmpCommandRejectsInvalidDisplayRows(t *testing.T) {
	if err := icmp.IcmpCommandMain(false, 0, "localhost", 1, 8, false, 1, 0); err == nil {
		t.Fatal("expected invalid display row error")
	}
}

// test ProbingFunc
func Test_IcmpCommandMain(t *testing.T) {
	if os.Getenv("NT_RUN_LIVE_TESTS") != "1" {
		t.Skip("set NT_RUN_LIVE_TESTS=1 to run live network tests")
	}
	t.Chdir(t.TempDir())

	// initial test vars
	recording := true
	displayRow := 13
	destHost := "192.168.1.1"
	count := 4
	size := 88
	timeout := 1
	interval := 1
	df := true

	// call the func IcmpProbingFunc
	err := icmp.IcmpCommandMain(recording, displayRow, destHost, count, size, df, timeout, interval)
	if err != nil {
		panic(err)
	}

}
