// *************************
// sudo go test -run ^Test_MtuCommandMain$
// *************************

package mtu_test

import (
	"fmt"
	"nt/pkg/cmd/mtu"
	"testing"
)

func Test_MtuCommandMain(t *testing.T) {

	destHost := "abc.com"

	ceilingSize := 1500

	err := mtu.MtuCommandMain(ceilingSize, destHost)

	if err != nil {
		fmt.Println(err)
	}
}
