package pipeline

import (
	"fmt"

	"github.com/wesleimp/rain/internal/pipeline/defaults"
	"github.com/wesleimp/rain/internal/pipeline/dist"
	"github.com/wesleimp/rain/internal/pipeline/docker"
	"github.com/wesleimp/rain/pkg/context"
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
	docker.Step{},
}
