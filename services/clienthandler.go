package services

import (
	"context"
	"errors"
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/sessions"
	"reflect"
)

var (
	ClientHandlerNotExistErr = errors.New("client handler not exist")
	RemoteServerNotExistErr  = errors.New("remote server not exist")
)

type clientHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var (
	_clientHandlers                  = make(map[string]*clientHandler)
	_remoteClientHandlers2ServerType = make(map[string]string)
	_sessionType                     = reflect.TypeOf((*sessions.Session)(nil))
)

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
	_clientHandlers[route] = handler
}

func HandleClientMsg(session *sessions.Session, message *codecs.Message) ([]byte, error) {
	handler, ok := _clientHandlers[message.Route]
	if !ok {
		return HandleRemoteClientMsg(session, message)
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

func HandleRemoteClientMsg(session *sessions.Session, message *codecs.Message) ([]byte, error) {
	if serverID, ok := session.Route2ServerId[message.Route]; ok {
		if conn, ok := rpc.ClientByServerID(serverID); ok {
			rpcClient := protos.NewRPCClient(conn)
			respond, err := rpcClient.Call(context.Background(), &protos.CallRequest{
				Route:    message.Route,
				Payload:  message.Payload,
				ServerID: serverID,
			})
			if err != nil {
				return nil, err
			}
			return respond.Data, nil
		} else {
			delete(session.Route2ServerId, message.Route)
		}
	}

	serverType, ok := _remoteClientHandlers2ServerType[message.Route]
	if !ok {
		return nil, ClientHandlerNotExistErr
	}

	server, err := env.SD.GetRandomServer(serverType)
	if err != nil {
		return nil, err
	}

	conn, ok := rpc.ClientByServerID(server.ID)
	if !ok {
		return nil, RemoteServerNotExistErr
	}

	session.Route2ServerId[message.Route] = server.ID
	rpcClient := protos.NewRPCClient(conn)
	respond, err := rpcClient.Call(context.Background(), &protos.CallRequest{
		Route:    message.Route,
		Payload:  message.Payload,
		ServerID: server.ID,
	})
	if err != nil {
		return nil, err
	}

	return respond.Data, nil
}
