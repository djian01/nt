package ntPinger

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// runProbes owns ProbeChan and closes it once when probing finishes.
func runProbes(p *Pinger, errChan chan<- error, probe func(int) (Packet, error)) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer close(p.ProbeChan)

	for seq := 0; p.InputVars.Count == 0 || seq < p.InputVars.Count; seq++ {
		if p.PingerEnd {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

		pkt, err := probe(seq)
		if err != nil {
			select {
			case errChan <- err:
			case <-ctx.Done():
			}
			return
		}

		p.StatsMu.RLock()
		previousAverage := p.Stat.AvgRtt
		p.StatsMu.RUnlock()
		p.UpdateStatistics(pkt)
		p.StatsMu.RLock()
		updatedStatistics := p.Stat
		p.StatsMu.RUnlock()
		pkt.UpdateStatistics(updatedStatistics)
		if CheckLatency(previousAverage, pkt.GetRtt()) {
			setPacketAdditionalInfo(pkt, "High_Latency")
		}
		select {
		case p.ProbeChan <- pkt:
		case <-ctx.Done():
			return
		}

		if p.InputVars.Count > 0 && seq+1 == p.InputVars.Count {
			return
		}

		delay := GetSleepTime(pkt.GetStatus(), p.InputVars.Interval, pkt.GetRtt())
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		}
	}
}

func setPacketAdditionalInfo(pkt Packet, info string) {
	switch packet := pkt.(type) {
	case *PacketTCP:
		packet.AdditionalInfo = info
	case *PacketICMP:
		packet.AdditionalInfo = info
	case *PacketHTTP:
		packet.AdditionalInfo = info
	case *PacketDNS:
		packet.AdditionalInfo = info
	}
}
