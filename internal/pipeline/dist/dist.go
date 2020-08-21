package dist

import (
	"fmt"
	"io/ioutil"
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
	files, err := ioutil.ReadDir(ctx.Config.Dist)
	if err != nil {
		return
	}
	if len(files) != 0 {
		log.Debugf("there are %d files on ./dist", len(files))
		return fmt.Errorf(
			"%s is not empty, remove it before running dim or use the --rm-dist flag",
			ctx.Config.Dist,
		)
	}
	log.Debug("./dist is empty")
	return mkdir(ctx)
}

func mkdir(ctx *context.Context) error {
	// #nosec
	return os.MkdirAll(ctx.Config.Dist, 0755)
}
