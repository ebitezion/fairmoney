package main

import (
	"fmt"
	"net/http"
)

// ErrorCode represents an error code with its corresponding meaning.
type ErrorCode struct {
	Code  string
	Value string
}

// Error codes
var (
	// Success codes
	Success        = ErrorCode{"00", "Success"}
	PartialSuccess = ErrorCode{"01", "PartialSuccess"}

	// Validation error codes
	ValidationError                  = ErrorCode{"10", "ValidationError"}
	MissingRequiredField             = ErrorCode{"11", "MissingRequiredField"}
	InvalidFormat                    = ErrorCode{"12", "InvalidFormat"}
	InvalidValue                     = ErrorCode{"13", "InvalidValue"}
	DuplicateEntry                   = ErrorCode{"14", "DuplicateEntry"}
	BVNValidationError               = ErrorCode{"15", "BVNValidationError"}
	ErrExistingAccountHolderResponse = ErrorCode{"16", " ErrExistingAccountHolderError"}
	InvalidTransactionPin            = ErrorCode{"17", " InvalidTransactionPin"}
	// Database error codes
	DatabaseError          = ErrorCode{"20", "DatabaseError"}
	ConnectionError        = ErrorCode{"21", "ConnectionError"}
	QueryError             = ErrorCode{"22", "QueryError"}
	DataIntegrityViolation = ErrorCode{"23", "DataIntegrityViolation"}
	Timeout                = ErrorCode{"24", "Timeout"}
	RecordNotFound         = ErrorCode{"25", "RecordNotFound"}
	FailedApiResponse      = ErrorCode{"26", "FailedApiResponse"}
	TransactionPINNotSet   = ErrorCode{"27", "TransactionPINNotSet"}

	// Network error codes
	NetworkError       = ErrorCode{"30", "NetworkError"}
	ConnectionTimeout  = ErrorCode{"31", "ConnectionTimeout"}
	DNSResolutionError = ErrorCode{"32", "DNSResolutionError"}
	ConnectionRefused  = ErrorCode{"33", "ConnectionRefused"}
	ProxyError         = ErrorCode{"34", "ProxyError"}
	BadRequest         = ErrorCode{"35", "BadRequest"}

	// Server error codes
	ServerError            = ErrorCode{"40", "ServerError"}
	InternalServerError    = ErrorCode{"41", "InternalServerError"}
	ServiceUnavailable     = ErrorCode{"42", "ServiceUnavailable"}
	GatewayTimeout         = ErrorCode{"43", "GatewayTimeout"}
	BandwidthLimitExceeded = ErrorCode{"44", "BandwidthLimitExceeded"}
	MethodNotAllowed       = ErrorCode{"45", "MethodNotAllowed"}

	// Bad request error codes
	BadRequestError         = ErrorCode{"50", "BadRequestError"}
	MalformedRequest        = ErrorCode{"51", "MalformedRequest"}
	InvalidParameters       = ErrorCode{"52", "InvalidParameters"}
	UnsupportedMediaType    = ErrorCode{"53", "UnsupportedMediaType"}
	CSRFTokenMissingInvalid = ErrorCode{"54", "CSRFTokenMissingInvalid"}

	// Authentication error codes
	AuthenticationError     = ErrorCode{"60", "AuthenticationError"}
	UnauthorizedAccess      = ErrorCode{"61", "UnauthorizedAccess"}
	ExpiredToken            = ErrorCode{"62", "ExpiredToken"}
	InvalidToken            = ErrorCode{"63", "InvalidToken"}
	InsufficientPermissions = ErrorCode{"64", "InsufficientPermissions"}
	UnauthorizedSource      = ErrorCode{"65", "UnauthorizedSource"}
	InvalidOTP              = ErrorCode{"66", "InvalidOTP"}

	// Authorization error codes
	AuthorizationError              = ErrorCode{"70", "AuthorizationError"}
	AccessDenied                    = ErrorCode{"71", "AccessDenied"}
	ForbiddenAccess                 = ErrorCode{"72", "ForbiddenAccess"}
	RoleBasedAccessControlViolation = ErrorCode{"73", "RoleBasedAccessControlViolation"}
	UnknownDevice                   = ErrorCode{"74", "UnknownDevice"}

	// File system error codes
	FileSystemError   = ErrorCode{"80", "FileSystemError"}
	FileNotFoundError = ErrorCode{"81", "FileNotFoundError"}
	PermissionDenied  = ErrorCode{"82", "PermissionDenied"}
	DiskFull          = ErrorCode{"83", "DiskFull"}

	// Miscellaneous error codes
	UnknownError         = ErrorCode{"90", "UnknownError"}
	UnsupportedOperation = ErrorCode{"91", "UnsupportedOperation"}
	DeprecatedFeature    = ErrorCode{"92", "DeprecatedFeature"}
	ResourceExhaustion   = ErrorCode{"93", "ResourceExhaustion"}
	CustomError          = ErrorCode{"99", "CustomError"}

	// KYC error codes
	KYCComplete          = ErrorCode{"100", "KYCComplete"}
	KYCEmailPhoneNotVer  = ErrorCode{"101", "KYCEmailPhoneNotVer"}
	KYCPhoneNotVer       = ErrorCode{"102", "KYCPhoneNotVer"}
	KYCNotVerified       = ErrorCode{"103", "KYCNotVerified"}
	KYCAccountUpgradeNot = ErrorCode{"104", "KYCAccountUpgradeNot"}
	KYCNoBVN             = ErrorCode{"105", "KYCNoBVN"}

	//Transfer error codes
	TransferSingleLimitExceeded = ErrorCode{"110", "Transfer amount exceeds single limit"}
	TransferDailyLimitExceeded  = ErrorCode{"111", "Transfer amount exceeds daily limit"}
	InvalidTransferPIN          = ErrorCode{"112", "The transfer PIN is invalid"}
	UnauthorizedAccountNo       = ErrorCode{"113", "The Senders AccountNo is not authorized"}
)

// The logError() method is a generic helper for logging an error message.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Using the interface{} type gives
// us more flexibility over the values that we can include in the response.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, NetworkStatus int, error_msg interface{}, ErrorCode ErrorCode, SpecificError interface{}) {
	// env := envelope{"error": message}
	env := envelope{

		"status_code": ErrorCode.Code,
		"status":      ErrorCode.Value,
		"error_msg":   error_msg,
		"error":       SpecificError,
		"data":        "",
	}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := app.writeJSON(w, NetworkStatus, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// The serverErrorResponse() method will be used when the application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message, ServerError, "")
}

// The notFoundResponse() method will be used to send a 404 Not Found status code
// and JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message, ServiceUnavailable, "")
}

// The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message, MethodNotAllowed, "")
}

// The badRequestResponse() method will be used to send a 400 Bad Request status code
// and JSON response to the client.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error(), BadRequest, "")
}

// The failedValidationResponse() method will be used to send a 422 Unprocessable Entity
// status code and JSON response to the client.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	message := ""

	app.errorResponse(w, r, http.StatusUnprocessableEntity, message, ValidationError, errors)

}

// The failedValidationResponse() method will be used to send a 422 Unprocessable Entity
// status code and JSON response to the client.
func (app *application) EmailUsernamefailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	message := "This email/username is already in use, please log in if it belongs to you."

	app.errorResponse(w, r, http.StatusUnprocessableEntity, message, ValidationError, errors)

}

// The BVNValidationError method will be used to send errors specifically for bvn
func (app *application) BVNValidationError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	message := "BVN validation error."

	app.errorResponse(w, r, http.StatusUnprocessableEntity, message, BVNValidationError, errors)

}

// The editConflictResponse() method will be used to send a 409 Conflict status code
// and JSON response to the client.
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message, DataIntegrityViolation, "")
}

// invalidCredentialsResponse() Appropriate for invalid tokens
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message, UnauthorizedAccess, "")
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message, AuthenticationError, "")
}
func (app *application) invalidOtpResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid otp provided"
	app.errorResponse(w, r, http.StatusUnauthorized, message, InvalidOTP, "")
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "You must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message, UnauthorizedAccess, "")
}
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "Your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message, AuthorizationError, "")
}

// Permissions
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message, ForbiddenAccess, "")
}
func (app *application) InvalidTransactionPin(w http.ResponseWriter, r *http.Request) {
	message := "Transaction pin provided is invalid"
	app.errorResponse(w, r, http.StatusForbidden, message, InvalidTransactionPin, "")
}

// API Key
func (app *application) invalidKey(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("WWW-Authenticate", "Api-Key")
	message := "Invalid API Key"
	app.errorResponse(w, r, http.StatusUnauthorized, message, InvalidToken, "")
}

func (app *application) noAPIKey(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("WWW-Authenticate", "Api-Key")
	message := "API Key is required"
	app.errorResponse(w, r, http.StatusUnauthorized, message, InvalidToken, "")
}
func (app *application) noBearerToken(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("WWW-Authenticate", "Api-Key")
	message := "Bearer Token is required"
	app.errorResponse(w, r, http.StatusUnauthorized, message, InvalidToken, "")
}

// RecordNotFound Record not found
func (app *application) RecordNotFound(w http.ResponseWriter, r *http.Request, err error) {
	message := "Record not found"
	app.errorResponse(w, r, http.StatusUnauthorized, message, RecordNotFound, err)
}

func (app *application) TransactionPINNotSet(w http.ResponseWriter, r *http.Request, err error) {
	message := "Transaction PIN Not Set"
	app.errorResponse(w, r, http.StatusUnauthorized, message, TransactionPINNotSet, err)
}
func (app *application) DatabaseError(w http.ResponseWriter, r *http.Request, err error) {
	message := "Database Error"
	app.errorResponse(w, r, http.StatusUnauthorized, message, DatabaseError, err)
}
func (app *application) FailedApiResponse(w http.ResponseWriter, r *http.Request, err interface{}) {
	message := "Request Failed" //"Request to subsequent api failed"
	app.errorResponse(w, r, http.StatusUnauthorized, message, FailedApiResponse, err)
}
func (app *application) FailedTransferResponse(w http.ResponseWriter, r *http.Request, err interface{}) {
	//	message := "Transfer Request Failed" //"Request to subsequent api failed"
	//response
	app.errorResponse(w, r, http.StatusUnauthorized, err, FailedApiResponse, err)
}

// Different Source Login
func (app *application) DifferentSourceResponse(w http.ResponseWriter, r *http.Request, err interface{}) {
	message := "Source Error: This login details is from source-Spectrumpay, device validation is required to continue"
	app.errorResponse(w, r, http.StatusUnauthorized, message, UnauthorizedSource, err)
}

func (app *application) SendersAccountNoUnauthorizedResponse(w http.ResponseWriter, r *http.Request, err interface{}) {
	message := "The account Number is not authorized to be used for this user"
	app.errorResponse(w, r, http.StatusUnauthorized, message, UnauthorizedAccountNo, err)
}

// KYC Errors

// AddressNotVerifiedResponse Address Not Verified Error Response
func (app *application) AddressNotVerifiedResponse(w http.ResponseWriter, r *http.Request) {
	message := "Home Address, Not Found"
	app.errorResponse(w, r, http.StatusUnauthorized, message, KYCNotVerified, "")
}

// StatusNotVerifiedResponse Status Not Verified
func (app *application) StatusNotVerifiedResponse(w http.ResponseWriter, r *http.Request) {
	message := "PhoneNo or Email, could not be validated"
	app.errorResponse(w, r, http.StatusUnauthorized, message, KYCNotVerified, "")
}

// StatusNotVerifiedResponse Status Not Verified
func (app *application) BVNNotVerifiedResponse(w http.ResponseWriter, r *http.Request, userID int64) {
	message := "BVN not validated"
	data := map[string]int64{"user_id": userID}
	app.errorResponse(w, r, http.StatusUnauthorized, message, KYCNotVerified, data)
}
func (app *application) ErrExistingAccountHolderResponse(w http.ResponseWriter, r *http.Request) {
	message := "Existing account number Has not be validated, submit OTP to validate."
	app.errorResponse(w, r, http.StatusUnauthorized, message, ErrExistingAccountHolderResponse, "")
}
