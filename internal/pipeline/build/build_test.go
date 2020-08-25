package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wesleimp/rain/pkg/config"
	"github.com/wesleimp/rain/pkg/context"
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
