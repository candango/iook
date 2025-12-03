package file

import (
	"io"
	"os"
)

type CopyOptions struct {
	Force bool
	Group int
	Mode  *os.FileMode
	User  int
}

type CopyOption func(*CopyOptions)

func WithForce(f bool) CopyOption {
	return func(o *CopyOptions) { o.Force = f }
}

func WithGroup(g int) CopyOption {
	return func(o *CopyOptions) { o.Group = g }
}

func WithFileMode(m os.FileMode) CopyOption {
	return func(o *CopyOptions) { o.Mode = &m }
}

// Copy copies a regular file from src to dst, preserving its content and
// setting the given permissions. If dst does not exist, it will be created. If
// it exists, it will be overwritten. Returns an error if any file operation
// fails.
//
// Example usage:
//
//	opts := file.WithMode(0654)
//	err := file.Copy("foo.txt", "bar.txt", opts)
//	if err != nil {
//	    log.Fatal(err)
//	}
func Copy(src, dst string, opts ...CopyOption) error {
	// Default options
	options := &CopyOptions{
		Mode:  nil,
		Force: false,
	}
	for _, opt := range opts {
		opt(options)
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return &os.PathError{Op: "Copy", Path: src, Err: os.ErrInvalid}
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	fi, err := in.Stat()
	if err != nil {
		return err
	}
	fileMode := fi.Mode()
	if options.Mode != nil {
		fileMode = *options.Mode
	}

	return os.Chmod(dst, fileMode)
}
