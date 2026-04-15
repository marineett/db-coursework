package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IPersonalDataService interface {
	GetPersonalData(id int64) (*types.PersonalData, error)
}

type PersonalDataService struct {
	personalDataRepository data_base.IPersonalDataRepository
}

func CreatePersonalDataService(personalDataRepository data_base.IPersonalDataRepository) *PersonalDataService {
	return &PersonalDataService{
		personalDataRepository: personalDataRepository,
	}
}

func (s *PersonalDataService) GetPersonalData(id int64) (*types.PersonalData, error) {
	return s.personalDataRepository.GetPersonalData(id)
}
