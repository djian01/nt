// *************************
// sudo go test -run ^Test_TcpCommandMain$
// *************************

package tcp_test

import (
	"nt/pkg/cmd/tcp"
	"testing"
)

// test ProbingFunc
func Test_TcpCommandMain(t *testing.T) {

	// initial test vars
	recording := true
	displayRow := 10
	destHost := "sina.com"
	destPort := 443
	count := 0
	size := 1550
	timeout := 1
	interval := 1

	// call the func IcmpProbingFunc
	err := tcp.TcpCommandMain(recording, displayRow, destHost, destPort, count, size, timeout, interval)
	if err != nil {
		panic(err)
	}

}
