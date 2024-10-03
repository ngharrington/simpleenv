package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestEnvironmentVariables_Set(t *testing.T) {
	envVars := &EnvironmentVariables{
		ID:   "test",
		vars: make(map[string]string),
	}

	envVars.Set("API_KEY", "12345")
	if value, exists := envVars.vars["API_KEY"]; !exists || value != "12345" {
		t.Errorf("Set failed: expected value '12345', got '%v'", value)
	}

	envVars.Set("DB_HOST", "localhost")
	if value, exists := envVars.vars["DB_HOST"]; !exists || value != "localhost" {
		t.Errorf("Set failed: expected value 'localhost', got '%v'", value)
	}

	envVars.Set("API_KEY", "67890")
	if value, exists := envVars.vars["API_KEY"]; !exists || value != "67890" {
		t.Errorf("Set failed: expected updated value '67890', got '%v'", value)
	}
}

func TestEnvironmentVariables_Get(t *testing.T) {
	envVars := &EnvironmentVariables{
		ID: "test",
		vars: map[string]string{
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		},
	}

	value, err := envVars.Get("DB_HOST")
	if err != nil {
		t.Errorf("Get failed: expected no error, got '%v'", err)
	}
	if value != "localhost" {
		t.Errorf("Get failed: expected 'localhost', got '%v'", value)
	}

	_, err = envVars.Get("NON_EXISTENT")
	if err == nil {
		t.Errorf("Get failed: expected error for non-existent key, got nil")
	}
	if err != errVarNotExists {
		t.Errorf("Get failed: expected '%v', got '%v'", errVarNotExists, err)
	}
}

func TestEnvironmentVariables_ListVariables(t *testing.T) {
	envVars := &EnvironmentVariables{
		ID: "test",
		vars: map[string]string{
			"API_KEY": "12345",
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		},
	}
	keys := envVars.ListVariables()
	expectedKeys := []string{"API_KEY", "DB_HOST", "DB_PORT"}

	sort.Strings(keys)
	sort.Strings(expectedKeys)

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("ListVariables failed: expected '%v', got '%v'", expectedKeys, keys)
	}
}

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

func TestMockInMemoryStorage_WriteAndRead(t *testing.T) {
	storage := NewMockInMemoryStorage()

	envVars := &EnvironmentVariables{
		ID: "test-env",
		vars: map[string]string{
			"API_KEY": "12345",
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		},
	}

	err := storage.Write(envVars)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	readVars, err := storage.Read("test-env")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !reflect.DeepEqual(readVars.vars, envVars.vars) {
		t.Errorf("Mismatch between written and read variables: expected %v, got %v", envVars.vars, readVars.vars)
	}
}

func TestMockInMemoryStorage_ReadNonExistent(t *testing.T) {
	storage := NewMockInMemoryStorage()

	_, err := storage.Read("non-existent-id")
	if err == nil {
		t.Fatalf("Expected error when reading non-existent environment variables, got nil")
	}
}

func TestPlainTextSerializerOrdering(t *testing.T) {
	serializer := &PlainTextSerializer{}

	envVars := &EnvironmentVariables{
		ID: "test-env",
		vars: map[string]string{
			"DB_HOST": "localhost",
			"API_KEY": "12345",
			"DB_PORT": "5432",
		},
	}

	serializedData, err := serializer.Serialize(envVars)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	expectedSerializedData := "API_KEY=12345\nDB_HOST=localhost\nDB_PORT=5432\n"

	if string(serializedData) != expectedSerializedData {
		t.Errorf("Serialized data ordering mismatch: expected %s, got %s", expectedSerializedData, serializedData)
	}
}
