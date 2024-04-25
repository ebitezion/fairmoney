package data

import (
	"database/sql"
	"errors"
)

var (
	// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
	// looking up a movie that doesn't exist in our database.
	ErrRecordNotFound       = errors.New("record not found")
	ErrTransactionPINNotSet = errors.New("Transaction PIN Not Set")
	// Define a custom ErrEditConflict error. We'll return this from our Update() method
	// when there is a data race.
	ErrEditConflict = errors.New("edit conflict")

	KYCLEVEL0 = "0" //no verification
	KYCLEVEL1 = "1" //Email,or phone verified
	KYCLEVEL2 = "2" //BVN verified
	KYCLEVEL3 = "3" //Address verified
	KYCLEVEL4 = "4" //Account Upgraded

	Pending   = "pending"
	Completed = "completed"
	Cancelled = "cancelled"
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	//AModel MyModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
	// VersionModel     VersionModel
	AccountModel AccountModel
	// MediaModel       MediaModel
	// ErrorModel       ErrorModel
	// VerifyModel      VerifyModel

}

// For ease of use, we also add a New() method which returns a Models struct containing
// the intitialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
		// VersionModel:     VersionModel{DB: db},
		AccountModel: AccountModel{DB: db},
		// MediaModel:       MediaModel{DB: db},
		// ErrorModel:       ErrorModel{DB: db},
		// VerifyModel:      VerifyModel{DB: db},
	}
}
