package build

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/rainproj/rain/internal/logger"
	"github.com/rainproj/rain/pkg/config"
	"github.com/rainproj/rain/pkg/context"
)

// Step for build
type Step struct{}

func (Step) String() string {
	return "building"
}

// Run the step
func (Step) Run(ctx *context.Context) error {
	for _, build := range ctx.Config.Builds {
		if build.Skip {
			log.WithField("name", build.Name).Info("Skip is set")
			continue
		}

		log.WithField("build", build).Debug("building...")
		err := run(ctx, build)
		if err != nil {
			return err
		}
	}

	return nil
}

func run(ctx *context.Context, build config.Build) error {
	log.Info(color.CyanString("%s", build.Command))

	var command = strings.Split(build.Command, " ")
	var env = append(ctx.Env.Strings(), build.Env...)

	var cmd = exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Env = env

	var entry = log.WithField("cmd", command).WithField("env", env)

	var b bytes.Buffer
	cmd.Stdout = io.MultiWriter(logger.NewWriter(entry), &b)
	cmd.Stderr = io.MultiWriter(logger.NewErrWriter(entry), &b)
	if build.Dir != "" {
		cmd.Dir = build.Dir
	}

	entry.Debug("running")
	if err := cmd.Run(); err != nil {
		entry.WithError(err).Debug("failed")
		return errors.Wrapf(err, "%q", b.String())
	}

	return nil
}
