package codec

// todo pool
type plain struct{}

func newPlain() Codec {
	return plain{}
}

func (p plain) Encode(m *Message) ([]byte, error) {
	data := make([]byte, 0, 1+len(m.Route)+len(m.Payload))
	data = append(data, byte(len(m.Route)))
	data = append(data, []byte(m.Route)...)
	data = append(data, m.Payload...)
	return data, nil
}

func (p plain) Decode(data []byte) (*Message, error) {
	if len(data) == 0 {
		return nil, MessageLenErr
	}

	offset := 0
	routeLen := int(data[offset])
	offset++
	if offset+routeLen > len(data) {
		return nil, MessageLenErr
	}

	route := data[offset : offset+routeLen]
	return NewMessage(string(route), data[offset+routeLen:]), nil
}
