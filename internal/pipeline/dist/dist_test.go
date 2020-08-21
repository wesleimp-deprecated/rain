package dist

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wesleimp/rain/pkg/config"
	"github.com/wesleimp/rain/pkg/context"
)

func TestDistDoesNotExist(t *testing.T) {
	folder, err := ioutil.TempDir("", "disttest")
	assert.NoError(t, err)
	var dist = filepath.Join(folder, "dist")
	assert.NoError(
		t,
		Step{}.Run(
			&context.Context{
				Config: config.Config{
					Dist: dist,
				},
			},
		),
	)
}

func TestPopulatedDistExists(t *testing.T) {
	folder, err := ioutil.TempDir("", "disttest")
	assert.NoError(t, err)
	var dist = filepath.Join(folder, "dist")
	assert.NoError(t, os.Mkdir(dist, 0755))
	_, err = os.Create(filepath.Join(dist, "mybin"))
	assert.NoError(t, err)
	var ctx = &context.Context{
		Config: config.Config{
			Dist: dist,
		},
	}
	assert.Error(t, Step{}.Run(ctx))
	ctx.RmDist = true
	assert.NoError(t, Step{}.Run(ctx))
	_, err = os.Stat(dist)
	assert.False(t, os.IsExist(err))
}

func TestEmptyDistExists(t *testing.T) {
	folder, err := ioutil.TempDir("", "disttest")
	assert.NoError(t, err)
	var dist = filepath.Join(folder, "dist")
	assert.NoError(t, os.Mkdir(dist, 0755))
	var ctx = &context.Context{
		Config: config.Config{
			Dist: dist,
		},
	}
	assert.NoError(t, Step{}.Run(ctx))
	_, err = os.Stat(dist)
	assert.False(t, os.IsNotExist(err))
}

func TestDescription(t *testing.T) {
	assert.NotEmpty(t, Step{}.String())
}
