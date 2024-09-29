package main

import (
	"reflect"
	"testing"
)

func TestPlainTextSerializer_Serialize(t *testing.T) {
	serializer := &PlainTextSerializer{}

	envVars := &EnvironmentVariables{
		vars: map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
		},
	}

	expected := "KEY1=value1\nKEY2=value2\n"

	data, err := serializer.Serialize(envVars)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != expected {
		t.Errorf("expected serialized data to be %q, got %q", expected, data)
	}
}

func TestPlainTextSerializer_Deserialize(t *testing.T) {
	serializer := &PlainTextSerializer{}

	data := []byte("KEY1=value1\nKEY2=value2\n")

	expectedVars := &EnvironmentVariables{
		vars: map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
		},
	}

	result, err := serializer.Deserialize(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expectedVars) {
		t.Errorf("expected deserialized data to be %+v, got %+v", expectedVars, result)
	}
}

func TestPlainTextSerializer_Deserialize_InvalidData(t *testing.T) {
	serializer := &PlainTextSerializer{}

	// Test with invalid data: missing "="
	data := []byte("KEY1value1\nKEY2=value2\n")

	_, err := serializer.Deserialize(data)
	if err == nil {
		t.Fatal("expected an error for invalid data, got none")
	}

	expectedErr := "invalid line in document: KEY1value1"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestPlainTextSerializer_SerializeAndDeserialize(t *testing.T) {
	serializer := &PlainTextSerializer{}

	envVars := &EnvironmentVariables{
		vars: map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
			"KEY3": "value3",
		},
	}

	// Serialize the environment variables
	data, err := serializer.Serialize(envVars)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Deserialize the data back to environment variables
	deserializedEnvVars, err := serializer.Deserialize(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Compare the original and deserialized environment variables
	if !reflect.DeepEqual(envVars, deserializedEnvVars) {
		t.Errorf("expected deserialized environment variables to match original, got %+v, want %+v", deserializedEnvVars, envVars)
	}
}
