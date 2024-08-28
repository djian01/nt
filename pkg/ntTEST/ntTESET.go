package ntTEST

import (
	"math/rand"
	"nt/pkg/sharedStruct"
	"time"
)

func ResultGenerate(count int, Type string, NtResultChain chan sharedStruct.NtResult) {

	// statistics
	var PacketsSent = 0
	var PacketsRecv = 0
	var PacketLoss float64
	var MinRtt time.Duration
	var AvgRtt time.Duration
	var MaxRtt time.Duration

	// initial loopCount
	loopCount := 0

	// endless loop
	if count == 0 {

		// create the input NtResult
		for {
			loopCount++
			PacketsSent++
			PacketsRecv++
			status := "OK"

			// generate RTT
			min := 200
			max := 700
			source := rand.NewSource(time.Now().UnixNano())
			ranRTT := time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

			// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
			if loopCount == 1 {
				MinRtt = ranRTT
				MaxRtt = ranRTT
				AvgRtt = ranRTT
				// after Initialization
			} else {
				if ranRTT > MaxRtt {
					MaxRtt = ranRTT
				} else if ranRTT < MinRtt {
					MinRtt = ranRTT
				}
				AvgRtt = time.Duration(((int64(AvgRtt)/1000000)*(int64(PacketsRecv-1))+(int64(ranRTT)/1000000))/int64(PacketsRecv)) * time.Millisecond
			}

			// generate statistic
			PacketLoss = (1 - float64(PacketsRecv)/float64(PacketsSent))

			NtResult := sharedStruct.NtResult{
				Seq:       loopCount - 1,
				Type:      Type,
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

			NtResultChain <- NtResult
			time.Sleep(1 * time.Second)
		}

	} else {
		// create the input NtResult
		for i := 0; i < count; i++ {

			PacketsSent++
			PacketsRecv++
			status := "OK"

			// generate RTT
			min := 200
			max := 700
			source := rand.NewSource(time.Now().UnixNano())
			ranRTT := time.Duration(rand.New(source).Intn(max-min+1)+min) * time.Millisecond

			// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
			if i == 0 {
				MinRtt = ranRTT
				MaxRtt = ranRTT
				AvgRtt = ranRTT
				// after Initialization
			} else {
				if ranRTT > MaxRtt {
					MaxRtt = ranRTT
				} else if ranRTT < MinRtt {
					MinRtt = ranRTT
				}
				AvgRtt = time.Duration(((int64(AvgRtt)/1000000)*(int64(PacketsRecv-1))+(int64(ranRTT)/1000000))/int64(PacketsRecv)) * time.Millisecond
			}

			// generate statistic
			PacketLoss = float64(PacketsRecv) / float64(PacketsSent)

			NtResult := sharedStruct.NtResult{
				Seq:       i,
				Type:      Type,
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

			NtResultChain <- NtResult
			time.Sleep(1 * time.Second)
		}

	}
}
