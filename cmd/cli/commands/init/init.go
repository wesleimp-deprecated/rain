package init

import (
	"os"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/rainproj/rain/internal/static"
)

// Command init
var Command = &cli.Command{
	Name:    "init",
	Usage:   "Generates a .rain.yml file",
	Aliases: []string{"i"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Value:   ".rain.yml",
			Aliases: []string{"c"},
			Usage:   "Load configuration from file",
		},
	},
	Action: func(c *cli.Context) error {
		var filename = c.String("config")
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		log.Infof(color.New(color.Bold).Sprintf("Generating %s file", filename))
		if _, err := f.WriteString(static.InitConfig); err != nil {
			return err
		}

		log.WithField("file", filename).Info("config created")
		return nil
	},
}
