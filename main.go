package main

import (
	"fmt"
	"github.com/remipassmoilesel/gitsearch/cli"
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/http"
	"github.com/remipassmoilesel/gitsearch/index"
	"log"
	"os"

	_ "github.com/remipassmoilesel/gitsearch/web_client"
)

func main() {
	fmt.Println()
	fmt.Println("🔎 Gitsearch ⚡")
	fmt.Println()

	err := applyCommand(os.Args)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func applyCommand(args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	gsIndex, err := index.NewIndex(cfg)
	if err != nil {
		return err
	}
	defer closeIndex(gsIndex)

	server := http.NewHttpServer(cfg, gsIndex)
	cliParser := cli.NewCliParser(cfg, gsIndex, server)

	return cliParser.ApplyCommand(args)
}

func closeIndex(index index.Index) {
	err := index.Close()
	if err != nil {
		fmt.Println("Error while closing index: ", err)
	}
}
