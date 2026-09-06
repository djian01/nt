// *************************
// sudo go test -run ^Test_MtuCommandMain$
// *************************

package mtu_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/djian01/nt/pkg/cmd/mtu"
)

func TestMtuCommandRejectsInvalidCeiling(t *testing.T) {
	for _, ceiling := range []int{28, 65536} {
		if err := mtu.MtuCommandMain(ceiling, "localhost"); err == nil {
			t.Fatalf("expected invalid ceiling error for %d", ceiling)
		}
	}
}

func Test_MtuCommandMain(t *testing.T) {
	if os.Getenv("NT_RUN_LIVE_TESTS") != "1" {
		t.Skip("set NT_RUN_LIVE_TESTS=1 to run live network tests")
	}

	destHost := "abc.com"

	ceilingSize := 1500

	err := mtu.MtuCommandMain(ceilingSize, destHost)

	if err != nil {
		fmt.Println(err)
	}
}
