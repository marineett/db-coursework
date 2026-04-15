package data_base

import (
	"data_base_project/types"
	"database/sql"
)

type ILessonRepository interface {
	InsertLesson(lesson types.Lesson) (int64, error)
	GetLessons(contractID int64, from int64, size int64) ([]types.Lesson, error)
}

type LessonRepository struct {
	db               *sql.DB
	lessonTable      string
	contractTable    string
	transactionTable string
}

func CreateLessonTable(db *sql.DB, lessonTable string, contractTable string, transactionTable string) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + lessonTable + ` (
		id SERIAL PRIMARY KEY,
		contract_id INTEGER NOT NULL,
		duration INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (contract_id) REFERENCES ` + contractTable + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	triggerFunc := `
		CREATE OR REPLACE FUNCTION create_transaction_on_lesson()
		RETURNS TRIGGER AS $$
		DECLARE
			price INTEGER;
			user_id INTEGER;
		BEGIN
			SELECT c.price, c.client_id INTO price, user_id
			FROM ` + contractTable + ` c
			WHERE c.id = NEW.contract_id;
			INSERT INTO ` + transactionTable + ` (
				user_id,
				amount,
				created_at,
                status,
                type
			) VALUES (
				user_id,
				(price * NEW.duration)::integer / 60,
				NEW.created_at,
                1,
                2
			);

			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	_, err = db.Exec(triggerFunc)
	if err != nil {
		return err
	}
	trigger := `
		DROP TRIGGER IF EXISTS lesson_transaction_trigger ON ` + lessonTable + `;
		CREATE TRIGGER lesson_transaction_trigger
		AFTER INSERT ON ` + lessonTable + `
		FOR EACH ROW
		EXECUTE FUNCTION create_transaction_on_lesson();
	`
	_, err = db.Exec(trigger)
	if err != nil {
		return err
	}

	return nil
}

func CreateLessonRepository(db *sql.DB, lessonTable string, contractTable string, transactionTable string) *LessonRepository {
	return &LessonRepository{
		db:               db,
		lessonTable:      lessonTable,
		contractTable:    contractTable,
		transactionTable: transactionTable,
	}
}

func (r *LessonRepository) InsertLesson(lesson types.Lesson) (int64, error) {
	query := `
		INSERT INTO ` + r.lessonTable + ` (contract_id, duration, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	row := r.db.QueryRow(query, lesson.ContractID, lesson.Duration, lesson.CreatedAt)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *LessonRepository) GetLessons(contractID int64, from int64, size int64) ([]types.Lesson, error) {
	query := `
		SELECT id, contract_id, duration, created_at
		FROM ` + r.lessonTable + `
		WHERE contract_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, contractID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lessons := []types.Lesson{}
	for rows.Next() {
		var lesson types.Lesson
		err := rows.Scan(&lesson.ID, &lesson.ContractID, &lesson.Duration, &lesson.CreatedAt)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}
