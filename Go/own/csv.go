package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Parses the given string as CSV
func Parse(csv string) [][]string {
	lines := strings.Split(csv, "\n")
	values := make([][]string, len(lines))

	for i, line := range lines {
		entries := strings.Split(line, ",")
		values[i] = make([]string, len(entries))

		for j, entry := range entries {
			values[i][j] = entry
		}
	}

	return values
}

// Converts the given CSV-Data into a map
func ToMap(parsed [][]string) []map[string]string {
	mm := make([]map[string]string, len(parsed))

	for i, row := range parsed[1:] {
		mm[i] = make(map[string]string)
		for j, col := range row {
			mm[i][parsed[0][j]] = col
		}
	}

	return mm
}

// Converts the given map to JSON (borked)
func ToJSON(mm map[string]string) string {
	s := "{\n"

	for k, v := range mm {
		fmt.Println(k)
		fmt.Println(v)
		s += "\"" + k
		s += "\": \"" + v
		s += "\",\n\n"
	}

	return s + "}"
}

func main() {
	// Let's use a command-line argument!
	if len(os.Args) < 2 {
		fmt.Print("Usage: csv <file>")
		os.Exit(-1)
	}

	// Read in the file
	bytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatalln(err)
	}

	// And now: parse the CSV!
	data := Parse(string(bytes))

	mm := ToMap(data)
	fmt.Print(ToJSON(mm[0]))

}
