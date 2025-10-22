// *************************
// go test -run ^Test_HttpCommandMain$
// *************************

package Http_test

import (
	"fmt"
	"testing"

	Http "github.com/djian01/nt/pkg/cmd/http"
	"github.com/djian01/nt/pkg/ntPinger"
)

func Test_http(t *testing.T) {

	a := "http://mywebsite.com:8080"

	b, _ := Http.ParseURL(a)

	fmt.Println(b)
}

func Test_HttpCommandMain(t *testing.T) {

	// initial InputVar
	HttpVarInput := Http.HttpVar{
		Scheme:   "https",
		Hostname: "tiktok.com",
		Port:     443,
		Path:     "",
	}

	// Initial Other Vars
	recording := true
	displayRow := 10
	HttpMethod := "GET"
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
	count := 3
	timeout := 2
	interval := 2
	HttpProxy := "http://user01:S%40cretPass@172.16.200.102:3128"

	// call the func IcmpProbingFunc
	err := Http.HttpCommandMain(recording, displayRow, HttpVarInput, HttpMethod, StatusCodes, count, timeout, interval, HttpProxy)

	if err != nil {
		panic(err)
	}

}
