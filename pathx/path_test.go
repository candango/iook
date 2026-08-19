package pathx

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "file")
	require.NoError(t, os.WriteFile(filePath, []byte("data"), 0o600))

	dirPath := filepath.Join(root, "directory")
	require.NoError(t, os.Mkdir(dirPath, 0o700))

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "regular file",
			path: filePath,
			want: true,
		},
		{
			name: "directory",
			path: dirPath,
			want: true,
		},
		{
			name: "missing path",
			path: filepath.Join(root, "missing"),
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exists, err := Exists(test.path)

			require.NoError(t, err)
			assert.Equal(t, test.want, exists)
		})
	}
}

func TestExistsFollowsSymlinks(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	require.NoError(t, os.WriteFile(target, []byte("data"), 0o600))

	link := filepath.Join(root, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symbolic links are unavailable: %v", err)
	}

	dangling := filepath.Join(root, "dangling")
	require.NoError(t, os.Symlink(filepath.Join(root, "missing"), dangling))

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "link to existing target",
			path: link,
			want: true,
		},
		{
			name: "dangling link",
			path: dangling,
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exists, err := Exists(test.path)

			require.NoError(t, err)
			assert.Equal(t, test.want, exists)
		})
	}
}

func TestLstat(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "file")
	require.NoError(t, os.WriteFile(filePath, []byte("data"), 0o600))

	dirPath := filepath.Join(root, "directory")
	require.NoError(t, os.Mkdir(dirPath, 0o700))

	t.Run("regular file", func(t *testing.T) {
		info, exists, err := Lstat(filePath)

		require.NoError(t, err)
		require.True(t, exists)
		require.NotNil(t, info)
		assert.True(t, info.Mode().IsRegular())
	})

	t.Run("directory", func(t *testing.T) {
		info, exists, err := Lstat(dirPath)

		require.NoError(t, err)
		require.True(t, exists)
		require.NotNil(t, info)
		assert.True(t, info.IsDir())
	})

	t.Run("missing path", func(t *testing.T) {
		info, exists, err := Lstat(filepath.Join(root, "missing"))

		require.NoError(t, err)
		assert.False(t, exists)
		assert.Nil(t, info)
	})
}

func TestLstatDoesNotFollowSymlinks(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	require.NoError(t, os.WriteFile(target, []byte("data"), 0o600))

	link := filepath.Join(root, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symbolic links are unavailable: %v", err)
	}

	dangling := filepath.Join(root, "dangling")
	require.NoError(t, os.Symlink(filepath.Join(root, "missing"), dangling))

	for _, path := range []string{link, dangling} {
		t.Run(filepath.Base(path), func(t *testing.T) {
			info, exists, err := Lstat(path)

			require.NoError(t, err)
			require.True(t, exists)
			require.NotNil(t, info)
			assert.NotZero(t, info.Mode()&fs.ModeSymlink)
		})
	}
}

func TestLookupErrorsArePreserved(t *testing.T) {
	path := "invalid\x00path"

	t.Run("Exists", func(t *testing.T) {
		_, statErr := os.Stat(path)
		if statErr == nil || errors.Is(statErr, fs.ErrNotExist) {
			t.Skip("platform does not expose a stable non-ErrNotExist stat error")
		}

		exists, err := Exists(path)

		assert.False(t, exists)
		require.Error(t, err)
		assert.False(t, errors.Is(err, fs.ErrNotExist))
		var pathErr *os.PathError
		assert.ErrorAs(t, err, &pathErr)
	})

	t.Run("Lstat", func(t *testing.T) {
		_, lstatErr := os.Lstat(path)
		if lstatErr == nil || errors.Is(lstatErr, fs.ErrNotExist) {
			t.Skip("platform does not expose a stable non-ErrNotExist lstat error")
		}

		info, exists, err := Lstat(path)

		assert.Nil(t, info)
		assert.False(t, exists)
		require.Error(t, err)
		assert.False(t, errors.Is(err, fs.ErrNotExist))
		var pathErr *os.PathError
		assert.ErrorAs(t, err, &pathErr)
	})
}
