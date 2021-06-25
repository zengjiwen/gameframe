# gameframe
game micro service frame

# feature
1. custom codec
2. custom marshaler
3. reflect handler
4. service discovery
5. rpc
6. group

in developing...

example:
```go
package main

import (
	"flag"
	"github.com/zengjiwen/gameframe"
	"github.com/zengjiwen/gameframe/services"
)

var _serverType = flag.String("servertype", "", "specify the server type")

func main() {
	flag.Parse()

	if *_serverType == "gateway" {
		services.RegisterClientHandler("hello", handleHello)
		gameframe.Run(*_serverType, "127.0.0.1:6666",
			gameframe.WithClientAddr("0.0.0.0:7777"),
			gameframe.WithConcurrentMode("actor"))
	} else if *_serverType == "game" {
		room := Room{}
		services.RegisterClientHandler("room.joinRoom", room.joinRoom)
		gameframe.Run(*_serverType, "127.0.0.1:8888",
			gameframe.WithConcurrentMode("csp"))
	} else {
		panic("incorrect server type!")
	}
}

func handleHello() {

}

type Room struct{}

func (r Room) joinRoom() {

}
```