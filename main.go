package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "write",
				Usage:  "Write environment variables to DigitalOcean Space",
				Action: writeCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "ID for the environment variables",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "vars",
						Usage:    "Environment variables in KEY=VALUE format",
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
