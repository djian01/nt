package ntTEST

import (
	"math"
	"math/rand"
	"nt/pkg/ntPinger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ResultGenerate(count int, Type string, probeChan *chan ntPinger.Packet) {

	// statistics
	var PacketsSent = 0
	var PacketsRecv = 0
	var PacketLoss float64
	var MinRtt time.Duration
	var AvgRtt time.Duration
	var MaxRtt time.Duration
	var status bool
	var AdditionalInfo string

	// initial RTT
	min := 5
	max := 1700
	var ranRTT time.Duration

	// random Source
	source := rand.NewSource(time.Now().UnixNano())

	// random Error Seed
	errorMax := 12
	errorMin := 2
	errorSeed := rand.New(source).Intn(errorMax-errorMin+1) + errorMin

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// endless loop
	if count == 0 {

		// forLoopFlag
		forLoopFlag := true

		// create the input NtResult
		for {

			// initial ranRTT
			ranRTT = time.Duration(0)

			// initial AdditionalInfo
			AdditionalInfo = ""

			// check forLoopFlag
			if !forLoopFlag {
				break
			}

			// PacketSent ++
			PacketsSent++

			// check Interrpution
			select {
			case <-interruptChan:
				// if doneChan <- true if interrupted
				forLoopFlag = false
			default:
				// pass
			}

			// error check
			if PacketsSent%errorSeed != 0 {
				PacketsRecv++
				status = true
			} else {
				status = false
				AdditionalInfo = "Timeout"
			}

			// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
			if PacketsSent == 1 {
				ranRTT = time.Duration(biasedRandom(min, max, rand.New(source))) * time.Millisecond
				MinRtt = ranRTT
				MaxRtt = ranRTT
				AvgRtt = ranRTT
				// after Initialization
			} else if status {
				// generate RTT
				ranRTT = time.Duration(biasedRandom(min, max, rand.New(source))) * time.Millisecond

				// RTT Statistics
				if ranRTT > MaxRtt {
					MaxRtt = ranRTT
				} else if ranRTT < MinRtt {
					MinRtt = ranRTT
				}
				AvgRtt = (AvgRtt*time.Duration(PacketsRecv-1) + ranRTT) / time.Duration(PacketsRecv)

			}

			// generate statistic
			PacketLoss = (1 - float64(PacketsRecv)/float64(PacketsSent))

			// Check Latency
			if ntPinger.CheckLatency(AvgRtt, ranRTT) {
				AdditionalInfo = "High_Latency"
			}

			switch Type {
			case "icmp":
				probeResult := ntPinger.PacketICMP{
					Seq:         PacketsSent - 1,
					Type:        Type,
					DestHost:    "google.com",
					DestAddr:    "1.2.3.4",
					PayLoadSize: 56,
					Status:      status,
					RTT:         ranRTT,
					SendTime:    time.Now(),

					PacketsSent:    PacketsSent,
					PacketsRecv:    PacketsRecv,
					PacketLoss:     PacketLoss,
					MinRtt:         MinRtt,
					MaxRtt:         MaxRtt,
					AvgRtt:         AvgRtt,
					AdditionalInfo: AdditionalInfo,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "tcp":
				destPort := 443

				probeResult := ntPinger.PacketTCP{
					Seq:         PacketsSent - 1,
					Type:        Type,
					DestHost:    "google.com",
					DestAddr:    "1.2.3.4",
					PayLoadSize: 56,
					Status:      status,
					RTT:         ranRTT,
					SendTime:    time.Now(),
					DestPort:    destPort,

					PacketsSent:    PacketsSent,
					PacketsRecv:    PacketsRecv,
					PacketLoss:     PacketLoss,
					MinRtt:         MinRtt,
					MaxRtt:         MaxRtt,
					AvgRtt:         AvgRtt,
					AdditionalInfo: AdditionalInfo,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "http":
				destPort := 443

				probeResult := ntPinger.PacketHTTP{
					Seq:         PacketsSent - 1,
					Type:        Type,
					DestHost:    "google.com",
					Status:      status,
					RTT:         ranRTT,
					SendTime:    time.Now(),
					DestPort:    destPort,
					Http_path: "",
					Http_scheme: "https",

					PacketsSent:    PacketsSent,
					PacketsRecv:    PacketsRecv,
					PacketLoss:     PacketLoss,
					MinRtt:         MinRtt,
					MaxRtt:         MaxRtt,
					AvgRtt:         AvgRtt,
					AdditionalInfo: AdditionalInfo,
				}
				if status {
					probeResult.Http_response_code = 200
					probeResult.Http_response = "OK"
				}else {
					probeResult.Http_response_code = 0
				}				

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)
			case "dns":

			}
		}

		// code finish
		close(*probeChan)

	} else {

		// forLoopFlag
		forLoopFlag := true

		// create the input NtResult
		for i := 0; i < count; i++ {

			// initial ranRTT
			ranRTT = time.Duration(0)

			// initial AdditionalInfo
			AdditionalInfo = ""

			// check forLoopFlag
			if !forLoopFlag {
				break
			}

			// PacketSent ++
			PacketsSent++

			// check Interrpution
			select {
			case <-interruptChan:
				// if doneChan <- true if interrupted
				forLoopFlag = false

			default:
				// pass
			}

			// error check
			if PacketsSent%errorSeed != 0 {
				PacketsRecv++
				status = true
			} else {
				status = false
				AdditionalInfo = "Timeout"
			}

			// setup the initial MinRtt, MaxRtt & AvgRtt. Based on the 1st "PING_OK" packet, set the MinRtt, MaxRtt & AvgRtt
			if PacketsSent == 1 {
				ranRTT = time.Duration(biasedRandom(min, max, rand.New(source))) * time.Millisecond
				MinRtt = ranRTT
				MaxRtt = ranRTT
				AvgRtt = ranRTT

				// if packet status true
			} else if status {
				// generate RTT
				ranRTT = time.Duration(biasedRandom(min, max, rand.New(source))) * time.Millisecond

				if ranRTT > MaxRtt {
					MaxRtt = ranRTT
				} else if ranRTT < MinRtt {
					MinRtt = ranRTT
				}
				AvgRtt = (AvgRtt*time.Duration(PacketsRecv-1) + ranRTT) / time.Duration(PacketsRecv)
			}

			// generate statistic
			PacketLoss = (1 - float64(PacketsRecv)/float64(PacketsSent))

			// Check Latency
			if ntPinger.CheckLatency(AvgRtt, ranRTT) {
				AdditionalInfo = "High_Latency"
			}

			switch Type {
			case "icmp":
				probeResult := ntPinger.PacketICMP{
					Seq:         PacketsSent - 1,
					Type:        Type,
					DestHost:    "google.com",
					DestAddr:    "1.2.3.4",
					PayLoadSize: 56,
					Status:      status,
					RTT:         ranRTT,
					SendTime:    time.Now(),

					PacketsSent:    PacketsSent,
					PacketsRecv:    PacketsRecv,
					PacketLoss:     PacketLoss,
					MinRtt:         MinRtt,
					MaxRtt:         MaxRtt,
					AvgRtt:         AvgRtt,
					AdditionalInfo: AdditionalInfo,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "tcp":
				destPort := 443

				probeResult := ntPinger.PacketTCP{
					Seq:         PacketsSent - 1,
					Type:        Type,
					DestHost:    "google.com",
					DestAddr:    "1.2.3.4",
					PayLoadSize: 56,
					Status:      status,
					RTT:         ranRTT,
					SendTime:    time.Now(),
					DestPort:    destPort,

					PacketsSent:    PacketsSent,
					PacketsRecv:    PacketsRecv,
					PacketLoss:     PacketLoss,
					MinRtt:         MinRtt,
					MaxRtt:         MaxRtt,
					AvgRtt:         AvgRtt,
					AdditionalInfo: AdditionalInfo,
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)

			case "http":
				destPort := 443

				probeResult := ntPinger.PacketHTTP{
					Seq:         PacketsSent - 1,
					Type:        Type,
					DestHost:    "google.com",
					Status:      status,
					RTT:         ranRTT,
					SendTime:    time.Now(),
					DestPort:    destPort,
					Http_path: "c/66dc2804-7f48-8011-88d8-c6bf57428b6a/c/66dc2804-7f48-8011-88d8-c6bf57428b6a",
					//Http_path: "web/login",
					Http_scheme: "https",

					PacketsSent:    PacketsSent,
					PacketsRecv:    PacketsRecv,
					PacketLoss:     PacketLoss,
					MinRtt:         MinRtt,
					MaxRtt:         MaxRtt,
					AvgRtt:         AvgRtt,
					AdditionalInfo: AdditionalInfo,
				}

				if status {
					probeResult.Http_response_code = 200
					probeResult.Http_response = "OK"
				}else {
					probeResult.Http_response_code = 0
				}

				*probeChan <- &probeResult
				time.Sleep(1 * time.Second)
			case "dns":
			}
		}

		// code finish
		close(*probeChan)
	}
}

// Biased random number generator using custom source. Bias towards smaller numbers
func biasedRandom(min, max int, r *rand.Rand) int {
	// Generate a random float in the range [0, 1]
	randomFloat := r.Float64()
	// ^3 or apply another bias to the random value to bias towards smaller numbers
	biased := math.Pow(randomFloat, 3) // You can adjust the power for more bias
	// Scale to the desired range and convert to int
	return min + int(biased*float64(max-min))
}
