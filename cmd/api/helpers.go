package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"maps"

	"github.com/julienschmidt/httprouter"
)

// A helper. Retrieves the "id" URL parameter from the current request context,
// then converts it to an integer and returns it. If the operation isn't successful,
// returns 0 and an error. A method of the application struct.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// When httprouter is parsing a request, any interpolated URL parameters will be
	// stored in the request context. We use the ParamsFromContext() function to
	// retrieve a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(r.Context())

	// Use the ByName() method to get the value of the "id" parameter from
	// the slice. We are retrieving unique positive integer ID, but
	// the value returned by ByName() is always a string. So we try to convert it to a
	// base 10 integer (with a bit size of 64). If the parameter couldn't be converted,
	// or is less than 1, we know the ID is invalid so we use the http.NotFound()
	// function to return a 404 Not Found response.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// Define an envelope type.
type envelope map[string]any

// A helper for sending responses. This takes the destination http.ResponseWriter,
// the HTTP status code to send, the enveloped data to encode to JSON, and a header map
// containing any additional HTTP headers we want to include in the response.
// A method of the application struct.
func (app *application) writeJSON(
	w http.ResponseWriter,
	status int,
	data envelope,
	headers http.Header,
) error {
	// Encode the data to JSON, returning the error if there was one.
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded
	// JSON. Here we use no line prefix ("") and tab indents ("\t") for each element.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	maps.Copy(w.Header(), headers)

	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
