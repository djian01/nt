package ntPinger

// func: icmpProbingRun
func icmpProbingRun(p *Pinger, errChan chan<- error) {
	payLoad := GeneratePayloadData(p.InputVars.PayLoadSize)
	runProbes(p, errChan, func(seq int) (Packet, error) {
		pkt, err := IcmpProbing(seq, p.DestAddr, p.InputVars.DestHost, p.InputVars.PayLoadSize, p.InputVars.Icmp_DF, p.InputVars.Timeout, payLoad)
		return &pkt, err
	})
}
