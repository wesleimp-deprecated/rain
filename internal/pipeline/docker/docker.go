package docker

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/rainproj/rain/internal/artifact"
	"github.com/rainproj/rain/internal/files"
	"github.com/rainproj/rain/internal/semerrgroup"
	"github.com/rainproj/rain/internal/tmpl"
	"github.com/rainproj/rain/pkg/config"
	"github.com/rainproj/rain/pkg/context"
)

// Step for dockers
type Step struct{}

func (Step) String() string {
	return "building dockers"
}

// Run docker step
func (Step) Run(ctx *context.Context) error {
	if len(ctx.Config.Dockers) == 0 || len(ctx.Config.Dockers[0].ImageTemplates) == 0 {
		return errors.New("docker section should be configured")
	}
	_, err := exec.LookPath("docker")
	if err != nil {
		return errors.New("docker not present in $PATH")
	}

	return run(ctx)
}

func run(ctx *context.Context) error {
	var g = semerrgroup.New(ctx.Parallelism)
	for _, docker := range ctx.Config.Dockers {
		docker := docker
		g.Go(func() error {
			tmp, err := ioutil.TempDir(ctx.Config.Dist, "raindocker-")
			if err != nil {
				return errors.Wrap(err, "failed to create temporary dir")
			}
			log.Debug("tempdir: " + tmp)

			images, err := parseTemplates(ctx, docker.ImageTemplates)
			if err != nil {
				return errors.Wrap(err, "error parsing image templates")
			}

			flags, err := parseTemplates(ctx, docker.BuildFlagTemplates)
			if err != nil {
				return errors.Wrap(err, "error parsing flag templates")
			}

			if err := linkFiles(ctx, tmp, docker); err != nil {
				return err
			}

			if err := build(ctx, tmp, images, flags); err != nil {
				return err
			}

			for _, img := range images {
				ctx.Artifacts.Add(&artifact.Artifact{
					Type: artifact.DockerImage,
					Name: img,
				})
			}

			return nil
		})
	}

	return g.Wait()
}

func parseTemplates(ctx *context.Context, templates []string) ([]string, error) {
	var out []string
	for _, template := range templates {
		output, err := tmpl.New(ctx).Apply(template)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to execute template %s", template)
		}

		out = append(out, output)
	}

	return out, nil
}

func linkFiles(ctx *context.Context, dir string, docker config.Docker) error {
	if err := os.Link(docker.Dockerfile, filepath.Join(dir, "Dockerfile")); err != nil {
		return errors.Wrap(err, "failed to link dockerfile")
	}

	files, err := files.Find(docker.Files)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(file)), 0755); err != nil {
			return errors.Wrapf(err, "failed to link extra file '%s'", file)
		}
		if err := link(file, filepath.Join(dir, file)); err != nil {
			return errors.Wrapf(err, "failed to link extra file '%s'", file)
		}
	}

	return nil
}

func link(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var dst = filepath.Join(dest, strings.Replace(path, src, "", 1))
		log.WithFields(log.Fields{
			"src": path,
			"dst": dst,
		}).Debug("extra file")
		if info.IsDir() {
			return os.MkdirAll(dst, info.Mode())
		}
		return os.Link(path, dst)
	})
}

func build(ctx *context.Context, dir string, images, flags []string) error {
	log.WithField("image", images[0]).
		Info("building docker image")

	var rest = mountCommand(images, flags)
	var cmd = exec.CommandContext(ctx, "docker", rest...)
	cmd.Dir = dir

	log.WithField("cmd", cmd.Args).
		WithField("cwd", cmd.Dir).
		Debug("running")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to build docker image: %s: \n%s", images[0], string(out))
	}
	log.Debugf("docker build output: \n%s", string(out))

	return nil
}

func mountCommand(images, flags []string) []string {
	base := []string{"build", "."}
	for _, image := range images {
		base = append(base, "-t", image)
	}
	base = append(base, flags...)
	return base
}
