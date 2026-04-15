package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"time"
)

type IContractRepository interface {
	InsertContract(contract types.Contract) (int64, error)
	BeginTx() (*sql.Tx, error)
	GetContract(id int64) (*types.Contract, error)
	GetContractsByRepetitorID(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error)
	GetContractsByClientID(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error)
	UpdateContractStatus(id int64, status types.ContractStatus) error
	UpdateContractRepetitorID(contractID int64, repetitorID int64) error
	UpdateContractPaymentStatus(id int64, paymentStatus types.PaymentStatus) error
	UpdateContractReviewClientID(id int64, reviewClientID int64) error
	UpdateContractReviewClientIDInSeq(tx *sql.Tx, id int64, reviewClientID int64) error
	UpdateContractReviewRepetitorID(id int64, reviewRepetitorID int64) error
	GetContractList(from int64, size int64, status types.ContractStatus) ([]types.Contract, error)
	GetAllContracts(from int64, size int64) ([]types.Contract, error)
}

func CreateContractTable(
	db *sql.DB,
	contractTableName string,
	userTableName string,
	reviewTableName string,
	repetitorTableName string,
	clientTableName string,
) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + contractTableName + ` (
		id SERIAL PRIMARY KEY,
		client_id INTEGER NOT NULL,
		repetitor_id INTEGER NOT NULL,
		review_client_id INTEGER NOT NULL,
		review_repetitor_id INTEGER NOT NULL,
		transaction_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		description TEXT NOT NULL,
		status INTEGER NOT NULL,
		payment_status INTEGER NOT NULL,

		price INTEGER NOT NULL,
		commission INTEGER NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,

		FOREIGN KEY (client_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", contractTableName, err)
	}

	query = `
	CREATE OR REPLACE FUNCTION update_repetitor_rating_and_count()
	RETURNS TRIGGER AS $$
	DECLARE
		new_rating INTEGER := 0;
	BEGIN
		SELECT rating INTO new_rating
		FROM ` + reviewTableName + `
		WHERE id = NEW.review_client_id;

		UPDATE ` + repetitorTableName + ` rt
		SET 
			summary_rating = summary_rating + new_rating,
			reviews_count = reviews_count + 1
		WHERE rt.id = NEW.repetitor_id;
		
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating function: %v", err)
	}

	query = `
	DROP TRIGGER IF EXISTS update_repetitor_rating_trigger ON ` + contractTableName + `;
	CREATE TRIGGER update_repetitor_rating_trigger
	AFTER UPDATE OF review_client_id ON ` + contractTableName + `
	FOR EACH ROW
	EXECUTE FUNCTION update_repetitor_rating_and_count();
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating trigger: %v", err)
	}
	query = `
	CREATE OR REPLACE FUNCTION update_client_rating_and_count()
	RETURNS TRIGGER AS $$
	DECLARE
		new_rating INTEGER := 0;
	BEGIN
		SELECT rating INTO new_rating
		FROM ` + reviewTableName + `
		WHERE id = NEW.review_repetitor_id;

		UPDATE ` + clientTableName + ` ct
		SET 
			summary_rating = summary_rating + new_rating,
			reviews_count = reviews_count + 1
		WHERE ct.id = NEW.client_id;
		
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating function: %v", err)
	}

	query = `
	DROP TRIGGER IF EXISTS update_client_rating_trigger ON ` + contractTableName + `;
	CREATE TRIGGER update_client_rating_trigger
	AFTER UPDATE OF review_repetitor_id ON ` + contractTableName + `
	FOR EACH ROW
	EXECUTE FUNCTION update_client_rating_and_count();
	`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating trigger: %v", err)
	}

	return nil
}

type ContractRepository struct {
	db            *sql.DB
	contractTable string
}

func CreateContractRepository(db *sql.DB, contractTable string) *ContractRepository {
	return &ContractRepository{
		db:            db,
		contractTable: contractTable,
	}
}

func (r *ContractRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *ContractRepository) InsertContract(contract types.Contract) (int64, error) {
	query := `
	INSERT INTO ` + r.contractTable + ` (client_id, repetitor_id, description, status, transaction_id, created_at, payment_status, review_client_id, review_repetitor_id, price, commission, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id
	`
	var insertedID int64
	err := r.db.QueryRow(query, contract.ClientID, contract.RepetitorID, contract.Description, contract.Status, contract.TransactionID, contract.CreatedAt, contract.PaymentStatus, 0, 0, contract.Price, contract.Commission, contract.StartDate, contract.EndDate).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ContractRepository) GetContract(id int64) (*types.Contract, error) {
	query := `
	SELECT id, client_id, repetitor_id, created_at, description, status, payment_status, review_client_id, review_repetitor_id, price, commission, start_date, end_date FROM ` + r.contractTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var contract types.Contract
	err := row.Scan(&contract.ID, &contract.ClientID, &contract.RepetitorID, &contract.CreatedAt, &contract.Description, &contract.Status, &contract.PaymentStatus, &contract.ReviewClientID, &contract.ReviewRepetitorID, &contract.Price, &contract.Commission, &contract.StartDate, &contract.EndDate)
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *ContractRepository) GetContractsByRepetitorID(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("SET ROLE repetitor")
	if err != nil {
		return nil, err
	}

	query := `
	SELECT 
		id, 
		client_id, 
		repetitor_id, 
		created_at, 
		description, 
		status, 
		payment_status, 
		review_client_id, 
		review_repetitor_id, 
		price, 
		commission, 
		start_date, 
		end_date FROM ` + r.contractTable + ` WHERE 
		repetitor_id = $1 AND status = $2 
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`
	rows, err := tx.Query(query, repetitorID, status, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.Contract{}
	for rows.Next() {
		var contract types.Contract
		err := rows.Scan(&contract.ID, &contract.ClientID, &contract.RepetitorID, &contract.CreatedAt, &contract.Description, &contract.Status, &contract.PaymentStatus, &contract.ReviewClientID, &contract.ReviewRepetitorID, &contract.Price, &contract.Commission, &contract.StartDate, &contract.EndDate)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	_, err = tx.Exec("RESET ROLE")
	if err != nil {
		return nil, err
	}
	return contracts, tx.Commit()
}

func (r *ContractRepository) GetContractsByClientID(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("SET ROLE client")
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, 
		client_id, 
		repetitor_id, 
		created_at, 
		description, 
		status, 
		payment_status, 
		review_client_id, 
		review_repetitor_id, 
		price, 
		commission, 
		start_date, 
		end_date 
		FROM ` + r.contractTable + ` WHERE client_id = $1 AND status = $2 
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`
	rows, err := tx.Query(query, clientID, status, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.Contract{}
	for rows.Next() {
		var contract types.Contract
		err := rows.Scan(&contract.ID, &contract.ClientID, &contract.RepetitorID, &contract.CreatedAt, &contract.Description, &contract.Status, &contract.PaymentStatus, &contract.ReviewClientID, &contract.ReviewRepetitorID, &contract.Price, &contract.Commission, &contract.StartDate, &contract.EndDate)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	_, err = tx.Exec("RESET ROLE")
	if err != nil {
		return nil, err
	}
	return contracts, tx.Commit()
}

func (r *ContractRepository) UpdateContractStatus(id int64, status types.ContractStatus) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("SET ROLE moderator")
	if err != nil {
		return err
	}

	query := `
	UPDATE ` + r.contractTable + ` SET status = $1 WHERE id = $2
	`
	result, err := tx.Exec(query, status, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	_, err = tx.Exec("RESET ROLE")
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ContractRepository) UpdateContractPaymentStatus(id int64, paymentStatus types.PaymentStatus) error {
	query := `
	UPDATE ` + r.contractTable + ` SET payment_status = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, paymentStatus, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *ContractRepository) UpdateContractReviewClientID(id int64, reviewClientID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	UPDATE ` + r.contractTable + ` SET review_client_id = $1 WHERE id = $2
	`
	result, err := tx.Exec(query, reviewClientID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}

	return tx.Commit()
}

func (r *ContractRepository) UpdateContractReviewRepetitorID(id int64, reviewRepetitorID int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET review_repetitor_id = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, reviewRepetitorID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *ContractRepository) UpdateContractPrice(id int64, price int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET price = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, price, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *ContractRepository) UpdateContractCommission(id int64, commission int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET commission = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, commission, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *ContractRepository) UpdateContractStartDate(id int64, startDate time.Time) error {
	query := `
	UPDATE ` + r.contractTable + ` SET start_date = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, startDate, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *ContractRepository) UpdateContractRepetitorID(contractID int64, repetitorID int64) error {

	query := `
	UPDATE ` + r.contractTable + ` SET repetitor_id = $1, status = $2 WHERE id = $3
	`

	result, err := r.db.Exec(query, repetitorID, types.ContractStatusActive, contractID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *ContractRepository) GetContractList(from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	query := `
	SELECT id, client_id, repetitor_id, created_at, description, status, payment_status, review_client_id, review_repetitor_id, price, commission, start_date, end_date FROM ` + r.contractTable + ` WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, status, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.Contract{}
	for rows.Next() {
		var contract types.Contract
		err := rows.Scan(&contract.ID, &contract.ClientID, &contract.RepetitorID, &contract.CreatedAt, &contract.Description, &contract.Status, &contract.PaymentStatus, &contract.ReviewClientID, &contract.ReviewRepetitorID, &contract.Price, &contract.Commission, &contract.StartDate, &contract.EndDate)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *ContractRepository) UpdateContractReviewClientIDInSeq(tx *sql.Tx, id int64, reviewClientID int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET review_client_id = $1 WHERE id = $2
	`
	result, err := tx.Exec(query, reviewClientID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *ContractRepository) GetAllContracts(from int64, size int64) ([]types.Contract, error) {
	query := `
	SELECT 
	id, 
	client_id, 
	repetitor_id, 
	created_at, 
	description, 
	status, 
	payment_status, 
	review_client_id, 
	review_repetitor_id, 
	price, 
	commission, 
	start_date, 
	end_date 
	FROM ` + r.contractTable + ` 
	ORDER BY created_at DESC 
	LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.Contract{}
	for rows.Next() {
		var contract types.Contract
		err := rows.Scan(&contract.ID, &contract.ClientID, &contract.RepetitorID, &contract.CreatedAt, &contract.Description, &contract.Status, &contract.PaymentStatus, &contract.ReviewClientID, &contract.ReviewRepetitorID, &contract.Price, &contract.Commission, &contract.StartDate, &contract.EndDate)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}
