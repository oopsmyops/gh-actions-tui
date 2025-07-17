package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("version should not be empty")
	}
}

func TestCommit(t *testing.T) {
	if commit == "" {
		t.Error("commit should not be empty")
	}
}

func TestDate(t *testing.T) {
	if date == "" {
		t.Error("date should not be empty")
	}
}