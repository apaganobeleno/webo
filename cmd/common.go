package cmd

import (
	"bytes"
	"go/format"
	"os"
	"text/template"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func writeTemplatedFile(templStr string, data interface{}) (error, []byte) {
	buf := new(bytes.Buffer)
	template.Must(template.New("some_template").Parse(templStr)).Execute(buf, data)

	formatted, err := format.Source(buf.Bytes())
	return err, formatted
}
