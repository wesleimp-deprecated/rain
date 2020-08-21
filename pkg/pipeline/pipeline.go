package pipeline

import (
	"fmt"

	"github.com/wesleimp/rain/internal/pipeline/defaults"
	"github.com/wesleimp/rain/internal/pipeline/docker"
	"github.com/wesleimp/rain/pkg/context"
)

// Pipeliner interface
type Pipeliner interface {
	fmt.Stringer

	Run(*context.Context) error
}

// PackPipeline execution
var BuildPipeline = []Pipeliner{
	defaults.Step{},
	docker.Step{},
}
