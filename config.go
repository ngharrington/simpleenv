package main

import (
	"fmt"
	"os"
)

type Config struct {
	SpaceName string
	Region    string
	AccessKey string
	SecretKey string
}

func LoadConfig() (*Config, error) {
	spaceName := os.Getenv("DO_SPACE_NAME")
	region := os.Getenv("DO_SPACE_REGION")
	accessKey := os.Getenv("DO_ACCESS_KEY")
	secretKey := os.Getenv("DO_SECRET_KEY")

	if spaceName == "" || region == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("DigitalOcean Space credentials must be set via environment variables: DO_SPACE_NAME, DO_SPACE_REGION, DO_ACCESS_KEY, DO_SECRET_KEY")
	}

	return &Config{
		SpaceName: spaceName,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}
