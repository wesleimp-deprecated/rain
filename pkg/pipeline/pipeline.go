package pipeline

import (
	"fmt"

	"github.com/wesleimp/rain/pkg/context"
)

// Pipeliner interface
type Pipeliner interface {
	fmt.Stringer

	Run(*context.Context) error
}

// PackPipeline execution
var PackPipeline = []Pipeliner{}
