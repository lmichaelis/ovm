package main

import (
	"fmt"
	"strconv"
)

func sqAsync(a <-chan int, b chan<- int) {
	for v := range a {
		b <- v * v
	}
}

// This is a program for testing asynchronous behavior
// using goroutines.
func main() {
	v := make(chan int)
	w := make(chan int)

	defer close(v)
	defer close(w)

	// Let's start this goroutine that checks for new values being
	// inserted into the channel v and then puts that number squared
	// into the channel w
	go sqAsync(v, w)

	for {
		// Read some user input (yes, this is unsafe)
		var tmp string
		fmt.Scanf("%s", &tmp)

		// ... and convert it to an integer
		k, err := strconv.Atoi(tmp)

		if err != nil {
			continue
		}

		// Then write it to the "input" channel
		v <- k

		// ... and retrieve it from the output channel
		fmt.Println(<-w)
	}
}
