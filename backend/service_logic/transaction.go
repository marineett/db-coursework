package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type ITransactionService interface {
	GetTransactionsList(user_id int64, from int64, size int64) ([]types.Transaction, error)
	CreateContractPaymentTransaction(amount int64, user_id int64) (int64, error)
	ChangeTransactionStatus(transaction_id int64, status types.TransactionStatus) error
	GetTransaction(transaction_id int64) (*types.Transaction, error)
	GetPendingContractPaymentTransaction() (*types.PendingContractPaymentTransaction, error)
	ApproveTransaction(transaction_id int64) error
}

type TransactionService struct {
	transactionRepository data_base.ITransactionRepository
}

func CreateTransactionService(transactionRepository data_base.ITransactionRepository) *TransactionService {
	return &TransactionService{transactionRepository: transactionRepository}
}

func (s *TransactionService) CreateContractPaymentTransaction(amount int64, user_id int64) (int64, error) {
	return s.transactionRepository.InsertTransaction(types.Transaction{
		UserID:    user_id,
		Amount:    amount,
		Status:    types.TransactionStatusPending,
		Type:      types.TransactionTypeContractPayment,
		CreatedAt: time.Now(),
	})
}

func (s *TransactionService) GetTransactionsList(user_id int64, from int64, size int64) ([]types.Transaction, error) {
	return s.transactionRepository.GetTransactionsList(user_id, from, size)
}

func (s *TransactionService) ChangeTransactionStatus(transaction_id int64, status types.TransactionStatus) error {
	return s.transactionRepository.UpdateTransactionStatus(transaction_id, status)
}

func (s *TransactionService) GetTransaction(transaction_id int64) (*types.Transaction, error) {
	return s.transactionRepository.GetTransaction(transaction_id)
}

func (s *TransactionService) GetPendingContractPaymentTransaction() (*types.PendingContractPaymentTransaction, error) {
	return s.transactionRepository.GetPendingContractPaymentTransaction()
}

func (s *TransactionService) ApproveTransaction(transaction_id int64) error {
	return s.transactionRepository.ApproveTransaction(transaction_id)
}
