package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IPersonalDataRepository interface {
	InsertPersonalDataInSeq(tx *sql.Tx, personalData types.PersonalData) (int64, error)
	InsertPersonalData(personalData types.PersonalData) (int64, error)
	GetPersonalData(id int64) (*types.PersonalData, error)
	UpdatePersonalData(id int64, personalData types.PersonalData) error
}

func CreatePersonalDataTable(db *sql.DB, personalDataTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + personalDataTableName + ` (
		id SERIAL PRIMARY KEY,
		telephone_number VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		passport_number VARCHAR(255) NOT NULL,
		passport_series VARCHAR(255) NOT NULL,
		passport_date TIMESTAMP NOT NULL,
		passport_issued_by VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		middle_name VARCHAR(255) NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", personalDataTableName, err)
	}

	query = `CREATE INDEX IF NOT EXISTS idx_` + personalDataTableName + `_email ON ` + personalDataTableName + ` (email);`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating index on %s: %v", personalDataTableName, err)
	}

	return nil
}

type PersonalDataRepository struct {
	db                *sql.DB
	personalDataTable string
}

func CreatePersonalDataRepository(db *sql.DB, personalDataTable string) *PersonalDataRepository {
	return &PersonalDataRepository{
		db:                db,
		personalDataTable: personalDataTable,
	}
}

func (r *PersonalDataRepository) InsertPersonalDataInSeq(tx *sql.Tx, personalData types.PersonalData) (int64, error) {
	query := `
	INSERT INTO ` + r.personalDataTable + ` (telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by, first_name, last_name, middle_name)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`
	row := tx.QueryRow(query, personalData.TelephoneNumber, personalData.Email, personalData.PassportNumber, personalData.PassportSeries, personalData.PassportDate, personalData.PassportIssuedBy, personalData.FirstName, personalData.LastName, personalData.MiddleName)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PersonalDataRepository) InsertPersonalData(personalData types.PersonalData) (int64, error) {
	query := `INSERT INTO ` + r.personalDataTable + ` (telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by, first_name, last_name, middle_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	var lastInsertId int
	err := r.db.QueryRow(query, personalData.TelephoneNumber, personalData.Email, personalData.PassportNumber, personalData.PassportSeries, personalData.PassportDate, personalData.PassportIssuedBy, personalData.FirstName, personalData.LastName, personalData.MiddleName).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return int64(lastInsertId), nil
}

func (r *PersonalDataRepository) GetPersonalData(id int64) (*types.PersonalData, error) {
	query := `
	SELECT * FROM ` + r.personalDataTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var personalData types.PersonalData
	err := row.Scan(&personalData.ID, &personalData.TelephoneNumber, &personalData.Email, &personalData.PassportNumber, &personalData.PassportSeries, &personalData.PassportDate, &personalData.PassportIssuedBy, &personalData.FirstName, &personalData.LastName, &personalData.MiddleName)
	if err != nil {
		return nil, err
	}
	return &personalData, nil
}

func (r *PersonalDataRepository) UpdatePersonalData(id int64, personalData types.PersonalData) error {
	query := `UPDATE ` + r.personalDataTable + ` SET telephone_number = $1, email = $2, passport_number = $3, passport_series = $4, passport_date = $5, passport_issued_by = $6, first_name = $7, last_name = $8, middle_name = $9 WHERE id = $10`
	_, err := r.db.Exec(query, personalData.TelephoneNumber, personalData.Email, personalData.PassportNumber, personalData.PassportSeries, personalData.PassportDate, personalData.PassportIssuedBy, personalData.FirstName, personalData.LastName, personalData.MiddleName, id)
	return err
}
