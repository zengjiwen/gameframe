package gameframe

import (
	"fmt"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"github.com/zengjiwen/gameframe/services"
	"github.com/zengjiwen/gameframe/services/proxy"
	"github.com/zengjiwen/gameframe/sessions"
	"github.com/zengjiwen/gamenet"
	"github.com/zengjiwen/gamenet/server"
	"os"
	"os/signal"
	"strings"
)

var _gameFrame *gameFrame

type gameFrame struct {
}

func Run(serverType, serviceAddr string, applies ...func(opts *options)) {
	for _, apply := range applies {
		apply(&_opts)
	}

	env.ServerType = serverType
	env.ServiceAddr = serviceAddr
	if _opts.codec != nil {
		env.Codec = _opts.codec
	}
	if _opts.marshaler != nil {
		env.Marshaler = _opts.marshaler
	}
	if _opts.sd != nil {
		env.SD = _opts.sd
	} else if _opts.sdAddr != "" {
		env.SD = servicediscovery.NewEtcd(_opts.sdAddr)
	} else {
		panic("please specify sd or sd addr!")
	}

	var tcpServer gamenet.Server
	if _opts.clientAddr != "" {
		cb := frontendEventCallback{}
		if strings.ToLower(_opts.concurrentMode) == "csp" {
			eventChan := make(chan func())
			tcpServer = server.NewServer("tcp", _opts.clientAddr, cb, server.WithEventChan(eventChan))
			go func() {
				for event := range eventChan {
					event()
				}
			}()
		} else if strings.ToLower(_opts.concurrentMode) == "actor" {
			tcpServer = server.NewServer("tcp", _opts.clientAddr, cb)
		} else {
			panic(fmt.Sprintf("incorrect concurrent mode: %s", _opts.concurrentMode))
		}

		go tcpServer.ListenAndServe()
	}
	if err := rpc.StartServer(services.NewStub()); err != nil {
		panic(err)
	}
	env.SD.AddServerListener(services.RemoteServerHandlers)
	env.SD.AddServerListener(services.RemoteClientHandlers)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case <-c:
		close(env.DieChan)
	case <-env.DieChan:
	}

	if tcpServer != nil {
		tcpServer.Shutdown()
	}
	rpc.StopServer()
}

type frontendEventCallback struct{}

func (frontendEventCallback) OnNewConn(conn gamenet.Conn) {
	frontendProxy := proxy.NewFrontend(conn)
	session := sessions.New(frontendProxy)
	conn.SetUserData(session)
}

func (frontendEventCallback) OnConnClosed(conn gamenet.Conn) {
	session, ok := conn.UserData().(*sessions.Session)
	if !ok {
		return
	}

	session.OnClosed()
}

func (frontendEventCallback) OnRecvData(conn gamenet.Conn, data []byte) {
	session, ok := conn.UserData().(*sessions.Session)
	if !ok {
		return
	}

	m, err := env.Codec.Decode(data)
	if err != nil {
		return
	}

	ret, err := services.HandleClientMsg(session, m)
	if err != nil {
		return
	}

	conn.Send(ret)
}
