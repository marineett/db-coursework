package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IClientRepository interface {
	InsertClient(client types.ClientData, personalData types.PersonalData, authData types.AuthData) (int64, error)
	GetClient(id int64) (*types.ClientData, error)
	UpdateClientPersonalData(clientID int64, personalData types.PersonalData) error
	UpdateClientPassword(clientID int64, authData types.AuthData, newPassword string) error
}

func CreateClientTable(db *sql.DB, clientTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + clientTableName + ` (
		id SERIAL PRIMARY KEY,
		summary_rating INTEGER NOT NULL,
		reviews_count INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", clientTableName, err)
	}
	return nil
}

type ClientRepository struct {
	db                     *sql.DB
	clientTable            string
	userRepository         IUserRepository
	personalDataRepository IPersonalDataRepository
	authRepository         IAuthRepository
}

func CreateClientRepository(db *sql.DB, personalDataTable string, userTable string, clientTable string, authTable string) *ClientRepository {
	return &ClientRepository{
		db:                     db,
		clientTable:            clientTable,
		userRepository:         CreateUserRepository(db, userTable),
		personalDataRepository: CreatePersonalDataRepository(db, personalDataTable),
		authRepository:         CreateAuthRepository(db, authTable),
	}
}

func (r *ClientRepository) InsertClient(client types.ClientData, personalData types.PersonalData, authData types.AuthData) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	client.UserData.PersonalDataID, err = r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	client.UserData.ID, err = r.userRepository.InsertUserInSeq(tx, client.UserData)
	if err != nil {
		return 0, err
	}
	authInfo := types.AuthInfo{
		UserID:   client.UserData.ID,
		UserType: types.Client,
		Login:    authData.Login,
		Password: authData.Password,
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, authInfo)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.clientTable + ` (id, summary_rating, reviews_count)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	row := tx.QueryRow(query, client.UserData.ID, 0, 0)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ClientRepository) GetClient(id int64) (*types.ClientData, error) {
	query := `
	SELECT * FROM ` + r.clientTable + ` WHERE id = $1
	`
	summaryRating, reviewsCount := 0, 0
	row := r.db.QueryRow(query, id)
	var client types.ClientData
	err := row.Scan(&client.ID, &summaryRating, &reviewsCount)
	if err != nil {
		return nil, err
	}
	userData, err := r.userRepository.GetUser(client.ID)
	if err != nil {
		return nil, err
	}
	client.UserData = *userData
	if reviewsCount > 0 {
		client.MeanRating = float64(summaryRating) / float64(reviewsCount)
	} else {
		client.MeanRating = 0
	}
	return &client, nil
}

func (r *ClientRepository) UpdateClientPersonalData(client_id int64, personalData types.PersonalData) error {
	client, err := r.GetClient(client_id)
	if err != nil {
		return err
	}
	return r.personalDataRepository.UpdatePersonalData(client.PersonalDataID, personalData)
}

func (r *ClientRepository) UpdateClientPassword(client_id int64, authData types.AuthData, newPassword string) error {
	return r.authRepository.ChangePassword(client_id, authData, newPassword)
}
