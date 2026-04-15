package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IUserRepository interface {
	InsertUser(user types.UserData) (int64, error)
	InsertUserInSeq(tx *sql.Tx, user types.UserData) (int64, error)
	GetUser(id int64) (*types.UserData, error)
}

func CreateUserTable(db *sql.DB, userTableName string, personalDataTable string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + userTableName + ` (
		id SERIAL PRIMARY KEY,
		registration_date TIMESTAMP NOT NULL,
		last_login_date TIMESTAMP NOT NULL,
		personal_data_id INTEGER NOT NULL,
		FOREIGN KEY (personal_data_id) REFERENCES ` + personalDataTable + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", userTableName, err)
	}

	query = `CREATE INDEX IF NOT EXISTS idx_` + userTableName + `_personal_data_id ON ` + userTableName + ` (personal_data_id);`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating index on %s: %v", userTableName, err)
	}

	return nil
}

type UserRepository struct {
	db        *sql.DB
	userTable string
}

func CreateUserRepository(db *sql.DB, userTable string) *UserRepository {
	return &UserRepository{
		db:        db,
		userTable: userTable,
	}
}

func (r *UserRepository) InsertUser(user types.UserData) (int64, error) {
	query := `
	INSERT INTO ` + r.userTable + ` (registration_date, last_login_date, personal_data_id)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	row := r.db.QueryRow(query, user.RegistrationDate, user.LastLoginDate, user.PersonalDataID)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) InsertUserInSeq(tx *sql.Tx, user types.UserData) (int64, error) {
	query := `
	INSERT INTO ` + r.userTable + ` (registration_date, last_login_date, personal_data_id)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	row := tx.QueryRow(query, user.RegistrationDate, user.LastLoginDate, user.PersonalDataID)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetUser(id int64) (*types.UserData, error) {
	query := `
	SELECT id, registration_date, last_login_date, personal_data_id FROM ` + r.userTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var user types.UserData
	err := row.Scan(&user.ID, &user.RegistrationDate, &user.LastLoginDate, &user.PersonalDataID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
