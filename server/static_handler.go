package server

import (
	"net/http"
	"path/filepath"
	"strings"
)

//StaticHandler Takes care of serving static files
type StaticHandler struct {
	Path      string
	ServePath string
}

func (sh StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, sh.ServePath)

	if len(path) < len(r.URL.Path) {
		path = applyHTMLExtension(path)
		r.URL.Path = filepath.Join(sh.Path, path)

		http.ServeFile(w, r, r.URL.Path)
		return
	}

	http.NotFound(w, r)
}

func applyHTMLExtension(path string) string {
	isDirectory := strings.HasSuffix(path, "/")
	if !isDirectory {
		parts := strings.Split(path, ".")
		if len(parts) == 1 {
			path = path + ".html"
		}
	}

	return path
}
