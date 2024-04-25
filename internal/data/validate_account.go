package data

import (
	"github.com/ebitezion/backend-framework/internal/validator"
)

func ValidateAccountData(v *validator.Validator, account *Account) {

	v.Check(account.Surname != "", "surname", "must be provided")
	v.Check(account.FirstName != "", "firstname", "must be provided")
	v.Check(account.HomeAddress != "", "homeAddress", "must be provided")
	v.Check(account.City != "", "city", "must be provided")
	v.Check(account.PhoneNumber != "", "phoneNumber", "must be provided")
	v.Check(account.BVN != "", "bvn", "must be provided")

	if len(account.PhoneNumber) != 11 {
		v.AddError("error", "phoneNumber should be 11 characters")
	}
	if len(account.BVN) != 11 {
		v.AddError("error", "bvn should be 11 characters")
	}

}

// func ValidateAccountHistoryData(v *validator.Validator, account *AccountHistory) {
// 	v.Check(account.AccountNumber != "", "accountNumber", "must be provided")
// 	v.Check(account.Pagination != "", "pagination", "must be provided")
// 	v.Check(account.StartDate != "", "startDate", "must be provided")
// 	v.Check(account.EndDate != "", "endDate", "must be provided")

// 	if len(account.AccountNumber) != 10 {
// 		v.AddError("error", "accountNumber should be 10 characters")
// 	}
// }

func ValidateAccountNumber(v *validator.Validator, data *AccountNumber) {
	v.Check(data.AccountNumber != "", "accountNumber", "must be provided")

	if len(data.AccountNumber) != 10 {
		v.AddError("error", "accountNumber should be 10 characters")
	}

}
func ValidateLimitRequestData(v *validator.Validator, request *UpgradeLimitRequest) {

	v.Check(request.UserID != "", "userID", "must be provided")
	v.Check(request.Single != "", "single", "must be provided")
	v.Check(request.Daily != "", "daily", "must be provided")
	v.Check(request.Type != "", "type", "must be provided")

	requestTypes := map[string]bool{"transfers": true, "ussd": true, "bills": true}

	_, requestType := requestTypes[request.Type]
	v.Check(requestType, "type", "must either be 'transfers' , 'ussd' or 'bills'")
}

func ValidateLimitUpgradeData(v *validator.Validator, limitUpgrade *UpgradeLimit) {
	v.Check(limitUpgrade.LimitID != "", "limitID", "must be provided")
	v.Check(limitUpgrade.UserID != "", "userID", "must be provided")
}

func ValidateAccountUpgradeData(v *validator.Validator, account *AccountUpgradeData) {
	v.Check(account.Title != "", "title", "must be provided")
	v.Check(account.FirstName != "", "firstName", "must be provided")
	v.Check(account.LastName != "", "lastName", "must be provided")
	v.Check(account.AccountNumber != "", "accountNumber", "must be provided")
	v.Check(account.BVN != "", "bvn", "must be provided")
	v.Check(account.DOB != "", "dob", "must be provided")
	v.Check(account.PhoneNumber != "", "phoneNumber", "must be provided")
	v.Check(account.Email != "", "email", "must be provided")
	v.Check(account.MaidenName != "", "maidenName", "must be provided")
	v.Check(account.Nationality != "", "nationality", "must be provided")
	v.Check(account.Country != "", "country", "must be provided")
	v.Check(account.Address != "", "address", "must be provided")
	v.Check(account.IdType != "", "idType", "must be provided")
	v.Check(account.Document != "", "document", "must be provided")
	v.Check(validator.Matches(account.Email, validator.EmailRX), "email", "must be a valid email address")

	if len(account.AccountNumber) != 10 {
		v.AddError("error", "accountNumber should be 10 characters")
	}
	if len(account.PhoneNumber) != 11 {
		v.AddError("error", "phoneNumber should be 11 characters")
	}
	if len(account.BVN) != 11 {
		v.AddError("error", "bvn should be 11 characters")
	}

}
