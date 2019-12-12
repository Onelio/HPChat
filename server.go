package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
)

var cRooms = make(map[string]*net.UDPAddr)

func main() {
	flag.Parse()
	fmt.Println("Listening...")
	address, err := net.ResolveUDPAddr("udp4", *addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var buffer [1024]byte
	for {
		length, address, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			continue
		}
		var request ChatRequest
		err = json.Unmarshal(buffer[:length], &request)
		if err != nil || request.Room != ServerAuth {
			fmt.Println(address.String(), "error", request.Room)
			continue
		}

		// Client is valid, proceed
		if pConn, ok := cRooms[request.Data]; ok {
			room := request.Data
			// Send first
			request.Data = address.String()
			response, _ := json.Marshal(&request)
			_, _ = conn.WriteToUDP(response, pConn)
			// Send second
			request.Data = pConn.String()
			response, _ = json.Marshal(&request)
			_, _ = conn.WriteToUDP(response, address)
			// Delete on complete
			delete(cRooms, room)
			fmt.Println(address.String(), "conn", request.Room)
			continue
		}
		cRooms[request.Data] = address
		fmt.Println(address.String(), "init", request.Data)
	}
}
