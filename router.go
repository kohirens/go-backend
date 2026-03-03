package backend

import (
	"net/http"
	"path/filepath"
)

// Route Is a function similar to http.HandleFunc, but returns an error.
type Route func(http.ResponseWriter, *http.Request, App) error

// Router Manages looking up a handler function (in the routes map) to respond
// to an HTTP request.
type Router struct {
	// routes A map of endpoints to handler functions. Unlike http.HandleFunc,
	// handlers also return an error. Which can slightly reduce error handling
	// code.
	routes          map[string]Route
	notFoundHandler Route
	h               http.HandlerFunc
}

type RouteManager interface {
	Add(route string, fn Route)
	Find(endpoint string) Route
	NotFound(f Route)
}

func NewRouteManager() RouteManager {
	return &Router{
		routes: make(map[string]Route),
		notFoundHandler: func(w http.ResponseWriter, r *http.Request, _ App) error {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
			return nil
		},
	}
}

// Add An endpoint that maps to a handler function.
func (router *Router) Add(route string, fn Route) {
	router.routes[route] = fn
}

// NotFound Return a 404 response when an endpoint does not map to a handler
// function.
func (router *Router) NotFound(f Route) {
	router.notFoundHandler = f
}

// Find a handler function based on the endpoint that it maps to. When no match
// can be found, then return the 404 handler.
//
//	Supported patterns:
//	* **exact match** - A pattern like `/api/sign-in-with-google` maps a single
//	page to a handler.
//	* **wildcard** - A patter like `*.html` maps any page that ends in `html`
//	to a handler.
func (router *Router) Find(endpoint string) Route {
	if len(router.routes) == 0 {
		Log.Fatf("%v", stderr.NoRoutes)
	}
	// Lookup the handler by endpoint
	fn, ok := router.routes[endpoint]
	if !ok { // or lookup by wildcard and an extension.
		ext := filepath.Ext(endpoint)
		Log.Dbugf("lookup extension *%v pattern", ext)

		fn, ok = router.routes["*"+ext]
		if !ok {
			fn = router.notFoundHandler
		}
	}

	return fn
}
