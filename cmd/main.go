package main

import (
	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/middleware"

	"log"
	"net/http"
)

type (
	App struct {
		htmx *htmx.HTMX
	}

	route struct {
		path    string
		handler http.Handler
	}
)

func main() {
	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()

	// wrap routes with the middleware
	wrapRoutes(mux, middleware.MiddleWare, []route{
		{path: "/", handler: http.HandlerFunc(app.Home)},
	})

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/index.html")
}


// wrapRoutes takes a mux, middleware, and a list of routes to apply the middleware to.
func wrapRoutes(mux *http.ServeMux, mw func(http.Handler) http.Handler, routes []route) {
	for _, r := range routes {
		mux.Handle(r.path, mw(r.handler))
	}
}

