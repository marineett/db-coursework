package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IContractService interface {
	CreateContract(initInfo types.ContractInitInfo) (int64, error)
	GetContract(contractID int64) (*types.Contract, error)
	UpdateContractStatus(contractID int64, status types.ContractStatus) error
	UpdateContractPaymentStatus(contractID int64, paymentStatus types.PaymentStatus) error
	CreateContractReviewClient(contractID int64, review types.Review) error
	CreateContractReviewRepetitor(contractID int64, review types.Review) error
	UpdateContractReviewClient(contractID int64, review types.Review) error
	UpdateContractReviewRepetitor(contractID int64, review types.Review) error
	UpdateContractRepetitorID(contractID int64, repetitorID int64) error
	GetClientContractList(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error)
	GetRepetitorContractList(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error)
	GetContractList(from int64, size int64, status types.ContractStatus) ([]types.Contract, error)
	GetAllContracts(from int64, size int64) ([]types.Contract, error)
}

type ContractService struct {
	contractRepository data_base.IContractRepository
	reviewRepository   data_base.IReviewRepository
}

func CreateContractService(contractRepository data_base.IContractRepository, reviewRepository data_base.IReviewRepository) IContractService {
	return &ContractService{
		contractRepository: contractRepository,
		reviewRepository:   reviewRepository,
	}
}

func (s *ContractService) transormInitContractData(initInfo types.ContractInitInfo) (*types.Contract, error) {
	return &types.Contract{
		ClientID:          initInfo.ClientID,
		Description:       initInfo.Description,
		Price:             int64(initInfo.Price),
		Commission:        initInfo.Commission,
		Status:            types.ContractStatusPending,
		PaymentStatus:     types.PaymentStatusNull,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		StartDate:         initInfo.StartDate,
		EndDate:           initInfo.StartDate.AddDate(0, 0, int(initInfo.Duration)),
	}, nil
}

func (s *ContractService) CreateContract(initInfo types.ContractInitInfo) (int64, error) {
	contract, err := s.transormInitContractData(initInfo)
	if err != nil {
		return 0, err
	}
	contract.CreatedAt = time.Now()
	return s.contractRepository.InsertContract(*contract)
}

func (s *ContractService) GetContract(contractID int64) (*types.Contract, error) {
	return s.contractRepository.GetContract(contractID)
}

func (s *ContractService) UpdateContractStatus(contractID int64, status types.ContractStatus) error {
	return s.contractRepository.UpdateContractStatus(contractID, status)
}

func (s *ContractService) UpdateContractPaymentStatus(contractID int64, paymentStatus types.PaymentStatus) error {
	return s.contractRepository.UpdateContractPaymentStatus(contractID, paymentStatus)
}

func (s *ContractService) CreateContractReviewClient(contractID int64, review types.Review) error {
	tx, err := s.contractRepository.BeginTx()
	if err != nil {
		return err
	}
	if tx != nil {
		defer tx.Rollback()
	}
	reviewID, err := s.reviewRepository.InsertReviewInSeq(tx, review)
	if err != nil {
		return err
	}
	err = s.contractRepository.UpdateContractReviewClientIDInSeq(tx, contractID, reviewID)
	if err != nil {
		return err
	}
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (s *ContractService) CreateContractReviewRepetitor(contractID int64, review types.Review) error {
	tx, err := s.contractRepository.BeginTx()
	if err != nil {
		return err
	}
	if tx != nil {
		defer tx.Rollback()
	}
	reviewID, err := s.reviewRepository.InsertReviewInSeq(tx, review)
	if err != nil {
		return err
	}
	err = s.contractRepository.UpdateContractReviewRepetitorID(contractID, reviewID)
	if err != nil {
		return err
	}
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (s *ContractService) UpdateContractReviewClient(contractID int64, review types.Review) error {
	contract, err := s.contractRepository.GetContract(contractID)
	if err != nil {
		return err
	}
	review.ID = contract.ReviewClientID
	return s.reviewRepository.UpdateReview(review)
}

func (s *ContractService) UpdateContractReviewRepetitor(contractID int64, review types.Review) error {
	contract, err := s.contractRepository.GetContract(contractID)
	if err != nil {
		return err
	}
	review.ID = contract.ReviewRepetitorID
	return s.reviewRepository.UpdateReview(review)
}

func (s *ContractService) GetRepetitorContractList(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	return s.contractRepository.GetContractsByRepetitorID(repetitorID, from, size, status)
}

func (s *ContractService) GetClientContractList(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	return s.contractRepository.GetContractsByClientID(clientID, from, size, status)
}

func (s *ContractService) UpdateContractRepetitorID(contractID int64, repetitorID int64) error {
	return s.contractRepository.UpdateContractRepetitorID(contractID, repetitorID)
}

func (s *ContractService) GetContractList(from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	return s.contractRepository.GetContractList(from, size, status)
}

func (s *ContractService) GetAllContracts(from int64, size int64) ([]types.Contract, error) {
	return s.contractRepository.GetAllContracts(from, size)
}
