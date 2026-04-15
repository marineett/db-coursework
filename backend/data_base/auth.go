package data_base

import (
	"data_base_project/types"
	"database/sql"
	"errors"
	"fmt"
)

type IAuthRepository interface {
	InsertAuthInSeq(tx *sql.Tx, auth types.AuthInfo) (int64, error)
	InsertAuth(auth types.AuthInfo) (int64, error)
	ChangePassword(userId int64, authData types.AuthData, newPassword string) error
	Authorize(auth_data types.AuthData) (types.AuthVerdict, error)
	CheckLogin(login string) (bool, error)
}

func CreateAuthTable(db *sql.DB, authTableName string, userTableName string) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + authTableName + ` (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		user_type INTEGER NOT NULL,
		login VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", authTableName, err)
	}
	return nil
}

func ApplyAuthIndex(db *sql.DB, authTableName string) error {
	query := `CREATE INDEX IF NOT EXISTS idx_` + authTableName + `_login ON ` + authTableName + ` (login);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating index on %s: %v", authTableName, err)
	}
	return nil
}

type AuthRepository struct {
	db        *sql.DB
	authTable string
}

func CreateAuthRepository(db *sql.DB, authTable string) *AuthRepository {
	return &AuthRepository{
		db:        db,
		authTable: authTable,
	}
}

func (r *AuthRepository) InsertAuthInSeq(tx *sql.Tx, auth types.AuthInfo) (int64, error) {
	ok, err := r.CheckLogin(auth.Login)
	if err != nil {
		return 0, err
	}
	if ok {
		return 0, errors.New("login already exists")
	}
	query := `
	INSERT INTO ` + r.authTable + ` (user_id, user_type, login, password) VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	row := tx.QueryRow(query, auth.UserID, auth.UserType, auth.Login, auth.Password)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthRepository) InsertAuth(auth types.AuthInfo) (int64, error) {
	query := `
	INSERT INTO ` + r.authTable + ` (user_id, user_type, login, password) VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, auth.UserID, auth.UserType, auth.Login, auth.Password)
	if err != nil {
		return 0, err
	}
	return auth.UserID, nil
}

func (r *AuthRepository) ChangePassword(userId int64, authData types.AuthData, newPassword string) error {
	verdict, err := r.Authorize(authData)
	if err != nil {
		return err
	}
	if verdict.UserID != userId {
		return errors.New("invalid user id")
	}
	query := `
	UPDATE ` + r.authTable + ` SET password = $1 WHERE user_id = $2
	`
	_, err = r.db.Exec(query, newPassword, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) Authorize(auth_data types.AuthData) (types.AuthVerdict, error) {
	query := `
	SELECT * FROM ` + r.authTable + ` WHERE login = $1
	`
	var auth types.AuthInfo
	err := r.db.QueryRow(query, auth_data.Login).Scan(&auth.ID, &auth.UserID, &auth.UserType, &auth.Login, &auth.Password)
	if err != nil {
		return types.AuthVerdict{}, err
	}
	if auth.Password != auth_data.Password {
		return types.AuthVerdict{}, errors.New("invalid password")
	}
	return types.AuthVerdict{
		UserID:   auth.UserID,
		UserType: auth.UserType,
	}, nil
}

func (r *AuthRepository) CheckLogin(login string) (bool, error) {
	query := `
	SELECT COUNT(*) FROM ` + r.authTable + ` WHERE login = $1
	`
	var count int
	err := r.db.QueryRow(query, login).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}
