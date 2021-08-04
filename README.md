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
	"log"
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
		gameframe.RegisterClientHandler("hello", handleHello)
		if err := gameframe.Run(*_serverID, *_serverType,
			gameframe.WithClientAddr(*_clientAddr),
			gameframe.WithConcurrentMode("actor"),
			gameframe.WithServiceAddr(*_serviceAddr)); err != nil {
			log.Println(err)
			return
		}
	} else if *_serverType == "game" {
		room := Room{}
		gameframe.RegisterClientHandler("room.joinRoom", room.joinRoom)
		if err := gameframe.Run(*_serverID, *_serverType,
			gameframe.WithConcurrentMode("csp"),
			gameframe.WithServiceAddr(*_serviceAddr)); err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println(errors.New("incorrect server type"))
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case sig := <-c:
		log.Println("receive signal: %s", sig.String())
	case err := <-gameframe.GetDieChan():
		log.Println(err)
	}

	if err := gameframe.Shutdown(); err != nil {
		log.Println(err)
	}
}

func handleHello(session *gameframe.Session) {

}

type Room struct{}

func (r Room) joinRoom(session *gameframe.Session) {

}
```