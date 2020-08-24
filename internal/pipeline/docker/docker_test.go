package docker

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wesleimp/rain/pkg/config"
	"github.com/wesleimp/rain/pkg/context"
)

var it = flag.Bool("it", false, "push images to docker hub")
var registry = "localhost:5000/"
var altRegistry = "localhost:5050/"

func TestMain(m *testing.M) {
	flag.Parse()
	if *it {
		registry = "docker.io/"
	}
	os.Exit(m.Run())
}

func start(t *testing.T) {
	if *it {
		return
	}
	if out, err := exec.Command(
		"docker", "run", "-d", "-p", "5000:5000", "--name", "registry", "registry:2",
	).CombinedOutput(); err != nil {
		t.Log("failed to start docker registry", string(out), err)
		t.FailNow()
	}
	if out, err := exec.Command(
		"docker", "run", "-d", "-p", "5050:5000", "--name", "alt_registry", "registry:2",
	).CombinedOutput(); err != nil {
		t.Log("failed to start alternate docker registry", string(out), err)
		t.FailNow()
	}
}

func killAndRm(t *testing.T) {
	if *it {
		return
	}
	t.Log("killing registry")
	_ = exec.Command("docker", "kill", "registry").Run()
	_ = exec.Command("docker", "rm", "registry").Run()
	_ = exec.Command("docker", "kill", "alt_registry").Run()
	_ = exec.Command("docker", "rm", "alt_registry").Run()
}

func TestRunPipe(t *testing.T) {
	type errChecker func(*testing.T, error)
	var shouldErr = func(msg string) errChecker {
		return func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), msg)
		}
	}
	var shouldNotErr = func(t *testing.T, err error) {
		assert.NoError(t, err)
	}
	type imageLabelFinder func(*testing.T, int)
	var shouldFindImagesWithLabels = func(image string, filters ...string) func(*testing.T, int) {
		return func(t *testing.T, count int) {
			for _, filter := range filters {
				output, err := exec.Command(
					"docker", "images", "-q", "*/"+image,
					"--filter", filter,
				).CombinedOutput()
				assert.NoError(t, err)
				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
				assert.Equal(t, count, len(lines))
			}
		}

	}
	var noLabels = func(t *testing.T, count int) {}

	var table = map[string]struct {
		dockers           []config.Docker
		env               map[string]string
		expect            []string
		assertImageLabels imageLabelFinder
		assertError       errChecker
		pubAssertError    errChecker
	}{
		"valid": {
			env: map[string]string{
				"FOO": "123",
			},
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/test_run_pipe:{{.Env.FOO}}",
						registry + "dim/test_run_pipe:latest",
						altRegistry + "dim/test_run_pipe:{{.Env.FOO}}",
						altRegistry + "dim/test_run_pipe:latest",
					},
					Dockerfile: "testdata/Dockerfile",
					BuildFlagTemplates: []string{
						"--label=org.label-schema.schema-version=1.0.0",
					},
					Files: []config.File{{Glob: "testdata/file.txt"}},
				},
			},
			expect: []string{
				registry + "dim/test_run_pipe:123",
				registry + "dim/test_run_pipe:latest",
				altRegistry + "dim/test_run_pipe:123",
				altRegistry + "dim/test_run_pipe:latest",
			},
			assertImageLabels: shouldFindImagesWithLabels(
				"label=org.label-schema.version=1.0.0",
			),
			assertError:    shouldNotErr,
			pubAssertError: shouldNotErr,
		},
		"multiple images with same dockerfile": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/test_run_pipe:latest",
					},
					Dockerfile: "testdata/Dockerfile",
					Files:      []config.File{{Glob: "testdata/file.txt"}},
				},
				{
					ImageTemplates: []string{
						registry + "dim/test_run_pipe2:latest",
					},
					Dockerfile: "testdata/Dockerfile",
					Files:      []config.File{{Glob: "testdata/file.txt"}},
				},
			},
			assertImageLabels: noLabels,
			expect: []string{
				registry + "dim/test_run_pipe:latest",
				registry + "dim/test_run_pipe2:latest",
			},
			assertError:    shouldNotErr,
			pubAssertError: shouldNotErr,
		},
		"valid_no_binaries": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/test_run_pipe:latest",
					},
					Dockerfile: "testdata/Dockerfile.bin",
					Files:      []config.File{{Glob: "testdata/file.txt"}},
					SkipPush:   true,
				},
			},
			expect: []string{
				registry + "dim/test_run_pipe:latest",
			},
			assertImageLabels: noLabels,
			assertError:       shouldNotErr,
		},
		"valid_skip_push": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/test_run_pipe:latest",
					},
					Dockerfile: "testdata/Dockerfile",
					Files:      []config.File{{Glob: "testdata/file.txt"}},
					SkipPush:   true,
				},
			},
			expect: []string{
				registry + "dim/test_run_pipe:latest",
			},
			assertImageLabels: noLabels,
			assertError:       shouldNotErr,
		},
		"valid build args": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/test_build_args:latest",
					},
					Dockerfile: "testdata/Dockerfile",
					Files:      []config.File{{Glob: "testdata/file.txt"}},
					BuildFlagTemplates: []string{
						"--label=foo=bar",
					},
				},
			},
			expect: []string{
				registry + "dim/test_build_args:latest",
			},
			assertImageLabels: noLabels,
			assertError:       shouldNotErr,
			pubAssertError:    shouldNotErr,
		},
		"bad build args": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/test_build_args:latest",
					},
					Dockerfile: "testdata/Dockerfile",
					Files:      []config.File{{Glob: "testdata/file.txt"}},
					BuildFlagTemplates: []string{
						"--bad-flag",
					},
				},
			},
			assertImageLabels: noLabels,
			assertError:       shouldErr("unknown flag: --bad-flag"),
		},
		"bad_dockerfile": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/bad_dockerfile:latest",
					},
					Dockerfile: "testdata/Dockerfile.invalid",
				},
			},
			assertImageLabels: noLabels,
			assertError:       shouldErr("pull access denied for none, repository does not exist"),
		},
		"missing_env_on_tag_template": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{
						registry + "dim/test_run_pipe:{{.Env.NONE}}",
					},
					Dockerfile: "testdata/Dockerfile",
					Files:      []config.File{{Glob: "testdata/file.txt"}},
				},
			},
			assertImageLabels: noLabels,
			assertError:       shouldErr(`template: tmpl:1:39: executing "tmpl" at <.Env.NONE>: map has no entry for key "NONE"`),
		},
		"no_permissions": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{"docker.io/none:latest"},
					Dockerfile:     "testdata/Dockerfile",
					Files:          []config.File{{Glob: "testdata/file.txt"}},
				},
			},
			expect: []string{
				"docker.io/none:latest",
			},
			assertImageLabels: noLabels,
			assertError:       shouldNotErr,
			pubAssertError:    shouldErr(`requested access to the resource is denied`),
		},
		"dockerfile_doesnt_exist": {
			dockers: []config.Docker{
				{
					ImageTemplates: []string{"whatever:latest"},
					Dockerfile:     "testdata/Dockerfilezzz",
				},
			},
			assertImageLabels: noLabels,
			assertError:       shouldErr(`failed to link dockerfile`),
		},
	}

	killAndRm(t)
	start(t)
	defer killAndRm(t)

	for name, docker := range table {
		t.Run(name, func(tt *testing.T) {
			folder, err := ioutil.TempDir("", "dockertest")
			assert.NoError(tt, err)
			var dist = filepath.Join(folder, "dist")
			assert.NoError(tt, os.Mkdir(dist, 0755))

			var ctx = context.New(config.Config{
				ProjectName: "test",
				Dist:        dist,
				Dockers:     docker.dockers,
			})

			ctx.Parallelism = 1
			ctx.Env = docker.env

			// this might fail as the image doesnt exist yet, so lets ignore the error
			for _, img := range docker.expect {
				_ = exec.Command("docker", "rmi", img).Run()
			}

			err = Step{}.Run(ctx)
			docker.assertError(tt, err)

			for _, d := range docker.dockers {
				docker.assertImageLabels(tt, len(d.ImageTemplates))
			}

			// this might should not fail as the image should have been created when
			// the step ran
			for _, img := range docker.expect {
				tt.Log("removing docker image", img)
				assert.NoError(tt, exec.Command("docker", "rmi", img).Run(), "could not delete image %s", img)
			}

		})
	}
}

func TestBuildCommand(t *testing.T) {
	images := []string{"dim/test_build_flag", "dim/test_multiple_tags"}
	tests := []struct {
		name   string
		flags  []string
		expect []string
	}{
		{
			name:   "no flags",
			flags:  []string{},
			expect: []string{"build", ".", "-t", images[0], "-t", images[1]},
		},
		{
			name:   "single flag",
			flags:  []string{"--label=foo"},
			expect: []string{"build", ".", "-t", images[0], "-t", images[1], "--label=foo"},
		},
		{
			name:   "multiple flags",
			flags:  []string{"--label=foo", "--build-arg=bar=baz"},
			expect: []string{"build", ".", "-t", images[0], "-t", images[1], "--label=foo", "--build-arg=bar=baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := mountCommand(images, tt.flags)
			assert.Equal(t, tt.expect, command)
		})
	}
}

func TestDescription(t *testing.T) {
	assert.NotEmpty(t, Step{}.String())
}

func TestNoDockers(t *testing.T) {
	assert.Error(t, Step{}.Run(context.New(config.Config{})))
}

func TestNoDockerWithoutImageName(t *testing.T) {
	assert.EqualError(t, Step{}.Run(context.New(config.Config{
		Dockers: []config.Docker{
			{},
		},
	})), "docker section should be configured")
}

func TestDockerNotInPath(t *testing.T) {
	var path = os.Getenv("PATH")
	defer func() {
		assert.NoError(t, os.Setenv("PATH", path))
	}()
	assert.NoError(t, os.Setenv("PATH", ""))

	var ctx = &context.Context{
		Config: config.Config{
			Dockers: []config.Docker{
				{
					ImageTemplates: []string{"a/b"},
				},
			},
		},
	}
	assert.EqualError(t, Step{}.Run(ctx), "docker not present in $PATH")
}
