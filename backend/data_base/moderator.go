package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IModeratorRepository interface {
	InsertModerator(moderator types.ModeratorData, personalData types.PersonalData, auth types.AuthData) (int64, error)
	GetModerator(moderatorID int64) (*types.ModeratorData, error)
	UpdateModeratorPersonalData(moderatorID int64, personalData types.PersonalData) error
	UpdateModeratorPassword(moderatorID int64, authData types.AuthData, newPassword string) error
	UpdateModeratorSalary(moderatorID int64, salary int64) error
	GetModerators() ([]int64, error)
}

func CreateModeratorTable(db *sql.DB, moderatorTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + moderatorTableName + ` (
		id SERIAL PRIMARY KEY,
		salary INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", moderatorTableName, err)
	}
	return nil
}

type ModeratorRepository struct {
	db                     *sql.DB
	moderatorTable         string
	userRepository         IUserRepository
	personalDataRepository IPersonalDataRepository
	authRepository         IAuthRepository
}

func CreateModeratorRepository(db *sql.DB, personalDataTable string, userTable string, moderatorTable string, authTable string) *ModeratorRepository {
	return &ModeratorRepository{
		db:                     db,
		moderatorTable:         moderatorTable,
		userRepository:         CreateUserRepository(db, userTable),
		personalDataRepository: CreatePersonalDataRepository(db, personalDataTable),
		authRepository:         CreateAuthRepository(db, authTable),
	}
}

func (r *ModeratorRepository) InsertModerator(moderator types.ModeratorData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	moderator.PersonalDataID, err = r.personalDataRepository.InsertPersonalDataInSeq(tx, personalData)
	if err != nil {
		return 0, err
	}
	moderator.ID, err = r.userRepository.InsertUserInSeq(tx, moderator.UserData)
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, types.AuthInfo{
		UserID:   moderator.ID,
		UserType: types.Moderator,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}

	query := `
	INSERT INTO ` + r.moderatorTable + ` (id, salary)
	VALUES ($1, $2)
	RETURNING id
	`
	var id int64
	err = tx.QueryRow(query, moderator.ID, moderator.Salary).Scan(&id)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ModeratorRepository) GetModerator(moderatorID int64) (*types.ModeratorData, error) {
	query := `
	SELECT * FROM ` + r.moderatorTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, moderatorID)
	var moderator types.ModeratorData
	err := row.Scan(&moderator.ID, &moderator.Salary)
	if err != nil {
		return nil, err
	}
	userData, err := r.userRepository.GetUser(moderator.UserData.ID)
	if err != nil {
		return nil, err
	}
	moderator.UserData = *userData
	return &moderator, nil
}

func (r *ModeratorRepository) UpdateModeratorPersonalData(moderator_id int64, personalData types.PersonalData) error {
	moderator, err := r.GetModerator(moderator_id)
	if err != nil {
		return err
	}
	return r.personalDataRepository.UpdatePersonalData(moderator.PersonalDataID, personalData)
}

func (r *ModeratorRepository) UpdateModeratorPassword(moderator_id int64, authData types.AuthData, newPassword string) error {
	return r.authRepository.ChangePassword(moderator_id, authData, newPassword)
}

func (r *ModeratorRepository) UpdateModeratorSalary(moderator_id int64, salary int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("SET ROLE admin")
	if err != nil {
		return err
	}
	query := `
	UPDATE ` + r.moderatorTable + ` SET salary = $1 WHERE id = $2
	`
	_, err = tx.Exec(query, salary, moderator_id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("RESET ROLE")
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ModeratorRepository) GetModerators() ([]int64, error) {
	query := `
	SELECT id FROM ` + r.moderatorTable + `
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
