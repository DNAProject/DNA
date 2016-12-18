package net

import (
	"log"
	"bytes"
	"errors"
	"encoding/binary"
)

const (
	MSGCMDLEN = 12
)

// The network and module communication message buffer
type Msg struct {
	magic	 uint32
	cmd	 [MSGCMDLEN]byte 	// the message command (message type)
	length   uint32
	checksum uint32
	payloader interface{}
}

type varStr struct {
	len uint
	buf []byte
}

type verACK struct {
	// No payload
}

type version struct {
	version		uint32
	services	uint64
	timeStamp	uint32
	port		uint16
	nonce		uint32
	userAgent	varStr
	startHeight	uint32
}

type headersReq struct {
	hashStart	[32]byte
	hashEnd		[32]byte
}

// Sample function, shoule be called from ledger module
func ledgerGetHeader() [32]byte {
	var t [32]byte
	return t
}

// TODO combine all of message alloc in one function via interface
func newMsg(p interface{}) (*Msg, error) {
	msg := new(Msg)
	switch t := p.(type) {
	case version:
		log.Printf("Port is %d", t.port)
	case verACK:
	case headersReq:
		t.hashStart = ledgerGetHeader()
	default:
		return nil, errors.New("Unknown message type")
	}

	msg.payloader = p
	return msg, nil
}

func newVersionMsg() *Msg {
	var ver version
	msg := new(Msg)
	msg.payloader = &ver
	return msg
}

func newVerackMsg() *Msg {
	var verACK verACK
	msg := new(Msg)
	msg.payloader = &verACK
	return msg
}

func newHeadersReqMsg() *Msg {
	var headersReq headersReq
	msg := new(Msg)
	msg.payloader = &headersReq
	return msg
}

func (msg Msg) serialization() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, msg)
	return buf.Bytes(), err
}

func (msg *Msg) deserialization(p []byte) error {
	var buf bytes.Buffer
	_, err := buf.Read(p)
	if (err != nil) {
		return err
	}
	err = binary.Read(&buf, binary.LittleEndian, msg)
	return err
}
