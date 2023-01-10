// this is used to hold models related to HTTP routes
package models

import (
	"github.com/gin-gonic/gin"
)

// Route model used to strictly define a route and its attributes
type Route struct {
	Path     string
	Handler  gin.HandlerFunc
	HttpVerb string
}
