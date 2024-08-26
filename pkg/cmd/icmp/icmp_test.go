package icmp_test

import (
	"nt/pkg/cmd/icmp"
	"testing"
)

// test ProbingFunc
func Test_ProbingFunc(t *testing.T) {
	err := icmp.IcmpProbingFunc("google.com", 15, 56, 1, false, "abc", 10)
	if err != nil {
		panic(err)
	}
}
