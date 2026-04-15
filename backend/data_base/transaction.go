package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type ITransactionRepository interface {
	InsertTransaction(transaction types.Transaction) (int64, error)
	UpdateTransactionStatus(transactionId int64, status types.TransactionStatus) error
	GetTransaction(transactionId int64) (*types.Transaction, error)
	GetTransactionsList(userId int64, from int64, size int64) ([]types.Transaction, error)
	GetPendingContractPaymentTransaction() (*types.PendingContractPaymentTransaction, error)
	InsertPendingContractPaymentTransaction(
		transactionPendingContractPayment types.PendingContractPaymentTransaction,
		transaction types.Transaction,
	) (int64, error)
	ApproveTransaction(transactionId int64) error
}

func CreateTransactionTable(db *sql.DB, transactionTableName string, userTableName string, pendingContractPaymentTransactionsTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + transactionTableName + ` (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		status INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		type INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", transactionTableName, err)
	}

	query = `
	CREATE OR REPLACE FUNCTION insert_pending_contract_payment_transaction()
	RETURNS TRIGGER AS $$
	BEGIN
		IF NEW.status = 1 AND NEW.type = 1 THEN
			INSERT INTO ` + pendingContractPaymentTransactionsTableName + ` (id, user_id, amount, created_at)
			VALUES (NEW.id, NEW.user_id, NEW.amount, NEW.created_at);
		END IF;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating trigger function: %v", err)
	}

	query = `
	DROP TRIGGER IF EXISTS after_insert_transaction ON ` + transactionTableName + `;
	CREATE TRIGGER after_insert_transaction
	AFTER INSERT ON ` + transactionTableName + `
	FOR EACH ROW
	EXECUTE FUNCTION insert_pending_contract_payment_transaction();
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating trigger: %v", err)
	}

	return nil
}

func CreatePendingContractPaymentTransactionsTable(
	db *sql.DB,
	pendingContractPaymentTransactionsTableName string,
	userTableName string,
	transactionTableName string,
) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + pendingContractPaymentTransactionsTableName + ` (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", pendingContractPaymentTransactionsTableName, err)
	}

	query = `
	CREATE OR REPLACE FUNCTION set_transaction_paid_on_pending_delete()
	RETURNS TRIGGER AS $$
	BEGIN
		RAISE NOTICE 'Deleting pending transaction with id: %', OLD.id;
		RAISE NOTICE 'Updating transaction with id: %', OLD.id;
		
		IF EXISTS (SELECT 1 FROM ` + transactionTableName + ` WHERE id = OLD.id) THEN
			UPDATE ` + transactionTableName + `
			SET status = 2 
			WHERE id = OLD.id AND status = 1;
			RAISE NOTICE 'Transaction status updated to Paid';
		ELSE
			RAISE NOTICE 'Transaction not found with id: %', OLD.id;
		END IF;
		
		RETURN OLD;
	END;
	$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating delete trigger function: %v", err)
	}

	query = `
	DROP TRIGGER IF EXISTS after_delete_pending_contract_payment ON ` + pendingContractPaymentTransactionsTableName + `;
	CREATE TRIGGER after_delete_pending_contract_payment
	AFTER DELETE ON ` + pendingContractPaymentTransactionsTableName + `
	FOR EACH ROW
	EXECUTE FUNCTION set_transaction_paid_on_pending_delete();
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating delete trigger: %v", err)
	}

	return nil
}

type TransactionRepository struct {
	db                                      *sql.DB
	transactionTable                        string
	pendingContractPaymentTransactionsTable string
}

func CreateTransactionRepository(db *sql.DB, transactionTable string, pendingContractPaymentTransactionsTable string) *TransactionRepository {
	return &TransactionRepository{db: db, transactionTable: transactionTable, pendingContractPaymentTransactionsTable: pendingContractPaymentTransactionsTable}
}

func (r *TransactionRepository) InsertTransaction(transaction types.Transaction) (int64, error) {
	query := `
	INSERT INTO ` + r.transactionTable + ` (user_id, amount, status, created_at, type)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var insertedId int64
	err := r.db.QueryRow(query, transaction.UserID, transaction.Amount, transaction.Status, transaction.CreatedAt, transaction.Type).Scan(&insertedId)
	if err != nil {
		return 0, err
	}
	return insertedId, nil
}

func (r *TransactionRepository) UpdateTransactionStatus(transactionId int64, status types.TransactionStatus) error {
	query := `
	UPDATE ` + r.transactionTable + ` SET status = $1 WHERE id = $2
	`
	_, err := r.db.Exec(query, status, transactionId)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) GetTransaction(id int64) (*types.Transaction, error) {
	query := `
	SELECT * FROM ` + r.transactionTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var transaction types.Transaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.CreatedAt, &transaction.Type)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetTransactionsList(userId int64, from int64, size int64) ([]types.Transaction, error) {
	query := `
	SELECT * FROM ` + r.transactionTable + ` WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userId, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []types.Transaction
	for rows.Next() {
		var transaction types.Transaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.CreatedAt, &transaction.Type)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) InsertPendingContractPaymentTransaction(
	transactionPendingContractPayment types.PendingContractPaymentTransaction,
	transaction types.Transaction,
) (int64, error) {
	query := `INSERT INTO ` + r.pendingContractPaymentTransactionsTable + ` (user_id, amount, created_at, transaction_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var lastInsertId int
	err := r.db.QueryRow(query, transactionPendingContractPayment.UserID, transactionPendingContractPayment.Amount, transactionPendingContractPayment.CreatedAt, transaction.ID).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return int64(lastInsertId), nil
}

func (r *TransactionRepository) GetPendingContractPaymentTransaction() (*types.PendingContractPaymentTransaction, error) {
	query := `
	SELECT * FROM ` + r.pendingContractPaymentTransactionsTable + ` LIMIT 1
	`
	row := r.db.QueryRow(query)
	var transaction types.PendingContractPaymentTransaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) ApproveTransaction(transactionId int64) error {
	query := `
	DELETE FROM ` + r.pendingContractPaymentTransactionsTable + ` WHERE id = $1
	`
	_, err := r.db.Exec(query, transactionId)
	if err != nil {
		return err
	}
	return nil
}
