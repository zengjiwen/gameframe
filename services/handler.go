package services

import (
	"errors"
	"github.com/zengjiwen/gameframe/codecs"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/sessions"
	"reflect"
)

type ClientHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var (
	_clientHandlers                  = make(map[string]*ClientHandler)
	_remoteClientHandlers2ServerType = make(map[string]int)
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

func HandleClientMsg(s *sessions.Session, m *codecs.Message) ([]byte, error) {
	handler, ok := _clientHandlers[m.Route]
	if !ok {
		serverType, ok := _remoteClientHandlers2ServerType[m.Route]
		if !ok {
			return nil, errors.New("handler isn't exist")
		}
		// todo
	}

	argi := reflect.New(handler.argt.Elem()).Interface()
	if err := env.Marshaler.Unmarshal(m.Payload, argi); err != nil {
		return nil, err
	}

	rets := handler.funv.Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(argi)})
	if len(rets) != 1 {
		return nil, errors.New("rets len isn't 1")
	}

	payLoad, err := env.Marshaler.Marshal(rets[0])
	if err != nil {
		return nil, err
	}

	m.Payload = payLoad
	retData, err := env.Codec.Encode(m)
	if err != nil {
		return nil, err
	}

	return retData, nil
}
