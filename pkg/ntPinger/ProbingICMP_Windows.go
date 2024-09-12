//go:build windows
// +build windows

package ntPinger

// func: IcmpProbing
func IcmpProbing(Seq int, destAddr string, desetHost string, PayLoadSize int, Icmp_DF bool, timeout int, payload []byte) (PacketICMP, error) {

	// Initial PacketICMP
	pkt := PacketICMP{
		Type:        "icmp",
		Status:      false,
		Seq:         Seq,
		DestAddr:    destAddr,
		DestHost:    desetHost,
		PayLoadSize: PayLoadSize,
		Icmp_DF:     Icmp_DF,
	}

	return pkt, nil
}
