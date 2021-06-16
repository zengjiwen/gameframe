package marshalers

type Marshaler interface {
	Marshal(arg interface{}) ([]byte, error)
	Unmarshal(payload []byte, arg interface{}) error
}
