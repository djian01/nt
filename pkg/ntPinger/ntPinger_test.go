// *************************
// sudo go test -run ^Test_pingerTCP$
// sudo go test -run ^Test_pingerICMP$
// sudo go test -run ^Test_ProbingICMP$
// sudo go test -run ^Test_ProbingHTTP$
// *************************

package ntPinger_test

import (
	"fmt"
	"nt/pkg/ntPinger"
	"testing"
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
		Type:     "icmp",
		Count:    5,
		Timeout:  1,
		Interval: 1,
		DestHost: "google.com",
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
	NBypes := 50
	Seq := 1
	size := 50
	payload := ntPinger.GeneratePayloadData(size)

	pkt, err := ntPinger.IcmpProbing(Seq, DestAddr, DestHost, NBypes, Timeout, payload)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pkt)
}



func Test_ProbingHTTP(t *testing.T) {

	// initial testing
	InputVar := ntPinger.InputVars{
		Type:     "http",
		Count:    0,
		Timeout:  4,
		Interval: 1,
		DestHost: "www.microsoft.com",
		DestPort: 443,
		Http_scheme: "https",
		Http_method: "GET",
		Http_path: "en-gb",
	}

	Seq := 0

	pkt, err := ntPinger.HttpProbing(Seq, InputVar.DestHost, InputVar.DestPort, InputVar.Http_path, InputVar.Http_scheme, InputVar.Http_method, InputVar.Timeout,)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pkt)
}
