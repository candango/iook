package pathx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_exists")
	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	assert.True(t, Exists(tmpPath), "Exists(%q) = false, want true", tmpPath)

	// Should return false for non-existent file
	nonExistent := tmpPath + "-does-not-exist"
	assert.False(t, Exists(nonExistent), "Exists(%q) = false, want true", nonExistent)

	// Should return true for existing directory
	tmpDir, err := os.MkdirTemp("", "test_exists_dir")
	if err != nil {
		t.Fatalf("Unable to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	assert.True(t, Exists(tmpDir), "Exists(%q) = false, want true", tmpDir)
}
