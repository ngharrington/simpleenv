//go:build integration
// +build integration

package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	testSpaceName = os.Getenv("DO_SPACE_NAME")
	testRegion    = os.Getenv("DO_SPACE_REGION")
	accessKey     = os.Getenv("DO_ACCESS_KEY")
	secretKey     = os.Getenv("DO_SECRET_KEY")
	testID        = "integration-test-env"
)

func TestDOSpaceStorage_Integration(t *testing.T) {
	t.Setenv("DO_SPACE_NAME", testSpaceName)
	t.Setenv("DO_SPACE_REGION", testRegion)
	t.Setenv("DO_ACCESS_KEY", accessKey)
	t.Setenv("DO_SECRET_KEY", secretKey)

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	storage, err := NewDOSpaceStorage(config.SpaceName, config.Region, config.AccessKey, config.SecretKey)
	if err != nil {
		t.Fatalf("failed to initialize DOSpaceStorage: %v", err)
	}

	envVars := &EnvironmentVariables{
		ID: "test-env-id",
		vars: map[string]string{
			"VAR1": "value1",
			"VAR2": "value2",
			"VAR3": "value3",
		},
	}

	err = storage.Write(envVars)
	if err != nil {
		t.Fatalf("failed to write environment variables to space: %v", err)
	}

	_, err = storage.Read(envVars.ID)
	if err != nil {
		t.Fatalf("failed to read environment variables from space: %v", err)
	}

	_, err = storage.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(config.SpaceName),
		Key:    aws.String(envVars.ID),
	})
	if err != nil {
		t.Fatalf("failed to delete test object from space: %v", err)
	}
}
