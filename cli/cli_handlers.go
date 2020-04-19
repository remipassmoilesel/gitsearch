package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/http"
	"github.com/remipassmoilesel/gitsearch/index"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type CliHandlers struct {
	config  config.Config
	index   index.Index
	server  http.HttpServer
	display CliDisplay
}

func NewCliHandlers(gsConfig config.Config, gsIndex index.Index, server http.HttpServer) CliHandlers {
	display := CliDisplay{gsConfig}
	return CliHandlers{gsConfig, gsIndex, server, display}
}

func (s *CliHandlers) BuildIndex() error {
	res, err := s.index.Build()
	if err != nil {
		return err
	}

	s.display.IndexBuild(res)
	return err
}

func (s *CliHandlers) CleanIndex() error {
	res, err := s.index.Clean()
	if err != nil {
		return err
	}

	s.display.IndexClean(res)
	return err
}

func (s *CliHandlers) Search(query string) error {
	err := s.checkIndex()
	if err != nil {
		return err
	}

	res, err := s.index.Search(query, 5, index.OutputAnsi)
	if err != nil {
		return err
	}

	s.display.Search(res)
	return nil
}

func (s *CliHandlers) StartServer() error {
	err := s.checkIndex()
	if err != nil {
		return err
	}

	serviceUrl := "http://" + s.config.Web.ListenAddress
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(serviceUrl)
	}()

	s.display.StartServer(serviceUrl)
	return s.server.Start()
}

func (s *CliHandlers) checkIndex() error {
	repoUpToDate, err := s.index.IsUpToDate()
	if err != nil {
		return errors.Wrap(err, "cannot check if index is up to date")
	}

	if !repoUpToDate {
		fmt.Println("Updating index ...")
		res, err := s.index.Build()
		if err != nil {
			return errors.Wrap(err, "cannot build index")
		}
		s.display.IndexBuild(res)
	}
	return nil
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

func leftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}

func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}
