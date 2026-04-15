package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IRepetitorService interface {
	CreateRepetitor(initData types.InitRepetitorData) error
	GetRepetitorData(userID int64) (*types.RepetitorData, error)
	GetRepetitorProfile(userID int64, reviewsOffset int64, reviewsLimit int64) (*types.RepetitorProfile, error)
	GetRepetitors(repetitorsOffset int64, repetitorsLimit int64) ([]*types.RepetitorView, error)
	UpdateRepetitorPersonalData(userID int64, personalData types.PersonalData) error
	UpdateRepetitorPassword(userID int64, authData types.AuthData, newPassword string) error
}

type RepetitorService struct {
	repetitorRepository    data_base.IRepetitorRepository
	personalDataRepository data_base.IPersonalDataRepository
	reviewRepository       data_base.IReviewRepository
	resumeRepository       data_base.IResumeRepository
}

func CreateRepetitorService(
	repetitorRepository data_base.IRepetitorRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	reviewRepository data_base.IReviewRepository,
	resumeRepository data_base.IResumeRepository,
) IRepetitorService {
	return &RepetitorService{
		repetitorRepository:    repetitorRepository,
		personalDataRepository: personalDataRepository,
		reviewRepository:       reviewRepository,
		resumeRepository:       resumeRepository,
	}
}
func (s *RepetitorService) transormInitRepetitorData(initData types.InitRepetitorData) (*types.RepetitorData, *types.PersonalData, *types.AuthData) {
	repetitor := types.RepetitorData{
		UserData: types.UserData{
			RegistrationDate: time.Now(),
			LastLoginDate:    time.Now(),
		},
		MeanRating: 0,
		ResumeID:   0,
		Reviews:    nil,
	}
	return &repetitor, &initData.InitUserData.PersonalData, &types.AuthData{
		Login:    initData.InitUserData.AuthData.Login,
		Password: initData.InitUserData.AuthData.Password,
	}
}

func (s *RepetitorService) CreateRepetitor(initData types.InitRepetitorData) error {
	repetitor, personalData, authData := s.transormInitRepetitorData(initData)
	_, err := s.repetitorRepository.InsertRepetitor(*repetitor, *personalData, *authData)
	if err != nil {
		return err
	}
	return nil
}

func (s *RepetitorService) GetRepetitorData(userID int64) (*types.RepetitorData, error) {
	return s.repetitorRepository.GetRepetitor(userID)
}

func (s *RepetitorService) UpdateRepetitorPersonalData(userID int64, personalData types.PersonalData) error {
	return s.repetitorRepository.UpdateRepetitorPersonalData(userID, personalData)
}

func (s *RepetitorService) UpdateRepetitorPassword(userID int64, authData types.AuthData, newPassword string) error {
	return s.repetitorRepository.UpdateRepetitorPassword(userID, authData, newPassword)
}

func (s *RepetitorService) GetRepetitorProfile(userID int64, reviewsOffset int64, reviewsLimit int64) (*types.RepetitorProfile, error) {
	repetitor, err := s.GetRepetitorData(userID)
	if err != nil {
		return nil, err
	}
	reviews, err := s.reviewRepository.GetReviewsByRepetitorID(userID, reviewsOffset, reviewsLimit)
	if err != nil {
		return nil, err
	}
	personalData, err := s.personalDataRepository.GetPersonalData(userID)
	if err != nil {
		return nil, err
	}
	resume, err := s.resumeRepository.GetResume(userID)
	if err != nil {
		resume = &types.Resume{
			ID:          0,
			RepetitorID: userID,
			Title:       "Very good Title",
			Description: "Very good Description",
			Prices:      map[string]int{"go": 100, "sql": 200, "dBeaver": 300},
		}
	}
	return &types.RepetitorProfile{
		FirstName:         personalData.FirstName,
		LastName:          personalData.LastName,
		MiddleName:        personalData.MiddleName,
		TelephoneNumber:   personalData.TelephoneNumber,
		Email:             personalData.Email,
		MeanRating:        repetitor.MeanRating,
		ResumeTitle:       resume.Title,
		ResumeDescription: resume.Description,
		ResumePrices:      resume.Prices,
		Reviews:           reviews,
	}, nil
}

func (s *RepetitorService) GetRepetitors(repetitorsOffset int64, repetitorsLimit int64) ([]*types.RepetitorView, error) {
	repetitorsIds, err := s.repetitorRepository.GetRepetitorsIds(repetitorsOffset, repetitorsLimit)
	if err != nil {
		return nil, err
	}
	repetitors := make([]*types.RepetitorView, len(repetitorsIds))
	for i, id := range repetitorsIds {
		repetitor, err := s.GetRepetitorData(id)
		if err != nil {
			return nil, err
		}
		personalData, err := s.personalDataRepository.GetPersonalData(id)
		if err != nil {
			return nil, err
		}
		repetitors[i] = &types.RepetitorView{
			FirstName:  personalData.FirstName,
			MeanRating: repetitor.MeanRating,
		}
	}
	return repetitors, nil
}
