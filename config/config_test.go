package config

import (
	"os"
	"path"
	"testing"
)

func TestSetup(t *testing.T) {
	var workingDir, _ = os.Getwd()

	Setup()

	if len(Stores) != 1 || Stores["global"] == nil {
		t.Errorf("Should have global store")
	}

	if err := os.Chdir(path.Join(workingDir, "mocks/project")); err != nil {
		t.Error(err)
	}

	Setup()

	if len(Stores) != 2 || Stores["global"] == nil || Stores["project"] == nil {
		t.Errorf("Should have global and project store")
	}

	if err := os.Chdir(path.Join(workingDir, "mocks/project/container")); err != nil {
		t.Error(err)
	}

	Setup()

	if len(Stores) != 3 || Stores["global"] == nil || Stores["project"] == nil || Stores["container"] == nil {
		t.Errorf("Should have global, project, and container store")
	}

	os.Chdir(workingDir)
}