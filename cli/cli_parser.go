package cli

import (
	"github.com/urfave/cli/v2"
	"gitlab.com/remipassmoilesel/gitsearch/config"
	"gitlab.com/remipassmoilesel/gitsearch/http"
	"gitlab.com/remipassmoilesel/gitsearch/index"
)

type CliParserImpl struct {
	cliHandlers CliHandlers
}

func NewCliParser(config config.Config, index index.Index, server http.HttpServer) CliParserImpl {
	handlers := NewCliHandlers(config, index, server)
	return CliParserImpl{cliHandlers: handlers}
}

const (
	FlagNumberOfResult = "number-of-results"
	FlagQuery          = "query"
	FlagNoPager        = "no-pager"
	FlagHash           = "hash"
)

func (s *CliParserImpl) ApplyCommand(args []string) error {
	app := &cli.App{
		Name:                 "gitsearch",
		Usage:                "🔎 Search in Git repositories ⚡",
		HelpName:             "gitsearch",
		Description:          "Git Search indexes versioned files in Git repositories and allows you to search using the command line or a web interface.",
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{
				Name:  "remipassmoilesel",
				Email: "r.passmoilesel@protonmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  FlagNoPager,
				Usage: "do not use pager to display output",
				Value: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "Search command",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     FlagQuery,
						Aliases:  []string{"q"},
						Usage:    "You can search term, search \"exact terms\", include +terms, exclude -terms etc ... See: https://blevesearch.com/docs/Query-String-Query/",
						Required: true,
					},
					&cli.IntFlag{
						Name:    FlagNumberOfResult,
						Aliases: []string{"n"},
						Usage:   "will return at most number of results",
						Value:   10,
					},
				},
				Action: func(context *cli.Context) error {
					query := context.String(FlagQuery)
					usePager := !context.Bool(FlagNoPager)
					numberOfResults := 50
					if context.Int(FlagNumberOfResult) > 0 {
						numberOfResults = context.Int(FlagNumberOfResult)
					}

					return s.cliHandlers.Search(query, numberOfResults, usePager)
				},
			},
			{
				Name:    "show-file",
				Aliases: []string{"f"},
				Usage:   "Show file with specified partial hash",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     FlagHash,
						Aliases:  []string{"ha"},
						Usage:    "Hash or partial hash of file",
						Required: true,
					},
				},
				Action: func(context *cli.Context) error {
					hash := context.String(FlagHash)
					usePager := !context.Bool(FlagNoPager)

					return s.cliHandlers.ShowFile(hash, usePager)
				},
			},
			{
				Name:    "web-ui",
				Aliases: []string{"w"},
				Usage:   "Start web user interface",
				Action: func(context *cli.Context) error {
					return s.cliHandlers.StartServer()
				},
			},
			{
				Name:    "index",
				Aliases: []string{"i"},
				Usage:   "Index commands",
				Subcommands: []*cli.Command{
					{
						Name:    "update",
						Aliases: []string{"u"},
						Usage:   "scan files from current git repository then index them",
						Action: func(c *cli.Context) error {
							return s.cliHandlers.UpdateIndex()
						},
					},
					{
						Name:    "clean",
						Aliases: []string{"c"},
						Usage:   "delete data from current git repository index",
						Action: func(c *cli.Context) error {
							return s.cliHandlers.CleanIndex()
						},
					},
				},
			},
		},
	}

	return app.Run(args)
}
