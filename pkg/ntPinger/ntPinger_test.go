// *************************
// sudo go test -run ^Test_pingerTCP$
// sudo go test -run ^Test_ProbingICMP$
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

func Test_ProbingICMP(t *testing.T) {

	// initial testing
	DestAddr := "4.2.2.2"
	DestHost := "google.com"
	Timeout := 4
	NBypes := 24
	df := false

	pkt, err := ntPinger.IcmpProbing(1, DestAddr, DestHost, NBypes, df, Timeout)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pkt)
}
