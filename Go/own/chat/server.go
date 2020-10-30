package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

var opSend rune = 0
var opRead rune = 1

var clients = make([]chan msg, 0)

type msg struct {
	Text   string
	Source net.Addr
}

// Handles incoming connections
// FIXME: client needs to be removed from the broadcast pool
func handle(conn net.Conn, in <-chan msg, out chan<- msg) {
	defer conn.Close()

	for {
		// Read what the client wants to do
		var op rune
		err := binary.Read(conn, binary.LittleEndian, &op)

		if err != nil {
			break
		}

		switch op {
		case opSend:
			// if it wants to send a message, read it and send it over to "handleTransmissions"
			// to broadcast it to all other clients.
			var len uint16
			binary.Read(conn, binary.LittleEndian, &len)

			m := make([]byte, len)
			conn.Read(m)

			out <- msg{fmt.Sprintf("%s >> %s", conn.RemoteAddr(), m), conn.RemoteAddr()}
		case opRead:
			// Otherwise, check the message input channel. There might be some other messages
			// left
			var c uint16 = uint16(len(in))
			binary.Write(conn, binary.LittleEndian, c)

			// Send then to the client!
			var i uint16
			for i = 0; i < c; i++ {
				s := <-in
				var length uint16 = uint16(len(s.Text))
				binary.Write(conn, binary.LittleEndian, length)
				conn.Write([]byte(s.Text))
			}
		}
	}

	// If the handler encounters an error, that probably means that the client
	// disconnected. Broadcast that to the other clients.
	out <- msg{fmt.Sprintf("[@] %s left.", conn.RemoteAddr()), conn.RemoteAddr()}

	for i, client := range clients {
		if client == in {
			clients[i] = nil
			return
		}
	}
}

// Handles all clients and broadcasting to them
func handleTransmissions(clientStream <-chan chan msg, input <-chan msg) {
	for {
		select {
		case s := <-clientStream:
			// We've connected a new client, let's add it to the list
			clients = append(clients, s)
		case i := <-input:
			fmt.Println(i.Text)

			// Now broadcast a new message to all connected clients
			for _, ch := range clients {
				if ch != nil {
					// (the transmission back to the client is handled by handle())
					ch <- i
				}
			}
		}
	}
}

// Runs the main chat server.
func main() {
	// Open a new socket on port 8102
	lst, err := net.Listen("tcp", ":8102")

	if err != nil {
		log.Fatalln("Could not open a listening socket on localhost:8102")
	}

	fmt.Println("Listening on localhost:8102")

	clients := make(chan chan msg)
	input := make(chan msg, 5)

	// This goroutine handles broadcasting incoming messages to all clients
	go handleTransmissions(clients, input)

	for {
		// Now, wait for a new client to connect.
		conn, err := lst.Accept()

		if err != nil {
			log.Println("Connection could not be established. Borked client?")
		} else {
			// Okay we're now fully connected. Let's create a buffered channel for all the messages.
			out := make(chan msg, 20)

			// And send it over to the handleTransmissions goroutine
			clients <- out

			// And notify all clients that a new one joined
			input <- msg{fmt.Sprintf("[@] %s joined.", conn.RemoteAddr()), conn.RemoteAddr()}

			// Now, turn the handling over to a seperate handler routine
			go handle(conn, out, input)
		}
	}
}
