package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type CommandHandler struct {
	config Config
	index  Index
	server HttpServer
}

func NewCommandHandler(config Config) (CommandHandler, error) {
	index, err := NewIndex(config)
	if err != nil {
		return *new(CommandHandler), errors.Wrap(err, "unable to create command handler")
	}

	server, err := NewHttpServer(config)
	if err != nil {
		return *new(CommandHandler), errors.Wrap(err, "unable to create command handler")
	}

	return CommandHandler{config, index, server}, nil
}

func (s *CommandHandler) CleanIndex() error {
	start := time.Now()
	err := s.index.Clean()
	elapsed := time.Since(start)

	fmt.Println("Clean took ", elapsed)
	return err
}

func (s *CommandHandler) BuildIndex() error {
	start := time.Now()
	err := s.index.Build()
	elapsed := time.Since(start)

	fmt.Println("Indexing took ", elapsed)
	return err
}

func (s *CommandHandler) Search(query string) error {
	results, err := s.index.Search(query)
	if err != nil {
		return err
	}

	for index, match := range results {
		fmt.Println(strconv.Itoa(index) + ": " + match.File.Path)
	}
	return nil
}

func (s *CommandHandler) StartServer() error {
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://" + s.config.Web.ListenAddress)
	}()
	return s.server.Start()
}

func (s *CommandHandler) Destroy() {
	err := s.index.Close()
	if err != nil {
		fmt.Println("Error while closing index: ", err)
	}
}

// See: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}
