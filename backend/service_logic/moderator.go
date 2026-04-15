package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IModeratorService interface {
	CreateModerator(init_data types.InitModeratorData) error
	GetModeratorData(user_id int64) (*types.ModeratorData, error)
	GetModeratorProfile(user_id int64) (*types.ModeratorProfile, error)
	UpdateModeratorPersonalData(user_id int64, personal_data types.PersonalData) error
	UpdateModeratorPassword(user_id int64, authData types.AuthData, newPassword string) error
	UpdateModeratorSalary(user_id int64, salary int64) error
	GetModerators() ([]*types.MoreratorProfileWithID, error)
	GetModeratorProfileWithId(moderator_id int64) (*types.MoreratorProfileWithID, error)
}

type ModeratorService struct {
	moderatorRepository    data_base.IModeratorRepository
	personalDataRepository data_base.IPersonalDataRepository
	departmentRepository   data_base.IDepartmentRepository
}

func CreateModeratorService(
	moderatorRepository data_base.IModeratorRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	departmentRepository data_base.IDepartmentRepository,
) IModeratorService {
	return &ModeratorService{
		moderatorRepository:    moderatorRepository,
		personalDataRepository: personalDataRepository,
		departmentRepository:   departmentRepository,
	}
}

func (s *ModeratorService) transormInitModeratorData(init_data types.InitModeratorData) (*types.ModeratorData, *types.PersonalData, *types.AuthData) {
	moderator := types.ModeratorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
		Salary: 0,
	}
	return &moderator, &init_data.InitUserData.PersonalData, &types.AuthData{
		Login:    init_data.InitUserData.AuthData.Login,
		Password: init_data.InitUserData.AuthData.Password,
	}
}

func (s *ModeratorService) CreateModerator(init_data types.InitModeratorData) error {
	moderator_data, personal_data, auth_info := s.transormInitModeratorData(init_data)
	_, err := s.moderatorRepository.InsertModerator(*moderator_data, *personal_data, *auth_info)
	if err != nil {
		return err
	}
	return nil
}

func (s *ModeratorService) GetModeratorData(user_id int64) (*types.ModeratorData, error) {
	return s.moderatorRepository.GetModerator(user_id)
}

func (s *ModeratorService) UpdateModeratorPersonalData(user_id int64, personal_data types.PersonalData) error {
	return s.moderatorRepository.UpdateModeratorPersonalData(user_id, personal_data)
}

func (s *ModeratorService) UpdateModeratorPassword(user_id int64, authData types.AuthData, newPassword string) error {
	return s.moderatorRepository.UpdateModeratorPassword(user_id, authData, newPassword)
}

func (s *ModeratorService) GetModeratorProfile(user_id int64) (*types.ModeratorProfile, error) {
	moderator_data, err := s.GetModeratorData(user_id)
	if err != nil {
		return nil, err
	}
	personal_data, err := s.personalDataRepository.GetPersonalData(user_id)
	if err != nil {
		return nil, err
	}
	departmentsIds, err := s.departmentRepository.GetUserDepartmentsIDs(user_id)
	if err != nil {
		return nil, err
	}
	departments := make([]string, len(departmentsIds))
	for _, departmentId := range departmentsIds {
		department, err := s.departmentRepository.GetDepartment(departmentId)
		departments = append(departments, department.Name)
		if err != nil {
			return nil, err
		}
	}
	return &types.ModeratorProfile{
		FirstName:       personal_data.FirstName,
		LastName:        personal_data.LastName,
		MiddleName:      personal_data.MiddleName,
		TelephoneNumber: personal_data.TelephoneNumber,
		Email:           personal_data.Email,
		Salary:          moderator_data.Salary,
		Departments:     departments,
	}, nil
}

func (s *ModeratorService) GetModerators() ([]*types.MoreratorProfileWithID, error) {
	moderatorsIDs, err := s.moderatorRepository.GetModerators()
	if err != nil {
		return nil, err
	}
	moderatorsWithID := make([]*types.MoreratorProfileWithID, len(moderatorsIDs))
	for i, moderatorID := range moderatorsIDs {
		moderatorProfile, err := s.GetModeratorProfile(moderatorID)
		if err != nil {
			return nil, err
		}
		moderatorsWithID[i] = &types.MoreratorProfileWithID{
			ModeratorProfile: *moderatorProfile,
			ID:               moderatorID,
		}
	}
	return moderatorsWithID, nil
}

func (s *ModeratorService) GetModeratorProfileWithId(moderator_id int64) (*types.MoreratorProfileWithID, error) {
	moderatorProfile, err := s.GetModeratorProfile(moderator_id)
	if err != nil {
		return nil, err
	}
	moderatorProfile.Departments = nil
	return &types.MoreratorProfileWithID{
		ModeratorProfile: *moderatorProfile,
		ID:               moderator_id,
	}, nil
}

func (s *ModeratorService) UpdateModeratorSalary(user_id int64, salary int64) error {
	return s.moderatorRepository.UpdateModeratorSalary(user_id, salary)
}
