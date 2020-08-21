package init

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

var (
	binaryName = "rain"
	run        = func(args []string) error {
		app := &cli.App{
			Name:     binaryName,
			Commands: []*cli.Command{Command},
		}
		return app.Run(args)
	}
)

func TestInit(t *testing.T) {

	assert := assert.New(t)

	var folder = mktemp(t)
	var path = filepath.Join(folder, "foo.yaml")

	err := run([]string{binaryName, "init", "-c", path})
	assert.NoError(err)
	assert.FileExists(path)
}

func TestInitFileExists(t *testing.T) {
	assert := assert.New(t)

	var folder = mktemp(t)
	var path = filepath.Join(folder, "foo.yaml")

	err := run([]string{binaryName, "init", "-c", path})
	assert.NoError(err)

	err = run([]string{binaryName, "init", "-c", path})
	assert.EqualError(err, "open "+path+": file exists")
	assert.FileExists(path)
}

func mktemp(t *testing.T) string {
	folder, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	return folder
}
