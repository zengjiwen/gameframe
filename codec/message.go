package codec

import "errors"

var MessageLenErr = errors.New("message len err")

type Message struct {
	Route   string
	Payload []byte
}

func NewMessage(route string, payload []byte) *Message {
	return &Message{
		Route:   route,
		Payload: payload,
	}
}
