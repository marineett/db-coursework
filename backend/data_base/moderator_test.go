package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateModeratorRepository(t *testing.T) {
	moderatorRepository := CreateModeratorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_moderator_table", "test_auth_table")
	if moderatorRepository == nil {
		t.Errorf("Failed to create moderator repository")
	}
}

func TestInsertModerator(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_moderator_table, test_auth_table CASCADE")
	moderatorRepository := CreateModeratorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_moderator_table", "test_auth_table")
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
	moderatorData := types.ModeratorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := moderatorRepository.InsertModerator(moderatorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert moderator: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	resultUserData := types.UserData{}
	globalDb.QueryRow("SELECT id, personal_data_id FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.ID, &resultUserData.PersonalDataID)
	globalDb.QueryRow("SELECT registration_date, last_login_date FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.RegistrationDate, &resultUserData.LastLoginDate)
	resultPersonalData := types.PersonalData{}
	globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1", resultUserData.PersonalDataID).Scan(&resultPersonalData.ID, &resultPersonalData.TelephoneNumber, &resultPersonalData.Email, &resultPersonalData.PassportData.PassportNumber, &resultPersonalData.PassportData.PassportSeries, &resultPersonalData.PassportData.PassportDate, &resultPersonalData.PassportData.PassportIssuedBy)
	resultModeratorData := types.ModeratorData{}
	globalDb.QueryRow("SELECT id FROM test_moderator_table WHERE id = $1", lastInsertedID).Scan(&resultModeratorData.ID)
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
	if resultModeratorData.ID != lastInsertedID {
		t.Errorf("Incorrect moderator ID: inserted %v, expected %v", resultModeratorData.ID, lastInsertedID)
	}
}

func TestGetModerator(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_moderator_table, test_auth_table CASCADE")
	moderatorRepository := CreateModeratorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_moderator_table", "test_auth_table")
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
	moderatorData := types.ModeratorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := moderatorRepository.InsertModerator(moderatorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert moderator: %v", err)
	}
	resultModeratorData, err := moderatorRepository.GetModerator(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get moderator: %v", err)
	}
	if resultModeratorData.ID != lastInsertedID {
		t.Errorf("Incorrect moderator ID: inserted %v, expected %v", resultModeratorData.ID, lastInsertedID)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id FROM test_personal_data_table WHERE id = $1", resultModeratorData.UserData.PersonalDataID).Scan(&resultPersonalData.ID)
	if err != nil {
		t.Errorf("Failed to get personal data: %v", err)
	}
	if resultPersonalData.ID != resultModeratorData.UserData.PersonalDataID {
		t.Errorf("Incorrect personal data ID: inserted %v, expected %v", resultPersonalData.ID, resultModeratorData.UserData.PersonalDataID)
	}
}

func TestUpdateModeratorPersonalData(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_moderator_table, test_auth_table CASCADE")
	moderatorRepository := CreateModeratorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_moderator_table", "test_auth_table")
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
	moderatorData := types.ModeratorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := moderatorRepository.InsertModerator(moderatorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert moderator: %v", err)
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
	err = moderatorRepository.UpdateModeratorPersonalData(lastInsertedID, newPersonalData)
	if err != nil {
		t.Errorf("Failed to update moderator personal data: %v", err)
	}
	resultModeratorData, err := moderatorRepository.GetModerator(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get moderator: %v", err)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1", resultModeratorData.UserData.PersonalDataID).Scan(&resultPersonalData.ID, &resultPersonalData.TelephoneNumber, &resultPersonalData.Email, &resultPersonalData.PassportData.PassportNumber, &resultPersonalData.PassportData.PassportSeries, &resultPersonalData.PassportData.PassportDate, &resultPersonalData.PassportData.PassportIssuedBy)
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

func TestUpdateModeratorPassword(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_moderator_table, test_auth_table CASCADE")
	moderatorRepository := CreateModeratorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_moderator_table", "test_auth_table")
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
	moderatorData := types.ModeratorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := moderatorRepository.InsertModerator(moderatorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert moderator: %v", err)
	}
	newPassword := "newpassword123"
	err = moderatorRepository.UpdateModeratorPassword(lastInsertedID, authData, newPassword)
	if err != nil {
		t.Errorf("Failed to update moderator password: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	if resultAuthInfo.Password != newPassword {
		t.Errorf("Incorrect password: inserted %v, expected %v", resultAuthInfo.Password, newPassword)
	}
}

func TestUpdateModeratorSalary(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_moderator_table, test_auth_table CASCADE")
	moderatorRepository := CreateModeratorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_moderator_table", "test_auth_table")
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
	moderatorData := types.ModeratorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
		Salary: 100000,
	}
	lastInsertedID, err := moderatorRepository.InsertModerator(moderatorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert moderator: %v", err)
	}
	newSalary := int64(150000)
	err = moderatorRepository.UpdateModeratorSalary(lastInsertedID, newSalary)
	if err != nil {
		t.Errorf("Failed to update moderator salary: %v", err)
	}
	resultModeratorData, err := moderatorRepository.GetModerator(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get moderator: %v", err)
	}
	if resultModeratorData.Salary != newSalary {
		t.Errorf("Incorrect salary: inserted %v, expected %v", resultModeratorData.Salary, newSalary)
	}
}
