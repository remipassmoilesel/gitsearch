package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type CliHandlers struct {
	config Config
	index  Index
	server HttpServer
}

func NewCliHandlers(config Config, index Index, server HttpServer) CliHandlers {
	return CliHandlers{config, index, server}
}

func (s *CliHandlers) BuildIndex() error {
	res, err := s.index.Build()
	if err != nil {
		return err
	}

	fmt.Println("Index build took ", res.TookSeconds)
	return err
}

func (s *CliHandlers) CleanIndex() error {
	res, err := s.index.Clean()
	if err != nil {
		return err
	}

	fmt.Println("Index clean took ", res.TookSeconds)
	return err
}

func (s *CliHandlers) Search(query string) error {
	res, err := s.index.Search(query)
	if err != nil {
		return err
	}

	// TODO: improve display
	for index, match := range res.Matches {
		fmt.Println(strconv.Itoa(index) + ": " + match.File.Path)
	}
	fmt.Println("Search took ", res.TookSeconds)
	return nil
}

func (s *CliHandlers) StartServer() error {
	var serviceUrl = "http://" + s.config.Web.ListenAddress
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(serviceUrl)
	}()
	fmt.Println("Listenning on " + serviceUrl)
	return s.server.Start()
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
