package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/hamcha/clessy/tg"
)

var clients []net.Conn

func startClientsServer(bind string) {
	listener, err := net.Listen("tcp", bind)
	assert(err)

	// Accept loop
	for {
		c, err := listener.Accept()
		if err != nil {
			log.Printf("Can't accept client: %s\n", err.Error())
			continue
		}
		clients = append(clients, c)
		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	b := bufio.NewReader(c)
	defer c.Close()

	// Start reading messages
	for {
		bytes, _, err := b.ReadLine()
		if err != nil {
			break
		}
		var cmd tg.ClientCommand
		err = json.Unmarshal(bytes, &cmd)
		if err != nil {
			log.Printf("[handleClient] Can't parse JSON: %s\r\n", err.Error())
			log.Printf("%s\n", string(bytes))
			continue
		}
		executeClientCommand(cmd)
	}
	removeCon(c)
}

func removeCon(c net.Conn) {
	for i, con := range clients {
		if c == con {
			clients = append(clients[:i], clients[i+1:]...)
		}
	}
}

func broadcast(message string) {
	for _, c := range clients {
		_, err := fmt.Fprintf(c, message+"\r\n")
		if err != nil {
			removeCon(c)
		}
	}
}
