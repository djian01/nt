// *************************
// sudo go test -run ^Test_IcmpCommandMain$
// *************************

package icmp_test

import (
	"nt/pkg/cmd/icmp"
	"testing"
)

// test ProbingFunc
func Test_IcmpCommandMain(t *testing.T) {

	// initial test vars
	recording := false
	path := "abc"
	displayRow := 10
	dest := "google.com"
	count := 10
	size := 24
	interval := 1

	// call the func IcmpProbingFunc
	err := icmp.IcmpCommandMain(recording, path, displayRow, dest, count, size, interval)
	if err != nil {
		panic(err)
	}

}
