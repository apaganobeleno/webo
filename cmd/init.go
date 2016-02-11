package cmd

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path"
	"text/template"

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
	err := createFolder()
	if err != nil {
		return
	}

	createSubfolders()
	createFiles()
}

func createFolder() error {
	projectPath := path.Join(initParams.DestionationPath, initParams.Name)

	if exists, _ := pathExists(projectPath); exists == true {
		log.Println("| ERROR: Folder already exist.")
		return errors.New("Folder already exist.")
	}

	//Create {Path}/{AppName} folder
	err := os.Mkdir(projectPath, 0777)
	if err != nil {
		log.Println("| ERROR: Could not create directory: ", err.Error())
		return err
	}

	return nil
}

func createSubfolders() {
	//Create folders inside {Path}/{AppName}
	projectPath := path.Join(initParams.DestionationPath, initParams.Name)
	folders := []string{
		path.Join(projectPath, "handlers"),
		path.Join(projectPath, "middlewares"),
	}

	for _, folder := range folders {
		err := os.Mkdir(folder, 0777)
		if err != nil {
			log.Println("| ERROR: Could not create directory: ", err.Error())
		}
	}
}

func createFiles() {
	projectPath := path.Join(initParams.DestionationPath, initParams.Name)
	//Create Files based on templates
	files := map[string]string{
		path.Join(projectPath, "handlers", "sample.go"):   handlersTemplate,
		path.Join(projectPath, "middlewares", ".gitkeep"): ".gitkeep",
		path.Join(projectPath, "routes.go"):               routesTemplate,
	}

	for path, content := range files {
		file, err := os.Create(path)
		if err != nil {
			log.Println("| ERROR: Could not create file: ", err.Error())
			continue
		}

		buf := new(bytes.Buffer)
		template.Must(template.New("sample").Parse(content)).Execute(buf, initParams)
		file.WriteString(buf.String())
	}
}

var routesTemplate = `package main
import "{{.Name}}/handlers"

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
