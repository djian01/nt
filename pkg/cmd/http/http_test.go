// *************************
// go test -run ^Test_HttpCommandMain$
// *************************

package Http_test

import (
	"fmt"
	Http "nt/pkg/cmd/http"
	"testing"
)



func Test_http (t *testing.T){

	a := "http://mywebsite.com:8080"

	b, _ := Http.ParseURL(a)

	fmt.Println(b)
}

func Test_HttpCommandMain(t *testing.T) {

	// initial InputVar
	HttpVarInput := Http.HttpVar {
		Scheme: "http",
		Hostname: "google.com",
		Port: 80,
		Path: "",
	}

	// Initial Other Vars
	recording := true
	displayRow := 10
	HttpMethod := "GET"
	count := 7
	timeout := 1
	interval := 1

	// call the func IcmpProbingFunc
	err := Http.HttpCommandMain(recording, displayRow, HttpVarInput, HttpMethod, count, timeout, interval)
	if err != nil {
		panic(err)
	}

}
