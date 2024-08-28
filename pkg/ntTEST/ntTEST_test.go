// *************************
// go test -run ^Test_ResultGenerate$
// *************************

package ntTEST_test

import (
	"fmt"
	"nt/pkg/ntTEST"
	"nt/pkg/sharedStruct"
	"testing"
)

func Test_ResultGenerate(t *testing.T) {

	count := 0

	NtResultChan := make(chan sharedStruct.NtResult, 1)
	defer close(NtResultChan)

	go ntTEST.ResultGenerate(count, "icmp", NtResultChan)

	if count == 0 {
		for r := range NtResultChan {
			fmt.Println(r)
		}
	} else {
		for i := 0; i < count; i++ {
			r := <-NtResultChan
			fmt.Println(r)
		}
	}

}
