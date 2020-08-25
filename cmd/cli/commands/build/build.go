package build

import (
	"time"

	"github.com/apex/log"
	"github.com/caarlos0/ctrlc"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/rainproj/rain/cmd/cli/config"
	"github.com/rainproj/rain/internal/middleware"
	"github.com/rainproj/rain/pkg/context"
	"github.com/rainproj/rain/pkg/pipeline"
	"github.com/urfave/cli/v2"
)

// Command build
var Command = &cli.Command{
	Name:    "build",
	Usage:   "Builds the current project",
	Aliases: []string{"b"},
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: ".rain.yml", Usage: "Load configuration from file"},
		&cli.IntFlag{Name: "parallelism", Value: 2, Usage: "number of tasks running concurrently"},
		&cli.DurationFlag{Name: "timeout", Value: 10 * time.Minute, Usage: "timeout to the entire push process", DefaultText: "10 minutes"},
		&cli.BoolFlag{Name: "rm-dist", Value: false, Usage: "remove the dist folder before build"},
	},
	Action: func(c *cli.Context) error {
		start := time.Now()
		log.Infof(color.New(color.Bold).Sprint("building..."))

		conf := c.String("config")
		cfg, err := config.Load(conf)
		if err != nil {
			return err
		}
		ctx, cancel := context.NewWithTimeout(cfg, c.Duration("timeout"))
		defer cancel()
		setupPushContext(ctx, c)
		err = ctrlc.Default.Run(ctx, func() error {
			for _, pipe := range pipeline.BuildPipeline {
				if err := middleware.Log(
					pipe.String(),
					middleware.ErrHandler(pipe.Run),
					middleware.DefaultInitialPadding,
				)(ctx); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, color.New(color.Bold).Sprintf("build failed after %0.2fs", time.Since(start).Seconds()))
		}

		log.Infof(color.New(color.Bold).Sprintf("build succeeded after %0.2fs", time.Since(start).Seconds()))
		return nil
	},
}

func setupPushContext(ctx *context.Context, c *cli.Context) *context.Context {
	ctx.Parallelism = c.Int("parallelism")
	log.Debugf("parallelism: %v", ctx.Parallelism)
	ctx.RmDist = c.Bool("rm-dist")
	ctx.Timeout = c.Duration("timeout")

	return ctx
}
