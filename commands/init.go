package commands

import (
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
)

type InitParams struct {
	Name             string
	DestionationPath string
}

var initParams = InitParams{}

var InitFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "path",
		Value:       ".",
		Usage:       "project Name",
		Destination: &initParams.DestionationPath,
	},
}

var InitAction = func(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Println("| ERROR: Please specify an app name.")
		return
	}

	initParams.Name = c.Args()[0]
	projectPath := path.Join(initParams.DestionationPath, initParams.Name)

	if exists, _ := pathExists(projectPath); exists == true {
		log.Println("| ERROR: Folder already exist.")
		return
	}

	err := os.Mkdir(projectPath, 0777)
	if err != nil {
		log.Println("| ERROR: Could not create directory: ", err.Error())
		return
	}

	//Create {Path}/{AppName} folder
	//Create folders inside {Path}/{AppName}
	//Create Files based on templates
}

func InitProject(c *cli.Context) error {
	return nil
}

var routesTemplate = `package main
import "{.Name}/handlers"

//defineRoutes gets called by the generated main.go, please define routes inside it.
func defineRoutes(r *Router){
  r.Get("/", handlers.Home)
}`

var handlersTemplate = `package handlers
import "http"
//Home renders a Hello message
func Home(rw http.ResponseWriter, req *http.Request){
  rw.Write([]byte("Hello From Webo!"))
}`
