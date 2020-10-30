package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Read 32 chars from stdin
	var code string
	fmt.Scanf("%32s", &code)

	// Make sure none of the chars are a letter or control-char
	for _, char := range code {
		if char < 48 || char > 57 {
			// If this char is not a digit, print "1", then exit the program
			// with exit code 1 (on unix the return code can be accessed by
			// executing "echo $?" in a terminal after running this program)
			fmt.Println("1")
			os.Exit(1)
		}
	}

	// Parse the needed substrings as integers.
	// These always succeed (we've made sure it only contains numbers above)
	ng, _ := strconv.Atoi(code[9:14])
	bg, _ := strconv.Atoi(code[18:24])

	if ng <= bg {
		fmt.Println("0")
		os.Exit(0)
	} else {
		fmt.Println("2")
		os.Exit(2)
	}
}
