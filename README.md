# gameframe
game micro service frame

# feature
1. custom codec
2. custom marshaler
3. reflect handler
4. custom service discovery
5. grpc
6. group
7. standalone mode
8. timeout, retry, circuit break

in developing...

example:

```go
package main

import (
	"errors"
	"flag"
	"github.com/zengjiwen/gameframe"
	"github.com/zengjiwen/gameframe/env"
	"github.com/zengjiwen/gameframe/service"
	"os"
	"os/signal"
)

var _serverType = flag.String("servertype", "", "specify the server type")
var _clientAddr = flag.String("clientaddr", "", "specify the client addr")
var _serviceAddr = flag.String("serviceaddr", "", "specify the service addr")

func main() {
	flag.Parse()

	if *_serverType == "gateway" {
		services.RegisterClientHandler("hello", handleHello)
		if err := gameframe.Run(*_serverType,
			gameframe.WithClientAddr(*_clientAddr),
			gameframe.WithConcurrentMode("actor"),
			gameframe.WithServiceAddr(*_serviceAddr)); err != nil {
			panic(err)
		}
	} else if *_serverType == "game" {
		room := Room{}
		services.RegisterClientHandler("room.joinRoom", room.joinRoom)
		if err := gameframe.Run(*_serverType,
			gameframe.WithConcurrentMode("csp"),
			gameframe.WithServiceAddr(*_serviceAddr)); err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("incorrect server type!"))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case <-c:
		close(env.DieChan)
	case <-env.DieChan:
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
```