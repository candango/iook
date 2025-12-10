package pathx

import "os"

// Exists checks whether the specified file or directory exists.
//
// It returns true if the file or directory at the given path exists,
// otherwise returns false.
//
// Example usage:
//
//	exists := pathx.Exists("/etc/passwd")
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
