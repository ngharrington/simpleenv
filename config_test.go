package main

import (
	"os"
	"testing"
)

func TestLoadConfig_Good(t *testing.T) {
	os.Setenv("DO_SPACE_NAME", "test-space")
	os.Setenv("DO_SPACE_REGION", "test-region")
	os.Setenv("DO_ACCESS_KEY", "test-access-key")
	os.Setenv("DO_SECRET_KEY", "test-secret-key")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.SpaceName != "test-space" {
		t.Errorf("expected space name to be %q, got %q", "test-space", config.SpaceName)
	}

	if config.Region != "test-region" {
		t.Errorf("expected region to be %q, got %q", "test-region", config.Region)
	}

	if config.AccessKey != "test-access-key" {
		t.Errorf("expected access key to be %q, got %q", "test-access-key", config.AccessKey)
	}

	if config.SecretKey != "test-secret-key" {
		t.Errorf("expected secret key to be %q, got %q", "test-secret-key", config.SecretKey)
	}
	os.Clearenv()
}

func TestLoadConfig_MissingSpaceName(t *testing.T) {
	os.Setenv("DO_SPACE_REGION", "test-region")
	os.Setenv("DO_ACCESS_KEY", "test-access-key")
	os.Setenv("DO_SECRET_KEY", "test-secret-key")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected an error for missing space name, got none")
	}
}
