package types

import "time"

type TransactionStatus int

const (
	TransactionStatusNull TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusPaid
	TransactionStatusRefunded
	TransactionStatusFailed
)

type TransactionType int

const (
	TransactionTypeNull TransactionType = iota
	TransactionTypeContractPayment
	TransactionTypeLessonPayment
)

type Transaction struct {
	ID        int64             `json:"id"`
	UserID    int64             `json:"user_id"`
	Amount    int64             `json:"amount"`
	Status    TransactionStatus `json:"status"`
	Type      TransactionType   `json:"type"`
	CreatedAt time.Time         `json:"created_at"`
}

type PendingContractPaymentTransaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	TransactionID int64     `json:"transaction_id"`
}
