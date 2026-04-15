package service_logic

import (
	service_test "data_base_project/tests/service_logic_tests"
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateTransaction(t *testing.T) {
	transactionRepository := service_test.CreateTestTransaction()
	transactionService := CreateTransactionService(&transactionRepository)
	_, err := transactionService.CreateContractPaymentTransaction(100, 10)
	if err != nil {
		t.Errorf("Error creating transaction: %v", err)
	}
}

func TestGetTransaction(t *testing.T) {
	transactionRepository := service_test.CreateTestTransaction()
	transactionService := CreateTransactionService(&transactionRepository)
	id, err := transactionService.CreateContractPaymentTransaction(100, 11)
	if err != nil {
		t.Errorf("Error creating first transaction: %v", err)
	}
	_, err = transactionService.CreateContractPaymentTransaction(200, 12)
	if err != nil {
		t.Errorf("Error creating second transaction: %v", err)
	}
	_, err = transactionService.CreateContractPaymentTransaction(300, 13)
	if err != nil {
		t.Errorf("Error creating third transaction: %v", err)
	}
	transaction, err := transactionService.GetTransaction(id)
	current_time := time.Now()
	if err != nil {
		t.Errorf("Error getting transaction: %v", err)
	}
	if transaction.ID != id {
		t.Errorf("Transaction ID is not %v: %v", id, transaction.ID)
	}
	if transaction.CreatedAt.After(current_time) {
		t.Errorf("Transaction created at is after current time: %v", transaction.CreatedAt)
	}
	if transaction.Amount != 100 {
		t.Errorf("Transaction amount is not 100: %v", transaction.Amount)
	}
	if transaction.Status != types.TransactionStatusPending {
		t.Errorf("Transaction status is not pending: %v", transaction.Status)
	}
	_, err = transactionService.GetTransaction(id + 3)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUpdateTransactionStatus(t *testing.T) {
	transactionRepository := service_test.CreateTestTransaction()
	transactionService := CreateTransactionService(&transactionRepository)
	id, err := transactionService.CreateContractPaymentTransaction(100, 10)
	if err != nil {
		t.Errorf("Error creating first transaction: %v", err)
	}
	_, err = transactionService.CreateContractPaymentTransaction(200, 11)
	if err != nil {
		t.Errorf("Error creating second transaction: %v", err)
	}
	_, err = transactionService.CreateContractPaymentTransaction(300, 12)
	if err != nil {
		t.Errorf("Error creating third transaction: %v", err)
	}
	err = transactionService.ChangeTransactionStatus(id, types.TransactionStatusPaid)
	if err != nil {
		t.Errorf("Error changing transaction status: %v", err)
	}
	transaction, err := transactionService.GetTransaction(id)
	if err != nil {
		t.Errorf("Error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPaid {
		t.Errorf("Transaction status is not paid: %v", transaction.Status)
	}
	err = transactionService.ChangeTransactionStatus(id, types.TransactionStatusFailed)
	if err != nil {
		t.Errorf("Error changing transaction status: %v", err)
	}
	transaction, err = transactionService.GetTransaction(id)
	if err != nil {
		t.Errorf("Error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusFailed {
		t.Errorf("Transaction status is not failed: %v", transaction.Status)
	}
	err = transactionService.ChangeTransactionStatus(id+3, types.TransactionStatusPaid)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetTransactionsList(t *testing.T) {
	transactionRepository := service_test.CreateTestTransaction()
	transactionService := CreateTransactionService(&transactionRepository)
	for i := range 10 {
		_, err := transactionService.CreateContractPaymentTransaction(100, int64(10+i%2))
		if err != nil {
			t.Errorf("Error creating transaction: %v", err)
		}
	}
	transactions, err := transactionService.GetTransactionsList(10, 1, 3)
	if err != nil {
		t.Errorf("Error getting transaction list: %v", err)
	}
	if len(transactions) != 3 {
		t.Errorf("Transaction list length is not 3: %v", len(transactions))
	}
	for _, transaction := range transactions {
		if transaction.UserID != 10 {
			t.Errorf("Transaction user ID if wrong: %v", transaction.UserID)
		}
	}
	for i := 1; i < len(transactions); i++ {
		if transactions[i].CreatedAt.After(transactions[i-1].CreatedAt) {
			t.Errorf("Transaction list is not sorted: %v", transactions[i].CreatedAt)
		}
	}
	new_transactions, err := transactionService.GetTransactionsList(10, 0, 2)
	if err != nil {
		t.Errorf("Error getting transaction list: %v", err)
	}
	if len(new_transactions) != 2 {
		t.Errorf("Transaction list length is not 2: %v", len(new_transactions))
	}
	if new_transactions[1].ID != transactions[0].ID {
		t.Errorf("Transaction windows is not correct: %v", new_transactions[1].ID)
	}
}
