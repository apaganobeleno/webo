package webo

import "net/http"

var commonValue = 0

func SampleHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
}

func unauthorizedMiddleware(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	commonValue = 1
	rw.WriteHeader(401)
}

func valueChangerMiddleware(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	commonValue = 4
	rw.WriteHeader(201)
}
