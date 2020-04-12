package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	err := applyCommand(os.Args)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func applyCommand(args []string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	index, err := NewIndex(config)
	if err != nil {
		return err
	}
	defer closeIndex(index)

	server := NewHttpServer(config, index)
	cliParser := NewCliParser(config, index, server)

	return cliParser.applyCommand(args)
}

func closeIndex(index Index) {
	err := index.Close()
	if err != nil {
		fmt.Println("Error while closing index: ", err)
	}
}
