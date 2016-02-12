package main

import (
	"os"
	"webo/cmd"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "webo"
	app.Usage = "Golang Web Development"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "creates webo file structure",
			Flags:   cmd.InitFlags,
			Action:  cmd.InitAction,
		},
	}

	app.Run(os.Args)
}
