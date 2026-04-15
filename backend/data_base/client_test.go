package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateClientRepository(t *testing.T) {
	clientRepository := CreateClientRepository(globalDb, "test_personal_data_table", "test_user_table", "test_client_table", "test_auth_table")
	if clientRepository == nil {
		t.Errorf("Failed to create client repository")
	}
}

func TestInsertClient(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_client_table, test_auth_table CASCADE")
	clientRepository := CreateClientRepository(globalDb, "test_personal_data_table", "test_user_table", "test_client_table", "test_auth_table")
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
	clientData := types.ClientData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := clientRepository.InsertClient(clientData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert client: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	resultUserData := types.UserData{}
	globalDb.QueryRow("SELECT id, personal_data_id FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.ID, &resultUserData.PersonalDataID)
	globalDb.QueryRow("SELECT registration_date, last_login_date FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.RegistrationDate, &resultUserData.LastLoginDate)
	resultPersonalData := types.PersonalData{}
	globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1", resultUserData.PersonalDataID).Scan(&resultPersonalData.ID, &resultPersonalData.TelephoneNumber, &resultPersonalData.Email, &resultPersonalData.PassportData.PassportNumber, &resultPersonalData.PassportData.PassportSeries, &resultPersonalData.PassportData.PassportDate, &resultPersonalData.PassportData.PassportIssuedBy)
	resultClientData := types.ClientData{}
	globalDb.QueryRow("SELECT id FROM test_client_table WHERE id = $1", lastInsertedID).Scan(&resultClientData.ID)
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
	if resultClientData.ID != lastInsertedID {
		t.Errorf("Incorrect client ID: inserted %v, expected %v", resultClientData.ID, lastInsertedID)
	}
}

func TestGetClient(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_client_table, test_auth_table CASCADE")
	clientRepository := CreateClientRepository(globalDb, "test_personal_data_table", "test_user_table", "test_client_table", "test_auth_table")
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
	clientData := types.ClientData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := clientRepository.InsertClient(clientData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert client: %v", err)
	}
	resultClientData, err := clientRepository.GetClient(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
	}
	if resultClientData.ID != lastInsertedID {
		t.Errorf("Incorrect client ID: inserted %v, expected %v", resultClientData.ID, lastInsertedID)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id FROM test_personal_data_table WHERE id = $1", resultClientData.PersonalDataID).Scan(&resultPersonalData.ID)
	if err != nil {
		t.Errorf("Failed to get personal data: %v", err)
	}
	if resultPersonalData.ID != resultClientData.PersonalDataID {
		t.Errorf("Incorrect personal data ID: inserted %v, expected %v", resultPersonalData.ID, resultClientData.UserData.PersonalDataID)
	}
}

func TestUpdateClientPersonalData(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_client_table, test_auth_table CASCADE")
	clientRepository := CreateClientRepository(globalDb, "test_personal_data_table", "test_user_table", "test_client_table", "test_auth_table")
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
	clientData := types.ClientData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := clientRepository.InsertClient(clientData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert client: %v", err)
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
	err = clientRepository.UpdateClientPersonalData(lastInsertedID, newPersonalData)
	if err != nil {
		t.Errorf("Failed to update client personal data: %v", err)
	}
	resultClientData, err := clientRepository.GetClient(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1", resultClientData.PersonalDataID).Scan(&resultPersonalData.ID, &resultPersonalData.TelephoneNumber, &resultPersonalData.Email, &resultPersonalData.PassportData.PassportNumber, &resultPersonalData.PassportData.PassportSeries, &resultPersonalData.PassportData.PassportDate, &resultPersonalData.PassportData.PassportIssuedBy)
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

func TestUpdateClientPassword(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_client_table, test_auth_table CASCADE")
	clientRepository := CreateClientRepository(globalDb, "test_personal_data_table", "test_user_table", "test_client_table", "test_auth_table")
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
	clientData := types.ClientData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := clientRepository.InsertClient(clientData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert client: %v", err)
	}
	newPassword := "newpassword123"
	err = clientRepository.UpdateClientPassword(lastInsertedID, authData, newPassword)
	if err != nil {
		t.Errorf("Failed to update client password: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	if resultAuthInfo.Password != newPassword {
		t.Errorf("Incorrect password: inserted %v, expected %v", resultAuthInfo.Password, newPassword)
	}
}
