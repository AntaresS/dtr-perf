package pull

import (
	"context"
	"fmt"

	"github.com/docker/perfkit/dtr/stress/pull"
	"github.com/docker/perfkit/dtr/stress/sharedutils"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

func NewCommands() []*cli.Command {
	cmd := []*cli.Command{
		{
			Name:  "pull",
			Usage: "execute docker pull command continuously",
			Action: func(c *cli.Context) error {
				cfg := &pull.Config{}
				in, err := sharedutils.ReadConfigFile(c.String("file"))
				if err != nil {
					return err
				}
				fmt.Println(string(in))
				err = yaml.Unmarshal(in, cfg)
				fmt.Printf("%+v", *cfg)
				j := pull.Job{
					Config: cfg,
				}
				return pull.StressPull(context.Background(), &j)
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "config file for pull simulation",
					Value:   "",
				},
			},
		},
	}
	return cmd
}
