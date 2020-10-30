package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// This is a public API key. You can't steal my information :D
var apiKey = "tw2BCBreQk2sfKMyZ8sZ8jNP4BbkjdxALFCh"

// Identifiers starting with a capital letter are
// exported (visible in other packages)
func GetRate(target string, out chan<- float64) {
	url := "https://currencyapi.net/api/v1/rates?key=" + apiKey

	// Open the TCP/HTTP stream
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("Failed to open HTTP stream :c")
	}

	// Close the body stream once this routine returns
	defer resp.Body.Close()

	// Read all bytes of the HTTP body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Failed to read HTTP body :c")
	}

	// Now unmarshal the body as a JSON object into a map (string -> any)
	var f map[string]interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Fatalln("Failed to parse JSON :c", err)
	}

	// Make use of some type assertions
	rates := f["rates"].(map[string]interface{})

	// Write the result to the channel
	out <- rates[target].(float64)
}

func main() {
	// Channels are used to send and receive data from goroutines. Normal channels can only
	// hold 1 piece of data at a time and they will block when reading from an empty channel
	// (until data is available) or writing to one that already contains data (until that data is
	// read)
	// These two channels are used to pass around 64-bit floating point values but they can hold any
	// type
	ch1 := make(chan float64)
	ch2 := make(chan float64)

	// We can run code once the function exits:
	defer close(ch1)
	defer close(ch2)

	fmt.Println("Fetching Conversion Rate USD -> EUR and USD -> CAD ...")

	// Launch the Goroutines
	go GetRate("EUR", ch1)
	go GetRate("CAD", ch2)

	// Fetch the incoming data
	fmt.Printf("$1 US is €%f\n", <-ch1)
	fmt.Printf("$1 US is $%f Canadian\n", <-ch2)
}
