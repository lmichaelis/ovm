package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// Struct members also have to be capitalized to be exported.
type Vector struct {
	X int32
	Y int32
}

func main() {
	// Open a new file ("test.bin")
	file, err := os.Create("test.bin")

	if err != nil {
		log.Fatalln("Failed to open \"test.bin\"")
	}

	// Write a struct of type Vector to the file by utilizing a
	// byte buffer and the builtin binary.Write function.
	test := Vector{100, 200}
	bbuf := new(bytes.Buffer)
	binary.Write(bbuf, binary.BigEndian, test)
	file.Write(bbuf.Bytes())
	file.Close()

	// Open "test.bin" for reading
	file, err = os.Open("test.bin")
	if err != nil {
		log.Fatalln("Failed to open \"test.bin\"")
	}

	// Get the file's size and read the bytes inside that file into
	// an array. make() is a special builtin function that will allocate
	// arrays, channels and maps on the heap
	stat, err := file.Stat()
	bbuf2 := make([]byte, stat.Size())
	file.Read(bbuf2)

	// Create a new byte reader
	reader := bytes.NewReader(bbuf2)

	// and read the struct from the file
	var r Vector
	binary.Read(reader, binary.BigEndian, &r)

	// Then print out it's contents
	fmt.Println(r)
	file.Close()
}
