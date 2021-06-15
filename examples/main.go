package main

import (
	"github.com/zengjiwen/gameframe"
	"github.com/zengjiwen/gameframe/services"
)

func main() {
	services.RegisterClientHandler("hello", handleHello)

	room := Room{}
	services.RegisterClientHandler("room.joinRoom", room.joinRoom)

	gameframe.Start("game", "127.0.0.1:6666", true,
		gameframe.WithConcurrentMode("csp"))
}

func handleHello() {

}

type Room struct {
}

func (r Room) joinRoom() {

}
