package services

import (
	"errors"
	"fmt"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"github.com/zengjiwen/gameframe/sessions"
	"reflect"
)

var (
	ClientHandlerNotExistErr = errors.New("client handler not exist")
	RemoteServerNotExistErr  = errors.New("remote server not exist")
)

var ClientHandlers = make(map[string]*clientHandler)

type clientHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var _remoteClientHandlers = make(remoteClientHandlers)

type remoteClientHandlers map[string]string

func (r remoteClientHandlers) OnAddServer(serverInfo *servicediscovery.ServerInfo) {
	for _, handler := range serverInfo.ClientHandlers {
		r[handler] = serverInfo.Type
	}
}

func (r remoteClientHandlers) OnRemoveServer(serverInfo *servicediscovery.ServerInfo) {
	for handler, serverType := range r {
		if serverType == serverInfo.Type {
			delete(r, handler)
		}
	}
}

var _sessionType = reflect.TypeOf((*sessions.Session)(nil))

func RegisterClientHandler(route string, ch interface{}) {
	ht := reflect.TypeOf(ch)
	if ht.Kind() != reflect.Func {
		return
	}
	if ht.NumIn() != 2 || ht.In(0) != _sessionType || ht.In(1).Kind() != reflect.Ptr {
		return
	}
	if ht.NumOut() != 1 || ht.Out(0).Kind() != reflect.Ptr {
		return
	}

	handler := &clientHandler{
		funv: reflect.ValueOf(ch),
		argt: ht.In(0),
	}
	ClientHandlers[route] = handler
}

func HandleClientMsg(session *sessions.Session, message *codec.Message) ([]byte, error) {
	handler, ok := ClientHandlers[message.Route]
	if !ok {
		return CallRemoteClientHandler(session, message)
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

func CallRemoteClientHandler(session *sessions.Session, message *codec.Message) ([]byte, error) {
	if serverID, ok := session.Route2ServerId[message.Route]; ok {
		if rpcConn, ok := rpc.GetConn(serverID); ok {
			rpcClient := protos.NewRPCClient(rpcConn)
			resp, err := rpc.TryBestCall(serverID, rpcClient, &protos.CallRequest{
				Route:    message.Route,
				Payload:  message.Payload,
				ServerID: env.ServerID,
			})
			if err != nil {
				return nil, err
			}
			return resp.Data, nil
		} else {
			delete(session.Route2ServerId, message.Route)
		}
	}

	serverType, ok := _remoteClientHandlers[message.Route]
	if !ok {
		return nil, ClientHandlerNotExistErr
	}

	serverInfo, ok := env.SD.GetRandomServer(serverType)
	if !ok {
		return nil, fmt.Errorf("get random server fail! server type:%s", serverType)
	}

	rpcConn, ok := rpc.GetConn(serverInfo.ID)
	if !ok {
		return nil, RemoteServerNotExistErr
	}
	rpcClient := protos.NewRPCClient(rpcConn)

	session.Route2ServerId[message.Route] = serverInfo.ID
	resp, err := rpc.TryBestCall(serverInfo.ID, rpcClient, &protos.CallRequest{
		Route:    message.Route,
		Payload:  message.Payload,
		ServerID: env.ServerID,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
