package main

import "github.com/zengjiwen/gameframe"

func main() {
	gameframe.RegisterHandler("hello", handleHello)
	room := Room{}
	gameframe.RegisterHandler("room.joinroom", room.joinRoom)
	gameframe.Start()
}

func handleHello() {

}

type Room struct {
}

func (r Room) joinRoom() {

}
