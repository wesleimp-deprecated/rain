package files

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"

	"github.com/rainproj/rain/pkg/config"
)

// Find resolves extra files globs et al into a map of names/paths or an error.
func Find(files []config.File) (map[string]string, error) {
	var result = map[string]string{}
	for _, file := range files {
		if file.Glob != "" {
			files, err := zglob.Glob(file.Glob)
			if err != nil {
				return result, errors.Wrapf(err, "globbing failed for pattern %s", file.Glob)
			}
			for _, file := range files {
				info, err := os.Stat(file)
				if err == nil && info.IsDir() {
					log.Debugf("ignoring directory %s", file)
					continue
				}
				var name = filepath.Base(file)
				if old, ok := result[name]; ok {
					log.Warnf("overriding %s with %s for name %s", old, file, name)
				}
				result[name] = file
			}
		}
	}
	return result, nil
}
