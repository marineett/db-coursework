package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateTransactionRepository(t *testing.T) {
	transactionRepository := CreateTransactionRepository(globalDb, "test_transaction_table", "test_pending_contract_payment_transactions_table")
	if transactionRepository == nil {
		t.Errorf("Failed to create transaction repository")
	}
}

func TestInsertTransaction(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_transaction_table, test_user_table, test_contract_table, test_auth_table, test_personal_data_table CASCADE")
	transactionRepository := CreateTransactionRepository(globalDb, "test_transaction_table", "test_pending_contract_payment_transactions_table")
	if transactionRepository == nil {
		t.Errorf("Failed to create transaction repository")
	}
	id, err := transactionRepository.InsertTransaction(types.Transaction{
		UserID:    1,
		Amount:    100,
		Status:    types.TransactionStatusPending,
		CreatedAt: time.Now(),
		Type:      types.TransactionTypeContractPayment,
	})
	if err != nil {
		t.Errorf("Failed to create transaction: %v", err)
	}
	resultTransaction := types.Transaction{}
	err = globalDb.QueryRow("SELECT * FROM test_transaction_table WHERE user_id = $1", id).Scan(&resultTransaction.ID, &resultTransaction.UserID, &resultTransaction.Amount, &resultTransaction.Status, &resultTransaction.CreatedAt, &resultTransaction.Type)
	if err != nil {
		t.Errorf("Failed to get transaction: %v", err)
	}
	if resultTransaction.ID != id {
		t.Errorf("Failed to get transaction: %v", resultTransaction.ID)
	}
	if resultTransaction.UserID != 1 {
		t.Errorf("Failed to get transaction: %v", resultTransaction.UserID)
	}
	if resultTransaction.Amount != 100 {
		t.Errorf("Failed to get transaction: %v", resultTransaction.Amount)
	}
	if resultTransaction.Status != types.TransactionStatusPending {
		t.Errorf("Failed to get transaction: %v", resultTransaction.Status)
	}
	if resultTransaction.CreatedAt.IsZero() {
		t.Errorf("Failed to get transaction: %v", resultTransaction.CreatedAt)
	}
	if resultTransaction.Type != types.TransactionTypeContractPayment {
		t.Errorf("Failed to get transaction: %v", resultTransaction.Type)
	}
}

func TestGetTransaction(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestContract(1, 1, 2)
	defer globalDb.Exec("TRUNCATE TABLE test_transaction_table, test_user_table, test_contract_table, test_auth_table, test_personal_data_table CASCADE")
	transactionRepository := CreateTransactionRepository(globalDb, "test_transaction_table", "test_pending_contract_payment_transactions_table")
	if transactionRepository == nil {
		t.Errorf("Failed to create transaction repository")
	}
	id, err := transactionRepository.InsertTransaction(types.Transaction{
		UserID:    1,
		Amount:    100,
		Status:    types.TransactionStatusPending,
		CreatedAt: time.Now(),
		Type:      types.TransactionTypeContractPayment,
	})
	if err != nil {
		t.Errorf("Failed to create transaction: %v", err)
	}
	resultTransaction, err := transactionRepository.GetTransaction(id)
	if err != nil {
		t.Errorf("Failed to get transaction: %v", err)
	}
	if resultTransaction.ID != id {
		t.Errorf("Failed to get transaction: %v", resultTransaction.ID)
	}
	if resultTransaction.UserID != 1 {
		t.Errorf("Failed to get transaction: %v", resultTransaction.UserID)
	}
	if resultTransaction.Amount != 100 {
		t.Errorf("Failed to get transaction: %v", resultTransaction.Amount)
	}
	if resultTransaction.Status != types.TransactionStatusPending {
		t.Errorf("Failed to get transaction: %v", resultTransaction.Status)
	}
	if resultTransaction.CreatedAt.IsZero() {
		t.Errorf("Failed to get transaction: %v", resultTransaction.CreatedAt)
	}
	if resultTransaction.Type != types.TransactionTypeContractPayment {
		t.Errorf("Failed to get transaction: %v", resultTransaction.Type)
	}
	_, err = transactionRepository.GetTransaction(id + 1)
	if err == nil {
		t.Errorf("Failed to get transaction: %v", err)
	}
}

func TestGetTransactionsList(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_transaction_table, test_user_table, test_contract_table, test_auth_table, test_personal_data_table CASCADE")
	transactionRepository := CreateTransactionRepository(globalDb, "test_transaction_table", "test_pending_contract_payment_transactions_table")
	if transactionRepository == nil {
		t.Errorf("Failed to create transaction repository")
	}
	transaction := types.Transaction{
		UserID:    1,
		Amount:    100,
		Status:    types.TransactionStatusPending,
		CreatedAt: time.Now(),
		Type:      types.TransactionTypeContractPayment,
	}
	transactionId, err := transactionRepository.InsertTransaction(transaction)
	if err != nil {
		t.Errorf("Failed to create transaction: %v", err)
	}
	transaction2 := types.Transaction{
		UserID:    2,
		Amount:    200,
		Status:    types.TransactionStatusPending,
		CreatedAt: time.Now(),
		Type:      types.TransactionTypeContractPayment,
	}
	_, err = transactionRepository.InsertTransaction(transaction2)
	if err != nil {
		t.Errorf("Failed to create transaction: %v", err)
	}
	transaction3 := types.Transaction{
		UserID:    1,
		Amount:    300,
		Status:    types.TransactionStatusPending,
		CreatedAt: time.Now(),
		Type:      types.TransactionTypeContractPayment,
	}
	transactionId3, err := transactionRepository.InsertTransaction(transaction3)
	if err != nil {
		t.Errorf("Failed to create transaction: %v", err)
	}
	transactions, err := transactionRepository.GetTransactionsList(1, 0, 10)
	if err != nil {
		t.Errorf("Failed to get transactions list: %v", err)
	}
	if len(transactions) != 2 {
		t.Errorf("Failed to get transactions list: %v", len(transactions))
	}
	if transactions[0].ID != transactionId3 {
		t.Errorf("Failed to get transactions list: %v", transactions[0].ID)
	}
	if transactions[0].Amount != 300 {
		t.Errorf("Failed to get transactions list: %v", transactions[0].Amount)
	}
	if transactions[0].UserID != 1 {
		t.Errorf("Failed to get transactions list: %v", transactions[0].UserID)
	}
	if transactions[0].Status != types.TransactionStatusPending {
		t.Errorf("Failed to get transactions list: %v", transactions[0].Status)
	}
	if transactions[1].ID != transactionId {
		t.Errorf("Failed to get transactions list: %v", transactions[1].ID)
	}
	if transactions[1].Amount != 100 {
		t.Errorf("Failed to get transactions list: %v", transactions[1].Amount)
	}
	if transactions[1].UserID != 1 {
		t.Errorf("Failed to get transactions list: %v", transactions[1].UserID)
	}
	if transactions[1].Status != types.TransactionStatusPending {
		t.Errorf("Failed to get transactions list: %v", transactions[1].Status)
	}
	if transactions[1].Type != types.TransactionTypeContractPayment {
		t.Errorf("Failed to get transactions list: %v", transactions[1].Type)
	}
	transactions, err = transactionRepository.GetTransactionsList(1, 1, 10)
	if err != nil {
		t.Errorf("Failed to get transactions list: %v", err)
	}
	if len(transactions) != 1 {
		t.Errorf("Failed to get transactions list: %v", len(transactions))
	}
	if transactions[0].ID != transactionId {
		t.Errorf("Failed to get transactions list: %v", transactions[0].ID)
	}
	if transactions[0].Amount != 100 {
		t.Errorf("Failed to get transactions list: %v", transactions[0].Amount)
	}
	if transactions[0].UserID != 1 {
		t.Errorf("Failed to get transactions list: %v", transactions[0].UserID)
	}
	if transactions[0].Status != types.TransactionStatusPending {
		t.Errorf("Failed to get transactions list: %v", transactions[0].Status)
	}
	if transactions[0].Type != types.TransactionTypeContractPayment {
		t.Errorf("Failed to get transactions list: %v", transactions[0].Type)
	}
	transactions, err = transactionRepository.GetTransactionsList(3, 0, 10)
	if err != nil {
		t.Errorf("Failed to get transactions list: %v", err)
	}
	if len(transactions) != 0 {
		t.Errorf("Found transactions for user with no transactions")
	}
}
