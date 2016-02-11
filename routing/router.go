package routing

import "net/http"

type Router struct {
	Routes []*Route
}

// RespondTo adds support for different HTTP methods
func (r *Router) RespondTo(method string, enpointPath string, handler http.HandlerFunc) *Route {
	route := NewRoute(enpointPath, method, handler)
	r.Routes = append(r.Routes, &route)
	return &route
}

func (r *Router) Group(prefix string, definition RouteGroupBlock) *RouteGroup {
	group := RouteGroup{Prefix: prefix, Router: r}
	definition(&group)
	return &group
}

func (r *Router) Get(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("GET", enpointPath, handler)
	return route
}

func (r *Router) Post(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("POST", enpointPath, handler)
	return route
}

func (r *Router) Put(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("PUT", enpointPath, handler)
	return route
}

func (r *Router) Delete(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("DELETE", enpointPath, handler)
	return route
}
