/*
 * API DoH
 *
 * This is a DNS-over-HTTP (DoH) resolver written in Go.
 *
 * API version: 0.1
 * Contact: dev@phunsites.net
 */

package api_doh

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "DoH Server")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"DNSQueryGet",
		strings.ToUpper("Get"),
		"/dns-query",
		DNSQueryGet,
	},

	Route{
		"DNSQueryPost",
		strings.ToUpper("Post"),
		"/dns-query",
		DNSQueryPost,
	},
}
