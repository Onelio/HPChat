package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

var port = flag.String("port", ":",
	"Set listen port")

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
				fmt.Print(">")
				continue
			}
			// On message
			//if address == peer { }
			fmt.Printf("\r%s: %s>", request.Room, request.Data)
		}
	}()
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		//conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		if peer == nil {
			continue
		}
		_, _ = conn.WriteToUDP([]byte(text), peer)
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
}
