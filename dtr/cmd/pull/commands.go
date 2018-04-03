package pull

import (
	"context"

	"github.com/docker/perfkit/dtr/stress"
	"github.com/docker/perfkit/dtr/stress/sharedutils"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

func NewCommands() []*cli.Command {
	cmd := []*cli.Command{
		{
			Name:  "pull",
			Usage: "execute docker pull command continuously",
			Action: func(c *cli.Context) error {
				if c.Bool("debug") {
					logrus.SetLevel(logrus.DebugLevel)
				}
				cfg := &stress.Config{}
				in, err := sharedutils.ReadConfigFile(c.String("file"))
				if err != nil {
					return err
				}
				err = yaml.Unmarshal(in, cfg)
				if err != nil {
					return err
				}
				j := stress.Job{
					Config: cfg,
				}
				return stress.StressPull(context.Background(), &j)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "config file for pull simulation",
					Value:   "",
				}, &cli.BoolFlag{
					Name:    "debug",
					Aliases: []string{"d"},
					Usage:   "debug mode",
					Value:   false,
				},
			},
		},
	}
	return cmd
}
