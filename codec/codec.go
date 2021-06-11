package codec

type Message struct {
	Route   string
	Payload []byte
}

type Codec interface {
	Encode(m *Message) ([]byte, error)
	Decode(data []byte) (*Message, error)
}
