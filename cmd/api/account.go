package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/nuban"

	//	fairmoney "github.com/ebitezion/backend-framework/internal/third_party"
	"github.com/ebitezion/backend-framework/internal/validator"
)

//CreateBankAccount creates bank account no and associates with a user
func (app *application) CreateBankAccount(w http.ResponseWriter, r *http.Request) {
	//get user request body input
	//mdata
	user_details := data.AccountDetails{}
	Account := data.Account{}
	err := app.readJSON(w, r, &Account)
	if err != nil {
		err1 := errors.New(err.Error() + "from: readJSON ")
		app.badRequestResponse(w, r, err1)
		return
	}

	//TODO : Validation of values
	v := validator.New()
	// Validate the Account struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateAccountData(v, &Account); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Details need for user_details table submission

	//get the user ID
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {

		r = app.contextSetUser(r, data.AnonymousUser)
	}
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	// Extract the actual authentication token from the header parts.
	token := headerParts[1]

	user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//Invoke package to generate account no
	nubanGenerator := nuban.NewNUBANGenerator()
	nuban := nubanGenerator.GenerateNUBAN()
	user_details.Account_number = nuban

	user_details.User_id = int(user.ID)

	//As a background task save these details on success to user_details table on DB
	app.background(func() {
		err := app.models.AccountModel.SaveCreatedAccountNo(&user_details)
		if err != nil {
			fmt.Println("Error from Saving Account No", err)
			return
		}
	})

	env := app.SuccessFormater(user_details, "Success")
	err = app.writeJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		err1 := errors.New(err.Error() + "from: WriteJSON")
		app.serverErrorResponse(w, r, err1)
		return
	}

}


