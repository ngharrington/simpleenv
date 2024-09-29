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

// Integration test example
func TestDOSpaceStorage_Integration(t *testing.T) {
	// Load configuration from environment variables
	spaceName := os.Getenv("DO_SPACE_NAME")
	region := os.Getenv("DO_SPACE_REGION")
	accessKey := os.Getenv("DO_ACCESS_KEY")
	secretKey := os.Getenv("DO_SECRET_KEY")

	if spaceName == "" || region == "" || accessKey == "" || secretKey == "" {
		t.Fatal("DigitalOcean Space credentials must be set via environment variables: DO_SPACE_NAME, DO_SPACE_REGION, DO_ACCESS_KEY, DO_SECRET_KEY")
	}

	// Initialize DigitalOcean Space storage
	storage, err := NewDOSpaceStorage(spaceName, region, accessKey, secretKey)
	if err != nil {
		t.Fatalf("failed to initialize DOSpaceStorage: %v", err)
	}

	// Define environment variables to test with
	envVars := &EnvironmentVariables{
		ID: "test-env-id",
		vars: map[string]string{
			"VAR1": "value1",
			"VAR2": "value2",
			"VAR3": "value3",
		},
	}

	// Write to DigitalOcean Space
	err = storage.Write(envVars)
	if err != nil {
		t.Fatalf("failed to write environment variables to space: %v", err)
	}

	// Read from DigitalOcean Space
	_, err = storage.Read(envVars.ID)
	if err != nil {
		t.Fatalf("failed to read environment variables from space: %v", err)
	}

	// Clean up: Delete the test object from the space
	_, err = storage.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(spaceName),
		Key:    aws.String(envVars.ID),
	})
	if err != nil {
		t.Fatalf("failed to delete test object from space: %v", err)
	}
}
