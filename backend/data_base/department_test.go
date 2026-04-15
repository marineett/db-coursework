package data_base

import (
	"data_base_project/types"
	"database/sql"
	"sort"
	"testing"
)

func TestCreateDepartmentRepository(t *testing.T) {
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	if departmentRepository == nil {
		t.Errorf("Failed to create department repository")
	}
}

func TestInsertDepartment(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_department_table CASCADE")
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	department := types.Department{
		Name:   "Test Department",
		HeadID: 0,
	}
	id, err := departmentRepository.InsertDepartment(department)
	if err != nil {
		t.Errorf("Failed to insert department: %v", err)
	}
	resultDepartment := types.Department{}
	globalDb.QueryRow("SELECT * FROM test_department_table WHERE id = $1", id).Scan(&resultDepartment.ID, &resultDepartment.Name, &resultDepartment.HeadID)
	if resultDepartment.ID != id {
		t.Errorf("Failed to insert department: %v", err)
	}
	if resultDepartment.Name != department.Name {
		t.Errorf("Failed to insert department: %v", err)
	}
	if resultDepartment.HeadID != department.HeadID {
		t.Errorf("Failed to insert department: %v", err)
	}
}

func TestGetDepartment(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_department_table CASCADE")
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	department := types.Department{
		Name:   "Test Department",
		HeadID: 0,
	}
	id, err := departmentRepository.InsertDepartment(department)
	if err != nil {
		t.Errorf("Failed to insert department: %v", err)
	}
	resultDepartment, err := departmentRepository.GetDepartment(id)
	if err != nil {
		t.Errorf("Failed to get department: %v", err)
	}
	if resultDepartment.ID != id {
		t.Errorf("Failed to get department: %v", err)
	}
	if resultDepartment.Name != department.Name {
		t.Errorf("Failed to get department: %v", err)
	}
	if resultDepartment.HeadID != department.HeadID {
		t.Errorf("Failed to get department: %v", err)
	}
}

func TestChangeDepartmentHead(t *testing.T) {
	defer globalDb.Exec("TRUNCATE TABLE test_hire_info_table CASCADE")
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	department := types.Department{
		Name:   "Test Department",
		HeadID: 0,
	}
	id, err := departmentRepository.InsertDepartment(department)
	if err != nil {
		t.Errorf("Failed to insert department: %v", err)
	}
	err = departmentRepository.ChangeDepartmentHead(id, 1)
	if err != nil {
		t.Errorf("Failed to change department head: %v", err)
	}
	resultDepartment, err := departmentRepository.GetDepartment(id)
	if err != nil {
		t.Errorf("Failed to get department: %v", err)
	}
	if resultDepartment.HeadID != 1 {
		t.Errorf("Failed to change department head: %v", err)
	}
	err = departmentRepository.ChangeDepartmentHead(id+5, 1)
	if err == nil {
		t.Errorf("Change department head of non-existent department should fail")
	}
}

func TestHireInfoInsert(t *testing.T) {
	InsertTestUser(1)
	InsertTestDepartment(1)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_personal_data_table, test_department_table, test_hire_info_table CASCADE")
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	hireInfo := types.HireInfo{
		DepartmentID: 1,
		UserID:       1,
	}
	err := departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	resultHireInfo := types.HireInfo{}
	err = globalDb.QueryRow("SELECT department_id, user_id FROM test_hire_info_table WHERE department_id = $1 AND user_id = $2", hireInfo.DepartmentID, hireInfo.UserID).Scan(&resultHireInfo.DepartmentID, &resultHireInfo.UserID)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	if resultHireInfo.DepartmentID != hireInfo.DepartmentID {
		t.Errorf("Failed to insert hire info: department id %v", resultHireInfo.DepartmentID)
	}
	if resultHireInfo.UserID != hireInfo.UserID {
		t.Errorf("Failed to insert hire info: user id %v", resultHireInfo.UserID)
	}
}

func TestHireInfoDelete(t *testing.T) {
	InsertTestUser(1)
	InsertTestDepartment(1)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_personal_data_table, test_department_table, test_hire_info_table CASCADE")
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	hireInfo := types.HireInfo{
		DepartmentID: 1,
		UserID:       1,
	}
	err := departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	err = departmentRepository.HireInfoDelete(hireInfo.UserID, hireInfo.DepartmentID)
	if err != nil {
		t.Errorf("Failed to delete hire info: %v", err)
	}
	resultHireInfo := types.HireInfo{}
	err = globalDb.QueryRow("SELECT * FROM test_hire_info_table WHERE department_id = $1 AND user_id = $2", hireInfo.DepartmentID, hireInfo.UserID).Scan(&resultHireInfo.DepartmentID, &resultHireInfo.UserID)
	if err != sql.ErrNoRows {
		t.Errorf("Failed to delete hire info: %v", err)
	}
	err = departmentRepository.HireInfoDelete(hireInfo.UserID+5, hireInfo.DepartmentID)
	if err == nil {
		t.Errorf("Failed to delete hire info: %v", err)
	}
}

func TestGetUserDepartmentsIDs(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestDepartment(1)
	InsertTestDepartment(2)
	InsertTestDepartment(3)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_personal_data_table, test_department_table, test_hire_info_table CASCADE")
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	hireInfo := types.HireInfo{
		DepartmentID: 1,
		UserID:       1,
	}
	err := departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	hireInfo = types.HireInfo{
		DepartmentID: 2,
		UserID:       1,
	}
	err = departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	hireInfo = types.HireInfo{
		DepartmentID: 3,
		UserID:       2,
	}
	err = departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	departments, err := departmentRepository.GetUserDepartmentsIDs(1)
	if err != nil {
		t.Errorf("Failed to get user departments: %v", err)
	}
	if len(departments) != 2 {
		t.Errorf("Failed to get user departments: %v", err)
	}
	sort.Slice(departments, func(i, j int) bool {
		return departments[i] < departments[j]
	})
	if departments[0] != 1 || departments[1] != 2 {
		t.Errorf("Failed to get user departments: %v", err)
	}
	departments, err = departmentRepository.GetUserDepartmentsIDs(3)
	if err != nil {
		t.Errorf("Failed to get user departments: %v", err)
	}
	if len(departments) != 0 {
		t.Errorf("Failed to get user departments: %v", err)
	}
}

func TestGetDepartmentUsersIDs(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestUser(3)
	InsertTestDepartment(1)
	InsertTestDepartment(2)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_personal_data_table, test_department_table, test_hire_info_table CASCADE")
	departmentRepository := CreateDepartmentRepository(globalDb, "test_department_table", "test_hire_info_table")
	hireInfo := types.HireInfo{
		DepartmentID: 1,
		UserID:       1,
	}
	err := departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	hireInfo = types.HireInfo{
		DepartmentID: 2,
		UserID:       1,
	}
	err = departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	hireInfo = types.HireInfo{
		DepartmentID: 2,
		UserID:       2,
	}
	err = departmentRepository.HireInfoInsert(hireInfo)
	if err != nil {
		t.Errorf("Failed to insert hire info: %v", err)
	}
	users, err := departmentRepository.GetDepartmentUsersIDs(2)
	if err != nil {
		t.Errorf("Failed to get department users: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Failed to get department users: %v", err)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i] < users[j]
	})
	if users[0] != 1 || users[1] != 2 {
		t.Errorf("Failed to get department users: %v", err)
	}
	users, err = departmentRepository.GetDepartmentUsersIDs(3)
	if err != nil {
		t.Errorf("Failed to get department users: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Failed to get department users: %v", err)
	}
}
