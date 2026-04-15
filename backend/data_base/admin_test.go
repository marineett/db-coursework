package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateAdminRepository(t *testing.T) {
	adminRepository := CreateAdminRepository(globalDb, "test_personal_data_table", "test_user_table", "test_admin_table", "test_auth_table")
	if adminRepository == nil {
		t.Errorf("Failed to create admin repository")
	}
}

func TestInsertAdmin(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_admin_table, test_auth_table CASCADE")
	adminRepository := CreateAdminRepository(globalDb, "test_personal_data_table", "test_user_table", "test_admin_table", "test_auth_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
	}
	authData := types.AuthData{
		Login:    "test123",
		Password: "test456",
	}
	adminData := types.AdminData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := adminRepository.InsertAdmin(adminData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert admin: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	resultUserData := types.UserData{}
	globalDb.QueryRow("SELECT id, personal_data_id FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.ID, &resultUserData.PersonalDataID)
	globalDb.QueryRow("SELECT registration_date, last_login_date FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.RegistrationDate, &resultUserData.LastLoginDate)
	resultPersonalData := types.PersonalData{}
	globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1", resultUserData.PersonalDataID).Scan(&resultPersonalData.ID, &resultPersonalData.TelephoneNumber, &resultPersonalData.Email, &resultPersonalData.PassportData.PassportNumber, &resultPersonalData.PassportData.PassportSeries, &resultPersonalData.PassportData.PassportDate, &resultPersonalData.PassportData.PassportIssuedBy)
	resultAdminData := types.AdminData{}
	globalDb.QueryRow("SELECT id FROM test_admin_table WHERE id = $1", lastInsertedID).Scan(&resultAdminData.ID)
	if resultAuthInfo.Login != authData.Login {
		t.Errorf("Incorrect login: inserted %v, expected %v", resultAuthInfo.Login, authData.Login)
	}
	if resultAuthInfo.Password != authData.Password {
		t.Errorf("Incorrect password: inserted %v, expected %v", resultAuthInfo.Password, authData.Password)
	}
	if resultAuthInfo.UserID != lastInsertedID {
		t.Errorf("Incorrect user ID: inserted %v, expected %v", resultAuthInfo.UserID, lastInsertedID)
	}
	if resultUserData.PersonalDataID != resultPersonalData.ID {
		t.Errorf("Incorrect personal data ID: inserted %v, expected %v", resultUserData.PersonalDataID, resultPersonalData.ID)
	}
	if resultPersonalData.TelephoneNumber != personalData.TelephoneNumber {
		t.Errorf("Incorrect telephone number: inserted %v, expected %v", resultPersonalData.TelephoneNumber, personalData.TelephoneNumber)
	}
	if resultPersonalData.Email != personalData.Email {
		t.Errorf("Incorrect email: inserted %v, expected %v", resultPersonalData.Email, personalData.Email)
	}
	if resultPersonalData.PassportData.PassportNumber != personalData.PassportData.PassportNumber {
		t.Errorf("Incorrect passport number: inserted %v, expected %v", resultPersonalData.PassportData.PassportNumber, personalData.PassportData.PassportNumber)
	}
	if resultPersonalData.PassportData.PassportSeries != personalData.PassportData.PassportSeries {
		t.Errorf("Incorrect passport series: inserted %v, expected %v", resultPersonalData.PassportData.PassportSeries, personalData.PassportData.PassportSeries)
	}
	if resultPersonalData.PassportData.PassportIssuedBy != personalData.PassportData.PassportIssuedBy {
		t.Errorf("Incorrect passport issued by: inserted %v, expected %v", resultPersonalData.PassportData.PassportIssuedBy, personalData.PassportData.PassportIssuedBy)
	}
	if resultAdminData.ID != lastInsertedID {
		t.Errorf("Incorrect admin ID: inserted %v, expected %v", resultAdminData.ID, lastInsertedID)
	}
}

func TestGetAdmin(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_admin_table, test_auth_table CASCADE")
	adminRepository := CreateAdminRepository(globalDb, "test_personal_data_table", "test_user_table", "test_admin_table", "test_auth_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
	}
	authData := types.AuthData{
		Login:    "test123",
		Password: "test456",
	}
	adminData := types.AdminData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := adminRepository.InsertAdmin(adminData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert admin: %v", err)
	}
	resultAdminData, err := adminRepository.GetAdmin(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get admin: %v", err)
	}
	if resultAdminData.ID != lastInsertedID {
		t.Errorf("Incorrect admin ID: inserted %v, expected %v", resultAdminData.ID, lastInsertedID)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id FROM test_personal_data_table WHERE id = $1", resultAdminData.UserData.PersonalDataID).Scan(&resultPersonalData.ID)
	if err != nil {
		t.Errorf("Failed to get personal data: %v", err)
	}
	if resultPersonalData.ID != resultAdminData.UserData.PersonalDataID {
		t.Errorf("Incorrect personal data ID: inserted %v, expected %v", resultPersonalData.ID, resultAdminData.UserData.PersonalDataID)
	}
}

func TestUpdateAdminPersonalData(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_admin_table, test_auth_table CASCADE")
	adminRepository := CreateAdminRepository(globalDb, "test_personal_data_table", "test_user_table", "test_admin_table", "test_auth_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
	}
	authData := types.AuthData{
		Login:    "test123",
		Password: "test456",
	}
	adminData := types.AdminData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := adminRepository.InsertAdmin(adminData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert admin: %v", err)
	}
	newPersonalData := types.PersonalData{
		TelephoneNumber: "+99999999999",
		Email:           "test2@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "0987654321",
			PassportDate:     time.Now(),
			PassportSeries:   "2048",
			PassportIssuedBy: "test2",
		},
	}
	err = adminRepository.UpdateAdminPersonalData(lastInsertedID, newPersonalData)
	if err != nil {
		t.Errorf("Failed to update admin personal data: %v", err)
	}
	resultAdminData, err := adminRepository.GetAdmin(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get admin: %v", err)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1", resultAdminData.UserData.PersonalDataID).Scan(&resultPersonalData.ID, &resultPersonalData.TelephoneNumber, &resultPersonalData.Email, &resultPersonalData.PassportData.PassportNumber, &resultPersonalData.PassportData.PassportSeries, &resultPersonalData.PassportData.PassportDate, &resultPersonalData.PassportData.PassportIssuedBy)
	if err != nil {
		t.Errorf("Failed to get personal data: %v", err)
	}
	if resultPersonalData.TelephoneNumber != newPersonalData.TelephoneNumber {
		t.Errorf("Incorrect telephone number: inserted %v, expected %v", resultPersonalData.TelephoneNumber, newPersonalData.TelephoneNumber)
	}
	if resultPersonalData.Email != newPersonalData.Email {
		t.Errorf("Incorrect email: inserted %v, expected %v", resultPersonalData.Email, newPersonalData.Email)
	}
	if resultPersonalData.PassportData.PassportNumber != newPersonalData.PassportData.PassportNumber {
		t.Errorf("Incorrect passport number: inserted %v, expected %v", resultPersonalData.PassportData.PassportNumber, newPersonalData.PassportData.PassportNumber)
	}
	if resultPersonalData.PassportData.PassportSeries != newPersonalData.PassportData.PassportSeries {
		t.Errorf("Incorrect passport series: inserted %v, expected %v", resultPersonalData.PassportData.PassportSeries, newPersonalData.PassportData.PassportSeries)
	}
	if resultPersonalData.PassportData.PassportIssuedBy != newPersonalData.PassportData.PassportIssuedBy {
		t.Errorf("Incorrect passport issued by: inserted %v, expected %v", resultPersonalData.PassportData.PassportIssuedBy, newPersonalData.PassportData.PassportIssuedBy)
	}
}

func TestUpdateAdminPassword(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_admin_table, test_auth_table CASCADE")
	adminRepository := CreateAdminRepository(globalDb, "test_personal_data_table", "test_user_table", "test_admin_table", "test_auth_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
	}
	authData := types.AuthData{
		Login:    "test123",
		Password: "test456",
	}
	adminData := types.AdminData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := adminRepository.InsertAdmin(adminData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert admin: %v", err)
	}
	newPassword := "newpassword123"
	err = adminRepository.UpdateAdminPassword(lastInsertedID, authData, newPassword)
	if err != nil {
		t.Errorf("Failed to update admin password: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	if resultAuthInfo.Password != newPassword {
		t.Errorf("Incorrect password: inserted %v, expected %v", resultAuthInfo.Password, newPassword)
	}
}

func TestUpdateAdminDepartment(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_admin_table, test_auth_table CASCADE")
	adminRepository := CreateAdminRepository(globalDb, "test_personal_data_table", "test_user_table", "test_admin_table", "test_auth_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
	}
	authData := types.AuthData{
		Login:    "test123",
		Password: "test456",
	}
	adminData := types.AdminData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
		DepartmentID: 1,
		Salary:       100000,
	}
	lastInsertedID, err := adminRepository.InsertAdmin(adminData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert admin: %v", err)
	}
	newDepartmentID := int64(2)
	err = adminRepository.UpdateAdminDepartment(lastInsertedID, newDepartmentID)
	if err != nil {
		t.Errorf("Failed to update admin department: %v", err)
	}
	resultAdminData, err := adminRepository.GetAdmin(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get admin: %v", err)
	}
	if resultAdminData.DepartmentID != newDepartmentID {
		t.Errorf("Incorrect department ID: inserted %v, expected %v", resultAdminData.DepartmentID, newDepartmentID)
	}
}

func TestUpdateAdminSalary(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_admin_table, test_auth_table CASCADE")
	adminRepository := CreateAdminRepository(globalDb, "test_personal_data_table", "test_user_table", "test_admin_table", "test_auth_table")
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
	}
	authData := types.AuthData{
		Login:    "test123",
		Password: "test456",
	}
	adminData := types.AdminData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
		DepartmentID: 1,
		Salary:       100000,
	}
	lastInsertedID, err := adminRepository.InsertAdmin(adminData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert admin: %v", err)
	}
	newSalary := int64(150000)
	err = adminRepository.UpdateAdminSalary(lastInsertedID, newSalary)
	if err != nil {
		t.Errorf("Failed to update admin salary: %v", err)
	}
	resultAdminData, err := adminRepository.GetAdmin(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get admin: %v", err)
	}
	if resultAdminData.Salary != newSalary {
		t.Errorf("Incorrect salary: inserted %v, expected %v", resultAdminData.Salary, newSalary)
	}
}
