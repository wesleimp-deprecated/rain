package files

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wesleimp/rain/pkg/config"
)

func TestShouldGetAllFiles(t *testing.T) {
	assert := assert.New(t)

	globs := []config.File{
		{Glob: "./testdata/file1.txt"},
	}

	files, err := Find(globs)
	assert.NoError(err)
	assert.Equal(1, len(files))

	path, ok := files["file1.txt"]
	assert.True(ok)
	assert.Equal(path, "./testdata/file1.txt")
}

func TestShouldGetFiles(t *testing.T) {
	assert := assert.New(t)

	globs := []config.File{
		{Glob: "test.txt"},
	}

	files, err := Find(globs)
	assert.NoError(err)
	assert.Equal(1, len(files))

	path, ok := files["test.txt"]
	assert.True(ok)
	assert.Equal(path, "test.txt")
}

func TestShouldGetAllFilesWithGoldenExtension(t *testing.T) {
	assert := assert.New(t)

	globs := []config.File{
		{Glob: "./testdata/*.txt"},
	}

	files, err := Find(globs)
	assert.NoError(err)
	assert.Equal(2, len(files))

	path, ok := files["file1.txt"]
	assert.True(ok)
	assert.Equal(path, "testdata/file1.txt")

	path, ok = files["file2.txt"]
	assert.True(ok)
	assert.Equal(path, "testdata/file2.txt")
}

func TestShouldGetAllFilesInsideTestdata(t *testing.T) {
	assert := assert.New(t)

	globs := []config.File{
		{Glob: "./testdata/**"},
	}

	files, err := Find(globs)
	assert.NoError(err)
	assert.Equal(3, len(files))

	path, ok := files["file1.txt"]
	assert.True(ok)
	assert.Equal(path, "testdata/file1.txt")

	path, ok = files["file2.txt"]
	assert.True(ok)
	assert.Equal(path, "testdata/file2.txt")

	path, ok = files["file3.text"]
	assert.True(ok)
	assert.Equal(path, "testdata/file3.text")
}
