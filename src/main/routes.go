// all route definitions should be placed here
// no complex logic should live within the route handler functions in this source - instead this should be delegated to
// a "service"
// route handlers should simply receive requests, extract input, call a service and return the response
// lastly, all routes should return sanitized errors to the user, actual backend errors should be logged for
// investigation later
package main

import (
	"fmt"
	"net/http"

	"harmonia-example.io/src/controllers"
	"harmonia-example.io/src/models"
	"harmonia-example.io/src/services/config"
	"harmonia-example.io/src/services/git"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// GetRoutes returns an array of `models.Route` representing all available routes
func GetRoutes() []models.Route {
	return []models.Route{
		// health routes
		{
			Path:     "/health",
			Handler:  getHealth,
			HttpVerb: http.MethodGet,
		},
		// swagger docs routes
		{
			Path:     "/",
			Handler:  swaggerRedirect,
			HttpVerb: http.MethodGet,
		},
		{
			Path:     "/index.html",
			Handler:  swaggerRedirect,
			HttpVerb: http.MethodGet,
		},
		{
			Path:     "/docs",
			Handler:  swaggerRedirect,
			HttpVerb: http.MethodGet,
		},
		{
			Path:     "/swagger/*any",
			Handler:  swagger,
			HttpVerb: http.MethodGet,
		},
		// rfc routes
		{
			Path:     "/submitRequest",
			Handler:  submitRequest,
			HttpVerb: http.MethodPost,
		},
		{
			Path:     "/updateRequest",
			Handler:  updateRequest,
			HttpVerb: http.MethodPost,
		},
		{
			Path:     "/reviewRequest",
			Handler:  reviewRequest,
			HttpVerb: http.MethodPost,
		},
		{
			Path:     "/mergeRequest",
			Handler:  mergeRequest,
			HttpVerb: http.MethodPost,
		},
		{
			Path:     "/loadRequest",
			Handler:  loadRequest,
			HttpVerb: http.MethodPost,
		},
		{
			Path:     "/status",
			Handler:  status,
			HttpVerb: http.MethodPost,
		},
		{
			Path:     "/getRfcs",
			Handler:  getRfcs,
			HttpVerb: http.MethodPost,
		},
		{
			Path:     "getRfcContents",
			Handler:  getRfcContents,
			HttpVerb: http.MethodPost,
		},
	}
}

// @Summary Health check
// @Description Simple health check used to determine if the service is healthy and responding
// @Tags Health
// @Produce json
// @Success 200 {object} models.Healthy "healthy response"
// @Router /health [get]
// getHealth returns a simple health check used to determine if the service is healthy and responding
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, &models.Healthy{Message: "healthy"})
}

// you don't see any openapi comments here because this is swagger itself
// swaggerRedirect redirects request to the swagger docs page
func swaggerRedirect(c *gin.Context) {
	c.Redirect(http.StatusFound, "/swagger/index.html")
}

// you don't see any openapi comments here because this is swagger itself
// swagger returns the swagger docs page with proper session information
func swagger(c *gin.Context) {
	docsHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	docsHandler(c)
}

// @description submit RFC
// @Tags RFC
// @Accept json
// @Produce json
// @Param RFC body models.RFC true "RFC JSON"
// @Response 200 {object} models.RFCIdentifier
// @Response 400 {object} models.Error
// @Response 403 {object} models.Error
// @Response 500 {object} models.Error
// @Router /submitRequest [post]
// submitRequest handles submitting an initial schema change request
func submitRequest(c *gin.Context) {
	RFC := new(models.RFC)
	// ensure the incoming request body conforms to the RFC model
	if err := c.ShouldBindBodyWith(RFC, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	} else {
		// initialize params for controller
		if accessToken, err := config.GetToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no token"})
		} else {
			// establish git client
			if github, err := git.NewGitHub(c, *accessToken); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git"})
			} else {
				// submit RFC
				if identifier, err := controllers.SubmitRequest(c, github, RFC); err != nil {
					c.JSON(http.StatusInternalServerError, &models.Error{Error: "Request creation error occurred"})
				} else {
					c.JSON(http.StatusOK, &models.RFCIdentifier{RFCIdentifier: *identifier})
				}
			}
		}
	}
}

// @description update RFC
// @Tags RFC
// @Accept json
// @Produce json
// @Param Update body models.Update true "Update JSON"
// @Response 200 {object} models.RFCIdentifier
// @Response 400 {object} models.Error
// @Response 403 {object} models.Error
// @Response 500 {object} models.Error
// @Router /updateRequest [post]
// updateRequest handles updating an existing schema change request
func updateRequest(c *gin.Context) {
	update := new(models.Update)
	// ensure the incoming request body conforms to the Update model
	if c.ShouldBindBodyWith(update, binding.JSON) == nil {
		// <this is a good point to augment logger with request metadata> //
		// initialize params for controller
		if accessToken, err := config.GetToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no token"})
		} else {
			// establish git client
			if github, err := git.NewGitHub(c, *accessToken); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git"})
			} else {
				// submit update request
				if identifier, err := controllers.UpdateRequest(c, github, update); err != nil {
					c.JSON(http.StatusInternalServerError, &models.Error{Error: "update request error occurred"})
				} else {
					c.JSON(http.StatusOK, &models.RFCIdentifier{RFCIdentifier: *identifier})
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	}
}

// @description review RFC
// @Tags RFC
// @Accept json
// @Produce json
// @Param Review body models.Review true "Review JSON"
// @Response 200 {object} models.Success
// @Response 400 {object} models.Error
// @Response 403 {object} models.Error
// @Response 500 {object} models.Error
// @Router /reviewRequest [post]
// reviewRequest handles all review actions: approval, requesting changes, or commenting. Requesting changes blocks
// merging, while the other events do not.
func reviewRequest(c *gin.Context) {
	review := new(models.Review)
	// ensure the incoming request body conforms to the Review model
	if c.ShouldBindBodyWith(review, binding.JSON) == nil {
		// <this is a good point to augment logger with request metadata> //
		// initialize params for controller
		if accessToken, err := config.GetToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no token"})
		} else {
			if machineAccessToken, err := config.GetMachineToken(); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{
					Error: "Configuration error occurred - no machine token"})
			} else {
				// establish git clients
				if github, err := git.NewGitHub(c, *accessToken); err != nil {
					c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git"})
				} else {
					if githubMachine, err := git.NewGitHub(c, *machineAccessToken); err != nil {
						c.JSON(http.StatusInternalServerError, &models.Error{
							Error: "Service error occurred - Git machine"})
					} else {
						// submit review
						if message, err := controllers.ReviewRequest(c, github, githubMachine, review); err != nil {
							c.JSON(http.StatusInternalServerError, &models.Error{
								Error: "Review submission error occurred"})
						} else {
							c.JSON(http.StatusOK, &models.Success{Success: *message})
						}
					}
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	}
}

// @description merge RFC
// @Tags RFC
// @Accept json
// @Produce json
// @Param Merge body models.Merge true "Merge JSON"
// @Response 200 {object} models.Success
// @Response 400 {object} models.Error
// @Response 403 {object} models.Error
// @Response 500 {object} models.Error
// @Router /mergeRequest [post]
// mergeRequest handles merging the given RFC and tagging it for tracking
func mergeRequest(c *gin.Context) {
	merge := new(models.Merge)
	// ensure the incoming request body conforms to the Merge model
	if c.ShouldBindBodyWith(merge, binding.JSON) == nil {
		// <this is a good point to augment logger with request metadata> //
		// initialize params for controller
		if machineAccessToken, err := config.GetMachineToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no machine token"})
		} else {
			// establish git clients
			if github, err := git.NewGitHub(c, *machineAccessToken); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git machine"})
			} else {
				// submit merge request
				if message, err := controllers.MergeRequest(c, github, merge); err != nil {
					c.JSON(http.StatusInternalServerError, &models.Error{Error: "Merge error occurred"})
				} else {
					c.JSON(http.StatusOK, &models.Success{Success: *message})
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	}
}

// @description load RFC
// @Tags RFC
// @Accept json
// @Produce json
// @Param Load body models.Load true "Load JSON"
// @Response 200 {object} models.Success
// @Response 400 {object} models.Error
// @Response 403 {object} models.Error
// @Response 500 {object} models.Error
// @Router /loadRequest [post]
// loadRequest handles loading the given RFC into the underlying datastore
func loadRequest(c *gin.Context) {
	load := new(models.Load)
	// ensure the incoming request body conforms to the Load model
	if c.ShouldBindBodyWith(load, binding.JSON) == nil {
		// <this is a good point to augment logger with request metadata> //
		// initialize params for controller
		if accessToken, err := config.GetToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no token"})
		} else {
			// establish git client
			if github, err := git.NewGitHub(c, *accessToken); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git"})
			} else {
				// submit load request
				// this only captures setup errors because the actual load is handled asynchronously
				if err = controllers.LoadRequest(c, github, load); err != nil {
					c.JSON(http.StatusInternalServerError, &models.Error{Error: "Load request error occurred"})
				} else {
					c.JSON(http.StatusOK, &models.LoadRequest{Message: fmt.Sprintf(
						"Submitted load request for RFC %s.You may query the load status through the /status endpoint.",
						load.RFCIdentifier)})
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	}
}

// @description status check
// @Tags RFC
// @Accept json
// @Produce json
// @Param Status body models.Status true "Load Status JSON"
// @Response 200 {object} models.Success
// @Response 400 {object} models.Error
// @Response 500 {object} models.Error
// @Router /status [post]
// status handles retrieving the load status of the given RFC
func status(c *gin.Context) {
	status := new(models.Status)
	// ensure the incoming request body conforms to the Status model
	if c.ShouldBindBodyWith(status, binding.JSON) == nil {
		// <this is a good point to augment logger with request metadata> //
		// operate as machine for status requests
		if machineAccessToken, err := config.GetMachineToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no machine token"})
		} else {
			// establish git clients
			if github, err := git.NewGitHub(c, *machineAccessToken); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git machine"})
			} else {
				// submit status request
				if loadStatus, err := controllers.Status(c, github, status); err != nil {
					c.JSON(http.StatusInternalServerError, &models.Error{Error: "Status error occurred"})
				} else {
					if loadStatus == nil {
						c.JSON(http.StatusOK, &models.StatusResponse{Status: "none"})
					} else {
						c.JSON(http.StatusOK, &models.StatusResponse{Status: *loadStatus})
					}
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	}
}

// @description get submitted RFCs
// @Tags RFC
// @Accept json
// @Produce json
// @Param Query body models.GetRfcs true "Query JSON"
// @Response 200 {object} models.RFCs
// @Response 400 {object} models.Error
// @Response 403 {object} models.Error
// @Response 500 {object} models.Error
// @Router /getRfcs [post]
// getRfcs queries the datastore for all RFCs with a given state, paginated output
func getRfcs(c *gin.Context) {
	request := new(models.GetRfcs)
	// ensure the incoming request body conforms to the request model
	if c.ShouldBindBodyWith(request, binding.JSON) == nil {
		// <this is a good point to augment logger with request metadata> //
		// operate as machine for credentials
		if machineAccessToken, err := config.GetMachineToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no machine token"})
		} else {
			// establish git clients
			if github, err := git.NewGitHub(c, *machineAccessToken); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git machine"})
			} else {
				// submit status request
				if results, err := controllers.GetRfcs(c, github, request); err != nil {
					fmt.Println(err)
					c.JSON(http.StatusInternalServerError, &models.Error{Error: "Error occurred when retrieving RFCs"})
				} else {
					count := len(results)
					if results == nil {
						c.JSON(http.StatusOK, &models.RFCs{RFCs: []map[string]string{}, Count: &count})
					} else {
						c.JSON(http.StatusOK, &models.RFCs{RFCs: results, Count: &count})
					}
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	}
}

// @description get submitted RFC contents
// @Tags RFC
// @Accept json
// @Produce json
// @Param RFC body models.GetRfcContents true "Query JSON"
// @Response 200 {object} models.RFCContents
// @Response 400 {object} models.Error
// @Response 403 {object} models.Error
// @Response 500 {object} models.Error
// @Router /getRfcContents [post]
// getRfcContents retrieves the body of a given RFC
func getRfcContents(c *gin.Context) {
	request := new(models.GetRfcContents)
	// ensure the incoming request body conforms to the request model
	if c.ShouldBindBodyWith(request, binding.JSON) == nil {
		// <this is a good point to augment logger with request metadata> //
		// operate as machine for status requests
		if machineAccessToken, err := config.GetMachineToken(); err != nil {
			c.JSON(http.StatusInternalServerError, &models.Error{Error: "Configuration error occurred - no machine token"})
		} else {
			// establish git clients
			if github, err := git.NewGitHub(c, *machineAccessToken); err != nil {
				c.JSON(http.StatusInternalServerError, &models.Error{Error: "Service error occurred - Git machine"})
			} else {
				// submit status request
				if contents, err := controllers.GetRfcContents(c, github, request); err != nil {
					c.JSON(http.StatusInternalServerError, &models.Error{
						Error: fmt.Sprintf("Error occurred when querying contents for RFC #%v", request.RFCIdentifier)})
				} else {
					if contents == nil {
						c.JSON(http.StatusOK, &models.RFCContents{Body: ""})
					} else {
						c.JSON(http.StatusOK, &models.RFCContents{Body: *contents})
					}
				}
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, &models.Error{Error: "Malformed request received"})
	}
}
