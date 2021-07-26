package marshaler

var _marshaler = newProtobuf()

type Marshaler interface {
	Marshal(arg interface{}) ([]byte, error)
	Unmarshal(payload []byte, arg interface{}) error
}

func Get() Marshaler {
	return _marshaler
}

func Set(marshaler Marshaler) {
	_marshaler = marshaler
}
