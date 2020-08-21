package middleware

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/fatih/color"
	"github.com/wesleimp/rain/pkg/context"
)

// Action is a function that takes a context and returns an error.
type Action func(ctx *context.Context) error

// Padding is a logging initial padding.
type Padding int

// DefaultInitialPadding is the default padding in the log library.
const DefaultInitialPadding Padding = 3

// ExtraPadding is the double of the DefaultInitialPadding.
const ExtraPadding Padding = DefaultInitialPadding * 2

// Log pretty prints the given action and its title.
func Log(title string, next Action, padding Padding) Action {
	return func(ctx *context.Context) error {
		defer func() {
			cli.Default.Padding = int(DefaultInitialPadding)
		}()

		cli.Default.Padding = int(padding)

		log.Infof(color.New(color.Bold).Sprint(title))
		cli.Default.Padding = int(padding + DefaultInitialPadding)

		return next(ctx)
	}
}

// ErrHandler handles an action error
func ErrHandler(action Action) Action {
	return func(ctx *context.Context) error {
		var err = action(ctx)
		if err == nil {
			return nil
		}

		return err
	}
}
