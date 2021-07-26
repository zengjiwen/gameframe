package main

import (
	"errors"
	"flag"
	"github.com/zengjiwen/gameframe"
	"github.com/zengjiwen/gameframe/service"
	"os"
	"os/signal"
)

var _serverID = flag.String("serverid", "", "specify the server id")
var _serverType = flag.String("servertype", "", "specify the server type")
var _clientAddr = flag.String("clientaddr", "", "specify the client addr")
var _serviceAddr = flag.String("serviceaddr", "", "specify the service addr")

func main() {
	flag.Parse()

	if *_serverType == "gateway" {
		service.RegisterClientHandler("hello", handleHello)
		if err := gameframe.Run(*_serverID, *_serverType,
			gameframe.WithClientAddr(*_clientAddr),
			gameframe.WithConcurrentMode("actor"),
			gameframe.WithServiceAddr(*_serviceAddr)); err != nil {
			panic(err)
		}
	} else if *_serverType == "game" {
		room := Room{}
		service.RegisterClientHandler("room.joinRoom", room.joinRoom)
		if err := gameframe.Run(*_serverID, *_serverType,
			gameframe.WithConcurrentMode("csp"),
			gameframe.WithServiceAddr(*_serviceAddr)); err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("incorrect server type"))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case <-c:
	case <-gameframe.GetDieChan():
	}

	if err := gameframe.Shutdown(); err != nil {
		panic(err)
	}
}

func handleHello() {

}

type Room struct{}

func (r Room) joinRoom() {

}
