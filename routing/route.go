package routing

import (
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	Path              string
	Method            string
	Handler           http.HandlerFunc
	Regexp            *regexp.Regexp
	BeforeMiddlewares []Middleware
}

//NewRoute creates a new Route and initializes the pattern
func NewRoute(routePath string, method string, handler http.HandlerFunc) Route {
	return Route{
		Path:    routePath,
		Method:  method,
		Handler: handler,
		Regexp:  buildPattern(routePath),
	}
}

//Matches returns if a route matches or not a request URI and method
func (r *Route) Matches(req *http.Request) bool {
	methodMatches := req.Method == r.Method
	routeMatches := req.URL != nil
	routeMatches = routeMatches && r.Regexp.MatchString(req.URL.RequestURI())

	return (methodMatches && routeMatches)
}

func (r *Route) Before(middleware Middleware) *Route {
	r.BeforeMiddlewares = append([]Middleware{middleware}, r.BeforeMiddlewares...)
	return r
}

func (r *Route) NextHandler(index int) http.HandlerFunc {
	if index < len(r.BeforeMiddlewares) {
		return func(rw http.ResponseWriter, req *http.Request) {
			r.BeforeMiddlewares[index](rw, req, r.NextHandler(index+1))
		}
	}

	return r.Handler
}

func (r *Route) Attend(rw http.ResponseWriter, req *http.Request) {
	r.NextHandler(0)(rw, req)
}

func buildPattern(routePath string) *regexp.Regexp {
	replaced := regexp.MustCompile("{(\\w*)?:?([^/]*)}").ReplaceAll([]byte(routePath), []byte("{?P<$1>$2}"))
	replaced = regexp.MustCompile("{([^/]*)}").ReplaceAll([]byte(replaced), []byte("($1)"))
	replaced = regexp.MustCompile("(\\<[\\w]*\\>)\\)").ReplaceAll([]byte(replaced), []byte("$1.*)"))
	replaced = []byte(strings.Replace(string(replaced), "?P<>", "", -1))
	replaced = []byte(string(replaced) + "$")
	return regexp.MustCompile(string(replaced))
}
