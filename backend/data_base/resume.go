package data_base

import (
	"data_base_project/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type IResumeRepository interface {
	InsertResume(resume types.Resume) (int64, error)
	InsertResumeInSeq(tx *sql.Tx, resume types.Resume) (int64, error)
	GetResume(id int64) (*types.Resume, error)
	UpdateResumeTitle(id int64, title string, updatedAt time.Time) error
	UpdateResumeDescription(id int64, description string, updatedAt time.Time) error
	UpdateResumePrices(id int64, prices map[string]int, updatedAt time.Time) error
	DeleteResume(id int64) error
}

func CreateResumeTable(db *sql.DB, resumeTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + resumeTableName + ` (
		id SERIAL PRIMARY KEY,
		repetitor_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		prices JSONB NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		FOREIGN KEY (repetitor_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", resumeTableName, err)
	}
	return nil
}

type ResumeRepository struct {
	db          *sql.DB
	resumeTable string
}

func CreateResumeRepository(db *sql.DB, resumeTable string) *ResumeRepository {
	return &ResumeRepository{db: db, resumeTable: resumeTable}
}

func (r *ResumeRepository) InsertResume(resume types.Resume) (int64, error) {
	query := `
	INSERT INTO ` + r.resumeTable + ` (repetitor_id, title, description, prices, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	pricesJSON, err := json.Marshal(resume.Prices)
	if err != nil {
		return 0, err
	}
	var insertedID int64
	err = r.db.QueryRow(query, resume.RepetitorID, resume.Title, resume.Description, pricesJSON, resume.CreatedAt, resume.UpdatedAt).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ResumeRepository) InsertResumeInSeq(tx *sql.Tx, resume types.Resume) (int64, error) {
	query := `
	INSERT INTO ` + r.resumeTable + ` (repetitor_id, title, description, prices, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	pricesJSON, err := json.Marshal(resume.Prices)
	if err != nil {
		return 0, err
	}
	var insertedID int64
	err = tx.QueryRow(query, resume.RepetitorID, resume.Title, resume.Description, pricesJSON, resume.CreatedAt, resume.UpdatedAt).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ResumeRepository) GetResume(id int64) (*types.Resume, error) {
	query := `
	SELECT id, repetitor_id, title, description, prices, created_at, updated_at FROM ` + r.resumeTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var resume types.Resume
	var pricesJSON []byte
	err := row.Scan(&resume.ID, &resume.RepetitorID, &resume.Title, &resume.Description, &pricesJSON, &resume.CreatedAt, &resume.UpdatedAt)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(pricesJSON, &resume.Prices)
	if err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) UpdateResumeTitle(id int64, title string, updatedAt time.Time) error {
	query := `
	UPDATE ` + r.resumeTable + ` SET title = $1, updated_at = $2 WHERE id = $3
	`
	result, err := r.db.Exec(query, title, updatedAt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("resume not found")
	}
	return nil
}

func (r *ResumeRepository) UpdateResumeDescription(id int64, description string, updatedAt time.Time) error {
	query := `
	UPDATE ` + r.resumeTable + ` SET description = $1, updated_at = $2 WHERE id = $3
	`
	result, err := r.db.Exec(query, description, updatedAt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("resume not found")
	}
	return nil
}

func (r *ResumeRepository) UpdateResumePrices(id int64, prices map[string]int, updatedAt time.Time) error {
	query := `
	UPDATE ` + r.resumeTable + ` SET prices = $1, updated_at = $2 WHERE id = $3
	`
	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(query, pricesJSON, updatedAt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("resume not found")
	}
	return nil
}

func (r *ResumeRepository) DeleteResume(id int64) error {
	query := `
	DELETE FROM ` + r.resumeTable + ` WHERE id = $1
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("resume not found")
	}
	return nil
}
