package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/mock"
	"github.com/shopspring/decimal"
)

func (app application) PaymentInitiation(w http.ResponseWriter, r *http.Request) {
	//get senders details

	// Retrieve token and validate
	token := app.GetBearerToken(w, r)
	if token == "" {
		return
	}

	// Read user input
	var Transaction data.Payment
	err := app.readJSON(w, r, &Transaction)
	if err != nil {
		err = errors.New(err.Error() + "from: Transaction ")
		app.badRequestResponse(w, r, err)
		return
	}

	// TODO: Validate user input
	//v := validator.New()
	// if data.ValidatePayment(v, &InternalTransaction); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }

	// Retrieve user details
	userDetail, err := app.models.Users.GetUserDetailsFromToken(data.ScopeAuthentication, token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.RecordNotFound(w, r, err)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	fmt.Println("User detail", userDetail)
	//simulate third party payment
	// response, err := paythirdparty(Tra) 
	// if err != nil {

	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	// //if response successful, debit or credit user based on request type
	// if response["status"] != "success" {
	// 	app.serverErrorResponse(w, r, err)
	// }

	// if Transaction.Type == "credit" {
	// 	//perform credit that is increase user account balance
	// }

	// if Transaction.Type == "debit" {
	// 	//perform debit that is decrease user account balance
	// }

	//save transaction

	//send response to user
}

func paythirdparty(amt decimal.Decimal, accountId, ref string) {
	// Start the mock server
	baseURL := mock.StartMockServer()

	// Define the request body
	requestBody := map[string]interface{}{
		"account_id": accountId,
		"reference":  ref,
		"amount":     amt,
	}

	// Convert request data to JSON
	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	// Perform the POST request
	resp, err := http.Post(baseURL+"/third-party/payments", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read and print the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response Body:", string(responseBody))
}

/*

// Retrieve token and validate
	token := app.GetBearerToken(w, r)
	if token == "" {
		return
	}

	// Read user input
	var InternalTransaction data.InternalTransaction
	err := app.readJSON(w, r, &InternalTransaction)
	if err != nil {
		err = errors.New(err.Error() + "from: InternalTransaction ")
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate user input
	if data.ValidateInternalTransfer(v, &InternalTransaction); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve user details
	userDetail, err := app.models.Users.GetUserDetailsAndPINFromToken(data.ScopeAuthentication, token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.RecordNotFound(w, r, err)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// Verify PIN
	if !VerifyPIN(userDetail.PIN, InternalTransaction.PIN) {
		app.errorResponse(w, r, http.StatusForbidden, InvalidTransferPIN.Code, InvalidTransferPIN, InvalidTransferPIN.Value)
		return
	}

	// Unmarshal user details JSON
	var userDetails UserDetails
	err = json.Unmarshal([]byte(userDetail.Limits), &userDetails)
	if err != nil {
		fmt.Println("Error unmarshalling user details JSON:", err)
		return
	}

	// Unmarshal counter JSON
	var counter Counter
	err = json.Unmarshal([]byte(userDetail.Counter), &counter)
	if err != nil {
		fmt.Println("Error unmarshalling counter JSON:", err)
		return
	}

	// Get single and daily limits for transfers
	singleLimit := userDetails.Transfers.Single
	dailyLimit := userDetails.Transfers.Daily

	// Convert InternalTransaction.Amount to decimal
	amount, err := decimal.NewFromString(InternalTransaction.Amount)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Calculate new counter value
	count := counter.Transfers.Add(amount)

	// Check transfer limits
	if amount.GreaterThan(singleLimit) {
		app.errorResponse(w, r, http.StatusForbidden, TransferSingleLimitExceeded.Code, TransferSingleLimitExceeded, "Transfer amount exceeds single limit")
		return
	}

	if count.GreaterThan(dailyLimit) {
		app.errorResponse(w, r, http.StatusForbidden, TransferDailyLimitExceeded.Code, TransferDailyLimitExceeded, "Transfer amount exceeds daily limit")
		return
	}

	// Prepare for transfer
	// InternalTransaction.SendersAccountNo = userDetail.AccountNumber
	// validate sendersaccount no
	var acct data.AccountNumber
	acct.AccountNumber = userDetail.AccountNumber
	isValidSenderAccountNo, err := ValidateSenderAccountno(&acct, InternalTransaction.SendersAccountNo, userDetail.AccountNumber)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("Cannot verify Account at the moment, try again later"))
		return
	}
	if !isValidSenderAccountNo {
		//account Invalid Error
		app.SendersAccountNoUnauthorizedResponse(w, r, " Senders AccountNo Unauthorized Response ")
		return
	}
	val := orbit.NewTransfer(InternalTransaction.SendersAccountNo, InternalTransaction.ReceiverAccountNo, InternalTransaction.Amount, "0", "0", InternalTransaction.Narration)

	// Initiate transfer
	resp, err := orbit.SpectrumTransfer(val)
	if err != nil {
		// Handle transfer failure
		app.handleTransferFailure(w, r, userDetail, InternalTransaction, count.String(), err)
		app.FailedTransferResponse(w, r, err)
		return
	}

	// Handle transfer success
	app.handleTransferSuccess(w, r, resp, userDetail, InternalTransaction, count.String())



*/
