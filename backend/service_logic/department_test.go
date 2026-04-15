package service_logic

import (
	service_test "data_base_project/tests/service_logic_tests"
	"data_base_project/types"
	"testing"
)

func TestCreateDepartment(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	authRepository := service_test.CreateTestAuthRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
}

func TestGetDepartment(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	authRepository := service_test.CreateTestAuthRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	newDepartment, err := departmentService.GetDepartment(1)
	if err != nil {
		t.Errorf("Error getting department: %v", err)
	}
	if newDepartment.Name != department.Name {
		t.Errorf("Department name mismatch: %v != %v", newDepartment.Name, department.Name)
	}
	if newDepartment.HeadID != department.HeadID {
		t.Errorf("Department head ID mismatch: %v != %v", newDepartment.HeadID, department.HeadID)
	}
}

func TestAssignAdminToDepartment(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	authRepository := service_test.CreateTestAuthRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	adminId := int64(1)
	err = departmentService.AssignAdminToDepartment(adminId, 1)
	if err != nil {
		t.Errorf("Error assigning admin to department: %v", err)
	}
	newDepartment, err := departmentService.GetDepartment(1)
	if err != nil {
		t.Errorf("Error getting department: %v", err)
	}
	if newDepartment.HeadID != adminId {
		t.Errorf("Department head ID mismatch: %v != %v", newDepartment.HeadID, adminId)
	}
}

func TestAssignModeratorToDepartment(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	authRepository := service_test.CreateTestAuthRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	moderatorId := int64(2)
	err = departmentService.AssignModeratorToDepartment(moderatorId, 1)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
}

func TestFireAdminFromDepartment(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	authRepository := service_test.CreateTestAuthRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	adminId := int64(1)
	err = departmentService.AssignAdminToDepartment(adminId, 1)
	if err != nil {
		t.Errorf("Error assigning admin to department: %v", err)
	}
	err = departmentService.FireAdminFromDepartment(adminId, 1)
	if err != nil {
		t.Errorf("Error firing admin from department: %v", err)
	}
	newDepartment, err := departmentService.GetDepartment(1)
	if err != nil {
		t.Errorf("Error getting department: %v", err)
	}
	if newDepartment.HeadID != 0 {
		t.Errorf("Department head ID mismatch: %v != %v", newDepartment.HeadID, 0)
	}
}

func TestFireModeratorFromDepartment(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	authRepository := service_test.CreateTestAuthRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	moderatorId := int64(2)
	err = departmentService.AssignModeratorToDepartment(moderatorId, 1)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	err = departmentService.FireModeratorFromDepartment(moderatorId, 1)
	if err != nil {
		t.Errorf("Error firing moderator from department: %v", err)
	}
}

func TestGetDepartmentUsersIDs(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	authRepository := service_test.CreateTestAuthRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	department = types.Department{
		Name: "Test Department2",
	}
	err = departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	moderatorIDS := []int64{1, 2, 3}
	departmentID1 := int64(1)
	departmentID2 := int64(2)
	departmentID3 := int64(3)
	err = departmentService.AssignModeratorToDepartment(moderatorIDS[0], departmentID1)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	err = departmentService.AssignModeratorToDepartment(moderatorIDS[1], departmentID2)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	err = departmentService.AssignModeratorToDepartment(moderatorIDS[2], departmentID1)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	usersIDs, err := departmentService.GetDepartmentUsersIDs(departmentID1)
	if err != nil {
		t.Errorf("Error getting department users IDs: %v", err)
	}
	if len(usersIDs) != 2 {
		t.Errorf("Department users IDs mismatch: %v != %v", usersIDs, []int64{1, 2})
	}
	usersIDs, err = departmentService.GetDepartmentUsersIDs(departmentID2)
	if err != nil {
		t.Errorf("Error getting department users IDs: %v", err)
	}
	if len(usersIDs) != 1 {
		t.Errorf("Department users IDs mismatch: %v != %v", usersIDs, []int64{3})
	}
	usersIDs, err = departmentService.GetDepartmentUsersIDs(departmentID3)
	if err != nil {
		t.Errorf("Error getting department users IDs: %v", err)
	}
	if len(usersIDs) != 0 {
		t.Errorf("Department users IDs mismatch: %v != %v", usersIDs, []int64{})
	}
}

func TestGetUserDepartmentsIDs(t *testing.T) {
	departmentRepository := service_test.CreateTestDepartmentRepository()
	authRepository := service_test.CreateTestAuthRepository()
	personalDataRepository := service_test.CreateTestPersonalDataRepository()
	userRepository := service_test.CreateTestUserRepository()
	moderatorRepository := service_test.CreateTestModeratorRepository(authRepository, personalDataRepository, userRepository)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	department := types.Department{
		Name: "Test Department",
	}
	err := departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	department = types.Department{
		Name: "Test Department2",
	}
	err = departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	department = types.Department{
		Name: "Test Department3",
	}
	err = departmentService.CreateDepartment(department)
	if err != nil {
		t.Errorf("Error creating department: %v", err)
	}
	userID1 := int64(1)
	userID2 := int64(2)
	userID3 := int64(3)
	departmentID1 := int64(1)
	departmentID2 := int64(2)
	err = departmentService.AssignModeratorToDepartment(userID1, departmentID1)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	err = departmentService.AssignModeratorToDepartment(userID1, departmentID2)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	err = departmentService.AssignModeratorToDepartment(userID2, departmentID2)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	departmentsIDs, err := departmentService.GetUserDepartmentsIDs(userID1)
	if err != nil {
		t.Errorf("Error getting user departments IDs: %v", err)
	}
	if len(departmentsIDs) != 2 {
		t.Errorf("User departments IDs mismatch: %v != %v", departmentsIDs, []int64{1, 2})
	}
	departmentsIDs, err = departmentService.GetUserDepartmentsIDs(userID2)
	if err != nil {
		t.Errorf("Error getting user departments IDs: %v", err)
	}
	if len(departmentsIDs) != 1 {
		t.Errorf("User departments IDs mismatch: %v != %v", departmentsIDs, []int64{2})
	}
	departmentsIDs, err = departmentService.GetUserDepartmentsIDs(userID3)
	if err != nil {
		t.Errorf("Error getting user departments IDs: %v", err)
	}
	if len(departmentsIDs) != 0 {
		t.Errorf("User departments IDs mismatch: %v != %v", departmentsIDs, []int64{})
	}
	departmentID3 := int64(3)
	err = departmentService.AssignModeratorToDepartment(userID3, departmentID3)
	if err != nil {
		t.Errorf("Error assigning moderator to department: %v", err)
	}
	departmentsIDs, err = departmentService.GetUserDepartmentsIDs(userID3)
	if err != nil {
		t.Errorf("Error getting user departments IDs: %v", err)
	}
	if len(departmentsIDs) != 1 {
		t.Errorf("User departments IDs mismatch: %v != %v", departmentsIDs, []int64{3})
	}
}
