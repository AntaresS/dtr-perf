package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/perfkit/dtr/cmd/pull"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:    "perfkit",
		Usage:   fmt.Sprintf("%s [options]", os.Args[0]),
		Version: "0.0.1",
		Flags:   []cli.Flag{},
	}
	app.Commands = append(app.Commands, pull.NewCommands()...)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
