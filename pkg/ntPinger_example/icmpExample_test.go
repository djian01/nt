// *************************
// sudo go test -run ^Test_Pinger$
// *************************

package ntPingerExample_test

import (
	ntPingerExample "nt/pkg/ntPinger_example"
	"testing"
)

func Test_Pinger(t *testing.T) {

	timeout := 1
	destAddr := "8.8.8.128"

	ntPingerExample.Pinger(destAddr, timeout)
}
