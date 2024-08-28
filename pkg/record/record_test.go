// *************************
// go test -run ^Test_RecordingFunc$
// *************************

package record_test

import (
	"nt/pkg/ntTEST"
	"nt/pkg/record"
	"nt/pkg/sharedStruct"
	"testing"
)

func Test_RecordingFunc(t *testing.T) {

	count := 0

	NtResultChan := make(chan sharedStruct.NtResult, 1)
	defer close(NtResultChan)

	recordingChan := make(chan sharedStruct.NtResult, 1)
	defer close(recordingChan)

	// go routine, Recording Func
	go record.RecordingFunc("icmp", "/home/uadmin/go/nt/record_test.csv", 10, recordingChan)

	// go routine, Result Generator
	go ntTEST.ResultGenerate(count, "icmp", NtResultChan)

	if count == 0 {
		for r := range NtResultChan {
			recordingChan <- r
		}
	} else {
		for i := 0; i < count; i++ {
			r := <-NtResultChan
			recordingChan <- r
		}
	}
}
