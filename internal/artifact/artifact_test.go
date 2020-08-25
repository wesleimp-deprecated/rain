package artifact

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestAdd(t *testing.T) {
	var g errgroup.Group
	var artifacts = New()
	for _, a := range []*Artifact{
		{
			Name: "dockerimage",
			Type: DockerImage,
		},
		{
			Name: "anotherimage",
			Type: DockerImage,
		},
	} {
		a := a
		g.Go(func() error {
			artifacts.Add(a)
			return nil
		})
	}
	assert.NoError(t, g.Wait())
	assert.Len(t, artifacts.List(), 2)
}
