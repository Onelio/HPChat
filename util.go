package main

import "flag"

var addr = flag.String("addr", "127.0.0.1:8080",
	"Set connection addr")
var room = flag.String("room", "ChatRoomNum_01",
	"Room access surname")

const ServerAuth = "0xCAFEBABE"

type ChatRequest struct {
	Room string
	Data string
}
