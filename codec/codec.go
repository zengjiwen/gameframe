package codec

var _codec = newPlain()

type Codec interface {
	Encode(m *Message) ([]byte, error)
	Decode(data []byte) (*Message, error)
}

func Get() Codec {
	return _codec
}

func Set(codec Codec) {
	_codec = codec
}
