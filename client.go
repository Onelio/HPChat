package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var port = flag.String("port", ":",
	"Set listen port")

var room = flag.String("room", "ChatRoomNum_01",
	"Room access surname")

var user = flag.String("user", "UnknownUserName",
	"Room peers username")

func readMessages(conn *net.UDPConn) {
	var buffer [1024]byte
	var peer *net.UDPAddr
	go func() {
		for {
			length, _, err := conn.ReadFromUDP(buffer[:])
			if err != nil {
				continue
			}
			var request ChatRequest
			err = json.Unmarshal(buffer[:length], &request)
			if err != nil {
				continue
			}
			// On server peer set
			if request.Room == ServerAuth {
				peer, err = net.ResolveUDPAddr("udp4", request.Data)
				if err != nil {
					panic(err)
				}
				// Set first message
				_, _ = conn.WriteToUDP([]byte("test"), peer)
				fmt.Printf("\rEnter room %s\n>", *room)
				continue
			}
			// On message
			fmt.Printf("\r%s: %s>", request.Room, request.Data)
		}
	}()
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		// On close
		if strings.HasPrefix(text, "/quit") {
			conn.Close()
			break
		}
		if peer == nil {
			continue
		}
		var request = ChatRequest{
			Room: *user,
			Data: text,
		}
		data, _ := json.Marshal(&request)
		_, _ = conn.WriteToUDP(data, peer)
		fmt.Print(">")
	}
}

func main() {
	flag.Parse()
	fmt.Println("Listening on", *port)
	address, err := net.ResolveUDPAddr("udp4", ":"+*port)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Write request
	server, err := net.ResolveUDPAddr("udp4", *addr)
	if err != nil {
		panic(err)
	}
	var request = ChatRequest{
		Room: ServerAuth,
		Data: *room,
	}
	data, _ := json.Marshal(&request)
	_, _ = conn.WriteToUDP(data, server)

	// Proceed to listen area
	readMessages(conn)
	fmt.Println("\nBye")
}
