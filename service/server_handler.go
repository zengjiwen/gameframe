package service

import (
	"errors"
	"fmt"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/marshaler"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"github.com/zengjiwen/gameframe/sessions"
	"google.golang.org/protobuf/proto"
	"reflect"
)

var ServerHandlers = make(map[string]*serverHandler)

type serverHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var _remoteServerHandlers = make(remoteServerHandlers)

type remoteServerHandlers map[string]string

func (s remoteServerHandlers) OnAddServer(serverInfo *servicediscovery.ServerInfo) {
	for _, handler := range serverInfo.ServerHandlers {
		s[handler] = serverInfo.Type
	}
}

func (s remoteServerHandlers) OnRemoveServer(serverInfo *servicediscovery.ServerInfo) {
	for handler, serverType := range s {
		if serverType == serverInfo.Type {
			delete(s, handler)
		}
	}
}

var _protoMsgType = reflect.TypeOf(proto.Message(nil)).Elem()

func RegisterServerHandler(route string, sh interface{}) {
	ht := reflect.TypeOf(sh)
	if ht.Kind() != reflect.Func {
		return
	}
	if ht.NumIn() != 1 || ht.In(0) != _protoMsgType {
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
	if err := marshaler.Get().Unmarshal(message.Payload, argi); err != nil {
		return nil, err
	}

	rets := handler.funv.Call([]reflect.Value{reflect.ValueOf(argi)})
	if len(rets) != 1 {
		return nil, errors.New("rets len isn't 1")
	}

	payLoad, err := marshaler.Get().Marshal(rets[0])
	if err != nil {
		return nil, err
	}

	message.Payload = payLoad
	retData, err := codec.Get().Encode(message)
	if err != nil {
		return nil, err
	}

	return retData, nil
}

func RPC(route string, arg proto.Message, ret proto.Message) error {
	serverType, ok := _remoteServerHandlers[route]
	if !ok {
		return fmt.Errorf("server handler not found! route:%s", route)
	}

	serverInfo, ok := servicediscovery.Get().GetRandomServer(serverType)
	if !ok {
		return fmt.Errorf("get random server fail! server type:%s", serverType)
	}

	rpcConn, ok := rpc.GetConn(serverInfo.ID)
	if !ok {
		return fmt.Errorf("rpc client not found! server id:%s", serverInfo.ID)
	}
	rpcClient := protos.NewRPCClient(rpcConn)

	argBytes, err := proto.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := rpc.TryBestCall(serverInfo.ID, rpcClient, &protos.CallRequest{
		Route:    route,
		Payload:  argBytes,
		ServerID: servicediscovery.GetServerInfo().ID,
	})
	if err != nil {
		return err
	}

	return proto.Unmarshal(resp.Data, ret)
}
