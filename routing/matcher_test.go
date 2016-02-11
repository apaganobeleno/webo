package routing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcherWithSimpleRoute(t *testing.T) {
	router := Router{}
	router.Get("/api", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello!"))
	})

	req, _ := http.NewRequest("GET", "/api", strings.NewReader(""))
	rw := httptest.NewRecorder()

	matcher := Matcher{}
	matcher.Match(rw, req, &router)

	assert.Equal(t, rw.Body.Bytes(), []byte("Hello!"))
}

func TestMatcherWithVariable(t *testing.T) {
	router := Router{}

	router.Get("/api/{id}", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello!"))
	})

	req, _ := http.NewRequest("GET", "/api/1", strings.NewReader(""))
	rw := httptest.NewRecorder()

	matcher := Matcher{}
	matcher.Match(rw, req, &router)

	assert.Equal(t, rw.Body.Bytes(), []byte("Hello!"))
}

func TestMatcherWithVariableExp(t *testing.T) {
	router := Router{}

	router.Get("/api/{id:[0-9]+}", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello With Regexp!"))
	})

	req, _ := http.NewRequest("GET", "/api/A", strings.NewReader(""))
	rw := httptest.NewRecorder()

	matcher := Matcher{}
	matcher.Match(rw, req, &router)

	assert.Equal(t, rw.Body.Bytes(), []byte(nil))
}

func TestMatcherWithTwoVariableExp(t *testing.T) {
	router := Router{}
	expectedResult := []byte("Hello With Two Regexp!")

	router.Get("/api/{id:[0-9]+}/other/{base}", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(expectedResult)
	})

	req, _ := http.NewRequest("GET", "/api/1/other/ouch", strings.NewReader(""))
	rw := httptest.NewRecorder()

	matcher := Matcher{}
	matcher.Match(rw, req, &router)

	assert.Equal(t, rw.Body.Bytes(), expectedResult)
}

func TestMatcherParamsInForm(t *testing.T) {
	router := Router{}

	router.Get("/api/{name}/other/{base}", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(req.Form.Get("name") + req.Form.Get("base")))
	})

	req, _ := http.NewRequest("GET", "/api/Tony/other/WWCO", strings.NewReader(""))
	rw := httptest.NewRecorder()

	matcher := Matcher{}
	matcher.Match(rw, req, &router)

	assert.Equal(t, string(rw.Body.Bytes()), "TonyWWCO")
}
