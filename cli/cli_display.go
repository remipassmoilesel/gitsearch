package cli

import (
	"fmt"
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/index"
	"github.com/ttacon/chalk"
	"strconv"
	"strings"
)

type CliDisplay struct {
	config config.Config
}

func (d CliDisplay) IndexBuild(res index.BuildOperationResult) {
	fmt.Println(fmt.Sprintf("Indexed %v/%v files in %v seconds", res.Files, res.TotalFiles, res.TookSeconds))
}

func (d CliDisplay) IndexClean(res index.CleanOperationResult) {
	fmt.Println(fmt.Sprintf("Index clean took %v ms", res.TookMs))
}

func (d CliDisplay) Search(res index.SearchResult) {
	for idx, match := range res.Matches {
		// header
		icon := "➡️"
		if idx == 0 {
			icon = "🎀️"
		}
		position := fmt.Sprintf("%v  %v: ", icon, strconv.Itoa(idx+1))
		commit := fmt.Sprintf("Commit: %v", string([]rune(match.File.Commit[0:15])))
		date := fmt.Sprintf("Date: %v", match.File.Date)
		id := fmt.Sprintf("Id: %v", match.File.Hash[0:15])
		header := position + match.File.Path + " - " + commit + " - " + date + " - " + id

		// body
		fragments := strings.Join(match.Fragments, "\n\n---\n\n")
		bodyLines := strings.Split(fragments, "\n")
		body := "\n    " + strings.Join(bodyLines, "\n    ")

		fmt.Println()
		fmt.Println(chalk.Cyan, header, chalk.Reset)
		fmt.Println(body)
		fmt.Println()
	}

	if len(res.Matches) < 1 {
		fmt.Println()
		fmt.Println("Nothing found. ")
	}

	fmt.Println(fmt.Sprintf("Search took %v ms", res.TookMs))
}

func (d CliDisplay) ShowFile(file index.IndexedFile) {
	commit := fmt.Sprintf("Commit: %v", string([]rune(file.Commit[0:15])))
	date := fmt.Sprintf("Date: %v", file.Date)
	id := fmt.Sprintf("Id: %v", file.Hash[0:15])
	header := file.Path + " - " + commit + " - " + date + " - " + id

	fmt.Println(chalk.Cyan, header, chalk.Reset)
	fmt.Println()
	fmt.Println(file.Content)
	fmt.Println()
}

func (d CliDisplay) StartServer(serviceUrl string) {
	fmt.Println("Listening on " + serviceUrl)
}
