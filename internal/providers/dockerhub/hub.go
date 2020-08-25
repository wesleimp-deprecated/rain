package dockerhub

import (
	"os/exec"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/rainproj/rain/internal/artifact"
	"github.com/rainproj/rain/pkg/config"
	"github.com/rainproj/rain/pkg/context"
)

// Provider for hub
type Provider struct{}

// Push images to hub
func (Provider) Push(ctx *context.Context, push config.Push) error {
	images := ctx.Artifacts.List()
	for _, image := range images {
		err := pushToHub(ctx, push, image)
		if err != nil {
			return err
		}
	}

	return nil
}

func pushToHub(ctx *context.Context, push config.Push, image *artifact.Artifact) error {
	var env = append(ctx.Env.Strings(), push.Env...)
	log.WithField("image", image.Name).Info("pushing docker image")
	var cmd = exec.CommandContext(ctx, "docker", "push", image.Name)
	cmd.Env = env

	log.WithField("cmd", cmd.Args).Debug("running")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to push docker image: \n%s", string(out))
	}
	return nil
}
