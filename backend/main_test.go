package main

import (
	"testing"
)

func TestValidateEnv(t *testing.T) {
	// This is a basic test - just ensure it doesn't panic
	validateEnv()
}

func TestGenerateRequestID(t *testing.T) {
	id := generateRequestID()
	if id == "" {
		t.Error("Request ID should not be empty")
	}
}
