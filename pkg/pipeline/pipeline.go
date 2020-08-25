package pipeline

import (
	"fmt"

	"github.com/rainproj/rain/internal/pipeline/build"
	"github.com/rainproj/rain/internal/pipeline/defaults"
	"github.com/rainproj/rain/internal/pipeline/dist"
	"github.com/rainproj/rain/internal/pipeline/docker"
	"github.com/rainproj/rain/pkg/context"
)

// Pipeliner interface
type Pipeliner interface {
	fmt.Stringer

	Run(*context.Context) error
}

// BuildPipeline execution
var BuildPipeline = []Pipeliner{
	defaults.Step{},
	dist.Step{},
	build.Step{},
	docker.Step{},
}
