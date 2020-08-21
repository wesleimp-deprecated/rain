package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

// File struct
type File struct {
	Glob string `yaml:",omitempty"`
}

// Docker contains all informarmation for build docker images
type Docker struct {
	Dockerfile         string   `yaml:",omitempty"`
	ImageTemplates     []string `yaml:"image_templates"`
	BuildFlagTemplates []string `yaml:"build_flag_templates,omitempty"`
	SkipPush           bool     `yaml:"skip_push,omitempty"`
	Files              []File   `yaml:"extra_files,omitempty"`
}

// Config contains the project configuration
type Config struct {
	ProjectName string   `yaml:"project_name,omitempty"`
	Version     string   `yaml:",omitempty"`
	Env         []string `yaml:",omitempty"`
	Dockers     []Docker `yaml:",omitempty"`
	Dist        string   `yaml:",omitempty"`
}

// Load config file.
func Load(file string) (config Config, err error) {
	f, err := os.Open(file) // #nosec
	if err != nil {
		return
	}
	defer f.Close()
	log.WithField("file", file).Info("loading config file")
	return LoadReader(f)
}

// LoadReader config via io.Reader.
func LoadReader(fd io.Reader) (config Config, err error) {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return config, err
	}
	err = yaml.UnmarshalStrict(data, &config)
	log.WithField("config", config).Debug("loaded config file")
	return config, err
}
