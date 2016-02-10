package webo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteDefinition(t *testing.T) {
	router := Router{}
	router.Get("/route", SampleHandler)

	assert.Equal(t, len(router.Routes), 1)
}

func TestRouteGroupDefinition(t *testing.T) {
	router := Router{}

	router.Group("/api", func(rg *RouteGroup) {
		rg.Get("/users", SampleHandler)
		rg.Get("/others", SampleHandler)
	})

	assert.Equal(t, len(router.Routes), 2)
	assert.Contains(t, router.Routes[0].Path, "api")
}

func TestRouteGroupNestedDefinition(t *testing.T) {
	router := Router{}

	router.Group("/api", func(rg *RouteGroup) {
		rg.Get("/users", SampleHandler)
		rg.Get("/others", SampleHandler)

		rg.Group("/user", func(rg *RouteGroup) {
			rg.Get("/{id}/data", SampleHandler)
		})
	})

	assert.Equal(t, len(router.Routes), 3)
	assert.Contains(t, router.Routes[0].Path, "api")
	assert.Equal(t, router.Routes[2].Path, "/api/user/{id}/data")
}

func TestEmptyNameGroup(t *testing.T) {
	router := Router{}

	router.Group("", func(rg *RouteGroup) {
		rg.Get("/users", SampleHandler)
	}).Before(unauthorizedMiddleware)

	assert.Equal(t, len(router.Routes), 1)
	assert.Equal(t, router.Routes[0].Path, "/users")
	assert.Equal(t, len(router.Routes[0].BeforeMiddlewares), 1)
}

func TestRouteGroupDefinitionWithCases(t *testing.T) {
	router := Router{}

	router.Group("/", func(rg *RouteGroup) {
		rg.Get("/AAA", SampleHandler)
	})

	router.Group("/api/", func(rg *RouteGroup) {
		rg.Get("/users", SampleHandler)
		rg.Get("/others", SampleHandler)
		rg.Get("/{id:[0-9]+}/name", SampleHandler)
		rg.Get("/{id}/name", SampleHandler)
	})

	assert.Equal(t, len(router.Routes), 5)
	assert.Equal(t, router.Routes[0].Path, "/AAA")
	assert.Equal(t, router.Routes[1].Path, "/api/users")
	assert.Equal(t, router.Routes[3].Path, "/api/{id:[0-9]+}/name")
	assert.Equal(t, router.Routes[4].Path, "/api/{id}/name")
}
