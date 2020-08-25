package build

import (
	"github.com/wesleimp/rain/pkg/context"
)

// Step for build
type Step struct{}

func (Step) String() string {
	return "building"
}

// Run build section
func (Step) Run(ctx *context.Context) error {
	return nil
}