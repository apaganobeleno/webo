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
	// 	{
	// 		Name:    "serve",
	// 		Aliases: []string{"s"},
	// 		Usage:   "runs the webo app",
	// 		Action: func(c *cli.Context) {
	// 			log.Println("Starting Webo App")
	// 		},
	// 	},
	// }
	app.Run(os.Args)
}
