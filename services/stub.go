package services

import (
	"context"
	"errors"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/services/proxy"
	"github.com/zengjiwen/gameframe/sessions"
)

var (
	HandlerNotExistErr   = errors.New("handler not exist!")
	ClientRpcNotExistErr = errors.New("client rpc not exist!")
	SessionNotFoundErr   = errors.New("session not found!")
)

type stub struct{}

func NewStub() *stub {
	return &stub{}
}

func (s *stub) Call(_ context.Context, request *protos.CallRequest) (*protos.CallRespond, error) {
	rpcClient, ok := rpc.ClientByServerID(request.ServerID)
	if !ok {
		return nil, ClientRpcNotExistErr
	}

	backendProxy := proxy.NewBackend(rpcClient)
	session := sessions.New(backendProxy)
	message := codec.NewMessage(request.Route, request.Payload)
	if _, ok := _clientHandlers[request.Route]; ok {
		retData, err := HandleClientMsg(session, message)
		return &protos.CallRespond{Data: retData}, err
	} else if _, ok := _serverHandlers[request.Route]; ok {
		retData, err := HandleServerMsg(session, message)
		return &protos.CallRespond{Data: retData}, err
	} else {
		return nil, HandlerNotExistErr
	}
}

func (s *stub) Send(_ context.Context, request *protos.SendRequest) (*protos.SendRespond, error) {
	session := sessions.SessionByID(request.SessionID)
	if session == nil {
		return nil, SessionNotFoundErr
	}

	return &protos.SendRespond{}, session.Send(request.Route, request.Payload)
}
