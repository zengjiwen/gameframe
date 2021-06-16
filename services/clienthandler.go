package services

import (
	"errors"
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/sd"
	"github.com/zengjiwen/gameframe/services/peers"
	"github.com/zengjiwen/gameframe/sessions"
	"reflect"
)

var (
	HandlerNotExistErr      = errors.New("handler not exist")
	RemoteServerNotExistErr = errors.New("remote server not exist")
)

type ClientHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var (
	_clientHandlers                  = make(map[string]*ClientHandler)
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

	handler := &ClientHandler{
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
	if serverId, ok := session.Route2ServerId[message.Route]; ok {
		if r, ok := peers.Remotes[serverId]; ok {
			respond, err := r.Call(message)
			if err != nil {
				return nil, err
			}
			return respond.Payload, nil
		} else {
			delete(session.Route2ServerId, message.Route)
		}
	}

	serverType, ok := _remoteClientHandlers2ServerType[message.Route]
	if !ok {
		return nil, HandlerNotExistErr
	}

	servers, err := sd.SD.GetServersByType(serverType)
	if err != nil {
		return nil, err
	}

	server := sd.GetMinLoadServer(servers)
	if server == nil {
		return nil, err
	}

	r, ok := peers.Remotes[server.ID]
	if !ok {
		return nil, RemoteServerNotExistErr
	}

	session.Route2ServerId[message.Route] = server.ID
	respond, err := r.Call(message)
	if err != nil {
		return nil, err
	}

	return respond.Payload, nil
}
