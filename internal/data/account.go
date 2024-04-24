package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"
)

type Email struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
}

type Account struct {
	Surname     string `json:"surname"`
	FirstName   string `json:"firstname"`
	HomeAddress string `json:"homeAddress"`
	City        string `json:"city"`
	PhoneNumber string `json:"phoneNumber"`
	BVN         string `json:"bvn"`
}

type AccountDetails struct {
	User_id         int    `json:"user_id"`
	Account_number  string `json:"account_number"`
	Transaction_pin string `json:"transaction_pin"`
	Created_at      string `json:"created_at"`
	Updated_at      string `json:"updated_at"`
	Limits          string `json:"limits"`
	Counter         string `json:"counter"`
}

type AccountHistory struct {
	AccountNumber string `json:"accountNumber"`
	Pagination    string `json:"pagination"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
}

type AccountNumber struct {
	AccountNumber string `json:"accountNumber"`
}

type AccountUpgradeData struct {
	Title         string `json:"title"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	MiddleName    string `json:"middleName"`
	AccountNumber string `json:"accountNumber"`
	BVN           string `json:"bvn"`
	DOB           string `json:"dob"`
	PhoneNumber   string `json:"phoneNumber"`
	Email         string `json:"email"`
	MaidenName    string `json:"maidenName"`
	Nationality   string `json:"nationality"`
	Country       string `json:"country"`
	Address       string `json:"address"`
	IdType        string `json:"idType"`
	Document      string `json:"document"`
}
type AccountProfile struct {
	TransferTag    *string `json:"transferTag"`
	AccountName    *string `json:"accountName"`
	AccountNumber  *string `json:"accountNumber"`
	PhoneNumber    *string `json:"phoneNumber"`
	Email          *string `json:"email"`
	NextOfKin      *string `json:"nextOfKin"`
	MarritalStatus *string `json:"marritalStatus"`
	KycLevel       *string `json:"kycLevel"`
	Username       *string `json:"username"`
}
type Transaction struct {
	ID                uint64   `json:"id"`
	UserID            uint64   `json:"user_id"`
	Type              string   `json:"type"`
	Source            string   `json:"source"`
	Narration         string   `json:"narration"`
	AccountNumber     string   `json:"accountNumber"`
	RequestID         string   `json:"requestId"`
	InternalReference string   `json:"internalReference"`
	ExternalReference *string  `json:"externalReference"`
	Amount            float64  `json:"amount"`
	CreatedAt         *string  `json:"createdAt"`
	UpdatedAt         *string  `json:"updatedAt"`
	Status            string   `json:"status"`
	Commission        *float64 `json:"commission"`
	BalanceAfter      *float64 `json:"balanceAfter"`
}
type UpgradeLimit struct {
	LimitID string `json:"limitID"`
	UserID  string `json:"userID"`
}
type UpgradeLimitRequest struct {
	Type   string `json:"type"`
	Single string `json:"single"`
	Daily  string `json:"daily"`
	UserID string `json:"userID"`
}
type Limits struct {
	Transfers TransferLimits `json:"transfers"`
	Bills     BillLimits     `json:"bills"`
	Ussd      UssdLimits     `json:"ussd"`
}

type TransferLimits struct {
	Single int64 `json:"single"`
	Daily  int64 `json:"daily"`
}

type BillLimits struct {
	Single int64 `json:"single"`
	Daily  int64 `json:"daily"`
}

type UssdLimits struct {
	Single int64 `json:"single"`
	Daily  int64 `json:"daily"`
}

type AccountModel struct {
	DB *sql.DB
}
type AccountHistoryResponse struct {
	AccountHistory  []AccountStatementData `json:"accountHistory"`
	ResponseCode    string                 `json:"responseCode"`
	ResponseMessage string                 `json:"responseMessage"`
}
type AccountStatementData struct {
	AccountType       string `json:"accountType"`
	BalanceAfter      string `json:"balanceAfter"`
	TransactionAmount string `json:"transactionAmount"`
	TransactionDate   string `json:"transactionDate"`
	TransactionDesc   string `json:"transactionDesc"`
	TransactionRef    string `json:"transactionRef"`
	TransactionType   string `json:"transactionType"`
}

// SortByDate sorts the transactions by transaction date.
func SortByDate(transactions []AccountStatementData) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].TransactionDate > transactions[j].TransactionDate
	})
}

type PdfData struct {
	AccountName          string `json:"accountName"`
	AccountNumber        string `json:"accountNumber"`
	TransactionStartDate string `json:"transactionDate"`
	TransactioEndDate    string
}

func (m AccountModel) UpdateLimitStatus(LimitID int) error {

	// Update the user's activated column in the database
	query := `
	UPDATE limit_upgrade_requests SET status = ? WHERE id = ? 
	`
	args := []interface{}{
		Completed,
		LimitID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m AccountModel) UpdateLimitInDB(limit string, userID string) error {
	// Update the user's activated column in the database
	query := `
	UPDATE user_details SET limits = ? WHERE user_id = ? 
	`
	args := []interface{}{

		limit,
		userID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m AccountModel) UpdateLimitCounterInDB(count string, userID string) error {
	// Update the user's activated column in the database
	query := `
	UPDATE user_details SET counter = ? WHERE user_id = ? 
	`
	var counter string = `{"transfers":` + count + `, "bills": 0, "ussd": 0, "ibank": 0}`
	args := []interface{}{
		counter,
		userID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m AccountModel) GetAccountLimits(userID string) (string, error) {

	// Set up the SQL query.
	query := `
    SELECT limits FROM user_details WHERE user_id = ?`

	// Create a slice containing the query arguments.
	args := []interface{}{userID}
	var limits string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&limits,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}
	// Return the matching user.
	return limits, nil
}

func (m AccountModel) GetUserAccountNoByID(userID string) (string, error) {

	// Set up the SQL query.
	query := `
    SELECT account_number FROM user_details WHERE user_id = ?`

	// Create a slice containing the query arguments.
	args := []interface{}{userID}
	var accountNo string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&accountNo,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}
	// Return the matching user.
	return accountNo, nil
}

func (m AccountModel) GetLimitUpgradeRequest(LimitID int) (*UpgradeLimitRequest, error) {

	// Set up the SQL query.
	query := `
    SELECT user_id,type,amount,daily FROM limit_upgrade_requests WHERE id = ?`

	// Create a slice containing the query arguments.
	args := []interface{}{LimitID}
	var limits UpgradeLimitRequest
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&limits.UserID,
		&limits.Type,
		&limits.Single,
		&limits.Daily,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &limits, nil
}

// NewAaccountUpgrade inserts account upgrade data into the database.
func (a AccountModel) CreateNewLimitRequest(data *UpgradeLimitRequest) error {

	query := `
	INSERT INTO limit_upgrade_requests (user_id,type,amount,daily)
	VALUES (?, ?, ?, ?)`
	args := []interface{}{
		data.UserID,
		data.Type,
		data.Single,
		data.Daily,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := a.DB.ExecContext(ctx, query, args...)
	return err
}
func (a AccountModel) NewMaillingList(email string) error {

	query := `
	INSERT INTO mailing_list (email)
	VALUES (?, ?)`
	args := []interface{}{
		email,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := a.DB.ExecContext(ctx, query, args...)
	return err
}

// NewAaccountUpgrade inserts account upgrade data into the database.
func (a AccountModel) NewAccountUpgrade(data *AccountUpgradeData) error {
	query := `
	INSERT INTO account_upgrade (title, first_name, last_name, middle_name, account_number, bvn, dob, phone_number, email, maiden_name, nationality, country, address, id_type, document)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	args := []interface{}{
		data.Title,
		data.FirstName,
		data.LastName,
		data.MiddleName,
		data.AccountNumber,
		data.BVN,
		data.DOB,
		data.PhoneNumber,
		data.Email,
		data.MaidenName,
		data.Nationality,
		data.Country,
		data.Address,
		data.IdType,
		data.Document,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := a.DB.ExecContext(ctx, query, args...)
	return err
}

func (a AccountModel) GetAccountHistory(accountNumber string, pagination string) ([]Transaction, error) {

	// Calculate the offset based on the pagination parameter.
	number, err := strconv.Atoi(pagination)
	if err != nil {
		return nil, err
	}
	offset := (number - 1) * 10

	query := "SELECT id, user_id, type, source, narration, account_number, request_id, internal_reference, external_reference, amount, created_at, updated_at, status, commission, balance_after FROM transactions WHERE account_number = ? ORDER BY created_at DESC LIMIT 10 OFFSET ?"

	var transactions []Transaction // Slice to hold multiple transaction records.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryContext method to execute the query, passing in the context.
	rows, err := a.DB.QueryContext(ctx, query, accountNumber, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and scan each row into a Transaction struct.
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.Type, &t.Source, &t.Narration, &t.AccountNumber, &t.RequestID, &t.InternalReference, &t.ExternalReference, &t.Amount, &t.CreatedAt, &t.UpdatedAt, &t.Status, &t.Commission, &t.BalanceAfter)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	if len(transactions) == 0 {
		return []Transaction{}, ErrRecordNotFound // Return empty slice if no transactions found.
	}
	return transactions, nil
}

// NewAaccountUpgrade inserts account upgrade data into the database.
func (a AccountModel) SaveCreatedAccountNo(data *AccountDetails) error {
	query := `
	INSERT INTO user_details(user_id, account_number, limits, created_at, updated_at, counter)
	VALUES (?, ?, COALESCE(?, '{"transfers":{"single":200000,"daily":600000},"bills":{"single":100000,"daily":200000},"ussd":{"single":10000,"daily":20000}}'), COALESCE(?, NOW()), COALESCE(?, NOW()),?)`

	// Set default values for limits, created_at, and updated_at if not provided
	if data.Limits == "" {
		data.Limits = `{"transfers":{"single":200000,"daily":600000},"bills":{"single":100000,"daily":200000},"ussd":{"single":10000,"daily":20000}}`
	}
	if data.Counter == "" {
		data.Counter = `{"transfers": 0, "bills": 0, "ussd": 0, "ibank": 0}`
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if data.Created_at == "" {
		data.Created_at = now
	}
	if data.Updated_at == "" {
		data.Updated_at = now
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query directly with the values
	re, err := a.DB.ExecContext(ctx, query, data.User_id, data.Account_number, data.Limits, data.Created_at, data.Updated_at, data.Counter)
	if err != nil {
		fmt.Println("ERRRRR:", err)
		return err
	}
	fmt.Println("Result:", re)
	return nil
}

// NewAccountUpgrade inserts account upgrade data into the database.
// func (a AccountModel) SaveTransactionDetails(transactions *Transaction) error {
// 	// Define the SQL query with correct placeholders created_at, updated_at,
// 	query := `
// 	INSERT INTO transactions(user_id,type, source, narration, account_number, request_id, internal_reference, external_reference, amount,  status)
// 	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?)`
// 	args := []interface{}{
// 		&transactions.UserID,
// 		&transactions.Type,
// 		&transactions.Source,
// 		&transactions.Narration,
// 		&transactions.AccountNumber,
// 		&transactions.RequestID,
// 		&transactions.InternalReference,
// 		&transactions.ExternalReference,
// 		&transactions.Amount,
// 		&transactions.Status,
// 	}

// 	//var transactions []Transaction
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Execute the query
// 	_, err := a.DB.ExecContext(ctx, query, args...)
// 	return err
// }

func (a AccountModel) SaveTransactionDetails(transaction *Transaction) error {
	// Define the SQL query with correct placeholders created_at, updated_at,
	query := `
	INSERT INTO transactions(user_id,type, source, narration, account_number, request_id, internal_reference, external_reference, amount,  status)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query with ExecContext.
	_, err := a.DB.ExecContext(ctx, query,
		transaction.UserID,
		transaction.Type,
		transaction.Source,
		transaction.Narration,
		transaction.AccountNumber,
		transaction.RequestID,
		transaction.InternalReference,
		transaction.ExternalReference,
		transaction.Amount,
		transaction.Status,
	)
	return err
}

// // NewAaccountUpgrade inserts account upgrade data into the database.
// func (a AccountModel) SaveCreatedAccountNo2(User_id, Account_number, Created_at, Updated_at, Limits string) error {
// 	query := `
// 	INSERT INTO user_details(user_id, account_number, created_at, updated_at, limits)
// 	VALUES (?, ?, ?, ?, ?, ?, ?)`
// 	args := []interface{}{
// 		User_id,
// 		Account_number,
// 		Created_at,
// 		Updated_at,
// 		Limits,
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	//set default value for limits, created_at,updated_at and counter,

// 	Limits = `{"transfers":{"single":200000,"daily":600000},"bills":{"single":100000,"daily":200000},"ussd":{"single":10000,"daily":20000}}`
// 	Created_at = time.Now().String()
// 	Updated_at = time.Now().String()

//		_, err := a.DB.ExecContext(ctx, query, args...)
//		return err
//	}
//
// NewAccountUpgrade inserts account upgrade data into the database.
func (a AccountModel) SaveCreatedAccountNo2(User_id, Account_number, Limits, Counter string, Created_at, Updated_at time.Time) error {
	// Define the SQL query with correct placeholders
	query := `
	INSERT INTO user_details(user_id, account_number, created_at, updated_at, limits, counter)
	VALUES (?, ?, ?, ?, ?,?)`
	args := []interface{}{
		User_id,
		Account_number,
		Created_at,
		Updated_at,
		Limits,
		Counter,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	_, err := a.DB.ExecContext(ctx, query, args...)
	return err
}

func (a AccountModel) GetAccountProfile(userID int64) (*AccountProfile, error) {
	query := `
		SELECT u.name, u.email, u.phone_number, u.kyc_level,u.username, t.transfer_tag,t.account_number
		FROM users u
		LEFT JOIN transfer_tag t ON u.id = t.user_id
		WHERE u.id = ?
	`

	var account AccountProfile

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := a.DB.QueryRowContext(ctx, query, userID).Scan(
		&account.AccountName,
		&account.Email,
		&account.PhoneNumber,
		&account.KycLevel,
		&account.Username,
		&account.TransferTag,
		&account.AccountNumber,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &account, nil
}
