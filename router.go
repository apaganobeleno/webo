package webo

import (
	"net/http"
	"path"
)

type Router struct {
	Routes []*Route
}

type RouteGroup struct {
	Prefix string
	Router *Router
	Routes []*Route
	Groups []*RouteGroup
}

type RouteGroupBlock func(rg *RouteGroup)

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

func (rg *RouteGroup) Before(mdwre Middleware) *RouteGroup {
	for _, route := range rg.Routes {
		route.Before(mdwre)
	}

	for _, rg := range rg.Groups {
		rg.Before(mdwre)
	}

	return rg
}

func (rg *RouteGroup) Group(prefix string, definition RouteGroupBlock) *RouteGroup {
	fullPrefix := path.Clean(rg.Prefix + "/" + prefix)
	group := RouteGroup{Prefix: fullPrefix, Router: rg.Router}
	rg.Groups = append(rg.Groups, &group)

	definition(&group)
	return rg
}

func (rg *RouteGroup) RespondTo(method string, enpointPath string, handler http.HandlerFunc) *Route {
	fullPath := path.Clean(rg.Prefix + enpointPath)
	route := rg.Router.RespondTo(method, fullPath, handler)
	rg.Routes = append(rg.Routes, route)
	return route
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

func (r *RouteGroup) Get(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("GET", enpointPath, handler)
	return route
}

func (r *RouteGroup) Post(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("POST", enpointPath, handler)
	return route
}

func (r *RouteGroup) Put(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("PUT", enpointPath, handler)
	return route
}

func (r *RouteGroup) Delete(enpointPath string, handler http.HandlerFunc) *Route {
	route := r.RespondTo("DELETE", enpointPath, handler)
	return route
}
