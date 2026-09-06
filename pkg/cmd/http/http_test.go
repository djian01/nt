// *************************
// go test -run ^Test_HttpCommandMain$
// *************************

package Http_test

import (
	"os"
	"testing"

	Http "github.com/djian01/nt/pkg/cmd/http"
	"github.com/djian01/nt/pkg/ntPinger"
)

func TestHttpCommandRejectsInvalidDisplayRows(t *testing.T) {
	input := Http.HttpVar{Scheme: "http", Hostname: "localhost", Port: 80}
	if err := Http.HttpCommandMain(false, 0, input, "GET", nil, 1, 1, 0, "none"); err == nil {
		t.Fatal("expected invalid display row error")
	}
}

func Test_http(t *testing.T) {
	parsed, err := Http.ParseURL("https://example.com/search?q=packet%20capture")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "https" || parsed.Hostname != "example.com" || parsed.Port != 443 || parsed.Path != "/search?q=packet%20capture" {
		t.Fatalf("unexpected parsed URL: %+v", parsed)
	}
	for _, input := range []string{"example.com", "ftp://example.com", "http://example.com:0"} {
		if _, err := Http.ParseURL(input); err == nil {
			t.Fatalf("expected invalid URL error for %q", input)
		}
	}
}

func Test_HttpCommandMain(t *testing.T) {
	if os.Getenv("NT_RUN_LIVE_TESTS") != "1" {
		t.Skip("set NT_RUN_LIVE_TESTS=1 to run live network tests")
	}
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
