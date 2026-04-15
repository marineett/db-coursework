package service_logic

import (
	service_test "data_base_project/tests/service_logic_tests"
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateClient(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	clientRepository := service_test.CreateTestClientRepository(authRepository, personalDataRepository, userRepository)
	reviewRepository := service_test.CreateTestReviewRepository()
	clientService := CreateClientService(clientRepository, personalDataRepository, reviewRepository)
	initData := types.InitClientData{
		InitUserData: types.InitUserData{
			PersonalData: types.PersonalData{
				TelephoneNumber: "88005553535",
				Email:           "admin@admin.com",
				PassportData: types.PassportData{
					PassportNumber: "1234567890",
					PassportSeries: "1234",
					PassportDate:   time.Now(),
				},
			},
			AuthData: types.AuthData{
				Login:    "aboba",
				Password: "1234",
			},
		},
	}
	err := clientService.CreateClient(initData)
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	client, err := clientRepository.GetClient(1)
	if err != nil {
		t.Errorf("Error getting client: %v", err)
	}
	if client.UserData.RegistrationDate.After(time.Now()) {
		t.Errorf("Registration date is zero")
	}
	if client.UserData.LastLoginDate.After(time.Now()) {
		t.Errorf("Last login date is zero")
	}
	if client.UserData.LastLoginDate.Before(client.UserData.RegistrationDate) {
		t.Errorf("Last login date is before registration date")
	}
	personalData, err := personalDataRepository.GetPersonalData(client.PersonalDataID)
	if err != nil {
		t.Errorf("Error getting personal data: %v", err)
	}
	if personalData.Email != "admin@admin.com" {
		t.Errorf("Email is not admin@admin.com")
	}
	if personalData.TelephoneNumber != "88005553535" {
		t.Errorf("Telephone number is not 88005553535")
	}
	if personalData.PassportData.PassportNumber != "1234567890" {
		t.Errorf("Passport number is not 1234567890")
	}
	if personalData.PassportData.PassportSeries != "1234" {
		t.Errorf("Passport series is not 1234")
	}
	if personalData.PassportData.PassportDate.After(time.Now()) {
		t.Errorf("Passport date is after current date")
	}
	_, err = authRepository.Authorize(types.AuthData{
		Login:    "aboba",
		Password: "1234",
	})
	if err != nil {
		t.Errorf("Error getting auth: %v", err)
	}
	if personalData.ID != client.PersonalDataID {
		t.Errorf("Personal data ID is not client.PersonalDataID")
	}
}

func TestGetClientData(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	clientRepository := service_test.CreateTestClientRepository(authRepository, personalDataRepository, userRepository)
	reviewRepository := service_test.CreateTestReviewRepository()
	clientService := CreateClientService(clientRepository, personalDataRepository, reviewRepository)
	initData := types.InitClientData{
		InitUserData: types.InitUserData{
			PersonalData: types.PersonalData{
				TelephoneNumber: "88005553535",
				Email:           "admin@admin.com",
				PassportData: types.PassportData{
					PassportNumber: "1234567890",
					PassportSeries: "1234",
					PassportDate:   time.Now(),
				},
			},
			AuthData: types.AuthData{
				Login:    "aboba",
				Password: "1234",
			},
		},
	}
	err := clientService.CreateClient(initData)
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	client, err := clientService.GetClientData(1)
	if err != nil {
		t.Errorf("Error getting client: %v", err)
	}
	if client.UserData.RegistrationDate.After(time.Now()) {
		t.Errorf("Registration date is zero")
	}
	if client.UserData.LastLoginDate.After(time.Now()) {
		t.Errorf("Last login date is zero")
	}
	if client.UserData.LastLoginDate.Before(client.UserData.RegistrationDate) {
		t.Errorf("Last login date is before registration date")
	}
	personalData, err := personalDataRepository.GetPersonalData(client.PersonalDataID)
	if err != nil {
		t.Errorf("Error getting personal data: %v", err)
	}
	if personalData.ID != client.PersonalDataID {
		t.Errorf("Personal data ID is not client.PersonalDataID")
	}
}

func TestUpdateClientPersonalData(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	clientRepository := service_test.CreateTestClientRepository(authRepository, personalDataRepository, userRepository)
	reviewRepository := service_test.CreateTestReviewRepository()
	clientService := CreateClientService(clientRepository, personalDataRepository, reviewRepository)
	initData := types.InitClientData{
		InitUserData: types.InitUserData{
			PersonalData: types.PersonalData{
				TelephoneNumber: "88005553535",
				Email:           "admin@admin.com",
				PassportData: types.PassportData{
					PassportNumber: "1234567890",
					PassportSeries: "1234",
					PassportDate:   time.Now(),
				},
			},
			AuthData: types.AuthData{
				Login:    "aboba",
				Password: "1234",
			},
		},
	}
	err := clientService.CreateClient(initData)
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	personalData := types.PersonalData{
		TelephoneNumber: "88005553535",
		Email:           "admin@admin.com",
		PassportData: types.PassportData{
			PassportNumber: "1234567890",
			PassportSeries: "1234",
			PassportDate:   time.Now(),
		},
	}
	err = clientService.UpdateClientPersonalData(1, personalData)
	if err != nil {
		t.Errorf("Error updating client personal data: %v", err)
	}
	client, err := clientService.GetClientData(1)
	if err != nil {
		t.Errorf("Error getting client: %v", err)
	}
	insertedPersonalData, err := personalDataRepository.GetPersonalData(client.PersonalDataID)
	if err != nil {
		t.Errorf("Error getting personal data: %v", err)
	}
	if insertedPersonalData.Email != personalData.Email {
		t.Errorf("Email is not personalData.Email")
	}
	if insertedPersonalData.TelephoneNumber != personalData.TelephoneNumber {
		t.Errorf("Telephone number is not personalData.TelephoneNumber")
	}
	if insertedPersonalData.PassportData.PassportNumber != personalData.PassportData.PassportNumber {
		t.Errorf("Passport number is not personalData.PassportData.PassportNumber")
	}
	if insertedPersonalData.PassportData.PassportSeries != personalData.PassportData.PassportSeries {
		t.Errorf("Passport series is not personalData.PassportData.PassportSeries")
	}
	if insertedPersonalData.PassportData.PassportDate.After(personalData.PassportData.PassportDate) {
		t.Errorf("Passport date is after personalData.PassportData.PassportDate")
	}
}

func TestUpdateClientPassword(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	clientRepository := service_test.CreateTestClientRepository(authRepository, personalDataRepository, userRepository)
	reviewRepository := service_test.CreateTestReviewRepository()
	clientService := CreateClientService(clientRepository, personalDataRepository, reviewRepository)
	initData := types.InitClientData{
		InitUserData: types.InitUserData{
			PersonalData: types.PersonalData{
				TelephoneNumber: "88005553535",
				Email:           "admin@admin.com",
				PassportData: types.PassportData{
					PassportNumber: "1234567890",
					PassportSeries: "1234",
					PassportDate:   time.Now(),
				},
			},
			AuthData: types.AuthData{
				Login:    "aboba",
				Password: "1234",
			},
		},
	}
	err := clientService.CreateClient(initData)
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	_, err = authRepository.Authorize(types.AuthData{
		Login:    "aboba",
		Password: "1234",
	})
	if err != nil {
		t.Errorf("Error getting auth: %v", err)
	}
	newPassword := "12345"
	client, err := clientService.GetClientData(1)
	if err != nil {
		t.Errorf("Error getting client: %v", err)
	}
	err = clientService.UpdateClientPassword(client.UserData.ID, types.AuthData{
		Login:    "aboba",
		Password: "1234",
	}, newPassword)
	if err != nil {
		t.Errorf("Error updating client password: %v", err)
	}
	_, err = authRepository.Authorize(types.AuthData{
		Login:    "aboba",
		Password: newPassword,
	})
	if err != nil {
		t.Errorf("Error getting auth: %v", err)
	}
}
