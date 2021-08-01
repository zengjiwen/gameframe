package gameframe

import (
	"errors"
	"fmt"
	"github.com/zengjiwen/gameframe/codec"
	"github.com/zengjiwen/gameframe/marshaler"
	"github.com/zengjiwen/gameframe/rpc"
	"github.com/zengjiwen/gameframe/service"
	"github.com/zengjiwen/gameframe/servicediscovery"
	"github.com/zengjiwen/gameframe/sessions"
	"github.com/zengjiwen/gameframe/sessions/proxy"
	"github.com/zengjiwen/gamenet"
	"github.com/zengjiwen/gamenet/server"
	"log"
	"strings"
)

var _gameFrame = &gameFrame{
	dieChan: make(chan error),
	opts: options{
		concurrentMode: "actor",
	},
}

type gameFrame struct {
	serverID       string
	serverType     string
	dieChan        chan error
	opts           options
	frontendServer gamenet.Server
}

type Session sessions.Session

func Run(serverID, serverType string, opts ...Option) error {
	_gameFrame.serverID = serverID
	_gameFrame.serverType = serverType

	for _, opt := range opts {
		opt.apply(&_gameFrame.opts)
	}

	gOpts := &_gameFrame.opts
	if !gOpts.standalone {
		if gOpts.serviceAddr == "" {
			return errors.New("must specify service addr if not standalone!")
		}

		servicediscovery.InitServerInfo(_gameFrame.serverID, _gameFrame.serverType, gOpts.serviceAddr)
		service.FillServerInfo()
		if servicediscovery.Get() == nil && gOpts.sdAddr != "" {
			servicediscovery.Set(servicediscovery.NewEtcd(gOpts.sdAddr, _gameFrame.dieChan))
			if err := servicediscovery.Get().Start(); err != nil {
				return err
			}
		} else {
			return errors.New("please specify service discovery or service discovery addr!")
		}

		if err := rpc.StartServer(gOpts.serviceAddr, service.NewService()); err != nil {
			return err
		}
		rpc.InitClients()
		rpc.WatchServer()
		service.WatchServer()
	}

	if gOpts.clientAddr != "" {
		cb := frontendEventCallback{}
		if strings.ToLower(gOpts.concurrentMode) == "csp" {
			eventChan := make(chan func())
			_gameFrame.frontendServer = server.NewServer("TCP", gOpts.clientAddr, cb, server.WithEventChan(eventChan))
			go func() {
				for event := range eventChan {
					event()
				}
			}()
		} else if strings.ToLower(gOpts.concurrentMode) == "actor" {
			_gameFrame.frontendServer = server.NewServer("TCP", gOpts.clientAddr, cb)
		} else {
			return fmt.Errorf("incorrect concurrent mode: %s", gOpts.concurrentMode)
		}

		go func() {
			_gameFrame.dieChan <- _gameFrame.frontendServer.ListenAndServe()
		}()
	}

	return nil
}

func Shutdown() error {
	var err error
	if _gameFrame.frontendServer != nil {
		err = _gameFrame.frontendServer.Shutdown()
	}

	if !_gameFrame.opts.standalone {
		if err == nil {
			err = rpc.CloseClients()
		} else {
			rpc.CloseClients()
		}
		rpc.StopServer()

		if err == nil {
			err = servicediscovery.Get().Close()
		} else {
			servicediscovery.Get().Close()
		}
	}
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
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	session, ok := conn.UserData().(*sessions.Session)
	if !ok {
		return
	}

	m, err := codec.Get().Decode(data)
	if err != nil {
		return
	}

	ret, err := service.HandleClientMsg(session, m)
	if err != nil {
		return
	}

	conn.Send(ret)
}

func RegisterClientHandler(route string, ch interface{}) {
	service.RegisterClientHandler(route, ch)
}

func RegisterServerHandler(route string, sh interface{}) {
	service.RegisterServerHandler(route, sh)
}

func RegisterCodec(cd codec.Codec) {
	codec.Set(cd)
}

func RegisterMarshaler(ma marshaler.Marshaler) {
	marshaler.Set(ma)
}

func RegisterServiceDiscovery(sd servicediscovery.ServiceDiscovery) {
	servicediscovery.Set(sd)
}

func GetDieChan() chan error {
	return _gameFrame.dieChan
}
