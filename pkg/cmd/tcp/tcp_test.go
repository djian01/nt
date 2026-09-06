// *************************
// go test -run ^Test_TcpCommandMain$
// *************************

package tcp_test

import (
	"os"
	"testing"

	"github.com/djian01/nt/pkg/cmd/tcp"
)

func TestTcpCommandRejectsInvalidDisplayRows(t *testing.T) {
	if err := tcp.TcpCommandMain(false, 0, "localhost", 443, 1, 0, 1, 0); err == nil {
		t.Fatal("expected invalid display row error")
	}
}

// test ProbingFunc
func Test_TcpCommandMain(t *testing.T) {
	if os.Getenv("NT_RUN_LIVE_TESTS") != "1" {
		t.Skip("set NT_RUN_LIVE_TESTS=1 to run live network tests")
	}
	t.Chdir(t.TempDir())

	// initial test vars
	recording := true
	displayRow := 10
	destHost := "google.com"
	destPort := 80
	count := 4
	size := 50
	timeout := 1
	interval := 1

	// call the func IcmpProbingFunc
	err := tcp.TcpCommandMain(recording, displayRow, destHost, destPort, count, size, timeout, interval)
	if err != nil {
		panic(err)
	}

}
