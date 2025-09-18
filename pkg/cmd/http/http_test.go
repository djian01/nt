// *************************
// go test -run ^Test_HttpCommandMain$
// *************************

package Http_test

import (
	"fmt"
	"testing"

	Http "github.com/djian01/nt/pkg/cmd/http"
)

func Test_http(t *testing.T) {

	a := "http://mywebsite.com:8080"

	b, _ := Http.ParseURL(a)

	fmt.Println(b)
}

func Test_HttpCommandMain(t *testing.T) {

	// initial InputVar
	HttpVarInput := Http.HttpVar{
		Scheme:   "http",
		Hostname: "google1231323.com",
		Port:     80,
		Path:     "",
	}

	// Initial Other Vars
	recording := true
	displayRow := 10
	HttpMethod := "GET"
	StatusCode := []string{"2xx", "3xx", "404"}
	count := 3
	timeout := 2
	interval := 2

	// call the func IcmpProbingFunc
	err := Http.HttpCommandMain(recording, displayRow, HttpVarInput, HttpMethod, StatusCode, count, timeout, interval)
	if err != nil {
		panic(err)
	}

}
