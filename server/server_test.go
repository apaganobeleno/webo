package routing

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/apaganobeleno/webo/routing"
	"github.com/stretchr/testify/assert"
)

func TestStatic(t *testing.T) {
	w := NewServer(func(r *routing.Router) {})
	w.AddStatic("/static", "/assets")

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
	w.AddStatic("../static", "/assets")

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assets/text.txt", nil)

	w.ServeHTTP(rw, req)
	responseStr := string(rw.Body.Bytes())
	assert.Equal(t, responseStr, "Hola Mundo!\n")
}

func TestStaticCallRootWithHTML(t *testing.T) {
	w := NewServer(func(r *routing.Router) {})
	w.AddStatic("../static", "/assets")

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assets/", nil)

	w.ServeHTTP(rw, req)
	responseStr := string(rw.Body.Bytes())
	assert.Contains(t, responseStr, "Wohoo")
}

func TestStaticCallWithCss(t *testing.T) {
	w := NewServer(func(r *routing.Router) {})
	w.AddStatic("../static", "/assets")

	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assets/css/styles.css", nil)

	w.ServeHTTP(rw, req)
	responseStr := string(rw.Body.Bytes())
	assert.Contains(t, responseStr, "body{")
}
