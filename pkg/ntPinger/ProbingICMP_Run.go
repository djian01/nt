package ntPinger

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// func: icmpProbingRun
func icmpProbingRun(p *Pinger, errChan chan<- error) {

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Sequence
	Seq := 0

	// forLoopEnds Flag
	forLoopEnds := false

	// generate payload
	payLoad := GeneratePayloadData(p.InputVars.PayLoadSize)

	// count
	if p.InputVars.Count == 0 {

		for {
			// Loop End Signal
			if forLoopEnds {
				break
			}

			// Pinger end Singal
			if p.PingerEnd {
				interruptChan <- os.Interrupt //send interrupt to interruptChan
			}

			pkt, err := IcmpProbing(Seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.PayLoadSize, p.InputVars.Icmp_DF, p.InputVars.Timeout, payLoad)
			if err != nil {
				errChan <- err
			}

			p.UpdateStatistics(&pkt)
			pkt.UpdateStatistics(p.Stat)
			p.ProbeChan <- &pkt
			Seq++

			// sleep for interval
			select {
			case <-time.After(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT)):
				// wait for SleepTime
			case <-interruptChan: // case interruptChan, close the channel & break the loop
				forLoopEnds = true
				close(p.ProbeChan)

			}
		}

	} else {
		for i := 0; i < p.InputVars.Count; i++ {
			if forLoopEnds {
				break
			}

			pkt, err := IcmpProbing(Seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.PayLoadSize, p.InputVars.Icmp_DF, p.InputVars.Timeout, payLoad)
			if err != nil {
				errChan <- err
			}

			p.UpdateStatistics(&pkt)
			pkt.UpdateStatistics(p.Stat)
			p.ProbeChan <- &pkt
			Seq++

			// check the last loop of the probing, close probeChan
			if i == (p.InputVars.Count - 1) {
				close(p.ProbeChan)
			}

			// sleep for interval
			select {
			case <-time.After(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT)):
				// wait for SleepTime
			case <-interruptChan: // case interruptChan, close the channel & break the loop
				forLoopEnds = true
				close(p.ProbeChan)
			}
		}
	}
}
