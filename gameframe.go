package gameframe

import (
	"errors"
	"fmt"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/marshaler"
	"github.com/zengjiwen/gamenet"
	"github.com/zengjiwen/gamenet/server"
	"os"
	"os/signal"
	"reflect"
)

type options struct {
}

type App struct {
	codec.Codec
	marshaler.Marshaler
	opts options
}

var _app = &App{}

func Start(applies ...func(opts *options)) {
	for _, apply := range applies {
		apply(&_app.opts)
	}

	eventChan := make(chan func())
	tcpServer := server.NewServer("tcp", "127.0.0.1:0", echoHandler{},
		server.WithEventChan(eventChan))
	go tcpServer.ListenAndServe()

	go func() {
		for event := range eventChan {
			event()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	tcpServer.Shutdown()
}

type echoHandler struct {
}

func (eh echoHandler) OnNewConn(c gamenet.Conn) {
	fmt.Println("OnNewConn")
}

func (eh echoHandler) OnConnClosed(c gamenet.Conn) {
	fmt.Println("OnConnClosed")
}

func (eh echoHandler) OnRecvData(c gamenet.Conn, data []byte) {
	fmt.Printf("OnRecvData: %s", data)
	m, err := _app.Codec.Decode(data)
	if err != nil {
		return
	}

	ret, err := ProcessHandlerMsg(m)
	if err != nil {
		return
	}

	c.Send(ret)
}

type Handler struct {
	funv reflect.Value
	argt reflect.Type
}

var _handlers = make(map[string]*Handler)

func RegisterHandler(route string, h interface{}) {
	ht := reflect.TypeOf(h)
	if ht.Kind() != reflect.Func {
		return
	}
	if ht.NumIn() != 1 || ht.In(0).Kind() != reflect.Ptr {
		return
	}
	if ht.NumOut() != 1 || ht.Out(0).Kind() != reflect.Ptr {
		return
	}

	handler := &Handler{
		funv: reflect.ValueOf(h),
		argt: ht.In(0),
	}
	_handlers[route] = handler
}

type Remote struct {
	fun  reflect.Value
	argt reflect.Type
}

var _remotes = make(map[string]*Remote)

func RegisterRemote(route string, r interface{}) {
	ht := reflect.TypeOf(r)
	if ht.Kind() != reflect.Func {
		panic(errors.New("error: remote isn't func!"))
	}

	remote := &Remote{}
	_remotes[route] = remote
}

var _route2ServerType = make(map[string]int)

func ProcessHandlerMsg(m *codec.Message) ([]byte, error) {
	handler, ok := _handlers[m.Route]
	if !ok {
		serverType, ok := _route2ServerType[m.Route]
		if !ok {
			return nil, errors.New("handler isn't exist")
		}
		// todo
	}

	argi := reflect.New(handler.argt.Elem()).Interface()
	if err := _app.Marshaler.Unmarshal(m.Payload, argi); err != nil {
		return nil, err
	}

	rets := handler.funv.Call([]reflect.Value{reflect.ValueOf(argi)})
	if len(rets) != 1 {
		return nil, errors.New("rets len isn't 1")
	}

	retData, err := _app.Marshaler.Marshal(rets[0])
	if err != nil {
		return nil, err
	}

	return retData, nil
}
