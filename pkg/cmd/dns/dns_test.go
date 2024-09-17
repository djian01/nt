// *************************
// go test -run ^Test_DnsCommandMain$
// *************************

package dns_test

import (
	"testing"

	"github.com/djian01/nt/pkg/cmd/dns"
)

// test ProbingFunc
func Test_DnsCommandMain(t *testing.T) {

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
