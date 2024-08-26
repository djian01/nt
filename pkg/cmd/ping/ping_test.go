package ping_test

import (
	"nt/pkg/cmd/ping"
	"testing"
)

// test ProbingFunc
func Test_ProbingFunc(t *testing.T){
	err := ping.ProbingFunc("google.com",15,56,1,false,"abc",10)
	if err != nil {
		panic(err)
	}
}