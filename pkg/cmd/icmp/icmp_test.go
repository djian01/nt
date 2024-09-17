// *************************
// sudo go test -run ^Test_IcmpCommandMain$
// *************************

package icmp_test

import (
	"testing"

	"github.com/djian01/nt/pkg/cmd/icmp"
)

// test ProbingFunc
func Test_IcmpCommandMain(t *testing.T) {

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
