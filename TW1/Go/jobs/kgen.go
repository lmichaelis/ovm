package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	// Make sure we actually generate random numbers
	rand.Seed(time.Now().UnixNano())

	// All IDs should begin with "KD"
	kn := "KD"
	qsum := 0

	// Generate a sequence of 8 random numbers
	for i := 0; i < 8; i++ {
		z := 1 + rand.Intn(8)
		qsum += z
		kn += strconv.Itoa(z)
	}

	// Make sure there are _two_ last digits
	if qsum < 10 {
		kn += "0"
	}

	// Append the digit to the number
	kn += strconv.Itoa(qsum)
	fmt.Println(kn)
}
