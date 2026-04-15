package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateRepetitorRepository(t *testing.T) {
	repetitorRepository := CreateRepetitorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_repetitor_table", "test_auth_table", "test_resume_table", "test_review_table")
	if repetitorRepository == nil {
		t.Errorf("Failed to create repetitor repository")
	}
}

func TestInsertRepetitor(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_repetitor_table, test_auth_table, test_resume_table, test_review_table CASCADE")
	repetitorRepository := CreateRepetitorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_repetitor_table", "test_auth_table", "test_resume_table", "test_review_table")
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
	repetitorData := types.RepetitorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := repetitorRepository.InsertRepetitor(repetitorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert repetitor: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	resultUserData := types.UserData{}
	globalDb.QueryRow("SELECT id, personal_data_id FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.ID, &resultUserData.PersonalDataID)
	globalDb.QueryRow("SELECT registration_date, last_login_date FROM test_user_table WHERE id = $1", lastInsertedID).Scan(&resultUserData.RegistrationDate, &resultUserData.LastLoginDate)
	resultPersonalData := types.PersonalData{}
	globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1",
		resultUserData.PersonalDataID).Scan(
		&resultPersonalData.ID,
		&resultPersonalData.TelephoneNumber,
		&resultPersonalData.Email,
		&resultPersonalData.PassportData.PassportNumber,
		&resultPersonalData.PassportData.PassportSeries,
		&resultPersonalData.PassportData.PassportDate,
		&resultPersonalData.PassportData.PassportIssuedBy)
	if err != nil {
		t.Errorf("Failed to get personal data: %v", err)
	}
	resultRepetitorData := types.RepetitorData{}
	summaryRating := 0
	reviewsCount := 0
	globalDb.QueryRow("SELECT id, resume_id, summary_rating, reviews_count FROM test_repetitor_table WHERE id = $1", lastInsertedID).Scan(&resultRepetitorData.ID, &resultRepetitorData.ResumeID, &summaryRating, &reviewsCount)

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
	if resultRepetitorData.ID != lastInsertedID {
		t.Errorf("Incorrect repetitor ID: inserted %v, expected %v", resultRepetitorData.ID, lastInsertedID)
	}
	if reviewsCount != 0 {
		t.Errorf("Incorrect reviews count: inserted %v, expected %v", reviewsCount, 0)
	}
	if summaryRating != 0 {
		t.Errorf("Incorrect summary rating: inserted %v, expected %v", summaryRating, 0)
	}
}

func TestGetRepetitor(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_repetitor_table, test_auth_table, test_resume_table, test_review_table CASCADE")
	repetitorRepository := CreateRepetitorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_repetitor_table", "test_auth_table", "test_resume_table", "test_review_table")
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
	repetitorData := types.RepetitorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := repetitorRepository.InsertRepetitor(repetitorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert repetitor: %v", err)
	}
	resultRepetitorData, err := repetitorRepository.GetRepetitor(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get repetitor: %v", err)
	}
	if resultRepetitorData.ID != lastInsertedID {
		t.Errorf("Incorrect repetitor ID: inserted %v, expected %v", resultRepetitorData.ID, lastInsertedID)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id FROM test_personal_data_table WHERE id = $1", resultRepetitorData.PersonalDataID).Scan(&resultPersonalData.ID)
	if err != nil {
		t.Errorf("Failed to get personal data: %v", err)
	}
	if resultPersonalData.ID != resultRepetitorData.PersonalDataID {
		t.Errorf("Incorrect personal data ID: inserted %v, expected %v", resultPersonalData.ID, resultRepetitorData.PersonalDataID)
	}
}

func TestUpdateRepetitorPersonalData(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_repetitor_table, test_auth_table, test_resume_table, test_review_table CASCADE")
	repetitorRepository := CreateRepetitorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_repetitor_table", "test_auth_table", "test_resume_table", "test_review_table")
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
	repetitorData := types.RepetitorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := repetitorRepository.InsertRepetitor(repetitorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert repetitor: %v", err)
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
	err = repetitorRepository.UpdateRepetitorPersonalData(lastInsertedID, newPersonalData)
	if err != nil {
		t.Errorf("Failed to update repetitor personal data: %v", err)
	}
	resultRepetitorData, err := repetitorRepository.GetRepetitor(lastInsertedID)
	if err != nil {
		t.Errorf("Failed to get repetitor: %v", err)
	}
	resultPersonalData := types.PersonalData{}
	err = globalDb.QueryRow("SELECT id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by FROM test_personal_data_table WHERE id = $1", resultRepetitorData.PersonalDataID).Scan(&resultPersonalData.ID, &resultPersonalData.TelephoneNumber, &resultPersonalData.Email, &resultPersonalData.PassportData.PassportNumber, &resultPersonalData.PassportData.PassportSeries, &resultPersonalData.PassportData.PassportDate, &resultPersonalData.PassportData.PassportIssuedBy)
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

func TestUpdateRepetitorPassword(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table, test_repetitor_table, test_auth_table, test_resume_table, test_review_table CASCADE")
	repetitorRepository := CreateRepetitorRepository(globalDb, "test_personal_data_table", "test_user_table", "test_repetitor_table", "test_auth_table", "test_resume_table", "test_review_table")
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
	repetitorData := types.RepetitorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	lastInsertedID, err := repetitorRepository.InsertRepetitor(repetitorData, personalData, authData)
	if err != nil {
		t.Errorf("Failed to insert repetitor: %v", err)
	}
	newPassword := "newpassword123"
	err = repetitorRepository.UpdateRepetitorPassword(lastInsertedID, authData, newPassword)
	if err != nil {
		t.Errorf("Failed to update repetitor password: %v", err)
	}
	resultAuthInfo := types.AuthInfo{}
	globalDb.QueryRow("SELECT login, password, user_id FROM test_auth_table WHERE login = $1", authData.Login).Scan(&resultAuthInfo.Login, &resultAuthInfo.Password, &resultAuthInfo.UserID)
	if resultAuthInfo.Password != newPassword {
		t.Errorf("Incorrect password: inserted %v, expected %v", resultAuthInfo.Password, newPassword)
	}
}
