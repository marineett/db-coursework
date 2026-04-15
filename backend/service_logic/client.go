package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IClientService interface {
	CreateClient(initData types.InitClientData) error
	GetClientData(UserID int64) (*types.ClientData, error)
	GetClientProfile(UserID int64, reviewsOffset int64, reviewsLimit int64) (*types.ClientProfile, error)
	UpdateClientPersonalData(UserID int64, personalData types.PersonalData) error
	UpdateClientPassword(UserID int64, authData types.AuthData, newPassword string) error
}

type ClientService struct {
	clientRepository       data_base.IClientRepository
	personalDataRepository data_base.IPersonalDataRepository
	reviewRepository       data_base.IReviewRepository
}

func CreateClientService(
	clientRepository data_base.IClientRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	reviewRepository data_base.IReviewRepository,
) IClientService {
	return &ClientService{
		clientRepository:       clientRepository,
		personalDataRepository: personalDataRepository,
		reviewRepository:       reviewRepository,
	}
}

func (s *ClientService) transormInitClientData(initData types.InitClientData) (*types.ClientData, *types.PersonalData, *types.AuthData) {
	client := types.ClientData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
	}
	return &client, &initData.InitUserData.PersonalData, &types.AuthData{
		Login:    initData.InitUserData.AuthData.Login,
		Password: initData.InitUserData.AuthData.Password,
	}
}

func (s *ClientService) CreateClient(initData types.InitClientData) error {
	clientData, personalData, authData := s.transormInitClientData(initData)
	_, err := s.clientRepository.InsertClient(*clientData, *personalData, *authData)
	if err != nil {
		return err
	}
	return nil
}

func (s *ClientService) GetClientData(userID int64) (*types.ClientData, error) {
	return s.clientRepository.GetClient(userID)
}

func (s *ClientService) UpdateClientPersonalData(userID int64, personalData types.PersonalData) error {
	return s.clientRepository.UpdateClientPersonalData(userID, personalData)
}

func (s *ClientService) UpdateClientPassword(userID int64, authData types.AuthData, newPassword string) error {
	return s.clientRepository.UpdateClientPassword(userID, authData, newPassword)
}

func (s *ClientService) GetClientProfile(userID int64, reviewsOffset int64, reviewsLimit int64) (*types.ClientProfile, error) {
	client, err := s.GetClientData(userID)
	if err != nil {
		return nil, err
	}
	personalData, err := s.personalDataRepository.GetPersonalData(userID)
	if err != nil {
		return nil, err
	}
	reviews, err := s.reviewRepository.GetReviewsByClientID(userID, reviewsOffset, reviewsLimit)
	if err != nil {
		return nil, err
	}
	return &types.ClientProfile{
		FirstName:       personalData.FirstName,
		LastName:        personalData.LastName,
		MiddleName:      personalData.MiddleName,
		TelephoneNumber: personalData.TelephoneNumber,
		Email:           personalData.Email,
		MeanRating:      client.MeanRating,
		Reviews:         reviews,
	}, nil
}
