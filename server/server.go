package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"strings"

	"github.com/apaganobeleno/webo/routing"
)

//Webo is a general struct that holds all the config
type Webo struct {
	Port           int
	Server         *http.Server
	Matcher        *routing.Matcher
	Router         *routing.Router
	staticHandlers map[string]http.Handler
}

//NewServer creates a new instance of the server
func NewServer(routesDefinitionFunction func(r *routing.Router)) *Webo {
	w := Webo{
		Router:         &routing.Router{},
		Matcher:        &routing.Matcher{},
		staticHandlers: map[string]http.Handler{},
	}

	routesDefinitionFunction(w.Router)
	return &w
}

//Start starts the Server on the passed port
func (w *Webo) Start(port string) {
	port = portToRun(port)
	log.Printf("| [webo] listening on port %s", port)

	w.Server = &http.Server{
		Addr:           ":" + port,
		Handler:        w,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(w.Server.ListenAndServe())
}

//AddStatic adds a static server handler for static files
func (w *Webo) AddStatic(dir, path string) {
	w.staticHandlers[path] = StaticHandler{
		Path:      dir,
		ServePath: path,
	}
}

func (w *Webo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("| Serving /%s", req.URL.Path[1:])
	defer recoverFromPanic(rw)

	//TODO: this should be a separate middleware.
	if req.Method == "OPTIONS" {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		rw.WriteHeader(200)

		rw.Write([]byte("OK"))
		return
	}

	found := w.Matcher.Match(rw, req, w.Router) || w.handleWithStaticFiles(rw, req)

	if !found {
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found"))
	}
}

func (w *Webo) handleWithStaticFiles(rw http.ResponseWriter, req *http.Request) bool {
	for key, handler := range w.staticHandlers {
		reqURL := req.URL.RequestURI()

		if strings.HasPrefix(reqURL, key) {
			log.Println("| found Handler!")
			handler.ServeHTTP(rw, req)
			return true
		}
	}

	return false
}

func portToRun(input string) string {
	port := os.Getenv("PORT")

	if port != "" {
		log.Printf("[webo] Serving from port %s since its specified in $PORT ENV variable.", port)
		return port
	}

	return input
}

func recoverFromPanic(rw http.ResponseWriter) {
	if err := recover(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println("| Error: ", err)
	}
}
