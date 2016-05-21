package server

import (
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/apaganobeleno/webo/routing"
	"github.com/stretchr/testify/assert"
)

func TestStatic(t *testing.T) {
	w := NewServer(func(r *routing.Router) {})
	w.AddStatic("static", "/assets")

	assert.Equal(t, len(w.staticHandlers), 1)
}

func TestStaticRepeated(t *testing.T) {
	w := NewServer(func(r *routing.Router) {})
	w.AddStatic("/static", "/assets")
	w.AddStatic("/ssss", "/assets")

	assert.Equal(t, len(w.staticHandlers), 1)
}

func TestStaticCall(t *testing.T) {
	w := NewServer(func(r *routing.Router) {})
	w.AddStatic("./static", "/assets")

	cases := map[string]string{
		"/assets/text.txt":       "Hola Mundo!",
		"/assets/":               "Wohoo!",
		"/assets/css/styles.css": "body{",
	}

	for path, expected := range cases {
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		w.ServeHTTP(rw, req)

		responseStr := string(rw.Body.Bytes())
		assert.Contains(t, responseStr, expected)
	}
}

func TestStaticCallWithRoute(t *testing.T) {
	w := NewServer(func(r *routing.Router) {
		r.Get("/api_endpoint", func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(422)
		})
	})

	w.AddStatic("static", "/")

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api_endpoint", nil)

	w.ServeHTTP(rw, req)
	assert.Equal(t, rw.Code, 422)
}

func TestPanicRecover(t *testing.T) {
	w := NewServer(func(r *routing.Router) {
		r.Get("/panic", func(rw http.ResponseWriter, req *http.Request) {
			var array []int
			println(array[0])
		})
	})

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	w.ServeHTTP(rw, req)

	assert.Equal(t, rw.Code, http.StatusInternalServerError)
}

func TestPort(t *testing.T) {
	os.Setenv("PORT", "12345")
	port := portToRun("123")
	assert.Equal(t, "12345", port)

	os.Setenv("PORT", "")
	port = portToRun("123")
	assert.Equal(t, "123", port)
}
