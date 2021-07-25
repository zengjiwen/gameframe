package gameframe

import (
	"errors"
	"fmt"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"github.com/zengjiwen/gameframe/services"
	"github.com/zengjiwen/gameframe/services/proxy"
	"github.com/zengjiwen/gameframe/sessions"
	"github.com/zengjiwen/gamenet"
	"github.com/zengjiwen/gamenet/server"
	"strings"
)

var _frontendServer gamenet.Server

func Run(serverType string, applies ...func(opts *options)) error {
	for _, apply := range applies {
		apply(&_opts)
	}

	env.ServerType = serverType
	if _opts.serviceAddr != "" {
		env.ServiceAddr = _opts.serviceAddr
	}
	if _opts.codec != nil {
		env.Codec = _opts.codec
	}
	if _opts.marshaler != nil {
		env.Marshaler = _opts.marshaler
	}

	if !_opts.standalone {
		if _opts.serviceAddr == "" {
			return errors.New("must specify service addr if not standalone!")
		}

		if _opts.sd != nil {
			env.SD = _opts.sd
		} else if _opts.sdAddr != "" {
			env.SD = servicediscovery.NewEtcd(_opts.sdAddr)
		} else {
			return errors.New("please specify service discovery or service discovery addr!")
		}

		if err := rpc.StartServer(services.NewService()); err != nil {
			return err
		}
		rpc.WatchServer()
		services.WatchServer()
	}

	if _opts.clientAddr != "" {
		cb := frontendEventCallback{}
		if strings.ToLower(_opts.concurrentMode) == "csp" {
			eventChan := make(chan func())
			_frontendServer = server.NewServer("tcp", _opts.clientAddr, cb, server.WithEventChan(eventChan))
			go func() {
				for event := range eventChan {
					event()
				}
			}()
		} else if strings.ToLower(_opts.concurrentMode) == "actor" {
			_frontendServer = server.NewServer("tcp", _opts.clientAddr, cb)
		} else {
			panic(fmt.Sprintf("incorrect concurrent mode: %s", _opts.concurrentMode))
		}

		go _frontendServer.ListenAndServe()
	}

	return nil
}

func Shutdown() error {
	var err error
	if _frontendServer != nil {
		err = _frontendServer.Shutdown()
	}
	rpc.StopServer()
	return err
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
