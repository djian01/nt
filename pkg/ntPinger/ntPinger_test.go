// *************************
// go test -run ^Test_pingerTCP$
// go test -run ^Test_pingerICMP$
// go test -run ^Test_ProbingICMP$
// go test -run ^Test_ProbingHTTP$
// go test -run ^Test_pingerHTTP$
// go test -run ^Test_ProbingDNS$
// go test -run ^Test_pingerDNS$
// *************************

package ntPinger_test

import (
	"fmt"
	"testing"

	"github.com/djian01/nt/pkg/ntPinger"
)

func Test_pingerTCP(t *testing.T) {

	InputVar := ntPinger.InputVars{
		Type:     "tcp",
		Count:    0,
		Timeout:  1,
		Interval: 1,
		DestHost: "sina.com",
		DestPort: 443,
	}

	// Channel - error
	errChan := make(chan error, 1)
	defer close(errChan)

	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		panic(err)
	}

	go p.Run(errChan)

	for pkt := range p.ProbeChan {
		fmt.Println(pkt)

	}
}

func Test_pingerICMP(t *testing.T) {

	InputVar := ntPinger.InputVars{
		Type:        "icmp",
		Count:       0,
		Timeout:     1,
		Interval:    1,
		DestHost:    "4.2.2.2",
		Icmp_DF:     true,
		PayLoadSize: 1480,
	}

	// Channel - error
	errChan := make(chan error, 1)
	defer close(errChan)

	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		panic(err)
	}

	go p.Run(errChan)

	for pkt := range p.ProbeChan {
		fmt.Println(pkt)

	}
}

func Test_ProbingICMP(t *testing.T) {

	// initial testing
	DestAddr := "8.8.8.8"
	DestHost := "google.com"
	Timeout := 1
	NBypes := 8505
	Seq := 1
	size := 50
	df := true
	payload := ntPinger.GeneratePayloadData(size)

	pkt, err := ntPinger.IcmpProbing(Seq, DestAddr, DestHost, NBypes, df, Timeout, payload)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pkt)
}

// go test -run ^Test_ProbingHTTP$
func Test_ProbingHTTP(t *testing.T) {

	// Http Status Codes
	StatusCodes := []ntPinger.HttpStatusCode{
		{
			LowerCode: 200,
			UpperCode: 299,
		},
		{
			LowerCode: 304,
			UpperCode: 304,
		},
	}

	// initial testing
	InputVar := ntPinger.InputVars{
		Type:     "http",
		Count:    0,
		Timeout:  4,
		Interval: 5,
		//DestHost: "httpbin.org/status/403",
		//DestHost: "dl.broadcom.com/%3CDownloadToken%3E/PROD/COMP/VCENTER/vmw/8.0.3.00500/package-pool/155c73ec8a8de71373b88686fddbe039e9a93649d8d1de5abe8257cc91f61d56.blo",
		DestHost: "www.youtube.com/watch?v=IQl8QcZzSKU",
		//DestHost: "www.dell.com",

		DestPort:         443,
		Http_scheme:      "https",
		Http_method:      "GET",
		Http_statusCodes: StatusCodes,
		Http_path:        "",
		Http_proxy:       "http://user01:S%40cretPass@172.16.200.102:3128",
		//Http_proxy: "http://172.16.200.102:3128",
	}

	Seq := 0

	pkt, err := ntPinger.HttpProbing(
		Seq,
		InputVar.DestHost,
		InputVar.DestPort,
		InputVar.Http_path,
		InputVar.Http_scheme,
		InputVar.Http_method,
		InputVar.Http_statusCodes,
		InputVar.Timeout,
		InputVar.Http_proxy,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pkt)
}

// go test -run ^Test_pingerHTTP$
func Test_pingerHTTP(t *testing.T) {

	// initial testing
	InputVar := ntPinger.InputVars{
		Type:     "http",
		Count:    0,
		Timeout:  4,
		Interval: 1,
		//DestHost: "google.com",
		DestHost:    "www.youtube.com",
		DestPort:    443,
		Http_scheme: "https",
		Http_method: "GET",
		Http_path:   "",
		//Http_path:  "/watch?v=IQl8QcZzSKU",
		Http_proxy: "http://user01:S%40cretPass@172.16.200.102:3128",
	}

	// Channel - error
	errChan := make(chan error, 1)
	defer close(errChan)

	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		panic(err)
	}

	go p.Run(errChan)

	for pkt := range p.ProbeChan {
		fmt.Println(pkt)

	}
}

func Test_ProbingDNS(t *testing.T) {

	// initial testing
	InputVar := ntPinger.InputVars{
		Type:         "dns",
		Count:        0,
		Timeout:      4,
		Interval:     5,
		DestHost:     "172.16.200.101",
		Dns_query:    "controller01.ocplab.homelab.local",
		Dns_Protocol: "udp", // "udp" or "tcp"
	}

	Seq := 0

	pkt, err := ntPinger.DnsProbing(Seq, InputVar.DestHost, InputVar.Dns_query, InputVar.Dns_Protocol, InputVar.Timeout)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pkt)
}

func Test_pingerDNS(t *testing.T) {

	// initial testing
	InputVar := ntPinger.InputVars{
		Type:         "dns",
		Count:        0,
		Timeout:      4,
		Interval:     1,
		DestHost:     "8.8.8.8",
		Dns_query:    "www.microsoft.com",
		Dns_Protocol: "udp", // "udp" or "tcp"
	}

	// Channel - error
	errChan := make(chan error, 1)
	defer close(errChan)

	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		panic(err)
	}

	go p.Run(errChan)

	for pkt := range p.ProbeChan {
		fmt.Println(pkt)

	}
}
