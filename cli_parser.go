package main

import (
	"github.com/urfave/cli/v2"
	"strings"
)

type CliParser struct {
	cliHandlers CliHandlers
}

func NewCliParser(config Config, index Index, server HttpServer) CliParser {
	handlers := NewCliHandlers(config, index, server)
	return CliParser{cliHandlers: handlers}
}

func (s *CliParser) applyCommand(args []string) error {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "index",
				Aliases: []string{"i"},
				Usage:   "Index commands",
				Subcommands: []*cli.Command{
					{
						Name:  "clean",
						Usage: "clean nearest index data",
						Action: func(c *cli.Context) error {
							return s.cliHandlers.CleanIndex()
						},
					},
					{
						Name:  "build",
						Usage: "build nearest index",
						Action: func(c *cli.Context) error {
							return s.cliHandlers.BuildIndex()
						},
					},
				},
			},
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
				Usage:   "Start web server",
				Action: func(context *cli.Context) error {
					return s.cliHandlers.StartServer()
				},
			},
		},
	}

	return app.Run(args)
}
