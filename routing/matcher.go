package routing

import "net/http"

type Matcher struct{}

//Match runs a route handler for specific request that matches it.
func (m *Matcher) Match(rw http.ResponseWriter, req *http.Request, router *Router) {
	var route *Route

	for _, r := range router.Routes {
		if r.Matches(req) {
			route = r
			break
		}
	}

	if route == nil {
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found"))
	} else {
		m.ApplyRouteParameters(req, route)
		route.Attend(rw, req)
	}
}

//ApplyRouteParamters adds named parameters into the form after parsing it.
func (m *Matcher) ApplyRouteParameters(req *http.Request, route *Route) {
	paramNames := route.Regexp.SubexpNames()
	result := route.Regexp.FindStringSubmatch(req.URL.RequestURI())
	req.ParseForm()

	for i, n := range result {
		if paramNames != nil && paramNames[i] != "" {
			req.Form.Add(paramNames[i], n)
		}
	}
}
