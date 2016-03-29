package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupAndTeardown(t *testing.T) {
	var workingDir, _ = os.Getwd()

	if Context != nil {
		t.Errorf("Expected config.Context to be null")
	}

	if len(Stores) != 0 {
		t.Errorf("Expected config.Stores to be an empty map")
	}

	Setup()

	if len(Stores) != 1 || Stores["global"] == nil {
		t.Errorf("Should have global store")
	}

	if err := os.Chdir(filepath.Join(workingDir, "mocks/project")); err != nil {
		t.Error(err)
	}

	Setup()

	if len(Stores) != 2 || Stores["global"] == nil || Stores["project"] == nil {
		t.Errorf("Should have global and project store")
	}

	if err := os.Chdir(filepath.Join(workingDir, "mocks/project/container")); err != nil {
		t.Error(err)
	}

	Setup()

	if len(Stores) != 3 {
		t.Error("Should have 3 config stores")
	}

	var list = []string{
		"global",
		"project",
		"container",
	}

	for _, k := range list {
		if Stores[k] == nil {
			t.Errorf("%v store missing", k)
		}
	}

	os.Chdir(workingDir)
	Teardown()

	if Context != nil {
		t.Errorf("Expected config.Context to be null")
	}

	if len(Stores) != 0 {
		t.Errorf("Expected config.Stores to be an empty map")
	}
}
