// Package handlers contains the full set of handler functions and routes
// supported by the http api
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/mistermoe/go-starter/framework"
	"github.com/mistermoe/go-starter/middleware"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) *framework.App {
	app := framework.NewApp(shutdown, middleware.Logger(log), middleware.Errors(log), middleware.Metrics(), middleware.Panics(log))

	readiness := readiness{
		log: log,
	}

	// attach all handlers here
	app.Handle(http.MethodGet, "/health", health)
	app.Handle(http.MethodGet, "/readiness", readiness.handle)

	return app
}
