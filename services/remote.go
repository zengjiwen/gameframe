package services

import (
	"reflect"
)

type ServerHandler struct {
	funv reflect.Value
	argt reflect.Type
}

var _serverHandlers = make(map[string]*ServerHandler)
var _remoteServerHandlers2ServerType = make(map[string]int)

func RegisterServerHandler(route string, sh interface{}) {
	ht := reflect.TypeOf(sh)
	if ht.Kind() != reflect.Func {
		return
	}
	if ht.NumIn() != 2 || ht.In(0) != _sessionType || ht.In(1).Kind() != reflect.Ptr {
		return
	}
	if ht.NumOut() != 1 || ht.Out(0).Kind() != reflect.Ptr {
		return
	}

	handler := &ServerHandler{
		funv: reflect.ValueOf(sh),
		argt: ht.In(0),
	}
	_serverHandlers[route] = handler
}
