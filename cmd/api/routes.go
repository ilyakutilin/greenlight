package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Defines routes used by the project. A method of the applcation struct.
// Returns a http.Handler instead of a *httprouter.Router due to the use of
// the panic recovery middleware.
func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for endpoints.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	// Return the httprouter instance.
	// Middlewares:
	// - Panic recovery middleware;
	// - Rate limit middleware - comes after our panic recovery middleware (so that any
	//   panics in rateLimit() are recovered), but otherwise we want it to be used as
	//   early as possible to prevent unnecessary work for our server.
	return app.recoverPanic(app.rateLimit(router))
}
