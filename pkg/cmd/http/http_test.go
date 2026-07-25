// *************************
// go test -run ^Test_HttpCommandMain$
// *************************

package Http_test

import (
	"fmt"
	"os"
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
	t.Chdir(t.TempDir())

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
	HttpProxy := os.Getenv("NT_HTTP_PROXY")

	// call the func IcmpProbingFunc
	err := Http.HttpCommandMain(recording, displayRow, HttpVarInput, HttpMethod, StatusCodes, count, timeout, interval, HttpProxy)

	if err != nil {
		panic(err)
	}

}
