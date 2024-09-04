package ntPingerExample

// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// taken from http://golang.org/src/pkg/net/ipraw_test.go

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	icmpv4EchoRequest = 8
	icmpv4EchoReply   = 0
	icmpv6EchoRequest = 128
	icmpv6EchoReply   = 129
)

// ************* Interface, Struct, Method ******************
// Struct - icmpMessage: Icmp Whole Message
type icmpMessage struct {
	Type     int             // type
	Code     int             // code
	Checksum int             // checksum
	Body     icmpMessageBody // body
}

// Interface - icmpMessageBody
type icmpMessageBody interface {
	Len() int
	Marshal() ([]byte, error)
	SetPayloadData(payLoadSize int)
}

// Struct - icmpEchoBody: Icmp request body. Satisfy Interface - icmpMessageBody
type icmpBody struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

// icmpEchoBody Method - Len()
func (b *icmpBody) Len() int {
	if b == nil {
		return 0
	}

	// the total Body len = 2 bytes (ID) + 2 bytes (Seq) + len(b.Data)
	return 4 + len(b.Data)
}

// icmpEchoBody Method - Marshal()
func (b *icmpBody) Marshal() ([]byte, error) {

	// len is len(b.Data) + 2 bytes (ID) + 2 bytes (Seq)
	binIcmpBody := make([]byte, 4+len(b.Data))

	// high byte -> ID, low byte -> bin[0]
	binIcmpBody[0], binIcmpBody[1] = byte(b.ID>>8), byte(b.ID&0xff)

	// high byte -> Seq, low byte -> bin[0]
	binIcmpBody[2], binIcmpBody[3] = byte(b.Seq>>8), byte(b.Seq&0xff)

	// copy the payload start from 5th byte
	copy(binIcmpBody[4:], b.Data)
	return binIcmpBody, nil
}

// icmpEchoBody Method - SetPayloadData()
func (b *icmpBody) SetPayloadData(payLoadSize int) {
	b.Data = make([]byte, payLoadSize)
	for i := 0; i < payLoadSize; i++ {
		b.Data[i] = byte(i) // Example payload data
	}
}

// icmpMessage Method - Len()
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

	// if the type is icmpv4EchoRequest or icmpv4EchoReply, move forward

	// Create Chechsum for the ICMP Message Bin
	csum := Checksum(binIcmpMsg)

	// Place checksum back in header
	binIcmpMsg[2] = byte(^csum & 0xff)
	binIcmpMsg[3] = byte(^csum >> 8)

	return binIcmpMsg, nil
}

// ************* Functions ******************

// func - parseICMPMessage, parses bin as an ICMP message binary []byte
func parseICMPMessage(bin []byte) (*icmpMessage, error) {

	msglen := len(bin)
	if msglen < 4 {
		return nil, errors.New("message too short")
	}

	icmpMsg := &icmpMessage{Type: int(bin[0]), Code: int(bin[1]), Checksum: int(bin[2])<<8 | int(bin[3])}

	if msglen > 4 {
		var err error
		switch icmpMsg.Type {
		case icmpv4EchoRequest, icmpv4EchoReply, icmpv6EchoRequest, icmpv6EchoReply:
			icmpMsg.Body, err = parseICMPEcho(bin[4:])
			if err != nil {
				return nil, err
			}
		}
	}
	return icmpMsg, nil
}

// parseICMPEcho parses b as an ICMP echo request or reply message body.
func parseICMPEcho(bin []byte) (*icmpBody, error) {

	bodylen := len(bin)

	icmpB := &icmpBody{ID: int(bin[0])<<8 | int(bin[1]), Seq: int(bin[2])<<8 | int(bin[3])}

	if bodylen > 4 {
		icmpB.Data = make([]byte, bodylen-4)
		copy(icmpB.Data, bin[4:])
	}
	return icmpB, nil
}

func ipv4Payload(b []byte) []byte {
	if len(b) < 20 {
		return b
	}
	hdrlen := int(b[0]&0x0f) << 2
	return b[hdrlen:]
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

// ******************** Pinger ***************************

func Ping(address string, timeout int) bool {
	err := Pinger(address, timeout)
	return err == nil
}

func Pinger(address string, timeout int) error {
	c, err := net.Dial("ip4:icmp", address)
	if err != nil {
		return err
	}
	c.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	defer c.Close()

	typ := icmpv4EchoRequest
	xid, xseq := os.Getpid()&0xffff, 1

	icmpB := icmpBody{
		ID:  xid,
		Seq: xseq,
	}
	icmpB.SetPayloadData(24)

	wb, err := (&icmpMessage{
		Type: typ,
		Code: 0,
		Body: &icmpB,
	}).Marshal()
	if err != nil {
		return err
	}

	// send
	fmt.Println(wb)
	if _, err = c.Write(wb); err != nil {
		return err
	}
	var m *icmpMessage
	rb := make([]byte, 20+len(wb))
	for {
		if _, err = c.Read(rb); err != nil {
			return err
		}
		rb = ipv4Payload(rb)
		if m, err = parseICMPMessage(rb); err != nil {
			return err
		}
		switch m.Type {
		case icmpv4EchoRequest, icmpv6EchoRequest:
			continue
		}
		break
	}
	fmt.Println(rb)
	fmt.Println(m)
	fmt.Println(m.Body)
	return nil
}
