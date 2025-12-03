package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	// Create a temp source file
	src, err := os.CreateTemp("", "srcfile")
	assert.NoError(t, err)
	defer os.Remove(src.Name())

	data := []byte("hello world")
	_, err = src.Write(data)
	assert.NoError(t, err)
	src.Close()

	// Create a temp destination file path
	dst, err := os.CreateTemp("", "dstfile")
	assert.NoError(t, err)
	dstPath := dst.Name()
	dst.Close()
	os.Remove(dstPath)

	t.Run("Should copy and keep origin file mode", func(t *testing.T) {
		defer os.Remove(dstPath)
		err = Copy(src.Name(), dstPath)
		assert.NoError(t, err)

		copied, err := os.ReadFile(dstPath)
		assert.NoError(t, err)
		assert.Equal(t, data, copied, "Copied file contents should match source")

		srcInfo, _ := os.Stat(src.Name())
		dstInfo, _ := os.Stat(dstPath)
		assert.Equal(t, srcInfo.Mode(), dstInfo.Mode(), "File modes should match")
	})

	t.Run("Should copy and change file to custom mode", func(t *testing.T) {
		defer os.Remove(dstPath)
		customMode := os.FileMode(0654)
		opts := WithFileMode(customMode)
		err = Copy(src.Name(), dstPath, opts)
		assert.NoError(t, err)

		copied, err := os.ReadFile(dstPath)
		assert.NoError(t, err)
		assert.Equal(t, data, copied, "Copied file contents should match source")

		dstInfo, _ := os.Stat(dstPath)
		assert.Equal(t, dstInfo.Mode(), customMode, "File modes should match")
	})
}
