package main

import "flag"

var addr = flag.String("addr", "127.0.0.1:8080",
	"Set connection addr")

const ServerAuth = "0xCAFEBABE"

type ChatRequest struct {
	Room string
	Data string
}
