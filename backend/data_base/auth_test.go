package data_base

import (
	"data_base_project/types"
	"testing"
)

func TestCreateAuthRepository(t *testing.T) {
	authRepository := CreateAuthRepository(globalDb, "test_auth_table")
	if authRepository == nil {
		t.Errorf("Failed to create auth repository")
	}
}

func TestInsertAuth(t *testing.T) {
	InsertTestUser(1)
	authRepository := CreateAuthRepository(globalDb, "test_auth_table")
	auth := types.AuthInfo{
		UserID:   1,
		Login:    "test123",
		Password: "test456",
	}
	_, err := authRepository.InsertAuth(auth)
	if err != nil {
		t.Errorf("Failed to insert auth: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_auth_table, test_personal_data_table, test_user_table CASCADE")
	result := globalDb.QueryRow("SELECT login, password FROM test_auth_table WHERE user_id = 1")
	if result == nil {
		t.Errorf("Expected result, got nil")
	}
	var login string
	var password string
	err = result.Scan(&login, &password)
	if err != nil {
		t.Errorf("Failed to scan result: %v", err)
	}
	if login != "test123" {
		t.Errorf("Expected username %s, got %s", "test123", login)
	}
	if password != "test456" {
		t.Errorf("Expected password %s, got %s", "test456", password)
	}
}

func TestInsertAuthInSeq(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_auth_table, test_personal_data_table, test_user_table CASCADE")
	tx, err := globalDb.Begin()
	if err != nil {
		t.Errorf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()
	authRepository := CreateAuthRepository(globalDb, "test_auth_table")
	auth := types.AuthInfo{
		UserID:   1,
		Login:    "test123",
		Password: "test456",
	}
	_, err = authRepository.InsertAuthInSeq(tx, auth)
	if err != nil {
		t.Errorf("Failed to insert auth: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("Failed to commit transaction: %v", err)
	}

	result := globalDb.QueryRow("SELECT login, password FROM test_auth_table WHERE user_id = 1")
	if result == nil {
		t.Errorf("Expected result, got nil")
	}

}

func TestAuthorize(t *testing.T) {
	InsertTestUser(1)
	authRepository := CreateAuthRepository(globalDb, "test_auth_table")
	auth := types.AuthInfo{
		UserID:   1,
		UserType: types.Unauthorized,
		Login:    "test123",
		Password: "test456",
	}
	_, err := authRepository.InsertAuth(auth)
	if err != nil {
		t.Errorf("Failed to insert auth: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_auth_table, test_personal_data_table, test_user_table CASCADE")
	result, err := authRepository.Authorize(types.AuthData{Login: "test123", Password: "test456"})
	if err != nil {
		t.Errorf("Failed to authorize: %v", err)
	}
	if result.UserID != 1 {
		t.Errorf("Expected user id %d, got %d", 1, result.UserID)
	}
	if result.UserType != types.Unauthorized {
		t.Errorf("Expected user type %d, got %d", types.Unauthorized, result.UserType)
	}
}

func TestChangePassword(t *testing.T) {
	InsertTestUser(1)
	authRepository := CreateAuthRepository(globalDb, "test_auth_table")
	auth := types.AuthInfo{
		UserID:   1,
		Login:    "test123",
		Password: "test456",
	}
	_, err := authRepository.InsertAuth(auth)
	if err != nil {
		t.Errorf("Failed to insert auth: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_auth_table, test_personal_data_table, test_user_table CASCADE")
	err = authRepository.ChangePassword(1, types.AuthData{Login: "test123", Password: "test456"}, "test789")
	if err != nil {
		t.Errorf("Failed to change password: %v", err)
	}
	result := globalDb.QueryRow("SELECT login, password FROM test_auth_table WHERE user_id = 1")
	if err != nil {
		t.Errorf("Failed to get auth: %v", err)
	}
	if result == nil {
		t.Errorf("Expected result, got nil")
	}
	var login string
	var password string
	err = result.Scan(&login, &password)
	if err != nil {
		t.Errorf("Failed to scan result: %v", err)
	}
	if login != "test123" {
		t.Errorf("Expected username %s, got %s", "test123", login)
	}
	if password != "test789" {
		t.Errorf("Expected password %s, got %s", "test789", password)
	}
}

func TestCheckLogin(t *testing.T) {
	InsertTestUser(1)
	defer globalDb.Exec("TRUNCATE TABLE test_auth_table, test_personal_data_table, test_user_table CASCADE")
	authRepository := CreateAuthRepository(globalDb, "test_auth_table")
	auth := types.AuthInfo{
		UserID:   1,
		Login:    "test123",
		Password: "test456",
	}
	_, err := authRepository.InsertAuth(auth)
	if err != nil {
		t.Errorf("Failed to insert auth: %v", err)
	}
	result, err := authRepository.CheckLogin("test123")
	if err != nil {
		t.Errorf("Failed to check login: %v", err)
	}
	if result != true {
		t.Errorf("Expected true, got false")
	}
	result, err = authRepository.CheckLogin("test456")
	if err != nil {
		t.Errorf("Failed to check login: %v", err)
	}
	if result != false {
		t.Errorf("Expected false, got true")
	}
}
