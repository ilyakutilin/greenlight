package main

import (
	"fmt"
	"net/http"
)

// A generic helper for logging an error message. Method of the application struct.
func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// A generic helper for sending JSON-formatted error messages to the client
// with a given status code. Method of the application struct.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// Used when our application encounters an unexpected problem at runtime.
// It logs the detailed error message, then uses the errorResponse() helper
// to send a 500 Internal Server Error status code and JSON response
// (containing a generic error message) to the client. Method of the application struct.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// Used to send a 404 Not Found status code and JSON response to the client.
// Method of the application struct.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// Used to send a 405 Method Not Allowed status code and JSON response to the client.
// Method of the application struct.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// Used to send a 400 Bad Request status code and JSON response to the client.
// Method of the application struct.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Used to send a 422 Unprocessable Entity status code and JSON response to the client.
// Method of the application struct.
// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in the Validator type.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// Used to send a 409 Conflict status code and JSON response to the client.
// Method of the application struct.
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

// Used to send a 429 Too Many Requests status code if the rate limit is exceeded
// and JSON response to the client. Method of the application struct.
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

// Used to send a 401 Unauthorized status code and JSON response to the client.
// Method of the application struct.
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Used to send a 401 Unauthorized status code and JSON response to the client.
// Method of the application struct.
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Used to send a 401 Unauthorized status code and JSON response to the client.
// Method of the application struct.
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Used to send a 403 Forbidden status code and JSON response to the client.
// Method of the application struct.
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
