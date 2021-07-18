package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/sessions"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type serverHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var ServerHandlers = make(map[string]*serverHandler)

// todo listener
var _remoteServerHandlers2ServerType = make(map[string]string)
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
	ServerHandlers[route] = handler
}

func HandleServerMsg(session *sessions.Session, message *codec.Message) ([]byte, error) {
	handler, ok := ServerHandlers[message.Route]
	if !ok {
		return nil, errors.New("server handler not exist!")
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

func RPC(route string, arg proto.Message, ret proto.Message) error {
	serverType, ok := _remoteServerHandlers2ServerType[route]
	if !ok {
		return fmt.Errorf("server handler not found! route:%s", route)
	}

	serverInfo, ok := env.SD.GetRandomServer(serverType)
	if !ok {
		return fmt.Errorf("get random server fail! server type:%s", serverType)
	}

	rpcClient, ok := rpc.Clients.ClientByServerID(serverInfo.ID)
	if !ok {
		return fmt.Errorf("rpc client not found! server id:%s", serverInfo.ID)
	}

	argBytes, err := proto.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := rpcClient.Call(context.Background(), &protos.CallRequest{
		Route:    route,
		Payload:  argBytes,
		ServerID: env.ServerID,
	})
	if err != nil {
		return err
	}

	return proto.Unmarshal(resp.Data, ret)
}
