package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/validator"
	"gopkg.in/guregu/null.v4"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Device_id   string `json:"device_id"`
		Device_os   string `json:"device_os"`
		Device_name string `json:"device_name"`
	}
	//Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:       input.Name,
		Username:   input.Username,
		Email:      input.Email,
		Activated:  null.BoolFrom(true),
		UserDevice: data.UserDevice{input.Device_id, input.Device_os, input.Device_name}, //TODO: Validate these fields
	}

	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		//TODO 3: Keep LOGS of failed Registration as background task
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmailOrUsername):
			v.AddError("email/username", "a user with this email/username address already exists")
			app.EmailUsernamefailedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	fmt.Println("User", ".....")
	// Add the "account:read" permission for the new user.
	// err = app.models.Permissions.AddForUser(user.ID, "account:read")
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }
	// After the user record has been created in the database, generate a new activation
	// token for the user.
	// token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	//TODO 4: Keep LOGS of successful Registration as background task

	// Call the Send() method on our Mailer, passing in the user's email address,
	// name of the template file, and the User struct containing the new user's data.

	// Launch a goroutine which runs an anonymous function that sends the welcome email.

	// app.background(func() {
	// Run a deferred function which uses recover() to catch any panic, and log an
	// error message instead of terminating the application.
	// As there are now multiple pieces of data that we want to pass to our email
	// templates, we create a map to act as a 'holding structure' for the data. This
	// contains the plaintext version of the activation token for the user, along
	// with their ID.
	// data := map[string]interface{}{
	// 	"activationToken": token.Plaintext,
	// 	"userID":          user.ID,
	// }
	// err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
	// if err != nil {
	// 	// Importantly, if there is an error sending the email then we use the
	// 	// app.logger.PrintError() helper to manage it, instead of the
	// 	// app.serverErrorResponse() helper like before.
	// 	app.logger.Println(err, nil)
	// }

	//	})

	resp_data := map[string]interface{}{"user": user}
	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	env := app.SuccessFormater(resp_data, "Success")

	err = app.writeJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	token := app.GetBearerToken(w, r)
	if token == "" {
		return
	}

	userDetails, err := app.models.Users.GetUserDetailsFromToken(data.ScopeAuthentication, token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):

			app.RecordNotFound(w, r, err)
			return
		case errors.Is(err, data.ErrTransactionPINNotSet):
			app.TransactionPINNotSet(w, r, err)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	env := app.SuccessFormater(userDetails, "Success")

	err = app.writeJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		err1 := errors.New(err.Error() + "from: WriteJSON")
		app.serverErrorResponse(w, r, err1)
		return
	}
}
func (app *application) GetBearerToken(w http.ResponseWriter, r *http.Request) string {
	authHeader := r.Header.Get("Authorization")

	// Check if the header is present
	if authHeader == "" {
		// Handle case where Authorization header is missing
		app.noBearerToken(w, r)
		return ""
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		app.invalidAuthenticationTokenResponse(w, r)
		return ""
	}
	// Extract the actual authentication token from the header parts.
	token := headerParts[1]

	return token

}
