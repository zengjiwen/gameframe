package marshaler

type Marshaler interface {
	Marshal(msg interface{}) ([]byte, error)
	Unmarshal(payload []byte, msg interface{}) error
}
