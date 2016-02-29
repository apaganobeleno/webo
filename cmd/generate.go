package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/serenize/snaker"
)

type GenHandlerParam struct {
	Name string
}

var genHanlerParams = GenHandlerParam{}

var GenHandlerAction = func(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Println("| ERROR: Please specify a handler name.")
		return
	}

	genHanlerParams.Name = strings.Replace(c.Args()[0], " ", "", -1)

	if exists, _ := pathExists("handlers"); exists == false {
		fmt.Println("| ERROR: handlers folder does not exist in current workign directory.")
		return
	}

	fileName := snaker.CamelToSnake(genHanlerParams.Name)
	filePath := path.Join("handlers", fileName+".go")

	if exists, _ := pathExists(filePath); exists == true {
		fmt.Printf("| ERROR: '%s' already exists \n", filePath)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("| ERROR: Could not create file: ", err.Error())
		return
	}

	err, dat := writeTemplatedFile(handlerGenTemplate, genHanlerParams)
	if err != nil {
		fmt.Println("| ERROR: Could not create handler: ", err.Error())
		os.Remove(filePath)
	}

	file.Write(dat)
}
