package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/jhillyerd/enmime"
)

func parse() {
	filename := os.Args[1]

	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file, %v", err)
	}

	message, err := enmime.ReadEnvelope(bytes.NewReader(filedata))
	if err != nil {
		log.Fatalf("failed to parse mail, %v", err)
	}

	log.Printf(message.HTML)
}
