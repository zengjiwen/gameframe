package services

import (
	"errors"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/sessions"
	"google.golang.org/protobuf/proto"
	"reflect"
)

var ServerHandlerNotExistErr = errors.New("server handler not exist!")

type serverHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var _serverHandlers = make(map[string]*serverHandler)
var _remoteServerHandlers2ServerType = make(map[string]int)
var _protoMsgType = reflect.TypeOf(proto.Message(nil)).Elem()

func RegisterServerHandler(route string, sh interface{}) {
	ht := reflect.TypeOf(sh)
	if ht.Kind() != reflect.Func {
		return
	}
	if ht.NumIn() != 2 || ht.In(0) != _sessionType || ht.In(1) != _protoMsgType {
		return
	}
	if ht.NumOut() != 1 || ht.Out(0) != _protoMsgType {
		return
	}

	handler := &serverHandler{
		funv: reflect.ValueOf(sh),
		argt: ht.In(0),
	}
	_serverHandlers[route] = handler
}

func HandleServerMsg(session *sessions.Session, message *codec.Message) ([]byte, error) {
	handler, ok := _serverHandlers[message.Route]
	if !ok {
		return nil, ServerHandlerNotExistErr
	}

	argi := reflect.New(handler.argt.Elem()).Interface()
	if err := env.Marshaler.Unmarshal(message.Payload, argi); err != nil {
		return nil, err
	}

	rets := handler.funv.Call([]reflect.Value{reflect.ValueOf(session), reflect.ValueOf(argi)})
	if len(rets) != 1 {
		return nil, errors.New("rets len isn't 1")
	}

	payLoad, err := env.Marshaler.Marshal(rets[0])
	if err != nil {
		return nil, err
	}

	message.Payload = payLoad
	retData, err := env.Codec.Encode(message)
	if err != nil {
		return nil, err
	}

	return retData, nil
}
