package internal

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// TODO add replies from server

type MsgType int

const (
	Identity MsgType = iota
	ListClients
	Message
)

type Response struct {
	Type MsgType
	Msg  string
}

type Request struct {
	Type    MsgType
	Message string // optional
}

func NewRequest() Request {
	return Request{}
}

func BuildIdentityRequest()    {}
func BuildListClientsRequest() {}
func BuildMessageRequest()     {}

func BuildIdentityResponse(id uint8) []byte {
	msg := Response{Type: Identity, Msg: fmt.Sprint(id)}
	buff := new(bytes.Buffer)

	gobobj := gob.NewEncoder(buff)
	gobobj.Encode(msg)

	return buff.Bytes()
}
