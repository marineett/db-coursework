package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IAdminRepository interface {
	InsertAdmin(admin types.AdminData, personalData types.PersonalData, auth types.AuthData) (int64, error)
	GetAdmin(id int64) (*types.AdminData, error)
	UpdateAdminPersonalData(adminId int64, personalData types.PersonalData) error
	UpdateAdminPassword(adminId int64, authData types.AuthData, newPassword string) error
	UpdateAdminDepartment(adminId int64, departmentId int64) error
	UpdateAdminSalary(adminId int64, salary int64) error
}

func CreateAdminTable(db *sql.DB, adminTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + adminTableName + ` (
		id SERIAL PRIMARY KEY,
		department_id INTEGER NOT NULL,
		salary INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", adminTableName, err)
	}
	return nil
}

type AdminRepository struct {
	db                     *sql.DB
	adminTable             string
	userRepository         IUserRepository
	personalDataRepository IPersonalDataRepository
	authRepository         IAuthRepository
}

func CreateAdminRepository(db *sql.DB, personalDataTable string, userTable string, adminTable string, authTable string) *AdminRepository {
	return &AdminRepository{
		db:                     db,
		adminTable:             adminTable,
		userRepository:         CreateUserRepository(db, userTable),
		personalDataRepository: CreatePersonalDataRepository(db, personalDataTable),
		authRepository:         CreateAuthRepository(db, authTable),
	}
}

func (r *AdminRepository) InsertAdmin(admin types.AdminData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	admin.PersonalDataID, err = r.personalDataRepository.InsertPersonalDataInSeq(tx, personalData)
	if err != nil {
		return 0, err
	}
	admin.ID, err = r.userRepository.InsertUserInSeq(tx, admin.UserData)
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, types.AuthInfo{
		UserID:   admin.ID,
		UserType: types.Admin,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.adminTable + ` (id, department_id, salary)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	_, err = tx.Exec(query, admin.ID, admin.DepartmentID, admin.Salary)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return admin.ID, nil
}

func (r *AdminRepository) GetAdmin(id int64) (*types.AdminData, error) {
	query := `
	SELECT id, department_id, salary FROM ` + r.adminTable + ` WHERE id = $1
	`
	var admin types.AdminData
	err := r.db.QueryRow(query, id).Scan(&admin.ID, &admin.DepartmentID, &admin.Salary)
	if err != nil {
		return nil, err
	}
	userData, err := r.userRepository.GetUser(admin.ID)
	if err != nil {
		return nil, err
	}
	admin.UserData = *userData
	return &admin, nil
}

func (r *AdminRepository) UpdateAdminPersonalData(adminId int64, personalData types.PersonalData) error {
	admin, err := r.GetAdmin(adminId)
	if err != nil {
		return err
	}
	return r.personalDataRepository.UpdatePersonalData(admin.PersonalDataID, personalData)
}

func (r *AdminRepository) UpdateAdminPassword(adminId int64, authData types.AuthData, newPassword string) error {
	return r.authRepository.ChangePassword(adminId, authData, newPassword)
}

func (r *AdminRepository) UpdateAdminDepartment(adminId int64, departmentId int64) error {
	query := `
	UPDATE ` + r.adminTable + ` SET department_id = $1 WHERE id = $2
	`
	_, err := r.db.Exec(query, departmentId, adminId)
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminRepository) UpdateAdminSalary(admin_id int64, salary int64) error {
	query := `
	UPDATE ` + r.adminTable + ` SET salary = $1 WHERE id = $2
	`
	_, err := r.db.Exec(query, salary, admin_id)
	if err != nil {
		return err
	}
	return nil
}
