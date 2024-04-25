package data

// TransactionType represents the type of transaction (credit or debit)
type TransactionType string

const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

// Transaction represents a debit or credit transaction
type Payment struct {
	AccountID string          `json:"account_id"`
	Reference string          `json:"reference"`
	Amount    int             `json:"amount"`
	Type      TransactionType `json:"type"`
}
