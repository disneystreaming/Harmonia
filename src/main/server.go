// Package main
// the main entrypoint of the API - it is in charge of delegating middleware and routes and ultimately running the
// gin server
package main

import (
	"net/http"

	"harmonia-example.io/src/main/docs"
	"harmonia-example.io/src/models"

	"github.com/gin-gonic/gin"
)

// harmoniaVersion is passed in from build and is used for swagger display
var harmoniaVersion string

// @title Harmonia
// @description Harmonia is a service for processing and accepting requests for schema changes

// @contact.name <your-contact-here>
// @contact.url <your-contact-here>
// @contact.email <your-contact-here>

// @license.name MIT

// @schemes https http

// main handles initializing the application and ultimately serving it
func main() {
	// initialize the gin engine
	engine := gin.Default()

	// < this is a good place to bind middleware > //

	// configure dynamic swagger documentation
	configureSwagger(harmoniaVersion)

	// create routes for app
	bindRoutes(engine, GetRoutes())

	// run application
	engine.Run(":8080")
}

// configureSwagger sets dynamic swagger configuration that is version/environment dependent
func configureSwagger(ver string) {
	// set display version (this is what is listed at the top of the swagger page)
	docs.SwaggerInfo.Version = ver
	// set host (where requests are routed to)
	docs.SwaggerInfo.Host = "localhost:8080"

}

// bindRoutes iterates over the provided routes array and adds the proper handlers to the given engine
func bindRoutes(engine *gin.Engine, routes []models.Route) {
	for _, route := range routes {
		// GET routes
		if route.HttpVerb == http.MethodGet {
			if route.Handler != nil {
				engine.GET(route.Path, route.Handler)
			}
			// POST ROUTES
		} else if route.HttpVerb == http.MethodPost {
			if route.Handler != nil {
				engine.POST(route.Path, route.Handler)
			}
		}
	}
}
