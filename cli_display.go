package main

import (
	"fmt"
	"github.com/ttacon/chalk"
	"strconv"
	"strings"
)

type CliDisplay struct {
	config Config
}

func (d CliDisplay) IndexClean(res BuildOperationResult) {
	fmt.Println(fmt.Sprintf("Index clean took %v seconds", res.TookSeconds))
}

func (d CliDisplay) IndexBuild(res CleanOperationResult) {
	fmt.Println(fmt.Sprintf("Index build took %v ms", res.TookMillis))
}

func (d CliDisplay) Search(res SearchResult) {
	fmt.Println()
	for index, match := range res.Matches {
		icon := "➡️"
		if index == 0 {
			icon = "🎀️"
		}
		position := fmt.Sprintf(" %v  %v:", icon, strconv.Itoa(index+1))
		commit := fmt.Sprintf("Commit: %v", string([]rune(match.File.Commit[0:15])))
		fmt.Println(chalk.Cyan, position, match.File.Path, " - ", commit, chalk.Reset)
		fmt.Println()
		fmt.Println(fmt.Println(strings.Join(match.Fragments, "\n\n...\n\n")))
		fmt.Println()
	}

	if len(res.Matches) < 1 {
		fmt.Println("Nothing found. ")
	}

	fmt.Println(fmt.Sprintf("Search took %v μs", res.TookUs))
}

func (d CliDisplay) StartServer(serviceUrl string) {
	fmt.Println("Listenning on " + serviceUrl)
}
