package routing

import "net/http"

type Middleware func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc)
