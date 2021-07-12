package codec

import "errors"

var (
	MessageLenErr = errors.New("message len err")
)

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

type Codec interface {
	Encode(m *Message) ([]byte, error)
	Decode(data []byte) (*Message, error)
}
