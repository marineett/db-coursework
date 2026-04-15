package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateUserRepository(t *testing.T) {
	userRepository := CreateUserRepository(globalDb, "test_user_table")
	if userRepository == nil {
		t.Errorf("Failed to create user repository")
	}
}

func TestInsertUser(t *testing.T) {
	InsertTestPesonalData(1)
	userRepository := CreateUserRepository(globalDb, "test_user_table")
	user := types.UserData{
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
		PersonalDataID:   1,
	}
	insertedID, err := userRepository.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert user: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table CASCADE")
	resultUser := types.UserData{}
	err = globalDb.QueryRow("SELECT * FROM test_user_table WHERE id = $1", insertedID).Scan(&resultUser.ID, &resultUser.RegistrationDate, &resultUser.LastLoginDate, &resultUser.PersonalDataID)
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	if resultUser.RegistrationDate.After(resultUser.LastLoginDate) {
		t.Errorf("Expected registration date to be before last login date")
	}
	if resultUser.PersonalDataID != 1 {
		t.Errorf("Expected personal data ID to be 1, got %d", resultUser.PersonalDataID)
	}
}

func TestInsertUserInSeq(t *testing.T) {
	InsertTestPesonalData(1)
	tx, err := globalDb.Begin()
	userRepository := CreateUserRepository(globalDb, "test_user_table")
	if err != nil {
		t.Errorf("Fail to begin transaction: %v", err)
	}
	user := types.UserData{
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
		PersonalDataID:   1,
	}
	insertedID, err := userRepository.InsertUserInSeq(tx, user)
	if err != nil {
		t.Errorf("Failed to insert user: %v", err)
	}
	tx.Commit()
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table CASCADE")
	resultUser := types.UserData{}
	err = globalDb.QueryRow("SELECT * FROM test_user_table WHERE id = $1", insertedID).Scan(&resultUser.ID, &resultUser.RegistrationDate, &resultUser.LastLoginDate, &resultUser.PersonalDataID)
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	if resultUser.RegistrationDate.After(resultUser.LastLoginDate) {
		t.Errorf("Expected registration date to be before last login date")
	}
	if resultUser.PersonalDataID != 1 {
		t.Errorf("Expected personal data ID to be 1, got %d", resultUser.PersonalDataID)
	}
}

func TestGetUser(t *testing.T) {
	InsertTestPesonalData(1)
	userRepository := CreateUserRepository(globalDb, "test_user_table")
	user := types.UserData{
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
		PersonalDataID:   1,
	}
	insertedID, err := userRepository.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert user: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_personal_data_table, test_user_table CASCADE")
	resultUser, err := userRepository.GetUser(insertedID)
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	if resultUser.ID != insertedID {
		t.Errorf("Expected user ID to be %d, got %d", insertedID, resultUser.ID)
	}
	if resultUser.PersonalDataID != 1 {
		t.Errorf("Expected personal data ID to be 1, got %d", resultUser.PersonalDataID)
	}
}
