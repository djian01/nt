package record_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/djian01/nt/pkg/ntPinger"
	"github.com/djian01/nt/pkg/record"
)

func TestRecordingFuncFlushesAndRenames(t *testing.T) {
	directory := t.TempDir()
	fileName := "Record_icmp_localhost.csv"
	recordingChan := make(chan ntPinger.Packet, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go record.RecordingFunc(directory, fileName, 10, recordingChan, &wg)

	recordingChan <- &ntPinger.PacketICMP{
		Type:        "icmp",
		Seq:         1,
		Status:      true,
		DestHost:    "localhost",
		DestAddr:    "127.0.0.1",
		PayLoadSize: 16,
		RTT:         time.Millisecond,
		SendTime:    time.Now(),
		PacketsSent: 1,
		PacketsRecv: 1,
	}
	close(recordingChan)
	wg.Wait()

	completedPath := filepath.Join(directory, fileName)
	contents, err := os.ReadFile(completedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "Type,Seq,Status") || !strings.Contains(string(contents), "icmp,1,true,localhost,127.0.0.1,16") {
		t.Fatalf("unexpected recording contents: %s", contents)
	}
	if _, err := os.Stat(filepath.Join(directory, "[recording...]_"+fileName)); !os.IsNotExist(err) {
		t.Fatalf("temporary recording still exists: %v", err)
	}
}
