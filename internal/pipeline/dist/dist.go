package dist

import (
	"os"

	"github.com/apex/log"
	"github.com/wesleimp/rain/pkg/context"
)

// Step for dist.
type Step struct{}

func (Step) String() string {
	return "checking ./dist"
}

// Run the pipe.
func (Step) Run(ctx *context.Context) (err error) {
	_, err = os.Stat(ctx.Config.Dist)
	if os.IsNotExist(err) {
		log.Debug("./dist doesn't exist, creating empty folder")
		return mkdir(ctx)
	}
	if ctx.RmDist {
		log.Info("--rm-dist is set, cleaning it up")
		err = os.RemoveAll(ctx.Config.Dist)
		if err == nil {
			err = mkdir(ctx)
		}
		return err
	}

	return nil
}

func mkdir(ctx *context.Context) error {
	// #nosec
	return os.MkdirAll(ctx.Config.Dist, 0755)
}
