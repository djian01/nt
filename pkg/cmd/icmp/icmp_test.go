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
	recording := true
	displayRow := 13
	destHost := "192.168.1.1"
	count := 10
	size := 88
	timeout := 1
	interval := 5
	//df := true

	// call the func IcmpProbingFunc
	err := icmp.IcmpCommandMain(recording, displayRow, destHost, count, size, timeout, interval)
	if err != nil {
		panic(err)
	}

}
