package main

import (
	"errors"
	"net/http"

	"greenlight.mazavrbazavr.ru/internal/data"
	"greenlight.mazavrbazavr.ru/internal/validator"
)

// Handler for the "POST /v1/users" endpoint. Method of the application struct.
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Launch a goroutine which runs an anonymous function that sends the welcome email.
	// The code in this background goroutine will be executed concurrently with the
	// subsequent code in our registerUserHandler, which means we are not waiting for
	// the email to be sent before we return a JSON response to the client. Most likely,
	// the background goroutine will still be executing its code long after the
	// registerUserHandler has returned. Use the background helper to execute an
	// anonymous function that sends the welcome email.
	app.background(func() {
		// Call the Send() method on our Mailer, passing in the user's email address,
		// name of the template file, and the User struct containing the new user's data
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper like before. This is because by the time
			// we encounter the errors, the client will probably have already been sent
			// a 202 Accepted response by our writeJSON() helper.
			app.logger.PrintError(err, nil)
		}
	})

	// Write a JSON response containing the user data along with a 202 Accepted status
	// code. This status code indicates that the request has been accepted for
	// processing, but the processing has not been completed.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
