// *************************
// sudo go test -run ^Test_Pinger$
// go test -run ^Test_bin$
// *************************

package ntPingerExample_test

import (
	"fmt"
	ntPingerExample "nt/pkg/ntPinger_example"
	"strconv"
	"testing"
)

func Test_Pinger(t *testing.T) {

	timeout := 1
	destAddr := "8.8.8.8"

	ntPingerExample.Pinger(destAddr, timeout)
}

func Test_bin(t *testing.T) {

	var a uint32 = 0b0110101111101011

	b := ^a // b is flipping all bits from a

	c := a ^ b // a XOR b, all the different bits will be "1"

	d := strconv.FormatInt(int64(c), 2) // output:  11111111111111111111111111111111

	fmt.Printf("%08s\n", d) // output: 01101011

	// fmt.Println(byte(a >> 8)) // 0b0110101111101011 -> 0b01101011 (move 8 bit right), puting the 1st 8 bits into a byte

	// fmt.Println(byte(a << 8)) // 0b0110101111101011 -> 0b00000000 (move 8 bit left, all 0)

	fmt.Println(byte(a & 0xff)) // 0b0110101111101011 -> 0b11101011, put the last 8 bits (on the right) into a byte

	// fmt.Println(a)

}
