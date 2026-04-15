package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IAdminService interface {
	CreateAdmin(initData types.InitAdminData) error
	GetAdminData(userID int64) (*types.AdminData, error)
	GetAdminProfile(userID int64) (*types.AdminProfile, error)
	UpdateAdminPersonalData(userID int64, personalData types.PersonalData) error
	UpdateAdminPassword(userID int64, authData types.AuthData, newPassword string) error
	UpdateAdminDepartment(userID int64, departmentID int64) error
	UpdateAdminSalary(userID int64, salary int64) error
}

type AdminService struct {
	adminRepository        data_base.IAdminRepository
	personalDataRepository data_base.IPersonalDataRepository
}

func CreateAdminService(
	adminRepository data_base.IAdminRepository,
	personalDataRepository data_base.IPersonalDataRepository,
) IAdminService {
	return &AdminService{
		adminRepository:        adminRepository,
		personalDataRepository: personalDataRepository,
	}
}

func (s *AdminService) transormInitAdminData(initData types.InitAdminData) (*types.AdminData, *types.PersonalData, *types.AuthData) {
	adminData := types.AdminData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
			PersonalDataID:   initData.PersonalData.ID,
		},
		Salary:       400_000,
		DepartmentID: 0,
	}
	return &adminData, &initData.PersonalData, &initData.AuthData
}
func (s *AdminService) CreateAdmin(initData types.InitAdminData) error {
	adminData, personalData, authData := s.transormInitAdminData(initData)
	_, err := s.adminRepository.InsertAdmin(*adminData, *personalData, *authData)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) GetAdminData(userID int64) (*types.AdminData, error) {
	return s.adminRepository.GetAdmin(userID)
}

func (s *AdminService) UpdateAdminPersonalData(userID int64, personalData types.PersonalData) error {
	return s.adminRepository.UpdateAdminPersonalData(userID, personalData)
}

func (s *AdminService) UpdateAdminPassword(userID int64, authData types.AuthData, newPassword string) error {
	return s.adminRepository.UpdateAdminPassword(userID, authData, newPassword)
}

func (s *AdminService) UpdateAdminDepartment(userID int64, departmentID int64) error {
	return s.adminRepository.UpdateAdminDepartment(userID, departmentID)
}

func (s *AdminService) UpdateAdminSalary(userID int64, salary int64) error {
	return s.adminRepository.UpdateAdminSalary(userID, salary)
}

func (s *AdminService) GetAdminProfile(userID int64) (*types.AdminProfile, error) {
	adminData, err := s.GetAdminData(userID)
	if err != nil {
		return nil, err
	}
	personalData, err := s.personalDataRepository.GetPersonalData(userID)
	if err != nil {
		return nil, err
	}
	return &types.AdminProfile{
		FirstName:       personalData.FirstName,
		LastName:        personalData.LastName,
		MiddleName:      personalData.MiddleName,
		TelephoneNumber: personalData.TelephoneNumber,
		Email:           personalData.Email,
		Salary:          adminData.Salary,
	}, nil
}
