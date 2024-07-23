package router

import (
	"net/http"
)

// pass a map of routes to handlers;
// router then invokes the correct handler
// for the route requested, otherwise
// sends back not found.
func Router(routes map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var head string
		head, r.URL.Path = ShiftPath(r.URL.Path)
		handler, ok := routes[head]
		if ok {
			handler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}
