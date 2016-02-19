package routing

import "net/http"

type Matcher struct{}

//Match runs a route handler for specific request that matches it,
//returns true or false depending if it finds or not the route match.
func (m *Matcher) Match(rw http.ResponseWriter, req *http.Request, router *Router) bool {
	var route *Route

	for _, r := range router.Routes {
		if r.Matches(req) {
			route = r
			break
		}
	}

	if route == nil {
		return false
	}

	m.ApplyRouteParameters(req, route)
	route.Attend(rw, req)
	return true
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
