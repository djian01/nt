package ntPinger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

// ICMP represents an ICMP message
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

// func: tcpProbing
func IcmpProbing(Seq int, destAddr string, desetHost string, nbytes int, df bool, timeout int) (PacketICMP, error) {

	// Initial PacketICMP
	pkt := PacketICMP{
		Type:           "icmp",
		Status:         false,
		Seq:            Seq,
		DestAddr:       destAddr,
		DestHost:       desetHost,
		NBytes:         nbytes,
		Icmp_dfragment: df,
	}

	// Convert the string to *net.IPAddr
	ipAddr, err := net.ResolveIPAddr("ip4", destAddr)
	if err != nil {
		return pkt, fmt.Errorf("failed to resolve IP address: %v", err)
	}

	// Create a raw socket with IPPROTO_ICMP
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		return pkt, fmt.Errorf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Set the "Don't Fragment" (DF) flag
	if df {
		rawConn, err := conn.SyscallConn()
		if err != nil {
			return pkt, fmt.Errorf("failed to get raw connection: %v", err)
		}
		rawConn.Control(func(fd uintptr) {
			syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MTU_DISCOVER, syscall.IP_PMTUDISC_DO)
		})
	}

	// Prepare ICMP message

	////  Add a custom payload
	payload := make([]byte, nbytes)
	for i := 0; i < nbytes; i++ {
		payload[i] = byte(i) // Example payload data
	}

	//// ICMP Body
	icmpBody := icmpBody{
		ID:   os.Getegid() & 0xffff,
		Seq:  Seq,
		Data: payload,
	}

	//// ICMP Message
	icmpMsg := icmpMessage{
		Type:     8, // Echo request
		Code:     0,
		Checksum: 0, // Will calculate later
		Body:     &icmpBody,
	}

	//// ICMP Message Marshal
	binIcmpMsg, err := icmpMsg.Marshal()
	if err != nil {
		return pkt, err
	}

	// TIMEOUT - Set timeout for the connection
	conn.SetDeadline(time.Now().Add(time.Duration(timeout)))

	// SEND - the ICMP packet
	startTime := time.Now()
	pkt.SendTime = startTime

	fmt.Println(binIcmpMsg)

	/// buffer.Bytes() -> conn.Write
	if _, err := conn.Write(binIcmpMsg); err != nil {
		return pkt, fmt.Errorf("failed to send packet: %v", err)
	}

	// RECEIVE -  the ICMP response
	reply := make([]byte, 1024)

	//// conn.Read -> reply
	_, err = conn.Read(reply)
	if err != nil {
		return pkt, fmt.Errorf("failed to receive reply: %v", err)
	}

	// RTT - Calculate RTT
	rtt := time.Since(startTime)

	// Extract the ICMP header from the response
	receivedICMP := ICMP{}
	buffer := *bytes.NewBuffer(reply[20:28])
	binary.Read(&buffer, binary.BigEndian, &receivedICMP)

	// check if the receive arrives
	if receivedICMP.SequenceNum == uint16(Seq) {
		pkt.Status = true
		pkt.RTT = rtt
	}

	// return
	return pkt, nil

}

// ============================

const (
	icmpv4EchoRequest = 8
	icmpv4EchoReply   = 0
	icmpv6EchoRequest = 128
	icmpv6EchoReply   = 129
)

type icmpMessageBody interface {
	Len() int
	Marshal() ([]byte, error)
}

// Struct - icmpEchoBody: Represenets an ICMP echo request or reply message body.
type icmpBody struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

func (b *icmpBody) Len() int {
	if b == nil {
		return 0
	}
	return 4 + len(b.Data) // the total Body len is len(b.Data) + 2 bytes (ID) + 2 bytes (Seq)
}

// ICMP Body Marshal
func (b *icmpBody) Marshal() ([]byte, error) {

	// len is len(b.Data) + 2 bytes (ID) + 2 bytes (Seq)
	bin := make([]byte, 4+len(b.Data))
	// high byte -> ID, low byte -> bin[0]
	bin[0], bin[1] = byte(b.ID>>8), byte(b.ID&0xff)
	// high byte -> Seq, low byte -> bin[0]
	bin[2], bin[3] = byte(b.Seq>>8), byte(b.Seq&0xff)
	// copy the payload start from 5th byte
	copy(bin[4:], b.Data)

	return bin, nil
}

// Struct - icmpMessage: Represenets the ICMP packet
type icmpMessage struct {
	Type     int             // type
	Code     int             // code
	Checksum int             // checksum
	Body     icmpMessageBody // body
}

// ICMP Message Marshal
func (icmpMsg *icmpMessage) Marshal() ([]byte, error) {

	// 4 x bytes of ICMP Header, 1st Byte: Type, 2nd Byte: Code, 3rd & 4th Bytes: Checksum
	binIcmpMsg := []byte{byte(icmpMsg.Type), byte(icmpMsg.Code), 0, 0}

	// if icmpMsg Body is not nil and Len is not 0, append ICMP Header & Body
	if icmpMsg.Body != nil && icmpMsg.Body.Len() != 0 {
		binIcmpBody, err := icmpMsg.Body.Marshal()
		if err != nil {
			return nil, err
		}
		binIcmpMsg = append(binIcmpMsg, binIcmpBody...)
	}

	// if the type is icmpv6EchoRequest or icmpv6EchoReply, return binIcmpMsg
	// IPv6 ICMP checksum is handled differently
	switch icmpMsg.Type {
	case icmpv6EchoRequest, icmpv6EchoReply:
		return binIcmpMsg, nil
	}

	// if the type is icmpv4EchoRequest or icmpv4EchoReply
	csumIcmpMsg := Checksum(binIcmpMsg) // checksum coverage (all bytes except for the last byte)

	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	binIcmpMsg[2] ^= byte(^csumIcmpMsg & 0xff)
	binIcmpMsg[3] ^= byte(^csumIcmpMsg >> 8)
	return binIcmpMsg, nil
}

// Checksum calculates the checksum for an ICMP message
func Checksum(data []byte) uint32 {

	csumCoverage := len(data) - 1 // checksum coverage

	csum := uint32(0)

	for i := 0; i < csumCoverage; i += 2 {
		csum += uint32(data[i+1])<<8 | uint32(data[i])
	}
	if csumCoverage&1 == 0 {
		csum += uint32(data[csumCoverage])
	}
	csum = csum>>16 + csum&0xffff
	csum = csum + csum>>16

	return csum
}
