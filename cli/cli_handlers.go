package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/http"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type CliHandlers interface {
	UpdateIndex() error
	CleanIndex() error
	Search(query string, numberOfResults int, usePager bool) error
	ShowFile(hash string, usePager bool) error
	StartServer() error
}

type CliHandlersImpl struct {
	config  config.Config
	index   index.Index
	server  http.HttpServer
	display CliDisplay
}

func NewCliHandlers(gsConfig config.Config, gsIndex index.Index, server http.HttpServer) CliHandlers {
	display := CliDisplayImpl{gsConfig}
	return &CliHandlersImpl{gsConfig, gsIndex, server, &display}
}

func (s *CliHandlersImpl) UpdateIndex() error {
	err := s.updateIndex()
	if err != nil {
		return err
	}

	return err
}

func (s *CliHandlersImpl) CleanIndex() error {
	res, err := s.index.Clean()
	if err != nil {
		return err
	}

	s.display.IndexClean(res)
	return err
}

func (s *CliHandlersImpl) Search(query string, numberOfResults int, usePager bool) error {
	err := s.updateIndex()
	if err != nil {
		return err
	}

	res, err := s.index.Search(query, numberOfResults, index.OutputAnsi)
	if err != nil {
		return err
	}

	if len(res.Matches) > numberOfResults {
		res.Matches = res.Matches[0:numberOfResults]
	}

	return s.display.Search(query, res, usePager)
}

// TODO: show files by path too
func (s *CliHandlersImpl) ShowFile(hash string, usePager bool) error {
	err := s.updateIndex()
	if err != nil {
		return err
	}

	res, err := s.index.FindDocumentById(hash)
	if err != nil {
		return err
	}

	return s.display.ShowFile(res, usePager)
}

func (s *CliHandlersImpl) StartServer() error {
	err := s.updateIndex()
	if err != nil {
		return err
	}

	addr, err := s.server.GetAvailableAddress()
	if err != nil {
		return err
	}

	serviceUrl := "http://" + addr
	go func() {
		time.Sleep(100 * time.Millisecond)
		err := openBrowser(serviceUrl)
		if err != nil {
			fmt.Println(err)
		}
	}()

	s.display.StartServer(serviceUrl)
	return s.server.Start(addr)
}

func (s *CliHandlersImpl) updateIndex() error {
	repoUpToDate, err := s.index.IsUpToDate()
	if err != nil {
		return errors.Wrap(err, "cannot check if index is up to date")
	}

	if !repoUpToDate {
		fmt.Println("Updating index ...")
		res, err := s.index.BuildWith(index.BuildOptionsSpacedBy())
		if err != nil {
			return errors.Wrap(err, "cannot build index")
		}
		s.display.IndexBuild(res)
	}
	return nil
}

// See: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openBrowser(url string) error {
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
	return err
}

func leftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}

func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}
