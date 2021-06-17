package marshalers

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

var (
	MsgTypeErr = errors.New("msg type err")
)

type Protobuf struct{}

func NewProtobuf() Protobuf {
	return Protobuf{}
}

func (p Protobuf) Marshal(msg interface{}) ([]byte, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return nil, MsgTypeErr
	}

	return proto.Marshal(protoMsg)
}

func (p Protobuf) Unmarshal(payload []byte, msg interface{}) error {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return MsgTypeErr
	}

	return proto.Unmarshal(payload, protoMsg)
}
