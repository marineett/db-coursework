package service_logic

import (
	service_test "data_base_project/tests/service_logic_tests"
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateAdmin(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	adminRepository := service_test.CreateTestAdminRepository(authRepository, personalDataRepository, userRepository)
	adminService := CreateAdminService(adminRepository, personalDataRepository)
	initData := types.InitAdminData{
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
		Salary: 100000,
	}
	err := adminService.CreateAdmin(initData)
	if err != nil {
		t.Errorf("Error creating admin: %v", err)
	}
	admin, err := adminRepository.GetAdmin(1)
	if err != nil {
		t.Errorf("Error getting admin: %v", err)
	}
	if admin.Salary != 100000 {
		t.Errorf("Salary is not 100000")
	}
	if admin.DepartmentID != 0 {
		t.Errorf("Department ID is not 0")
	}
	if admin.UserData.RegistrationDate.After(time.Now()) {
		t.Errorf("Registration date is zero")
	}
	if admin.UserData.LastLoginDate.After(time.Now()) {
		t.Errorf("Last login date is zero")
	}
	if admin.UserData.LastLoginDate.Before(admin.UserData.RegistrationDate) {
		t.Errorf("Last login date is before registration date")
	}
	personalData, err := personalDataRepository.GetPersonalData(admin.PersonalDataID)
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
	if personalData.ID != admin.PersonalDataID {
		t.Errorf("Personal data ID is not admin.PersonalDataID")
	}
}

func TestGetAdminData(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	adminRepository := service_test.CreateTestAdminRepository(authRepository, personalDataRepository, userRepository)
	adminService := CreateAdminService(adminRepository, personalDataRepository)
	initData := types.InitAdminData{
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
		},
		Salary: 100000,
	}
	err := adminService.CreateAdmin(initData)
	if err != nil {
		t.Errorf("Error creating admin: %v", err)
	}
	admin, err := adminService.GetAdminData(1)
	if err != nil {
		t.Errorf("Error getting admin: %v", err)
	}
	if admin.Salary != 100000 {
		t.Errorf("Salary is not 100000")
	}
	if admin.DepartmentID != 0 {
		t.Errorf("Department ID is not 0")
	}
	if admin.UserData.RegistrationDate.After(time.Now()) {
		t.Errorf("Registration date is zero")
	}
	if admin.UserData.LastLoginDate.After(time.Now()) {
		t.Errorf("Last login date is zero")
	}
	if admin.UserData.LastLoginDate.Before(admin.UserData.RegistrationDate) {
		t.Errorf("Last login date is before registration date")
	}
	personalData, err := personalDataRepository.GetPersonalData(admin.PersonalDataID)
	if err != nil {
		t.Errorf("Error getting personal data: %v", err)
	}
	if personalData.ID != admin.PersonalDataID {
		t.Errorf("Personal data ID is not admin.PersonalDataID")
	}
}

func TestUpdateAdminPersonalData(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	adminRepository := service_test.CreateTestAdminRepository(authRepository, personalDataRepository, userRepository)
	adminService := CreateAdminService(adminRepository, personalDataRepository)
	initData := types.InitAdminData{
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
		},
		Salary: 100000,
	}
	err := adminService.CreateAdmin(initData)
	if err != nil {
		t.Errorf("Error creating admin: %v", err)
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
	err = adminService.UpdateAdminPersonalData(1, personalData)
	if err != nil {
		t.Errorf("Error updating admin personal data: %v", err)
	}
	admin, err := adminService.GetAdminData(1)
	if err != nil {
		t.Errorf("Error getting admin: %v", err)
	}
	insertedPersonalData, err := personalDataRepository.GetPersonalData(admin.PersonalDataID)
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

func TestUpdateAdminPassword(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	adminRepository := service_test.CreateTestAdminRepository(authRepository, personalDataRepository, userRepository)
	adminService := CreateAdminService(adminRepository, personalDataRepository)
	initData := types.InitAdminData{
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
		Salary: 100000,
	}
	err := adminService.CreateAdmin(initData)
	if err != nil {
		t.Errorf("Error creating admin: %v", err)
	}
	authRepository.Authorize(types.AuthData{
		Login:    "aboba",
		Password: "1234",
	})
	if err != nil {
		t.Errorf("Error getting auth: %v", err)
	}
	newPassword := "12345"
	admin, err := adminService.GetAdminData(1)
	if err != nil {
		t.Errorf("Error getting admin: %v", err)
	}
	err = adminService.UpdateAdminPassword(admin.UserData.ID, types.AuthData{
		Login:    "aboba",
		Password: "1234",
	}, newPassword)
	if err != nil {
		t.Errorf("Error updating admin password: %v", err)
	}
	_, err = authRepository.Authorize(types.AuthData{
		Login:    "aboba",
		Password: newPassword,
	})
	if err != nil {
		t.Errorf("Error getting auth: %v", err)
	}
}

func TestUpdateAdminDepartment(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	adminRepository := service_test.CreateTestAdminRepository(authRepository, personalDataRepository, userRepository)
	adminService := CreateAdminService(adminRepository, personalDataRepository)
	initData := types.InitAdminData{
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
		Salary: 100000,
	}
	err := adminService.CreateAdmin(initData)
	if err != nil {
		t.Errorf("Error creating admin: %v", err)
	}
	newDepartmentID := int64(11)
	err = adminService.UpdateAdminDepartment(1, newDepartmentID)
	if err != nil {
		t.Errorf("Error updating admin department: %v", err)
	}
	admin, err := adminService.GetAdminData(1)
	if err != nil {
		t.Errorf("Error getting admin: %v", err)
	}
	if admin.DepartmentID != newDepartmentID {
		t.Errorf("Department ID is not newDepartmentID")
	}
}

func TestUpdateAdminSalary(t *testing.T) {
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	adminRepository := service_test.CreateTestAdminRepository(authRepository, personalDataRepository, userRepository)
	adminService := CreateAdminService(adminRepository, personalDataRepository)
	initData := types.InitAdminData{
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
		Salary: 100000,
	}
	err := adminService.CreateAdmin(initData)
	if err != nil {
		t.Errorf("Error creating admin: %v", err)
	}
	newSalary := int64(110000)
	err = adminService.UpdateAdminSalary(1, newSalary)
	if err != nil {
		t.Errorf("Error updating admin salary: %v", err)
	}
	admin, err := adminService.GetAdminData(1)
	if err != nil {
		t.Errorf("Error getting admin: %v", err)
	}
	if admin.Salary != newSalary {
		t.Errorf("Salary is not newSalary")
	}
}
