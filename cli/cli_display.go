//go:generate mockgen -package mock -destination ../mocks/mocks_CliDisplay.go gitlab.com/remipassmoilesel/gitsearch/cli CliDisplay
package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ttacon/chalk"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type CliDisplay interface {
	IndexBuild(res index.BuildOperationResult)
	IndexClean(res index.CleanOperationResult)
	Search(query string, res index.SearchResult, usePager bool) error
	ShowFile(file index.IndexedFile, usePager bool) error
	StartServer(serviceUrl string)
}

type CliDisplayImpl struct {
	config config.Config
}

func (d *CliDisplayImpl) IndexBuild(res index.BuildOperationResult) {
	fmt.Println(fmt.Sprintf("Indexed %v/%v files in %v seconds. Oldest commit: %v", res.Files, res.TotalFiles, res.TookSeconds, res.OldestCommit))
}

func (d *CliDisplayImpl) IndexClean(res index.CleanOperationResult) {
	fmt.Println(fmt.Sprintf("Index clean took %d ms", res.TookMs))
}

func (d *CliDisplayImpl) Search(query string, res index.SearchResult, usePager bool) error {
	output := fmt.Sprintf("Query: %s\n", query)
	for idx, match := range res.Matches {
		// header
		icon := "➡️"
		if idx == 0 {
			icon = "🎀️"
		}
		rank := fmt.Sprintf("%v  %v: ", icon, strconv.Itoa(idx+1))
		commit := fmt.Sprintf("Commit: %v", string([]rune(match.File.Commit[0:15])))
		date := fmt.Sprintf("Date: %v", match.File.Date)
		id := fmt.Sprintf("Id: %v", match.File.Hash[0:15])
		header := rank + match.File.Path + " - " + commit + " - " + date + " - " + id

		// body
		fragments := strings.Join(match.Fragments, "\n\n---\n\n")
		bodyLines := strings.Split(fragments, "\n")
		body := "\n    " + strings.Join(bodyLines, "\n    ")

		output += fmt.Sprintf("\n%s%s%s", chalk.Cyan, header, chalk.Reset)
		output += body
		output += "\n"
	}

	if len(res.Matches) < 1 {
		output += "\nNothing found."
	}

	output += fmt.Sprintf("Search took %v ms", res.TookMs)
	if usePager {
		return lessPipe(output)
	} else {
		fmt.Println(output)
		return nil
	}
}

func (d *CliDisplayImpl) ShowFile(file index.IndexedFile, usePager bool) error {
	commit := fmt.Sprintf("Commit: %v", string([]rune(file.Commit[0:15])))
	date := fmt.Sprintf("Date: %v", file.Date)
	id := fmt.Sprintf("Id: %v", file.Hash[0:15])
	header := file.Path + " - " + commit + " - " + date + " - " + id

	output := fmt.Sprintf("\n%s%s%s", chalk.Cyan, header, chalk.Reset)
	output += file.Content
	output += "\n"

	if usePager {
		return lessPipe(output)
	} else {
		fmt.Println(output)
		return nil
	}
}

func (d *CliDisplayImpl) StartServer(serviceUrl string) {
	fmt.Println("Listening on " + serviceUrl)
	fmt.Println("Press CTRL+C to quit")
}

func lessPipe(content string) error {
	pager := "/usr/bin/less"
	if len(os.Getenv("PAGER")) > 0 {
		pager = os.Getenv("PAGER")
	}
	cmd := exec.Command(pager)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("*** You should install a pager and set $PAGER variable in your shell ***")
		return errors.Wrap(err, "cannot pipe output to less")
	}
	return nil
}
