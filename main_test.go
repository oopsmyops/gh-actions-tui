package main

import (
	"testing"
)

func TestVersionDefaults(t *testing.T) {
	// During testing, version variables have default values
	if version != "dev" && version == "" {
		t.Error("version should have a default value")
	}
	if commit != "none" && commit == "" {
		t.Error("commit should have a default value")
	}
	if date != "unknown" && date == "" {
		t.Error("date should have a default value")
	}
}

func TestMainFunctionExists(t *testing.T) {
	// This test ensures the package compiles correctly
	// We test that the version variables are properly initialized
	t.Logf("Version: %s, Commit: %s, Date: %s", version, commit, date)
}