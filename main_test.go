package main

import (
	"bytes"
	"testing"
)

func TestOutputVars(t *testing.T) {
	vars := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	writer := bytes.Buffer{}
	outputVars(vars, false, &writer)
	expected := "KEY1=value1\nKEY2=value2\n"
	if writer.String() != expected {
		t.Errorf("expected output to be %q, got %q", expected, writer.String())
	}
}

func TestOutputVars_Source(t *testing.T) {
	vars := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	writer := bytes.Buffer{}
	outputVars(vars, true, &writer)
	expected := "export KEY1=value1\nexport KEY2=value2\n"
	if writer.String() != expected {
		t.Errorf("expected output to be %q, got %q", expected, writer.String())
	}
}
