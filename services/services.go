package services

import (
	"context"
	"errors"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/rpc/protos"
	"github.com/zengjiwen/gameframe/services/proxy"
	"github.com/zengjiwen/gameframe/sessions"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Call(_ context.Context, request *protos.CallRequest) (*protos.CallRespond, error) {
	conn, ok := rpc.GetConn(request.ServerID)
	if !ok {
		return nil, errors.New("client rpc not exist!")
	}

	backendProxy := proxy.NewBackend(request.ServerID, conn)
	session := sessions.New(backendProxy)
	message := codec.NewMessage(request.Route, request.Payload)
	if _, ok := ClientHandlers[request.Route]; ok {
		retData, err := HandleClientMsg(session, message)
		return &protos.CallRespond{Data: retData}, err
	} else if _, ok := ServerHandlers[request.Route]; ok {
		retData, err := HandleServerMsg(session, message)
		return &protos.CallRespond{Data: retData}, err
	} else {
		return nil, errors.New("handler not exist!")
	}
}

func (s *service) Send(_ context.Context, request *protos.SendRequest) (*protos.SendRespond, error) {
	session := sessions.SessionByID(request.SessionID)
	if session == nil {
		return nil, errors.New("session not found!")
	}

	return &protos.SendRespond{}, session.Send(request.Route, request.Payload)
}

func WatchServer() {
	env.SD.AddServerWatcher(_remoteServerHandlers)
	env.SD.AddServerWatcher(_remoteClientHandlers)
}
