package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
)

const mainTemplate = `package main
import (
  "log"
  "net/http"
  "webo/routing"
  "{{.Package}}/config"
)

func main(){
  router := routing.Router{}
  config.DefineRoutes(&router)
	log.Println("| [{{.Package}}] Running")

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
    log.Printf("| Serving /%s", r.URL.Path[1:] )

    matcher := routing.Matcher{}
    matcher.Match(w, r, &router)
  })

  http.ListenAndServe(":8080", nil)
	defer log.Printf("| {{.Package}} Closing")
}
`

func RunAction(c *cli.Context) {
	tempFolder := os.TempDir()
	tempFolder = path.Join(tempFolder, ".webo")
	setupTmp(tempFolder)

	workingDirectory, _ := os.Getwd()
	folders := strings.Split(workingDirectory, "/")
	packageName := folders[len(folders)-1]
	packageFolder := path.Join(tempFolder, "src", packageName)

	copySources(packageFolder)
	addMain(packageFolder, packageName)
	runMain(tempFolder, packageName)
}

func setupTmp(tmpFolder string) {
	cleanupTmp(tmpFolder)

	os.Mkdir(tmpFolder, 0777)
	os.Mkdir(path.Join(tmpFolder, "src"), 0777)
}

func copySources(destFolder string) {
	workingDirectory, _ := os.Getwd()

	filepath.Walk(workingDirectory, func(strPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(strPath, ".tmp") || strings.Contains(strPath, ".git") { //TODO add other non-needed files
			return nil
		}

		newPath := strings.Replace(strPath, workingDirectory, "", 1)
		newPath = path.Join(destFolder, newPath)

		if !info.IsDir() {
			data, _ := ioutil.ReadFile(strPath)
			err = ioutil.WriteFile(newPath, data, 0777)
		} else {
			os.Mkdir(newPath, 0777)
		}

		return nil
	})
}

func addMain(destFolder, packageName string) {
	mainPath := path.Join(destFolder, "main.go")
	file, _ := os.Create(mainPath)
	buf := new(bytes.Buffer)

	template.Must(template.New("run").Parse(mainTemplate)).Execute(buf, struct{ Package string }{packageName})
	file.WriteString(buf.String())
}

func runMain(tempFolder, packageName string) {
	mainPath := path.Join(tempFolder, "src", packageName, "main.go")
	execGopath := os.Getenv("GOPATH") + ":" + path.Join(tempFolder)

	command := exec.Command("go", "run", mainPath)
	command.Env = []string{"GOPATH=" + execGopath}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	command.Run()
	command.Wait()
}

func cleanupTmp(tmpFolder string) {
	var directoryList []string

	filepath.Walk(tmpFolder, func(strPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directoryList = append(directoryList, strPath)
		} else {
			os.Remove(strPath)
		}
		return nil
	})

	for _, dir := range directoryList {
		os.Remove(dir)
	}
}
