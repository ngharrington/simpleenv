package main

import (
	"fmt"
	"io"
	"os"
	"sort"
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
					&cli.BoolFlag{
						Name:  "source",
						Usage: "Output in 'source' format (e.g. export KEY=VALUE)",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		// we send error messages to stderr because we rely on stdout for
		// e.g. output that can be used by "source" to set env variables in the shell.
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
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
	id := c.String("id")
	sourceFlag := c.Bool("source")

	if spaceName == "" || region == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("DigitalOcean Space credentials must be set via environment variables: DO_SPACE_NAME, DO_SPACE_REGION, DO_ACCESS_KEY, DO_SECRET_KEY")
	}

	storage, err := NewDOSpaceStorage(spaceName, region, accessKey, secretKey)
	if err != nil {
		return fmt.Errorf("failed to initialize DOSpaceStorage: %v", err)
	}

	return executeReadCommand(storage, id, sourceFlag)
}

func executeReadCommand(storage Storage, id string, sourceFlag bool) error {
	envVars, err := storage.Read(id)
	if err != nil {
		return fmt.Errorf("failed to read environment variables from storage: %v", err)
	}
	outputVars(envVars.vars, sourceFlag, os.Stdout)
	return nil
}

func outputVars(vars map[string]string, sourceFlag bool, w io.Writer) {
	orderedVars := make([]string, 0, len(vars))
	for k := range vars {
		orderedVars = append(orderedVars, k)
	}
	sort.Strings(orderedVars)
	for _, key := range orderedVars {
		if sourceFlag {
			fmt.Fprintf(w, "export %s=%s\n", key, vars[key])
		} else {
			fmt.Fprintf(w, "%s=%s\n", key, vars[key])
		}
	}
}
