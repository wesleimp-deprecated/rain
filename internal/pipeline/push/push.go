package push

import (
	"github.com/rainproj/rain/internal/semerrgroup"
	"github.com/rainproj/rain/pkg/context"
	"github.com/rainproj/rain/pkg/provider"
)

// Step for push
type Step struct{}

func (Step) String() string {
	return "pushing images"
}

// Run push
func (Step) Run(ctx *context.Context) error {
	var g = semerrgroup.New(ctx.Parallelism)
	for _, push := range ctx.Config.Pushes {
		push := push
		g.Go(func() error {
			p, err := provider.Get(push.Provider)
			if err != nil {
				return err
			}

			return p.Push(ctx, push)
		})
	}

	return g.Wait()
}
