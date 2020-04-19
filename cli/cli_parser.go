package cli

import (
	"github.com/remipassmoilesel/gitsearch/config"
	"github.com/remipassmoilesel/gitsearch/http"
	"github.com/remipassmoilesel/gitsearch/index"
	"github.com/urfave/cli/v2"
	"strings"
)

type CliParser struct {
	cliHandlers CliHandlers
}

func NewCliParser(config config.Config, index index.Index, server http.HttpServer) CliParser {
	handlers := NewCliHandlers(config, index, server)
	return CliParser{cliHandlers: handlers}
}

func (s *CliParser) ApplyCommand(args []string) error {
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
		Commands: []*cli.Command{
			{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "Search command",
				Action: func(context *cli.Context) error {
					query := strings.Join(context.Args().Slice(), " ")
					return s.cliHandlers.Search(query)
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
						Name:    "build",
						Aliases: []string{"b"},
						Usage:   "scan files from current git repository then index them",
						Action: func(c *cli.Context) error {
							return s.cliHandlers.BuildIndex()
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
