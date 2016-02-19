package routing

import (
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/apaganobeleno/webo/routing"
)

type Webo struct {
	Port           int
	Server         *http.Server
	Matcher        *routing.Matcher
	Router         *routing.Router
	staticHandlers map[string]http.Handler
}

func NewServer(routesDefinitionFunction func(r *routing.Router)) *Webo {
	w := Webo{
		Router:         &routing.Router{},
		Matcher:        &routing.Matcher{},
		staticHandlers: map[string]http.Handler{},
	}

	routesDefinitionFunction(w.Router)
	return &w
}

func (w *Webo) Start(port string) {
	log.Printf("| [webo] listening on port " + port)
	w.Server = &http.Server{
		Addr:           ":" + port,
		Handler:        w,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(w.Server.ListenAndServe())
}

func (w *Webo) AddStatic(dir, path string) {
	fs := http.FileServer(http.Dir(dir))
	w.staticHandlers[path] = http.StripPrefix(path, fs)
}

func (w *Webo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("| Serving /%s", req.URL.Path[1:])
	found := w.Matcher.Match(rw, req, w.Router)

	if !found {
		for key, handler := range w.staticHandlers {
			reqURL := req.URL.RequestURI()

			if strings.HasPrefix(reqURL, key) {
				handler.ServeHTTP(rw, req)
				return
			}
		}

		//TODO: move this to a separate class/function
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found"))
	}
}
