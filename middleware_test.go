package webo

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeforeMiddlewareSet(t *testing.T) {
	router := Router{}
	route := router.Get("/endpoint", func(rw http.ResponseWriter, req *http.Request) {
		log.Println("Here!")
		rw.WriteHeader(201)
	}).Before(unauthorizedMiddleware)

	assert.Equal(t, len(route.BeforeMiddlewares), 1)
}

func TestBeforeMiddlewareCalled(t *testing.T) {
	router := Router{}
	router.Get("/endpoint", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(201)
	}).Before(unauthorizedMiddleware)

	req, _ := http.NewRequest("GET", "/endpoint", strings.NewReader(""))
	rw := httptest.NewRecorder()
	matcher := Matcher{}

	matcher.Match(rw, req, &router)
	assert.Equal(t, rw.Code, 401)
}

func TestBeforeMiddlewarePasses(t *testing.T) {
	router := Router{}
	router.Get("/endpoint", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(201)
	}).Before(func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		next(rw, req)
	})

	req, _ := http.NewRequest("GET", "/endpoint", strings.NewReader(""))
	rw := httptest.NewRecorder()
	matcher := Matcher{}

	matcher.Match(rw, req, &router)
	assert.Equal(t, rw.Code, 201)
}

func TestBeforeMiddlewareGroup(t *testing.T) {
	router := Router{}
	var route *Route

	router.Group("/api", func(rg *RouteGroup) {
		route = rg.Get("/endpoint", func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(201)
		}).Before(valueChangerMiddleware)
	}).Before(unauthorizedMiddleware)

	req, _ := http.NewRequest("GET", "/api/endpoint", strings.NewReader(""))
	rw := httptest.NewRecorder()
	matcher := Matcher{}

	matcher.Match(rw, req, &router)
	assert.Equal(t, len(route.BeforeMiddlewares), 2)
	assert.Equal(t, rw.Code, 401)
	assert.Equal(t, 1, commonValue)

}

func TestBeforeMiddlewareNestedGroup(t *testing.T) {
	router := Router{}
	var route *Route

	router.Group("/api", func(rg *RouteGroup) {
		rg.Get("/endpoint", func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(201)
		}).Before(valueChangerMiddleware)

		rg.Group("/users", func(rg *RouteGroup) {
			rg.Group("/active", func(rg *RouteGroup) {
				route = rg.Get("/{id}", SampleHandler)
			}).Before(unauthorizedMiddleware)
		})

	}).Before(unauthorizedMiddleware)

	req, _ := http.NewRequest("GET", "/api/endpoint", strings.NewReader(""))
	rw := httptest.NewRecorder()
	matcher := Matcher{}

	matcher.Match(rw, req, &router)
	assert.Equal(t, len(route.BeforeMiddlewares), 2)
	assert.Equal(t, rw.Code, 401)
	assert.Equal(t, 1, commonValue)

}
