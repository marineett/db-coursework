package service_logic

import (
	service_test "data_base_project/tests/service_logic_tests"
	"data_base_project/types"
	"testing"
)

func TestAuthorize(t *testing.T) {
	authRepository := service_test.CreateTestAuthRepository()
	authService := CreateAuthService(authRepository)
	userId := int64(1)
	authData := types.AuthData{
		Login:    "aboba",
		Password: "123456",
	}
	_, err := authRepository.InsertAuth(types.AuthInfo{
		UserID:   userId,
		Login:    authData.Login,
		Password: authData.Password,
	})
	if err != nil {
		t.Errorf("Error authorizing: %v", err)
	}
	verdict, err := authService.Authorize(authData)
	if err != nil {
		t.Errorf("Error authorizing: %v", err)
	}
	if verdict.UserID != userId {
		t.Errorf("Expected id %d, got %d", userId, verdict.UserID)
	}
	authData.Password = "1234567"
	_, err = authService.Authorize(authData)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	authData.Password = "123456"
	authData.Login = "abobus"
	_, err = authService.Authorize(authData)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
func TestCheckLogin(t *testing.T) {
	authRepository := service_test.CreateTestAuthRepository()
	authService := CreateAuthService(authRepository)
	userId := int64(1)
	authData := types.AuthData{
		Login:    "aboba",
		Password: "123456",
	}
	_, err := authRepository.InsertAuth(types.AuthInfo{
		UserID:   userId,
		Login:    authData.Login,
		Password: authData.Password,
	})
	if err != nil {
		t.Errorf("Error authorizing: %v", err)
	}
	ok, err := authService.CheckLogin(authData.Login)
	if err != nil {
		t.Errorf("Error checking login: %v", err)
	}
	if !ok {
		t.Errorf("Expected login to be true, got false")
	}
	ok, err = authService.CheckLogin("abobus")
	if err != nil {
		t.Errorf("Error checking login: %v", err)
	}
	if ok {
		t.Errorf("Expected login to be false, got true")
	}
}
