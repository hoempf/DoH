/*
 * DoH Service - HTTP Router
 *
 * This is the collection for the HTTP request router
 *
 * Contact: dev@phunsites.net
 */

package dohservice

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

// routers defines set of HTTP handler routes
var routers = routes{
	route{
		"Index",
		"GET",
		"/",
		rootIndex,
	},

	route{
		"Status",
		"GET",
		"/status",
		status,
	},

	route{
		"DNSQueryGet",
		strings.ToUpper("Get"),
		"/dns-query",
		DNSQueryGet,
	},

	route{
		"DNSQueryPost",
		strings.ToUpper("Post"),
		"/dns-query",
		DNSQueryPost,
	},
}

// NewRouter initializes an HTTP multiplexer for the webservice
func NewRouter(chanTelemetry chan uint) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routers {
		var handler http.Handler
		//handler = route.HandlerFunc
		handler = httpHandler(route.HandlerFunc, route.Name, chanTelemetry)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

		ConsoleLogger(LogInform, fmt.Sprintf("Registered HTTP handler: method=%s, path=%s", route.Method, route.Pattern), false)
	}

	return router
}

// httpHandler wraps the http request handler and logging routine.
func httpHandler(inner http.Handler, name string, chanTelemetry chan uint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// add some extra verbosity before we handle the request
		ConsoleLogger(LogDebug, fmt.Sprintf("Client Requested URL: %s", r.URL), false)
		ConsoleLogger(LogDebug, fmt.Sprintf("Client Request Headers: %s", r.Header), false)

		// Telemetry: Logging HTTP request type
		chanTelemetry <- TelemetryValues[r.Method]
		ConsoleLogger(LogDebug, fmt.Sprintf("Logging HTTP Telemetry for %s request.", r.Method), false)

		// serve the HTTP request
		inner.ServeHTTP(w, r)

		// Logging HTTP request in verbose mode
		ConsoleLogger(LogInform, fmt.Sprintf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		), false)
	})
}
