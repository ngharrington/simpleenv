package main

import (
	"fmt"
	"os"
	"strings"

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
			{
				Name:   "read",
				Usage:  "Read environment variables from DigitalOcean Space",
				Action: readCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "ID for the environment variables",
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

func writeCommand(c *cli.Context) error {
	spaceName := os.Getenv("DO_SPACE_NAME")
	region := os.Getenv("DO_SPACE_REGION")
	accessKey := os.Getenv("DO_ACCESS_KEY")
	secretKey := os.Getenv("DO_SECRET_KEY")

	if spaceName == "" || region == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("DigitalOcean Space credentials must be set via environment variables: DO_SPACE_NAME, DO_SPACE_REGION, DO_ACCESS_KEY, DO_SECRET_KEY")
	}
	storage, err := NewDOSpaceStorage(spaceName, region, accessKey, secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize DOSpaceStorage: %v", err)
	}
	envVars := &EnvironmentVariables{
		ID:   c.String("id"),
		vars: make(map[string]string),
	}
	for _, v := range c.StringSlice("vars") {
		parts := strings.Split(v, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid environment variable: %s", v)
		}
		envVars.vars[parts[0]] = parts[1]
	}
	err = storage.Write(envVars)
	if err != nil {
		return fmt.Errorf("failed to write environment variables to space: %v", err)
	}
	fmt.Println("Environment variables written successfully")
	return nil
}

func readCommand(c *cli.Context) error {
	spaceName := os.Getenv("DO_SPACE_NAME")
	region := os.Getenv("DO_SPACE_REGION")
	accessKey := os.Getenv("DO_ACCESS_KEY")
	secretKey := os.Getenv("DO_SECRET_KEY")

	if spaceName == "" || region == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("DigitalOcean Space credentials must be set via environment variables: DO_SPACE_NAME, DO_SPACE_REGION, DO_ACCESS_KEY, DO_SECRET_KEY")
	}
	storage, err := NewDOSpaceStorage(spaceName, region, accessKey, secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize DOSpaceStorage: %v", err)
	}
	envVars, err := storage.Read(c.String("id"))
	if err != nil {
		return fmt.Errorf("failed to read environment variables from space: %v", err)
	}
	fmt.Println("Environment variables read successfully")
	for k, v := range envVars.vars {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
