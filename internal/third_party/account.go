package thirdparty

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/go-resty/resty/v2"
)

type Account struct {
	Surname     string
	FirstName   string
	HomeAddress string
	City        string
	PhoneNumber string
	BVN         string
}

const (
	ACCOUNT_CREATION_URL                = ""
	ACCOUNT_HISTORY_WITH_PAGINATION_URL = ""
	BALANCE_ENQUIRY_URL                 = ""

	ACCOUNT_DETAILS_URL = ""
)

// initiatize and populate the Account_creation struct
func New(Surname, FirstName, HomeAddress, City, PhoneNumber, BVN string) *Account {
	return &Account{
		Surname:     Surname,
		FirstName:   FirstName,
		HomeAddress: HomeAddress,
		City:        City,
		PhoneNumber: PhoneNumber,
		BVN:         BVN,
	}
}

// CreateBankAccount creates an account for app users using orbit API
func CreateBankAccount(account *Account) (*resty.Response, error) {
	requestData := map[string]interface{}{
		"surname":     account.Surname,
		"firstName":   account.FirstName,
		"homeAddress": account.HomeAddress,
		"city":        account.City,
		"phoneNumber": account.PhoneNumber,
		"bvn":         account.BVN,
	}

	// Making HTTP request using resty library
	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(requestData).
		Post(ACCOUNT_CREATION_URL)

	return response, err
}

// Transaction represents a single transaction in the account history
type Transaction struct {
	AccountType       string  `json:"accountType"`
	TransactionDesc   string  `json:"transactionDesc"`
	TransactionAmount string  `json:"transactionAmount"`
	BalanceAfter      string  `json:"balanceAfter"`
	TransactionDate   string  `json:"transactionDate"`
	TransactionType   string  `json:"transactionType"`
	TransactionRef    *string `json:"transactionRef"`
}

// Response represents the API response structure
type Response struct {
	ResponseCode    string        `json:"responseCode"`
	ResponseMessage string        `json:"responseMessage"`
	AccountHistory  []Transaction `json:"accountHistory"`
}

func BalanceEnquiryApi(data *data.AccountNumber) (*resty.Response, error) {
	requestData := map[string]interface{}{
		"accountNumber": data.AccountNumber,
	}

	// Making HTTP request using resty library
	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(requestData).
		Post(BALANCE_ENQUIRY_URL)

	return response, err

}

func AccountDetailsApi(data *data.AccountNumber) (*resty.Response, error) {

	requestData := map[string]interface{}{
		"accountNumber": data.AccountNumber,
	}

	// Making HTTP request using resty library
	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(requestData).
		Post(ACCOUNT_DETAILS_URL)

	return response, err

}
func AccountDetails2Api(account string) (*resty.Response, error) {

	requestData := map[string]interface{}{
		"accountNumber": account,
	}

	// Making HTTP request using resty library
	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(requestData).
		Post(ACCOUNT_DETAILS_URL)

	return response, err

}


func CreateBankAccount2(account *Account) ([]byte, error) {

	requestData := map[string]string{
		"surname":     account.Surname,
		"firstName":   account.FirstName,
		"homeAddress": account.HomeAddress,
		"city":        account.City,
		"phoneNumber": account.PhoneNumber,
		"bvn":         account.BVN,
	}
	// Convert request data to JSON
	reqBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	// Perform HTTP POST request using the global HTTP client
	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(reqBody).
		Post(ACCOUNT_CREATION_URL)
	if err != nil {
		return nil, err
	}

	// Check for errors in the response status code
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("Unexpected status code")
	}

	// Return the response body
	return resp.Body(), nil

}
