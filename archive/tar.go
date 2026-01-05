package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Tar creates a tar archive from the given source file or directory and writes
// it to w. It preserves file structure and metadata.
func Tar(w io.Writer, src string) error {
	src = filepath.Clean(src)
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(w)
	defer tw.Close()

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)
	}

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		relPath := path
		if baseDir != "" {
			relPath, err = filepath.Rel(filepath.Dir(src), path)
			if err != nil {
				return err
			}
		} else {
			relPath = filepath.Base(path)
		}
		hdr.Name = relPath
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	}

	if info.IsDir() {
		return filepath.Walk(src, walkFn)
	}
	return walkFn(src, info, nil)
}

// Untar extracts a tar archive from r into the destination directory dest.
// It performs security checks to prevent path traversal and preserves file
// structure and metadata.
func Untar(r io.Reader, dest string) error {
	destDir := filepath.Clean(dest)
	f, err := os.Open(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Sanitize path to prevent path traversal (security check)
		name := filepath.Clean(header.Name)
		if strings.Contains(name, "..") || filepath.IsAbs(name) {
			return fmt.Errorf(
				"security error: tar entry %q contains an invalid or unsafe "+
					"path (possible path traversal attempt), extraction "+
					"aborted", name)
		}
		target := filepath.Join(destDir, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
			// Set timestamps
			if err := os.Chtimes(target, header.AccessTime, header.ModTime); err != nil {
				// Not fatal, continue
			}
			// TODO: Set file ownership (UID/GID) if needed and running as root
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		case tar.TypeLink:
			linkTarget := filepath.Join(destDir, header.Linkname)
			if err := os.Link(linkTarget, target); err != nil {
				return err
			}
		}
	}
	return nil
}
