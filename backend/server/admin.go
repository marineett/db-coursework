package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"net/http"
	"strconv"
)

func SetupAdminRouter(
	adminService service_logic.IAdminService,
	departmentService service_logic.IDepartmentService,
	moderatorService service_logic.IModeratorService,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(ADMIN_GET_PROFILE, AdminGetProfileHandler(adminService))
	router.HandleFunc(ADMIN_CREATE_DEPARTMENT, AdminCreateDepartmentHandler(adminService, departmentService))
	router.HandleFunc(ADMIN_GET_DEPARTMENTS, AdminGetDepartmentsHandler(departmentService, moderatorService))
	router.HandleFunc(ADMIN_GET_MODERATORS, AdminGetModeratorsHandler(moderatorService))
	router.HandleFunc(ADMIN_HIRE_MODERATOR, AdminHireModeratorHandler(departmentService))
	router.HandleFunc(ADMIN_FIRE_MODERATOR, AdminFireModeratorHandler(departmentService))
	router.HandleFunc(ADMIN_CHANGE_MODERATOR_SALARY, AdminChangeModeratorSalaryHandler(moderatorService))
	return router
}

func AdminGetProfileHandler(adminService service_logic.IAdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		adminProfile, err := adminService.GetAdminProfile(userID)
		if err != nil {
			http.Error(w, "Failed to get admin profile", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(adminProfile)
	}
}

func AdminCreateDepartmentHandler(
	adminService service_logic.IAdminService,
	departmentService service_logic.IDepartmentService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		departmentName := r.URL.Query().Get("name")
		if departmentName == "" {
			http.Error(w, "Department name is required", http.StatusBadRequest)
			return
		}
		department := types.Department{
			Name:   departmentName,
			HeadID: userID,
		}
		err = departmentService.CreateDepartment(department)
		if err != nil {
			http.Error(w, "Failed to create department", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func AdminGetDepartmentsHandler(departmentService service_logic.IDepartmentService, moderatorService service_logic.IModeratorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		departments, err := departmentService.GetDepartmentsByHeadID(userID)
		if err != nil {
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		completeDepartments := make([]types.CompleteDepartmentInfo, len(departments))
		for i, department := range departments {
			moderatorsIDs, err := departmentService.GetDepartmentUsersIDs(department.ID)
			if err != nil {
				http.Error(w, "Failed to get complete department info", http.StatusInternalServerError)
				return
			}
			moderators := make([]types.MoreratorProfileWithID, len(moderatorsIDs))
			for j, moderatorID := range moderatorsIDs {
				moderator, err := moderatorService.GetModeratorProfileWithId(moderatorID)
				if err != nil {
					http.Error(w, "Failed to get complete department info", http.StatusInternalServerError)
					return
				}
				moderators[j] = *moderator
			}
			completeDepartments[i] = types.CompleteDepartmentInfo{
				Department: department,
				Moderators: moderators,
			}
		}
		json.NewEncoder(w).Encode(completeDepartments)
	}
}

func AdminGetDepartmentHandler(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		departmentIDStr := r.URL.Query().Get("id")
		departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		department, err := departmentService.GetDepartment(departmentID)
		if err != nil {
			http.Error(w, "Failed to get department", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(department)
	}
}

func AdminGetModeratorsHandler(moderatorService service_logic.IModeratorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		moderators, err := moderatorService.GetModerators()
		if err != nil {
			http.Error(w, "Failed to get moderators", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(moderators)
	}
}

func AdminHireModeratorHandler(
	departmentService service_logic.IDepartmentService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		adminIDStr := r.URL.Query().Get("id")
		adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid admin ID", http.StatusBadRequest)
			return
		}
		departmentIDStr := r.URL.Query().Get("d_id")
		departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		department, err := departmentService.GetDepartment(departmentID)
		if err != nil {
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		departments, err := departmentService.GetUserDepartmentsIDs(moderatorID)
		if err != nil {
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		for _, departmentID := range departments {
			if departmentID == department.ID {
				http.Error(w, "Moderator is already in this department", http.StatusBadRequest)
				return
			}
		}
		if department.HeadID != adminID {
			http.Error(w, "You are not the head of this department", http.StatusBadRequest)
			return
		}
		err = departmentService.AssignModeratorToDepartment(moderatorID, departmentID)
		if err != nil {
			http.Error(w, "Failed to hire moderator", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func AdminFireModeratorHandler(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		adminIDStr := r.URL.Query().Get("id")
		adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid admin ID", http.StatusBadRequest)
			return
		}
		departmentIDStr := r.URL.Query().Get("d_id")
		departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		departments, err := departmentService.GetDepartment(departmentID)
		if err != nil {
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		if departments.HeadID != adminID {
			http.Error(w, "You are not the head of this department", http.StatusBadRequest)
			return
		}
		err = departmentService.FireModeratorFromDepartment(moderatorID, departmentID)
		if err != nil {
			http.Error(w, "Failed to fire moderator", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func AdminChangeModeratorSalaryHandler(
	moderatorService service_logic.IModeratorService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		newSalaryStr := r.URL.Query().Get("salary")
		newSalary, err := strconv.ParseInt(newSalaryStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		err = moderatorService.UpdateModeratorSalary(moderatorID, newSalary)
		if err != nil {
			http.Error(w, "Failed to change moderator salary", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
