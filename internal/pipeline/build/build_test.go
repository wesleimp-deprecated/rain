package build

import (
	"testing"

	"github.com/rainproj/rain/pkg/config"
	"github.com/rainproj/rain/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	assert.NotEmpty(t, Step{}.String())
}

func TestBuild(t *testing.T) {
	ctx := context.New(config.Config{
		Builds: []config.Build{{
			Name:    "my build",
			Command: "touch testdata/file.txt",
		}},
	})

	assert.NoError(t, Step{}.Run(ctx))
	assert.FileExists(t, "testdata/file.txt")
}
