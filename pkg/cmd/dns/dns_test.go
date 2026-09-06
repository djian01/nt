// *************************
// go test -run ^Test_DnsCommandMain$
// *************************

package dns_test

import (
	"os"
	"testing"

	"github.com/djian01/nt/pkg/cmd/dns"
)

func TestDnsCommandRejectsInvalidDisplayRows(t *testing.T) {
	if err := dns.DnsCommandMain(false, 0, "127.0.0.1", "localhost", "udp", 1, 1, 0); err == nil {
		t.Fatal("expected invalid display row error")
	}
}

// test ProbingFunc
func Test_DnsCommandMain(t *testing.T) {
	if os.Getenv("NT_RUN_LIVE_TESTS") != "1" {
		t.Skip("set NT_RUN_LIVE_TESTS=1 to run live network tests")
	}
	t.Chdir(t.TempDir())

	// initial test vars
	recording := true
	displayRow := 10
	destHost := "8.8.8.8"
	Dns_query := "www.youtube.com"
	Dns_protocol := "udp"
	count := 4
	timeout := 1
	interval := 1

	// call the func IcmpProbingFunc
	err := dns.DnsCommandMain(recording, displayRow, destHost, Dns_query, Dns_protocol, count, timeout, interval)
	if err != nil {
		panic(err)
	}

}
