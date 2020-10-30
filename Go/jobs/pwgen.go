package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Make sure we actually generate random numbers
	rand.Seed(time.Now().UnixNano())
	pwd := ""

	// Generate 10 random chars and append the to pwd
	for i := 0; i < 10; i++ {
		// the printable ASCII chars begin at char 32 (space)
		// and end at char 126 (~)
		char := 32 + rand.Intn(95)
		pwd += string(rune(char))
	}

	fmt.Printf("%s", pwd)
}
