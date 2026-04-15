package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"log"
)

type IDepartmentRepository interface {
	InsertDepartment(department types.Department) (int64, error)
	GetDepartmentIdByName(name string) (int64, error)
	GetDepartmentsByHeadID(headID int64) ([]types.Department, error)
	GetDepartment(id int64) (*types.Department, error)
	ChangeDepartmentHead(id int64, headID int64) error
	HireInfoInsert(hireInfo types.HireInfo) error
	HireInfoDelete(userId int64, departmentId int64) error
	GetUserDepartmentsIDs(userId int64) ([]int64, error)
	GetDepartmentUsersIDs(departmentId int64) ([]int64, error)
}

func CreateDepartmentTable(db *sql.DB, departmentTable string, hireInfoTable string, userTableName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
	CREATE TABLE IF NOT EXISTS ` + departmentTable + ` (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		head_id INTEGER NOT NULL
	)
	`
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}
	query = `
	CREATE INDEX IF NOT EXISTS idx_department_name ON ` + departmentTable + ` (name)
	`
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}
	query = `
	CREATE TABLE IF NOT EXISTS ` + hireInfoTable + ` (
		id SERIAL PRIMARY KEY,
		department_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (department_id) REFERENCES ` + departmentTable + `(id),
		FOREIGN KEY (user_id) REFERENCES ` + userTableName + `(id)
	)
	`
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}
	return tx.Commit()
}

type DepartmentRepository struct {
	db              *sql.DB
	departmentTable string
	hireInfoTable   string
}

func CreateDepartmentRepository(db *sql.DB, departmentTable string, hireInfoTable string) *DepartmentRepository {
	return &DepartmentRepository{db: db, departmentTable: departmentTable, hireInfoTable: hireInfoTable}
}

func (r *DepartmentRepository) InsertDepartment(department types.Department) (int64, error) {
	query := `
	INSERT INTO ` + r.departmentTable + ` (name, head_id) VALUES ($1, $2)
	RETURNING id
	`
	var id int64
	err := r.db.QueryRow(query, department.Name, department.HeadID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *DepartmentRepository) GetDepartment(id int64) (*types.Department, error) {
	query := `
	SELECT id, name, head_id FROM ` + r.departmentTable + ` WHERE id = $1
	`
	var department types.Department
	err := r.db.QueryRow(query, id).Scan(&department.ID, &department.Name, &department.HeadID)
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *DepartmentRepository) GetDepartmentsByHeadID(headID int64) ([]types.Department, error) {
	query := `
	SELECT id, name, head_id FROM ` + r.departmentTable + ` WHERE head_id = $1
	`
	var departments []types.Department
	rows, err := r.db.Query(query, headID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var department types.Department
		err := rows.Scan(&department.ID, &department.Name, &department.HeadID)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, nil
}

func (r *DepartmentRepository) ChangeDepartmentHead(id int64, headID int64) error {
	query := `
	UPDATE ` + r.departmentTable + ` SET head_id = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, headID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("department not found")
	}
	return nil
}

func (r *DepartmentRepository) HireInfoInsert(hireInfo types.HireInfo) error {
	query := `
	INSERT INTO ` + r.hireInfoTable + ` (department_id, user_id) VALUES ($1, $2)
	`
	_, err := r.db.Exec(query, hireInfo.DepartmentID, hireInfo.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DepartmentRepository) HireInfoDelete(userId int64, departmentId int64) error {
	query := `
	DELETE FROM ` + r.hireInfoTable + ` WHERE user_id = $1 AND department_id = $2
	`
	result, err := r.db.Exec(query, userId, departmentId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("hire info not found")
	}
	return nil
}

func (r *DepartmentRepository) GetUserDepartmentsIDs(userId int64) ([]int64, error) {
	query := `
	SELECT department_id FROM ` + r.hireInfoTable + ` WHERE user_id = $1
	`
	var departments []int64
	rows, err := r.db.Query(query, userId)
	if err != nil {
		log.Printf("Error getting user departments IDs: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var departmentID int64
		err := rows.Scan(&departmentID)
		if err != nil {
			log.Printf("Error getting user departments IDs: %v", err)
			return nil, err
		}
		departments = append(departments, departmentID)
	}
	return departments, nil
}

func (r *DepartmentRepository) GetDepartmentUsersIDs(departmentId int64) ([]int64, error) {
	query := `
	SELECT user_id FROM ` + r.hireInfoTable + ` WHERE department_id = $1
	`
	var users []int64
	rows, err := r.db.Query(query, departmentId)
	if err != nil {
		log.Printf("Error getting department users IDs: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userID int64
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		users = append(users, userID)
	}
	return users, nil
}

func (r *DepartmentRepository) GetDepartmentIdByName(name string) (int64, error) {
	query := `
	SELECT id FROM ` + r.departmentTable + ` WHERE name = $1
	`
	var id int64
	err := r.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
