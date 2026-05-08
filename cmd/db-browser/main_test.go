package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestServeRejectsFlagLookingStringValue(t *testing.T) {
	cmd := newServeCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--db", "--scripts-dir"})
	cmd.SetContext(context.Background())

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "usually means the intended value was empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServeRejectsUnexpectedPositionalArgs(t *testing.T) {
	cmd := newServeCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"unexpected"})
	cmd.SetContext(context.Background())

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown command") && !strings.Contains(err.Error(), "accepts 0 arg") {
		t.Fatalf("unexpected error: %v", err)
	}
}
