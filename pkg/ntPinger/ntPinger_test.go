// *************************
// sudo go test -run ^Test_pingerTCP$
// *************************

package ntPinger_test

import (
	"nt/pkg/ntPinger"
	"testing"
)

func Test_pingerTCP(t *testing.T) {

	InputVar := ntPinger.InputVars{
		Type:     "tcp",
		Count:    10,
		Timeout:  1,
		Interval: 1,
		DestHost: "google.com",
		DestPort: 443,
	}

	p, err := ntPinger.NewPinger(InputVar)
	if err != nil {
		panic(err)
	}

	p.Run()
}
