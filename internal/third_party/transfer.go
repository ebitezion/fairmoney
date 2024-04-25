package thirdparty

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	// "github.com/shopspring/decimal"
)

type InternalTransaction struct {
	SendersAccountNo  string
	ReceiverAccountNo string
	Amount            string
	Charge            string
	Tax               string
	Narration         string
}

type EasyPayTransfer struct {
}

type FundsTransferCreditRequest struct {
	NameEnquiryRef                    string `json:"nameEnquiryRef"`
	DestinationInstitutionCode        string `json:"destinationInstitutionCode"`
	ChannelCode                       string `json:"channelCode"`
	BeneficiaryAccountName            string `json:"beneficiaryAccountName"`
	BeneficiaryAccountNumber          string `json:"beneficiaryAccountNumber"`
	BeneficiaryBankVerificationNumber string `json:"beneficiaryBankVerificationNumber"`
	BeneficiaryKYCLevel               string `json:"beneficiaryKYCLevel"`
	OriginatorAccountName             string `json:"originatorAccountName"`
	OriginatorAccountNumber           string `json:"originatorAccountNumber"`
	OriginatorBankVerificationNumber  string `json:"originatorBankVerificationNumber"`
	OriginatorKYCLevel                string `json:"originatorKYCLevel"`
	TransactionLocation               string `json:"transactionLocation"`
	Narration                         string `json:"narration"`
	PaymentReference                  string `json:"paymentReference"`
	Amount                            string `json:"amount"`
}

type NameEnquiryRequest struct {
	AccountNumber              string `json:"accountNumber"`
	DestinationInstitutionCode string `json:"destinationInstitutionCode"`
	ChannelCode                string `json:"channelCode"`
}

// Internal Transfer
const (
	ORBIT_INTERNAL_DEBIT_TRANSFER  = ""
	ORBIT_EXTERNAL_DEBIT_TRANSFER  = ""
	ORBIT_INTERNAL_CREDIT_TRANSFER = ""
)

// External Transfer
const (
	BaseURL                     = ""
	FundsTransferCreditEndpoint = "/fundsTransferCredit"
	GetBanksPath                = "/getBanks"
	NameEnquiryPath             = "/nameEnquiry"
)

// initiatize and populate the Account_creation struct
func NewTransfer(SendersAccountNo, ReceiverAccountNo, Amount, Charge, Tax, Narration string) *InternalTransaction {
	return &InternalTransaction{
		SendersAccountNo:  SendersAccountNo,
		ReceiverAccountNo: ReceiverAccountNo,
		Amount:            Amount,
		Charge:            Charge,
		Narration:         Narration,
		Tax:               Tax,
	}
}

var httpClient = resty.New() // Global HTTP client

// Internal Transfer
func SpectrumTransfer(tnx *InternalTransaction) ([]byte, error) {
	requestData := map[string]string{
		"fromAccountNumber":      tnx.SendersAccountNo,
		"toAccountNumber":        tnx.ReceiverAccountNo,
		"transactionAmount":      tnx.Amount,
		"transactionDescription": tnx.Narration,
		"chargeAmount":           tnx.Charge,
		"taxAmount":              tnx.Tax,
	}

	body, err := postRequest(requestData)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Internal Transfer To Ledger, credit ledger, debit user
func SpectrumTransferToLedger(Amount, SendersAccountNo, Narration string) ([]byte, error) {
	requestData := map[string]string{
		"fromAccountNumber":      SendersAccountNo,
		"toAccountNumber":        "",
		"transactionAmount":      Amount,
		"transactionDescription": Narration,
		"chargeAmount":           "",
		"taxAmount":              "",
	}
	body, err := postRequest(requestData)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Internal Transfer To Ledger, credit ledger, debit user
func SpectrumTransferFromLedgerToAccountNo(Amount string, AccountNo, Narration string) ([]byte, error) {

	requestData := map[string]interface{}{
		"accountNumber": AccountNo,
		"amount":        Amount,
		"narration":     Narration,
		"chargeAmount":  "0",
		"taxAmount":     "0",
	}
	// requestData := map[string]string{
	// 	"fromAccountNumber":      SendersAccountNo,
	// 	"toAccountNumber":        "",
	// 	"transactionAmount":      Amount,
	// 	"transactionDescription": Narration,
	// 	"chargeAmount":           "",
	// 	"taxAmount":              "",
	// }
	body, err := postCreditRequest(requestData)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func OutwardTransfer(tnx *FundsTransferCreditRequest) ([]byte, error) {
	requestData := map[string]string{

		"nameEnquiryRef":                    tnx.NameEnquiryRef,
		"destinationInstitutionCode":        tnx.DestinationInstitutionCode,
		"channelCode":                       tnx.ChannelCode,
		"beneficiaryAccountName":            tnx.BeneficiaryAccountName,
		"beneficiaryAccountNumber":          tnx.BeneficiaryAccountNumber,
		"beneficiaryBankVerificationNumber": tnx.BeneficiaryBankVerificationNumber,
		"beneficiaryKYCLevel":               tnx.BeneficiaryKYCLevel,
		"originatorAccountName":             tnx.OriginatorAccountName,
		"originatorAccountNumber":           tnx.OriginatorAccountNumber,
		"originatorBankVerificationNumber":  tnx.OriginatorBankVerificationNumber,
		"originatorKYCLevel":                tnx.OriginatorKYCLevel,
		"transactionLocation":               "6.625800, 3.334580",
		"narration":                         tnx.Narration,
		"paymentReference":                  tnx.PaymentReference,
		"amount":                            tnx.Amount,
	}

	body, err := postExternalRequest(requestData)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// External Transfer
func FundsTransferCredit(requestData *FundsTransferCreditRequest) ([]byte, error) {
	// Convert request data to JSON
	reqBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}
	fmt.Println("Request Body:", requestData)

	// Perform HTTP POST request using the global HTTP client
	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(reqBody).
		Post(BaseURL + FundsTransferCreditEndpoint)
	if err != nil {
		fmt.Println("Logging error from ext. tnx", err)
		return nil, err
	}

	fmt.Println("Logging response from ext. tnx", resp)
	// Check for errors in the response status code
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("Unexpected status code")
	}

	// Return the response body
	return resp.Body(), nil
}

// Get List of Banks and code
func GetBanks() ([]byte, error) {
	// Perform HTTP GET request using the global HTTP client
	resp, err := httpClient.R().
		SetHeader("Accept", "application/json").
		Get(BaseURL + GetBanksPath)
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

func NameEnquiry(requestData *NameEnquiryRequest) ([]byte, error) {
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
		Post(BaseURL + NameEnquiryPath)
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

// Helpers
func postRequest(data map[string]string) ([]byte, error) {
	// Convert request data to JSON
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Perform HTTP POST request using the global HTTP client
	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(reqBody).
		Post(ORBIT_INTERNAL_DEBIT_TRANSFER)
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

func postExternalRequest(data map[string]string) ([]byte, error) {
	// Convert request data to JSON
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Perform HTTP POST request using the global HTTP client
	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(reqBody).
		Post(ORBIT_EXTERNAL_DEBIT_TRANSFER)
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
func postCreditRequest(data map[string]interface{}) ([]byte, error) {
	// Convert request data to JSON
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Perform HTTP POST request using the global HTTP client
	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(reqBody).
		Post(ORBIT_INTERNAL_CREDIT_TRANSFER)
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
