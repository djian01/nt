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
	recording := false
	displayRow := 10
	destHost := "google.com"
	destPort := 443
	count := 0
	size := 50
	timeout := 1
	interval := 1

	// call the func IcmpProbingFunc
	err := tcp.TcpCommandMain(recording, displayRow, destHost, destPort, count, size, timeout, interval)
	if err != nil {
		panic(err)
	}

}
