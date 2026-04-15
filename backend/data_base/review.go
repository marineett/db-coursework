package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IReviewRepository interface {
	InsertReview(review types.Review) (int64, error)
	InsertReviewInSeq(tx *sql.Tx, review types.Review) (int64, error)
	UpdateReview(review types.Review) error
	GetReview(id int64) (*types.Review, error)
	GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.Review, error)
	GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.Review, error)
}

func CreateReviewTable(db *sql.DB, reviewTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + reviewTableName + ` (
		id SERIAL PRIMARY KEY,
		client_id INTEGER NOT NULL,
		repetitor_id INTEGER NOT NULL,
		rating INTEGER NOT NULL,
		comment TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (client_id) REFERENCES ` + userTableName + `(id),
		FOREIGN KEY (repetitor_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", reviewTableName, err)
	}
	return nil
}

type ReviewRepository struct {
	db          *sql.DB
	reviewTable string
}

func CreateReviewRepository(db *sql.DB, reviewTable string) IReviewRepository {
	return &ReviewRepository{
		db:          db,
		reviewTable: reviewTable,
	}
}

func (r *ReviewRepository) InsertReview(review types.Review) (int64, error) {
	query := `
	INSERT INTO ` + r.reviewTable + ` (client_id, repetitor_id, rating, comment, created_at) VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var insertedID int64
	err := r.db.QueryRow(query, review.ClientID, review.RepetitorID, review.Rating, review.Comment, review.CreatedAt).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ReviewRepository) InsertReviewInSeq(tx *sql.Tx, review types.Review) (int64, error) {
	query := `
	INSERT INTO ` + r.reviewTable + ` (client_id, repetitor_id, rating, comment, created_at) VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var insertedID int64
	err := tx.QueryRow(query, review.ClientID, review.RepetitorID, review.Rating, review.Comment, review.CreatedAt).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ReviewRepository) GetReview(id int64) (*types.Review, error) {
	query := `
	SELECT * FROM ` + r.reviewTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var review types.Review
	err := row.Scan(&review.ID, &review.ClientID, &review.RepetitorID, &review.Rating, &review.Comment, &review.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) UpdateReview(review types.Review) error {
	query := `
	UPDATE ` + r.reviewTable + ` SET rating = $1, comment = $2 WHERE id = $3
	`
	result, err := r.db.Exec(query, review.Rating, review.Comment, review.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

func (r *ReviewRepository) GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.Review, error) {
	query := `
	SELECT * FROM ` + r.reviewTable + ` WHERE repetitor_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, repetitorID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := []types.Review{}
	for rows.Next() {
		var review types.Review
		err := rows.Scan(&review.ID, &review.ClientID, &review.RepetitorID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.Review, error) {
	query := `
	SELECT * FROM ` + r.reviewTable + ` WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, clientID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := []types.Review{}
	for rows.Next() {
		var review types.Review
		err := rows.Scan(&review.ID, &review.ClientID, &review.RepetitorID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}
