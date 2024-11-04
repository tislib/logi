package main

import (
	"github.com/tislib/logi/pkg/parser/logi"
	"log"
	"os"
)

func main() {
	var file = os.Args[1]

	log.Print("Processing file: ", file)

	content, err := os.ReadFile(file)

	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	parsed, err := logi.ParseFullWithMacro(string(content), ``)

	if err != nil {
		log.Fatal("Error parsing file: ", err)
	}

	log.Print("File content: ", parsed)
}
