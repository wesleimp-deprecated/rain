package config

import (
	"errors"
	"os"

	"github.com/wesleimp/rain/pkg/config"
)

// Load config file
func Load(path string) (config.Config, error) {
	if path != "" {
		return config.Load(path)
	}
	for _, f := range [4]string{
		".rain.yml",
		".rain.yaml",
		"rain.yml",
		"rain.yaml",
	} {
		proj, err := config.Load(f)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		return proj, err
	}

	return config.Config{}, errors.New("could nor find a config file")
}
