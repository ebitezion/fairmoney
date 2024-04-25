package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/ebitezion/backend-framework/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var (
	ErrDeviceIDNotFound error = errors.New("DeviceID not found")
	//ErrKYCPhoneNotVerified error = errors.New("KYC Error: Phone No. Not Verified")
	//ErrKYCEmailNotVerified error = errors.New("Email Error: Email Not Verified")
	ErrKYCStatusNotVerified  error = errors.New("KYC Error: Status Not Verified")
	ErrKYCAddressNotVerified error = errors.New("KYC Error: Address Not Verified")
	ErrKYCAccountUpgraded    error = errors.New("KYC Error: Account Not Upgraded")
	ErrKYCBVNNotVerified     error = errors.New("KYC Error: BVN Not Verified")
	ErrDifferentSource       error = errors.New("Source Error: Spectrumpay User")
	ErrExistingAccountHolder error = errors.New("Status Error: User Has Not Verified Existing Account Number.")
)

type User struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          password  `json:"-"`
	Status            string    `json:"-"`
	Activated         null.Bool `json:"activated" `
	CreatedAt         *string   `json:"-"` //leave this as a pnter to avoid null nightmare
	PhoneNumber       string    `json:"phone_number"`
	Version           int       `json:"-"`
	DeviceID          string    `json:"device_id"`
	DeviceOS          string    `json:"device_os"`
	DeviceName        string    `json:"device_name"`
	Address_verified  null.Bool `json:"address_verified"`
	BVN_verified      null.Bool `json:"bvn_verified"`
	Account_upgraded  null.Bool `json:"account_upgrade"`
	KYC_level         int       `json:"kyc_level"`
	Source            string    `json:"-"`
	Is_spectrum_extra null.Bool `json:"-" `

	UserDevice UserDevice `json:"user_device"` //TODO : REMOVE !!!
}

type UserDetails struct {
	UserID        string `json:"userID"`
	PIN           string `json:"-"`
	AccountNumber string `json:"accountNumber"`
	Limits        string `json:"limits"`
	Counter       string `json:"count"`
}
type UserDetailsForLimits struct {
	UserID        string      `json:"userID"`
	PIN           string      `json:"-"`
	AccountNumber string      `json:"accountNumber"`
	Limits        Limits      `json:"limits"`
	Counter       LimitCounts `json:"count"`
	//TransactionPIN string `json:`
}
type LimitCounts struct {
	Transfers int `json:"transfers"`
	Bills     int `json:"bills"`
	USSD      int `json:"ussd"`
	IBank     int `json:"ibank"`
}

type UserDetailsWithPIN struct {
	UserID        string `json:"userID"`
	PIN           string `json:"-"`
	AccountNumber string `json:"accountNumber"`
	Limits        string `json:"limits"`
	Counter       string `json:"count"`
}
type UserDevice struct { //TODO : REMOVE !!!
	DeviceID   string `json:"device_id"`
	DeviceOS   string `json:"device_os"`
	DeviceName string `json:"device_name"`
}
type ResetPassword struct {
	Otp      string `json:"otp"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Hash     []byte `json:"-"`
}
type SetPinData struct {
	UserID string `json:"user_id"`
	PIN    string `json:"pin"`
}
type UpdatePinData struct {
	UserID string `json:"-"`
	PIN    string `json:"pin"`
	OTP    string `json:"otp"`
}

// Declare a new AnonymousUser variable.
var AnonymousUser = &User{}

// Check if a User instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Define a custom ErrDuplicateEmail error.
var (
	ErrDuplicateEmailOrUsername = errors.New("Duplicate email or username")
)

// Create a custom password type which is a struct containing the plaintext and hashed
// versions of the password for a user. The plaintext field is a *pointer* to a string,
// so that we're able to distinguish between a plaintext password not being present in
// the struct at all, versus a plaintext password which is the empty string "".
type password struct {
	plaintext *string
	hash      []byte
}

// Create a UserModel struct which wraps the connection pool.
type UserModel struct {
	DB *sql.DB
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (m UserModel) SetPassword(plaintextPassword string) ([]byte, error) {
	var Password password
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	Password.plaintext = &plaintextPassword
	Password.hash = hash

	return Password.hash, nil
}

// The Rese method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func HashPassword(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil

}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	fmt.Println("GOTTT TO Passsword Match area")
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidateResetPassword(v *validator.Validator, ResetPassword *ResetPassword) {

	v.Check(ResetPassword.Email != "", "email", "must be provided")
	v.Check(validator.Matches(ResetPassword.Email, validator.EmailRX), "email", "must be a valid email address")
	v.Check(ResetPassword.Password != "", "password", "must be provided")
	v.Check(len(ResetPassword.Password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(ResetPassword.Password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateDeviceID(v *validator.Validator, deviceID string) {
	v.Check(deviceID != "", "deviceID", "must be provided")
	//	v.Check(len(deviceID) >= 5, "deviceID", "must be at least 5 bytes long")
	v.Check(len(deviceID) <= 72, "deviceID", "must not be more than 72 bytes long")
}

func ValidateDeviceName(v *validator.Validator, deviceName string) {
	v.Check(deviceName != "", "deviceName", "must be provided")

	v.Check(len(deviceName) <= 72, "deviceName", "must not be more than 72 bytes long")
}

func ValidateDeviceOS(v *validator.Validator, deviceOS string) {
	v.Check(deviceOS != "", "deviceOS", "must be provided")
	v.Check(len(deviceOS) <= 72, "deviceOS", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)
	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		fmt.Println("missing password hash for user")
	}
}
func ValidateSetPinData(v *validator.Validator, user *SetPinData) {

	//v.Check(user.UserID != "", "user_id", "must be provided")
	v.Check(user.PIN != "", "pin", "must be provided")
	if len(user.PIN) != 4 {
		v.AddError("error", "pin should be 4 characters")
	}

}
func ValidateUpdatePinData(v *validator.Validator, user *UpdatePinData) {

	v.Check(user.PIN != "", "pin", "must be provided")
	v.Check(user.OTP != "", "otp", "must be provided")
	if len(user.PIN) != 4 {
		v.AddError("error", "pin should be 4 characters")
	}

}
func (m UserModel) Insert(user *User) error {

	insertQuery := `
	INSERT INTO users (name,username, email, password, activated, device_id, device_os, device_name)
	VALUES (?, ?, ? ,?, ?, ?, ?, ?)`
	args := []interface{}{user.Name, user.Username, user.Email, user.Password.hash, user.Activated, user.UserDevice.DeviceID, user.UserDevice.DeviceOS, user.UserDevice.DeviceName}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the INSERT query
	result, err := m.DB.ExecContext(ctx, insertQuery, args...)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Duplicate entry"):
			return ErrDuplicateEmailOrUsername
		default:
			return err
		}
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = lastInsertID

	return nil
}

func (m UserModel) Insert_Existing_Account_Holder(user *User) error {
	var is_spectrum_extra bool = true
	var status string = "existing"
	insertQuery := `
	INSERT INTO users (name,username, email, password,status, activated, device_id, device_os, device_name, is_spectrum_extra)
	VALUES (?, ?, ? ,?, ?, ?, ?, ?, ?, ?)`
	args := []interface{}{user.Name, user.Username, user.Email, user.Password.hash, status, user.Activated, user.UserDevice.DeviceID, user.UserDevice.DeviceOS, user.UserDevice.DeviceName, is_spectrum_extra}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the INSERT query
	result, err := m.DB.ExecContext(ctx, insertQuery, args...)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Duplicate entry"):
			return ErrDuplicateEmailOrUsername
		default:
			return err
		}
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = lastInsertID

	return nil
}

func (m UserModel) UpdateUserDevice(Email, DeviceID, DeviceName, DeviceOS string) error {
	query := `
	UPDATE users
	SET device_id  = ?, device_os  = ?, device_name  = ?
	WHERE email = ?
	`
	fmt.Println("Users Output", Email)
	args := []interface{}{DeviceID, DeviceOS, DeviceName, Email}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the INSERT query
	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Duplicate entry"):
			return ErrDuplicateEmailOrUsername
		default:
			return err
		}
	}

	// Retrieve the last inserted ID
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// Retrieve the User details from the database based on the user's email address.
// Because we have a UNIQUE constraint on the email column, this SQL query will only
// return one record (or none at all, in which case we return a ErrRecordNotFound error).
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, name, username, email, activated,password, created_at, version
	FROM users
	WHERE email = ?`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Activated,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// Retrieve the User details from the database based on the user's phone number .
// Because we have a UNIQUE constraint on the phonenumer column, this SQL query will only
// return one record (or none at all, in which case we return a ErrRecordNotFound error).
func (m UserModel) GetByPhoneNumber(phoneNumber string) (*User, error) {
	query := `
	SELECT id, name, username, email, activated,password, created_at, version
	FROM users
	WHERE phone_number = ?`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Activated,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// Retrieve the User details from the database based on the user's email address and device ID.
// Because we have a UNIQUE constraint on the email column, this SQL query will only
// return one record (or none at all, in which case we return a ErrRecordNotFound error).
func (m UserModel) GetByEmailAndDeviceID(email string, deviceID, deviceName, deviceOS string) (*User, error) {
	var phoneNumber sql.NullString
	query := `
    SELECT id, name, username, email, status, activated, password, created_at, phone_number, source, device_id, device_name,version, device_os, address_verified, bvn_verified, account_upgraded, kyc_level, is_spectrum_extra
    FROM users
    WHERE email = ? or username = ?`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Status,
		&user.Activated,
		&user.Password.hash,
		&user.CreatedAt,
		&phoneNumber, //using null string of sql
		&user.Source,
		&user.DeviceID,
		&user.DeviceName,
		&user.Version,
		&user.DeviceOS,
		&user.Address_verified,
		&user.BVN_verified,
		&user.Account_upgraded,
		&user.KYC_level,
		&user.Is_spectrum_extra,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	if phoneNumber.Valid {
		user.PhoneNumber = phoneNumber.String
	}

	//Existing Users from old Spectrumpay app
	//TODO: Add source for spectrumextra on registration

	if !user.Is_spectrum_extra.Bool {

		return &user, ErrDifferentSource
	}
	if user.Status == "existing" {
		fmt.Println("Existing User error")
		return &user, ErrExistingAccountHolder
	}
	//device validation
	if user.DeviceID != deviceID {
		return &user, ErrDeviceIDNotFound
	}
	if user.DeviceName != deviceName {
		return &user, ErrDeviceIDNotFound
	}
	if user.DeviceOS != deviceOS {
		return &user, ErrDeviceIDNotFound
	}

	//KYC levels
	// if !user.Address_verified {
	// 	return nil, ErrKYCAddressNotVerified
	// }
	if user.BVN_verified.Valid && !user.BVN_verified.Bool {
		return &user, ErrKYCBVNNotVerified
	}
	// if !user.Account_upgraded {
	// 	return nil, ErrKYCAccountUpgraded
	// }
	if user.Activated.Valid && !user.Activated.Bool {
		// User is not activated
		return nil, ErrKYCStatusNotVerified
	}

	return &user, nil
}

func (m UserModel) UpdatePassword(password *ResetPassword) error {

	query := `
	UPDATE users
	SET password = ?
	WHERE email = ?
	`
	args := []interface{}{
		password.Hash,
		password.Email,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername
		default:
			return err
		}
	}

	return nil
}

// Update the details for a specific user. Notice that we check against the version
// field to help prevent any race conditions during the request cycle, just like we did
// when updating a movie. And we also check for a violation of the "users_email_key"
// constraint when performing the update, just like we did when inserting the user
// record originally.
func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET name = ?, username = ?,email = ?, password = ?, activated = ?, version = version + 1
	WHERE id = ? AND version = ?
	RETURNING version`
	args := []interface{}{
		user.Name,
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetUserIdByToken(tokenScope, tokenPlaintext string) (int64, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the SQL query.
	query := `
	SELECT tokens.user_id
	FROM tokens
	WHERE tokens.hash = ?
	AND tokens.scope = ?
	AND tokens.expiry > ?`
	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}
	var ID int64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&ID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrRecordNotFound
		default:
			return 0, err
		}
	}
	// Return the matching user.
	return ID, nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the SQL query.
	query := `
	SELECT users.id, users.created_at, users.name, users.username, users.email, users.password, users.phone_number, users.activated, users.version
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = ?
	AND tokens.scope = ?
	AND tokens.expiry > NOW()`

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query.
	var user User
	var phoneNumber sql.NullString
	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], tokenScope).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&phoneNumber,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	if phoneNumber.Valid {
		user.PhoneNumber = phoneNumber.String
	}

	// Return the matching user.
	return &user, nil
}

func (m UserModel) GetUserDetailsFromToken(tokenScope, tokenPlaintext string) (*UserDetailsForLimits, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the SQL query.
	query := `
    SELECT user_details.user_id, user_details.account_number, user_details.limits, user_details.counter, user_details.transaction_pin
    FROM user_details
    INNER JOIN tokens
    ON user_details.user_id = tokens.user_id
    WHERE tokens.hash = ?
    AND tokens.scope = ?
    AND tokens.expiry > NOW()`

	// Create a slice containing the query arguments.
	args := []interface{}{tokenHash[:], tokenScope}
	var user UserDetailsForLimits
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct.
	var limitsJSON, limitCount string // Assuming limitsJSON is retrieved from the database
	var transactionPIN sql.NullString
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID,
		&user.AccountNumber,
		&limitsJSON, // Updated to scan limitsJSON
		&limitCount,
		&transactionPIN,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	fmt.Println("Check TPIN", transactionPIN.String)
	if !transactionPIN.Valid || transactionPIN.String == "" {
		fmt.Println("I got into Transaction Pin Validity Check", transactionPIN.String)
		return &user, ErrTransactionPINNotSet
	}

	// Unmarshal limitsJSON into the Limits struct
	err = json.Unmarshal([]byte(limitsJSON), &user.Limits)
	if err != nil {
		return nil, err
	}
	//Unmarshal limitCount into the Limitcount struct
	err = json.Unmarshal([]byte(limitCount), &user.Counter)
	if err != nil {
		return nil, err
	}

	// Return the matching user.
	return &user, nil
}

func (m UserModel) GetUserDetailsAndPINFromToken(tokenScope, tokenPlaintext string) (*UserDetailsWithPIN, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the SQL query.
	query := `
		SELECT user_id, account_number, limits, counter, transaction_pin
		FROM user_details
		WHERE EXISTS (
			SELECT 1
			FROM tokens
			WHERE user_details.user_id = tokens.user_id
			AND tokens.hash = ? AND tokens.scope = ? AND tokens.expiry > NOW()
		)`

	// Create a slice containing the query arguments.
	args := []interface{}{tokenHash[:], tokenScope}
	var user UserDetailsWithPIN
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID,
		&user.AccountNumber,
		&user.Limits,
		&user.Counter,
		&user.PIN,
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
	return &user, nil
}

func (m UserModel) GetUserByUserID(UserId string) (*User, error) {

	// Set up the SQL query.
	query := `
	SELECT users.id, users.created_at, users.name, users.username, users.email, users.password	, users.activated, users.version
	FROM users WHERE id = ?`

	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []interface{}{UserId}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
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
	return &user, nil
}

// Retrieve the User details from the database based on the user's id.
// Because we have a UNIQUE constraint on the email column, this SQL query will only
// return one record (or none at all, in which case we return a ErrRecordNotFound error).
func (m UserModel) GetAccountNoByID(id string) (string, error) {
	fmt.Println("From GetAccountNoByID(id string) ID is", id)
	query := `
	SELECT account_number
	FROM user_details
	WHERE user_id = ?`
	var user UserDetails
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.AccountNumber,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}
	return user.AccountNumber, nil
}

// Update the details for a specific user. Notice that we check against the version
// field to help prevent any race conditions during the request cycle, just like we did
// when updating a movie. And we also check for a violation of the "users_email_key"
// constraint when performing the update, just like we did when inserting the user
// record originally.
func (m UserModel) UpdateActivated(user *User) error {
	query := `
	UPDATE users
	SET activated = ?, version = version + 1
	WHERE id = ? AND version = ?
	RETURNING version`
	args := []interface{}{
		user.Name,
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) activateUser(user *User) error {
	// Update the user's activated column in the database
	query := `
	UPDATE users SET activated = true , version = version + 1 WHERE id = ? AND version = ?
	`
	args := []interface{}{

		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// updateExistingAccountNoUserStatus updates the status from exixting to active neccessary for login-Do this after Verification
func (m UserModel) updateExistingAccountNoUserStatus(user *User) error {
	// Update the user's activated column in the database
	query := `
	UPDATE users SET status = "active" , version = version + 1 WHERE id = ? AND version = ?
	`
	args := []interface{}{

		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m UserModel) UpdateBvnVerified(user *User) error {
	query := `
	UPDATE users
	SET bvn_verified = ?
	WHERE id = ? 
	`
	verified := 1
	args := []interface{}{verified, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// UpdateStatus_PhoneNumberVerified
func (m UserModel) UpdateStatus_PhoneNumberVerified(user *User) error {
	var status string = "active"
	query := `
	UPDATE users
	SET phone_number_verified = ?, status = ?
	WHERE id = ? 
	`
	verified := 1
	args := []interface{}{verified, status, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (m UserModel) UpdatePhoneNumberVerified(user *User) error {
	query := `
	UPDATE users
	SET phone_number_verified = ? , activated = ?
	WHERE id = ? 
	`
	verified := 1
	args := []interface{}{verified, verified, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) UpdateEmailVerified(user *User) error {
	query := `
	UPDATE users
	SET email_verified = ? , activated = ?
	WHERE id = ? 
	`
	verified := 1
	args := []interface{}{verified, verified, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	re, err := m.DB.ExecContext(ctx, query, args...)
	fmt.Println("Result, err :::", re, err)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) UpdateEmailVerified_for_Spectrumpay_Users(user *User) error {
	query := `
	UPDATE users
	SET email_verified = ?, is_spectrum_extra = ?
	WHERE id = ? 
	`
	verified := 1
	args := []interface{}{verified, verified, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) SetFirstTimePIN(user *UserDetails) error {
	query := `
	UPDATE user_details
	SET transaction_pin = ?
	WHERE user_id = ? 
	`

	args := []interface{}{user.PIN, user.UserID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_id_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// this function updats the kyc level of a user based on the kyc level provided
func (m UserModel) UpdateKycLevel(user *User, level string) error {
	query := `
	UPDATE users
	SET kyc_level = ?
	WHERE id = ? 
	`

	args := []interface{}{level, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) UpdateActivatedIfVerified(user *User) error {
	query := `
		UPDATE users
		SET activated = 1
		WHERE id = ? AND email_verified = 1 OR phone_number_verified = 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, user.ID)
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
func (m UserModel) UpdateTransactionPin(data *SetPinData) error {
	query := `
	UPDATE user_details
	SET transaction_pin = ?
	WHERE user_id = ?
	`
	args := []interface{}{
		&data.PIN,
		&data.UserID,
	}
	fmt.Println("Waht args shows", args)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	re, err := m.DB.ExecContext(ctx, query, data.PIN, data.UserID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername
		default:
			return err
		}
	}
	ire, _ := re.LastInsertId()
	fmt.Println("Result", ire)
	return nil
}
func (m UserModel) UpdateTransactionPin2(pin, userid string) error {
	query := `
	UPDATE user_details
	SET transaction_pin = ?
	WHERE user_id = ?
	`

	fmt.Println("Waht args shows", pin, userid)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	re, err := m.DB.ExecContext(ctx, query, pin, userid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmailOrUsername
		default:
			return err
		}
	}
	ire, _ := re.LastInsertId()
	fmt.Println("Result", ire)
	return nil
}

func (m UserModel) IsAuthTokenForUserID(tokenPlaintext, UserID string) bool {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the SQL query. user.ID
	query := `
	SELECT users.id, users.created_at, users.name, users.username, users.email, users.password	, users.activated, users.version
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = ?
	AND tokens.expiry > ?`
	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []interface{}{tokenHash[:], time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		//Handle Error Here
		fmt.Println(err)
	}
	if fmt.Sprint(user.ID) == UserID {
		return true
	}
	// Return the matching user.
	return false
}
