//go:generate mockgen -package mock -destination ../test/mock/mocks_CliHandlers.go gitlab.com/remipassmoilesel/gitsearch/cli CliHandlers
package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/http"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	"gitlab.com/remipassmoilesel/gitsearch/utils"
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
	utils   utils.Utils
}

func NewCliHandlers(gsConfig config.Config, gsIndex index.Index, server http.HttpServer) CliHandlers {
	display := NewCliDisplay(gsConfig)
	return &CliHandlersImpl{gsConfig, gsIndex, server, display, utils.NewUtils()}
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

	return s.display.Display(s.display.IndexClean(res), false)
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

	return s.display.Display(s.display.Search(res), usePager)
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

	return s.display.Display(s.display.ShowFile(res), usePager)
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
		// Here we wait a little for server start before opening browser
		time.Sleep(100 * time.Millisecond)
		err := s.utils.OpenWebBrowser(serviceUrl)
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = s.display.Display(s.display.StartServer(serviceUrl), false)
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
		return s.display.Display(s.display.IndexBuild(res), false)
	}
	return nil
}
