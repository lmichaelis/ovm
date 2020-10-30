package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

var opSend rune = 0
var opRead rune = 1

// Just read messages from stdin and write them to
// the "out" channel.
func readin(out chan<- string) {
	// I don't know why I have to use a bufio.Reader for this but meh
	reader := bufio.NewReader(os.Stdin)

	for {
		s, _ := reader.ReadString('\n')
		out <- s[:len(s)-1]
	}
}

func main() {
	// Connect to localhost:8102 where the server should be running.
	// If this fails, make sure the server is running.
	conn, err := net.Dial("tcp", ":8102")

	if err != nil {
		log.Fatalln("Failed to connect to localhost:8102")
	}

	// Let's read user input asynchronously
	in := make(chan string)
	go readin(in)

	for {

		// We are going to make use of select-statements here. It tries to read a value from
		// "in" but can't so it just continues to the default branch
		select {
		case d := <-in:
			// If there is some data in the "in" channel, that means that
			// the user sent a message. Let's send that message to the server
			binary.Write(conn, binary.LittleEndian, opSend)
			var size uint16 = uint16(len(d))

			binary.Write(conn, binary.LittleEndian, size)
			conn.Write([]byte(d))
		default:
			// Otherwise the user must be waiting for new messages. Let's ask the server for
			// some. This should probably be done using a timeout since it just spams the server with 0's.
			binary.Write(conn, binary.LittleEndian, opRead)
			var c uint16
			binary.Read(conn, binary.LittleEndian, &c)

			var i uint16
			var l uint16

			// And now retrieve the data and print it for the user to see.
			for i = 0; i < c; i++ {
				binary.Read(conn, binary.LittleEndian, &l)
				v := make([]byte, l)
				conn.Read(v)

				fmt.Println(string(v))
			}
		}
	}
}
