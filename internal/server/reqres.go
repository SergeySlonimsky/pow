package server

import (
	"github.com/SergeySlonimsky/pow/pkg/protocol"
)

type Response interface {
	ToString() string
}

type Request interface {
	GetMessage() protocol.Message
	GetAddr() string
}

type messageRequest struct {
	message protocol.Message
	addr    string
}

func (r messageRequest) GetMessage() protocol.Message {
	return r.message
}

func (r messageRequest) GetAddr() string {
	return r.addr
}
