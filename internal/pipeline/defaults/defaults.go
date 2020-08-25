package defaults

import (
	"errors"

	"github.com/rainproj/rain/pkg/context"
	"github.com/rainproj/rain/pkg/defaults"
)

// Step for defaults
type Step struct{}

func (Step) String() string {
	return "setting defaults"
}

// Run defaults step
func (Step) Run(ctx *context.Context) error {
	if ctx.Config.Dist == "" {
		ctx.Config.Dist = "dist"
	}

	for _, def := range defaults.Defaulters {
		err := def.Default(ctx)
		if err != nil {
			return errors.New("error setting defaults")
		}
	}

	return nil
}
