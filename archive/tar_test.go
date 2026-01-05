package archive

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestTar(t *testing.T) {
	t.Run("should tar and untar", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "justify.txt")
		content := []byte("justify this test")
		if err := os.WriteFile(filePath, content, 0600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		// Tar the file
		var buf bytes.Buffer
		if err := Tar(&buf, filePath); err != nil {
			t.Fatalf("Tar: %v", err)
		}

		// Untar to a new directory
		extractDir := t.TempDir()
		if err := Untar(bytes.NewReader(buf.Bytes()), extractDir); err != nil {
			t.Fatalf("Untar: %v", err)
		}

		// Check the extracted file
		extracted := filepath.Join(extractDir, "justify.txt")
		data, err := os.ReadFile(extracted)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if string(data) != string(content) {
			t.Errorf("Extracted content mismatch: got %q, want %q", string(data), string(content))
		}
	})
}
