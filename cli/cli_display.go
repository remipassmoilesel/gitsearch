//go:generate mockgen -package mock -destination ../test/mock/mocks_CliDisplay.go gitlab.com/remipassmoilesel/gitsearch/cli CliDisplay
package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ttacon/chalk"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/index"
	"gitlab.com/remipassmoilesel/gitsearch/utils"
	"strconv"
	"strings"
)

type CliDisplay interface {
	IndexBuild(res index.BuildOperationResult) string
	IndexClean(res index.CleanOperationResult) string
	Search(res index.SearchResult) string
	ShowFile(file index.IndexedFile) string
	StartServer(serviceUrl string) string
	Display(output string, withPager bool) error
}

func NewCliDisplay(config config.Config) CliDisplay {
	return &CliDisplayImpl{
		config: config,
		utils:  utils.NewUtils(),
	}
}

type CliDisplayImpl struct {
	config config.Config
	utils  utils.Utils
}

func (d *CliDisplayImpl) Display(output string, withPager bool) error {
	if withPager {
		err := d.utils.Pager(output)
		if err != nil {
			d.utils.PrintLn("*** You should install a pager and set $PAGER variable in your shell ***")
			d.utils.PrintLn(output)
		}
		return errors.Wrap(err, "cannot pipe output to less")
	} else {
		d.utils.PrintLn(output)
		return nil
	}
}

func (d *CliDisplayImpl) IndexBuild(res index.BuildOperationResult) string {
	return fmt.Sprintf("Indexed %v/%v files in %v seconds. Oldest commit: %v", res.Files, res.TotalFiles, res.TookSeconds, res.OldestCommit)
}

func (d *CliDisplayImpl) IndexClean(res index.CleanOperationResult) string {
	return fmt.Sprintf("Index clean took %d ms", res.TookMs)
}

func (d *CliDisplayImpl) Search(res index.SearchResult) string {
	output := fmt.Sprintf("Query: %s\n", res.Query)
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
	return output
}

func (d *CliDisplayImpl) ShowFile(file index.IndexedFile) string {
	commit := fmt.Sprintf("Commit: %v", string([]rune(file.Commit[0:15])))
	date := fmt.Sprintf("Date: %v", file.Date)
	id := fmt.Sprintf("Id: %v", file.Hash[0:15])
	header := file.Path + " - " + commit + " - " + date + " - " + id

	output := fmt.Sprintf("\n%s%s%s", chalk.Cyan, header, chalk.Reset)
	output += "\n" + file.Content
	output += "\n"

	return output
}

func (d *CliDisplayImpl) StartServer(serviceUrl string) string {
	output := "Listening on " + serviceUrl
	output += "\nPress CTRL+C to quit"
	return output
}
