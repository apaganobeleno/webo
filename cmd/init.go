package cmd

import (
	"errors"
	"fmt"
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
		fmt.Println("| ERROR: Please specify an app name.")
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
		fmt.Println("| ERROR: Folder already exist.")
		return errors.New("Folder already exist.")
	}
	err := os.Mkdir(projectPath, 0777)
	if err != nil {
		fmt.Println("| ERROR: Could not create directory: ", err.Error())
		return err
	}

	return nil
}

func createSubfolders() {
	projectPath := path.Join(initParams.DestionationPath, initParams.Name)
	folders := []string{
		path.Join(projectPath, "handlers"),
		path.Join(projectPath, "middlewares"),
		path.Join(projectPath, "config"),
	}

	for _, folder := range folders {
		err := os.Mkdir(folder, 0777)
		if err != nil {
			fmt.Println("| ERROR: Could not create directory: ", err.Error())
		}
	}
}

func createFiles() {
	projectPath := path.Join(initParams.DestionationPath, initParams.Name)

	files := map[string]string{
		path.Join(projectPath, "handlers", "sample.go"):   handlersTemplate,
		path.Join(projectPath, "middlewares", ".gitkeep"): ".gitkeep",
		path.Join(projectPath, "config", "routes.go"):     routesTemplate,
		path.Join(projectPath, "main.go"):                 mainTemplate,
	}

	for path, content := range files {
		file, err := os.Create(path)

		if err != nil {
			fmt.Println("| ERROR: Could not create file: ", err.Error())
			continue
		}

		_, data := writeTemplatedFile(content, initParams)
		file.Write(data)
	}
}
