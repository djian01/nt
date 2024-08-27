// *************************
// sudo go test -run ^Test_ProbingFunc$
// *************************

package icmp_test

import (
	"nt/pkg/cmd/icmp"
	"testing"
)

// test ProbingFunc
func Test_ProbingFunc(t *testing.T) {

	// defer func() to capture the panic & debug stack messages
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println(r)
	// 	}
	// }()

	// call the func IcmpProbingFunc
	err := icmp.IcmpProbingFunc("4.2.2.2", 15, 56, 1, false, "abc", 10)
	if err != nil {
		panic(err)
	}

}
