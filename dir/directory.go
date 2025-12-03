package dir

import (
	"os"
	"path/filepath"

	"github.com/candango/iook/file"
)

// CopyAll recursively copies the contents of the src directory into the dst directory.
//
// It preserves file and directory permissions. Existing files or directories in the destination
// with the same name will be overwritten. Symbolic links and special files are not handled.
//
// Parameters:
//
//	src - Path to the source directory.
//	dst - Path to the destination directory.
//
// Returns an error if any operation fails.
//
// Example usage:
//
//	err := dir.CopyAll("source", "destination")
//	if err != nil {
//	    log.Fatal(err)
//	}
func CopyAll(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo,
		err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return file.Copy(path, targetPath)
	})
}

// Exists reports whether the specified path exists and is a directory.
//
// It returns true if the path exists and is a directory, and false otherwise.
// This function does not follow symlinks; it checks the path as provided.
//
// Example usage:
//
//	if dir.Exists("/tmp/mydir") {
//	    // Directory exists
//	}
func Exists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
