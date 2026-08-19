package pathx

import (
	"errors"
	"io/fs"
	"os"
)

// Exists reports whether path resolves to an existing filesystem object.
// It follows symbolic links in the same way as [os.Stat]. A dangling symbolic
// link is therefore reported as missing.
//
// Exists returns false with a nil error when the lookup fails with
// [fs.ErrNotExist]. It preserves every other lookup error.
//
// The result describes only this lookup. It does not protect a later
// filesystem operation from changes to the path.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Lstat reports whether path names an existing directory entry and returns
// its metadata. It does not follow the final symbolic link, matching
// [os.Lstat]. An existing symbolic link, including a dangling one, is therefore
// reported as existing and its link metadata is returned.
//
// Lstat returns nil, false, nil when the lookup fails with [fs.ErrNotExist]. It
// preserves every other lookup error.
//
// The returned metadata belongs to this lookup. It avoids a second metadata
// query but does not make a later filesystem operation safe from path changes.
func Lstat(path string) (fs.FileInfo, bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return info, true, nil
}
