package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"log"
)

type IRepetitorRepository interface {
	InsertRepetitor(repetitor types.RepetitorData, personalData types.PersonalData, auth types.AuthData) (int64, error)
	GetRepetitor(repetitorID int64) (*types.RepetitorData, error)
	UpdateRepetitorPersonalData(repetitorID int64, personalData types.PersonalData) error
	UpdateRepetitorPassword(repetitorID int64, authData types.AuthData, newPassword string) error
	GetRepetitorsIds(repetitorsOffset int64, repetitorsLimit int64) ([]int64, error)
}

func CreateRepetitorTable(db *sql.DB, repetitorTableName string, userTableName string, resumeTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + repetitorTableName + ` (
		id SERIAL PRIMARY KEY,
		resume_id INTEGER NOT NULL,
		summary_rating INTEGER NOT NULL,
		reviews_count INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id),
		FOREIGN KEY (resume_id) REFERENCES ` + resumeTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", repetitorTableName, err)
	}
	return nil
}

type RepetitorRepository struct {
	db                     *sql.DB
	userRepository         IUserRepository
	authRepository         IAuthRepository
	personalDataRepository IPersonalDataRepository
	resumeRepository       IResumeRepository
	repetitorTable         string
	reviewTable            string
}

func CreateRepetitorRepository(db *sql.DB, personalDataTable string, userTable string, repetitorTable string, authTable string, resumeTable string, reviewTable string) *RepetitorRepository {
	return &RepetitorRepository{
		db:                     db,
		userRepository:         CreateUserRepository(db, userTable),
		authRepository:         CreateAuthRepository(db, authTable),
		personalDataRepository: CreatePersonalDataRepository(db, personalDataTable),
		resumeRepository:       CreateResumeRepository(db, resumeTable),
		repetitorTable:         repetitorTable,
		reviewTable:            reviewTable,
	}
}
func (r *RepetitorRepository) InsertRepetitor(repetitor types.RepetitorData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	repetitor.PersonalDataID, err = r.personalDataRepository.InsertPersonalDataInSeq(tx, personalData)
	if err != nil {
		return 0, err
	}
	repetitor.UserData.ID, err = r.userRepository.InsertUserInSeq(tx, repetitor.UserData)
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, types.AuthInfo{
		UserID:   repetitor.UserData.ID,
		UserType: types.Repetitor,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	resumeID, err := r.resumeRepository.InsertResumeInSeq(tx, types.Resume{
		RepetitorID: repetitor.UserData.ID,
	})
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.repetitorTable + ` (id, summary_rating, reviews_count, resume_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	row := tx.QueryRow(query, repetitor.UserData.ID, 0, 0, resumeID)
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

func (r *RepetitorRepository) GetRepetitor(id int64) (*types.RepetitorData, error) {
	log.Println("GetRepetitor", id)
	query := `
	SELECT id, resume_id, summary_rating, reviews_count FROM ` + r.repetitorTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var repetitor types.RepetitorData
	var summaryRating float64
	var reviewsCount int64
	err := row.Scan(&repetitor.ID, &repetitor.ResumeID, &summaryRating, &reviewsCount)
	if err != nil {
		return nil, err
	}
	userData, err := r.userRepository.GetUser(repetitor.ID)
	if err != nil {
		return nil, err
	}
	repetitor.UserData = *userData
	if reviewsCount == 0 {
		repetitor.MeanRating = 0
	} else {
		repetitor.MeanRating = summaryRating / float64(reviewsCount)
	}
	return &repetitor, nil
}

func (r *RepetitorRepository) UpdateRepetitorPersonalData(repetitor_id int64, personalData types.PersonalData) error {
	repetitor, err := r.GetRepetitor(repetitor_id)
	if err != nil {
		return err
	}
	return r.personalDataRepository.UpdatePersonalData(repetitor.PersonalDataID, personalData)
}

func (r *RepetitorRepository) UpdateRepetitorPassword(repetitor_id int64, authData types.AuthData, newPassword string) error {
	return r.authRepository.ChangePassword(repetitor_id, authData, newPassword)
}

func (r *RepetitorRepository) GetRepetitorsIds(repetitorsOffset int64, repetitorsLimit int64) ([]int64, error) {
	query := `
	SELECT id FROM ` + r.repetitorTable + ` ORDER BY CASE WHEN reviews_count = 0 THEN 0 ELSE summary_rating / reviews_count END DESC OFFSET $1 LIMIT $2
	`
	rows, err := r.db.Query(query, repetitorsOffset, repetitorsLimit)
	if err != nil {
		return nil, err
	}
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
