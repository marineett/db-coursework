package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IReviewService interface {
	GetReview(id int64) (*types.Review, error)
	GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.Review, error)
	GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.Review, error)
}

type ReviewService struct {
	reviewRepository data_base.IReviewRepository
}

func CreateReviewService(reviewRepository data_base.IReviewRepository) IReviewService {
	return &ReviewService{reviewRepository: reviewRepository}
}

func (s *ReviewService) GetReview(id int64) (*types.Review, error) {
	return s.reviewRepository.GetReview(id)
}

func (s *ReviewService) GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.Review, error) {
	return s.reviewRepository.GetReviewsByRepetitorID(repetitorID, from, size)
}

func (s *ReviewService) GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.Review, error) {
	return s.reviewRepository.GetReviewsByClientID(clientID, from, size)
}
