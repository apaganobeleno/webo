package routing

import (
	"log"
	"net/http"
	"time"
	"webo/routing"
)

type Webo struct {
	Port    int
	Server  *http.Server
	Matcher *routing.Matcher
	Router  *routing.Router
}

func NewServer(routesDefinitionFunction func(r *routing.Router)) *Webo {
	w := Webo{
		Router:  &routing.Router{},
		Matcher: &routing.Matcher{},
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

func (w *Webo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("| Serving /%s", req.URL.Path[1:])
	w.Matcher.Match(rw, req, w.Router)
}
