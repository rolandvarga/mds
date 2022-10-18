package internal

import (
	"encoding/binary"
	"fmt"
)

// TODO add replies from server

type MessageType int

const (
	Identity MessageType = iota
	ListClients
	Message
)

type Request struct {
	Type    MessageType
	Message string // optional
}

func NewRequest() Request {
	return Request{}
}

func BuildIdentityRequest()    {}
func BuildListClientsRequest() {}
func BuildMessageRequest()     {}

func BuildIdentityResponse(id uint8) []byte {
	// TODO cleanup
	out := make([]byte, 100)
	msg := []byte(fmt.Sprintf("%d", id))
	length := len(msg)

	fmt.Printf("msg: %v length: %d\n", msg, length)

	binary.LittleEndian.PutUint16(out, uint16(length))

	return append(out, msg...)
}
