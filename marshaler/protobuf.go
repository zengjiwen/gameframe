package marshaler

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

var (
	MsgTypeErr = errors.New("msg type err")
)

type protobuf struct{}

func newProtobuf() Marshaler {
	return protobuf{}
}

func (p protobuf) Marshal(msg interface{}) ([]byte, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return nil, MsgTypeErr
	}

	return proto.Marshal(protoMsg)
}

func (p protobuf) Unmarshal(payload []byte, msg interface{}) error {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return MsgTypeErr
	}

	return proto.Unmarshal(payload, protoMsg)
}
