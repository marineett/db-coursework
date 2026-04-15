package service_logic

import (
	service_test "data_base_project/tests/service_logic_tests"
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateContract(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	contract, err := contractRepository.GetContract(contract_id)
	if err != nil {
		t.Errorf("Error getting contract: %v", err)
	}
	if contract.ID != contract_id {
		t.Errorf("Contract ID is not %d", contract_id)
	}
	if contract.ClientID != contractInitInfo.ClientID {
		t.Errorf("Client ID is not %d", contractInitInfo.ClientID)
	}
	if contract.RepetitorID != 0 {
		t.Errorf("Repetitor ID is not 0")
	}
	if contract.TransactionID != 0 {
		t.Errorf("Transaction ID is not 0")
	}
	if contract.Status != types.ContractStatusPending {
		t.Errorf("Status is not %d", types.ContractStatusPending)
	}
	if contract.PaymentStatus != types.PaymentStatusNull {
		t.Errorf("Payment status is not %d", types.PaymentStatusNull)
	}
	if contract.ReviewClientID != 0 {
		t.Errorf("Review client ID is not 0")
	}
	if contract.ReviewRepetitorID != 0 {
		t.Errorf("Review repetitor ID is not 0")
	}
	if contract.Price != int64(contractInitInfo.Price) {
		t.Errorf("Price is not %d", contractInitInfo.Price)
	}
	if contract.Commission != contractInitInfo.Commission {
		t.Errorf("Commission is not %d", contractInitInfo.Commission)
	}
	if contract.Description != contractInitInfo.Description {
		t.Errorf("Description is not %s", contractInitInfo.Description)
	}
	if contract.StartDate.After(time.Now()) {
		t.Errorf("Start date is in the future")
	}
}

func TestGetContract(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	contract, err := contractService.GetContract(contract_id)
	if err != nil {
		t.Errorf("Error getting contract: %v", err)
	}
	if contract.ID != contract_id {
		t.Errorf("Contract ID is not %d", contract_id)
	}
	if contract.ClientID != contractInitInfo.ClientID {
		t.Errorf("Client ID is not %d", contractInitInfo.ClientID)
	}
	if contract.RepetitorID != 0 {
		t.Errorf("Repetitor ID is not 0")
	}
	if contract.TransactionID != 0 {
		t.Errorf("Transaction ID is not 0")
	}
	if contract.Status != types.ContractStatusPending {
		t.Errorf("Status is not %d", types.ContractStatusPending)
	}
	if contract.PaymentStatus != types.PaymentStatusNull {
		t.Errorf("Payment status is not %d", types.PaymentStatusNull)
	}
	if contract.ReviewClientID != 0 {
		t.Errorf("Review client ID is not 0")
	}
	if contract.ReviewRepetitorID != 0 {
		t.Errorf("Review repetitor ID is not 0")
	}
	if contract.Price != int64(contractInitInfo.Price) {
		t.Errorf("Price is not %d", contractInitInfo.Price)
	}
	if contract.Commission != contractInitInfo.Commission {
		t.Errorf("Commission is not %d", contractInitInfo.Commission)
	}
	if contract.Description != contractInitInfo.Description {
		t.Errorf("Description is not %s", contractInitInfo.Description)
	}
	if contract.StartDate.After(time.Now()) {
		t.Errorf("Start date is in the future")
	}
	_, err = contractService.GetContract(contract_id + 1)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
func TestUpdateContractStatus(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	err = contractService.UpdateContractStatus(contract_id, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Error updating contract status: %v", err)
	}
	contract, err := contractRepository.GetContract(contract_id)
	if err != nil {
		t.Errorf("Error getting contract: %v", err)
	}
	if contract.Status != types.ContractStatusActive {
		t.Errorf("Status is not %d", types.ContractStatusActive)
	}
}

func TestUpdateContractPaymentStatus(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	err = contractService.UpdateContractPaymentStatus(contract_id, types.PaymentStatusPaid)
	if err != nil {
		t.Errorf("Error updating contract payment status: %v", err)
	}
	contract, err := contractRepository.GetContract(contract_id)
	if err != nil {
		t.Errorf("Error getting contract: %v", err)
	}
	if contract.PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Payment status is not %d", types.PaymentStatusPaid)
	}
}

func TestUpdateContractReviewClientID(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	err = contractService.UpdateContractRepetitorID(contract_id, 2)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "everithing is good",
	}
	err = contractService.CreateContractReviewClient(contract_id, review)
	if err != nil {
		t.Errorf("Error updating contract review client ID: %v", err)
	}
	contract, err := contractRepository.GetContract(contract_id)
	if err != nil {
		t.Errorf("Error getting contract: %v", err)
	}
	if contract.ReviewClientID == 0 {
		t.Errorf("Review client ID is still 0")
	}
}

func TestUpdateContractReviewRepetitorID(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	err = contractService.UpdateContractRepetitorID(contract_id, 2)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	review := types.Review{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      5,
		Comment:     "everithing is good",
	}
	err = contractService.CreateContractReviewRepetitor(contract_id, review)
	if err != nil {
		t.Errorf("Error updating contract review repetitor ID: %v", err)
	}
	contract, err := contractRepository.GetContract(contract_id)
	if err != nil {
		t.Errorf("Error getting contract: %v", err)
	}
	if contract.ReviewRepetitorID == 0 {
		t.Errorf("Review repetitor ID is still 0")
	}
}

func TestUpdateContractRepetitorID(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	err = contractService.UpdateContractRepetitorID(contract_id, 2)
	if err != nil {
		t.Errorf("Error updating contract repetitor ID: %v", err)
	}
	contract, err := contractRepository.GetContract(contract_id)
	if err != nil {
		t.Errorf("Error getting contract: %v", err)
	}
	if contract.RepetitorID != 2 {
		t.Errorf("Repetitor ID is not %d", 2)
	}
}

func TestGetClientContractList(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	contracts, err := contractService.GetClientContractList(1, 0, 10, types.ContractStatusPending)
	if err != nil {
		t.Errorf("Error getting contract list: %v", err)
	}
	if len(contracts) != 1 {
		t.Errorf("Contract list length is not %d", 1)
	}
	if contracts[0].ID != contract_id {
		t.Errorf("Contract ID is not %d", contract_id)
	}
	contracts, err = contractService.GetClientContractList(1, 0, 10, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Error getting contract list: %v", err)
	}
	if len(contracts) != 0 {
		t.Errorf("Contract list length is not %d", 0)
	}
}

func TestGetRepetitorContractList(t *testing.T) {
	contractRepository := service_test.CreateTestContractRepository()
	reviewRepository := service_test.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractInitInfo := types.ContractInitInfo{
		ClientID:         1,
		ContractCategory: types.ContractCategoryProgramming,
		ContractSubcategories: []types.ContractSubcategory{
			types.ContractSubcategoryProgramming,
		},
		Description: "test contract",
		Price:       100,
		Commission:  5,
	}
	contract_id, err := contractService.CreateContract(contractInitInfo)
	if err != nil {
		t.Errorf("Error creating contract: %v", err)
	}
	err = contractService.UpdateContractRepetitorID(contract_id, 2)
	if err != nil {
		t.Errorf("Error updating contract repetitor ID: %v", err)
	}
	contracts, err := contractService.GetRepetitorContractList(2, 0, 10, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Error getting contract list: %v", err)
	}
	if len(contracts) != 1 {
		t.Errorf("Contract list length is not %d", 1)
	}
	if contracts[0].ID != contract_id {
		t.Errorf("Contract ID is not %d", contract_id)
	}
	contracts, err = contractService.GetRepetitorContractList(2, 0, 10, types.ContractStatusPending)
	if err != nil {
		t.Errorf("Error getting contract list: %v", err)
	}
	if len(contracts) != 0 {
		t.Errorf("Contract list length is not %d", 0)
	}
}
