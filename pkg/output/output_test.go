package output_test

import (
	"fmt"
	"math/rand"
	"nt/pkg/output"
	"nt/pkg/sharedStruct"
	"testing"
	"time"
)

// test - func Output
func Test_Output(t *testing.T) {

	// statistics
	var PacketsSent int
	var PacketsRecv int
	var PacketLoss float64
	var MinRtt time.Duration
	var AvgRtt time.Duration
	var MaxRtt time.Duration

	// create a NtResult Channel
	c := make(chan sharedStruct.NtResult)
	defer close(c)

	// starts func SliceProcessing
	go output.OutputFunc(c, 10)

	// create the input NtResult
	for i := 0; i < 20; i++ {

		PacketsSent++

		status := ""
		if i%2 == 0 {
			status = "PING_Failed"
		} else {
			status = "PING_OK"
			PacketsRecv++
		}

		// generate RTT
		min := 200
		max := 700
		source := rand.NewSource(time.Now().UnixNano())
		ranRTT := time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

		// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
		if MaxRtt == time.Duration(0) && status == "PING_OK" {
			MinRtt = ranRTT
			MaxRtt = ranRTT
			AvgRtt = ranRTT
			// after Initialization
		} else if status == "PING_OK" {
			if ranRTT > MaxRtt {
				MaxRtt = ranRTT
			} else if ranRTT < MinRtt {
				MinRtt = ranRTT
			}

			AvgRtt = time.Duration(((int64(AvgRtt)/1000000)*(int64(PacketsRecv-1))+(int64(ranRTT)/1000000))/int64(PacketsRecv)) * time.Millisecond

		} else {
			ranRTT = time.Duration(0)
		}

		// generate statistic
		PacketLoss = float64(PacketsRecv) / float64(PacketsSent)

		NtResult := sharedStruct.NtResult{
			Seq:       i,
			Type:      "ping",
			HostName:  "google.com",
			IP:        "1.2.3.4",
			Size:      56,
			Status:    status,
			RTT:       ranRTT,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),

			PacketsSent: PacketsSent,
			PacketsRecv: PacketsRecv,
			PacketLoss:  PacketLoss,
			MinRtt:      MinRtt,
			MaxRtt:      MaxRtt,
			AvgRtt:      AvgRtt,
		}

		c <- NtResult
		time.Sleep(1 * time.Second)
	}

}

// test - func GetAvailableSliceItem
func Test_GetAvailableSliceItem(t *testing.T) {
	// create NtResult Slice
	testSlice := []sharedStruct.NtResult{}

	len := 10

	for i := 0; i < len; i++ {
		testSlice = append(testSlice, sharedStruct.NtResult{})
	}

	for i := 0; i < 10; i++ {
		NtResult := sharedStruct.NtResult{
			Seq:       i,
			Type:      "ping",
			HostName:  "google.com",
			IP:        "1.2.3.4",
			Size:      56,
			Status:    "PING_OK",
			RTT:       45 * time.Microsecond,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		}
		testSlice[i] = NtResult
	}

	// call test function
	idx := output.GetAvailableSliceItem(&testSlice)

	fmt.Println(idx)

	fmt.Println(testSlice)

	// check output
	expectedIdx := 10
	if idx != expectedIdx {
		t.Errorf("Expected output: %v, but got: %v", expectedIdx, idx)
	}
}
