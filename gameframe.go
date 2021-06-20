package gameframe

import (
	"fmt"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/services"
	"github.com/zengjiwen/gameframe/services/proxy"
	"github.com/zengjiwen/gameframe/sessions"
	"github.com/zengjiwen/gamenet"
	"github.com/zengjiwen/gamenet/server"
	"os"
	"os/signal"
	"strings"
)

func Run(serverType, addr string, isFrontend bool, applies ...func(opts *options)) {
	for _, apply := range applies {
		apply(&_opts)
	}

	env.ServerType = serverType
	env.Addr = addr
	env.IsFrontend = isFrontend

	if _opts.concurrentMode != "" {
		env.ConcurrentMode = _opts.concurrentMode
	}
	if _opts.codec != nil {
		env.Codec = _opts.codec
	}
	if _opts.marshaler != nil {
		env.Marshaler = _opts.marshaler
	}

	var tcpServer gamenet.Server
	if isFrontend {
		if strings.ToLower(env.ConcurrentMode) == "csp" {
			eventChan := make(chan func())
			tcpServer = server.NewServer("tcp", addr, frontendEventCallback{},
				server.WithEventChan(eventChan))

			go func() {
				for event := range eventChan {
					event()
				}
			}()
		} else if strings.ToLower(env.ConcurrentMode) == "actor" {
			tcpServer = server.NewServer("tcp", "127.0.0.1:0", frontendEventCallback{})
		} else {
			panic(fmt.Sprintf("incorrect concurrent mode: %s", env.ConcurrentMode))
		}

		go tcpServer.ListenAndServe()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	if isFrontend {
		tcpServer.Shutdown()
	}
}

type frontendEventCallback struct{}

func (fecb frontendEventCallback) OnNewConn(conn gamenet.Conn) {
	clientProxy := proxy.NewClient(conn)
	session := sessions.New(clientProxy)
	conn.SetUserData(session)
}

func (fecb frontendEventCallback) OnConnClosed(conn gamenet.Conn) {
	session, ok := conn.UserData().(*sessions.Session)
	if !ok {
		return
	}

	session.OnClosed()
}

func (fecb frontendEventCallback) OnRecvData(conn gamenet.Conn, data []byte) {
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
