package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func main() {
	err := applyCommand(os.Args)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func applyCommand(args []string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	commandHandler, err := NewCommandHandler(config)
	if err != nil {
		return err
	}
	defer commandHandler.Destroy()

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
							return commandHandler.CleanIndex()
						},
					},
					{
						Name:  "build",
						Usage: "build nearest index",
						Action: func(c *cli.Context) error {
							return commandHandler.BuildIndex()
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
					return commandHandler.Search(query)
				},
			},
			{
				Name:    "web",
				Aliases: []string{"w"},
				Usage:   "Start web server",
				Action: func(context *cli.Context) error {
					return commandHandler.StartServer()
				},
			},
		},
	}

	return app.Run(args)
}
