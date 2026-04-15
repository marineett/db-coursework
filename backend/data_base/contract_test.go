package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateContractRepository(t *testing.T) {
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	if contractRepository == nil {
		t.Errorf("Failed to create contract repository")
	}
}

func TestInsertContract(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")

	resultContract := types.Contract{}
	err = globalDb.QueryRow("SELECT * FROM test_contract_table WHERE id = $1", insertedID).Scan(&resultContract.ID, &resultContract.ClientID, &resultContract.RepetitorID, &resultContract.ReviewClientID, &resultContract.ReviewRepetitorID, &resultContract.TransactionID, &resultContract.CreatedAt, &resultContract.Description, &resultContract.Status, &resultContract.PaymentStatus, &resultContract.Price, &resultContract.Commission, &resultContract.StartDate, &resultContract.EndDate)
	if err != nil {
		t.Errorf("Failed to get contract: %v", err)
	}
	if resultContract.ClientID != 1 {
		t.Errorf("Contract ID is not 1")
	}
	if resultContract.RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if resultContract.ReviewClientID != 0 {
		t.Errorf("Contract Review Client ID is not 0")
	}
	if resultContract.ReviewRepetitorID != 0 {
		t.Errorf("Contract Review Repetitor ID is not 0")
	}
	if resultContract.TransactionID != 0 {
		t.Errorf("Contract Transaction ID is not 0")
	}
	if resultContract.Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if resultContract.Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
}

func TestGetContract(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	resultContract, err := contractRepository.GetContract(insertedID)
	if err != nil {
		t.Errorf("Failed to get contract: %v", err)
	}
	if resultContract.ClientID != 1 {
		t.Errorf("Contract ID is not 1")
	}
	if resultContract.RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if resultContract.ReviewClientID != 0 {
		t.Errorf("Contract Review Client ID is not 0")
	}
	if resultContract.ReviewRepetitorID != 0 {
		t.Errorf("Contract Review Repetitor ID is not 0")
	}
	if resultContract.TransactionID != 0 {
		t.Errorf("Contract Transaction ID is not 0")
	}
	if resultContract.Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if resultContract.Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
}

func TestGetContractsByRepetitorID(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestUser(3)
	InsertTestUser(4)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	contract = types.Contract{
		ClientID:          1,
		RepetitorID:       3,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	_, err = contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	contract = types.Contract{
		ClientID:          4,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	insertedID3, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	contracts, err := contractRepository.GetContractsByRepetitorID(2, 0, 10, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Failed to get contracts: %v", err)
	}
	if len(contracts) != 2 {
		t.Errorf("Number of contracts is not 1")
	}
	if contracts[0].ID != insertedID3 {
		t.Errorf("Contract ID is not %d", insertedID3)
	}
	if contracts[0].ClientID != 4 {
		t.Errorf("Contract Client ID is not 4")
	}
	if contracts[0].RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if contracts[0].ReviewClientID != 0 {
		t.Errorf("Contract Review Client ID is not 0")
	}
	if contracts[0].ReviewRepetitorID != 0 {
		t.Errorf("Contract Review Repetitor ID is not 0")
	}
	if contracts[0].Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if contracts[0].Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
	if contracts[0].Status != types.ContractStatusActive {
		t.Errorf("Contract Status is not %v", types.ContractStatusActive)
	}
	if contracts[0].PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Contract Payment Status is not %v", types.PaymentStatusPaid)
	}

	if contracts[1].ID != insertedID {
		t.Errorf("Contract ID is not %d", insertedID)
	}
	if contracts[1].ClientID != 1 {
		t.Errorf("Contract Client ID is not 1")
	}
	if contracts[1].RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if contracts[1].Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if contracts[1].Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
	if contracts[1].Status != types.ContractStatusActive {
		t.Errorf("Contract Status is not %v", types.ContractStatusActive)
	}
	if contracts[1].PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Contract Payment Status is not %v", types.PaymentStatusPaid)
	}
	if contracts[1].ReviewClientID != 0 {
		t.Errorf("Contract Review Client ID is not 0")
	}
	if contracts[1].ReviewRepetitorID != 0 {
		t.Errorf("Contract Review Repetitor ID is not 0")
	}
	if contracts[1].TransactionID != 0 {
		t.Errorf("Contract Transaction ID is not 0")
	}
	contracts, err = contractRepository.GetContractsByRepetitorID(2, 1, 5, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Failed to get contracts: %v", err)
	}
	if len(contracts) != 1 {
		t.Errorf("Number of contracts is not 1")
	}
	if contracts[0].ID != insertedID {
		t.Errorf("Contract ID is not %d", insertedID)
	}
	if contracts[0].ClientID != 1 {
		t.Errorf("Contract Client ID is not 1")
	}
	if contracts[0].RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if contracts[0].Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if contracts[0].Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
	if contracts[0].Status != types.ContractStatusActive {
		t.Errorf("Contract Status is not %v", types.ContractStatusActive)
	}
	if contracts[0].PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Contract Payment Status is not %v", types.PaymentStatusPaid)
	}
	contracts, err = contractRepository.GetContractsByRepetitorID(5, 0, 5, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Failed to get contracts: %v", err)
	}
	if len(contracts) != 0 {
		t.Errorf("Number of contracts is not 0")
	}
}

func TestGetContractsByClientID(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestUser(3)
	InsertTestUser(4)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	contract = types.Contract{
		ClientID:          1,
		RepetitorID:       3,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	insertedID2, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	contract = types.Contract{
		ClientID:          4,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
		CreatedAt:         time.Now(),
		StartDate:         time.Now(),
		EndDate:           time.Now(),
	}
	_, err = contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	contracts, err := contractRepository.GetContractsByClientID(1, 0, 10, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Failed to get contracts: %v", err)
	}
	if len(contracts) != 2 {
		t.Errorf("Number of contracts is not 1")
	}
	if contracts[0].ID != insertedID2 {
		t.Errorf("Contract ID is not %d", insertedID2)
	}
	if contracts[0].ClientID != 1 {
		t.Errorf("Contract Client ID is not 1")
	}
	if contracts[0].RepetitorID != 3 {
		t.Errorf("Contract Repetitor ID is not 3")
	}
	if contracts[0].ReviewClientID != 0 {
		t.Errorf("Contract Review Client ID is not 0")
	}
	if contracts[0].ReviewRepetitorID != 0 {
		t.Errorf("Contract Review Repetitor ID is not 0")
	}
	if contracts[0].Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if contracts[0].Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
	if contracts[0].Status != types.ContractStatusActive {
		t.Errorf("Contract Status is not %v", types.ContractStatusActive)
	}
	if contracts[0].PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Contract Payment Status is not %v", types.PaymentStatusPaid)
	}

	if contracts[1].ID != insertedID {
		t.Errorf("Contract ID is not %d", insertedID)
	}
	if contracts[1].ClientID != 1 {
		t.Errorf("Contract Client ID is not 1")
	}
	if contracts[1].RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if contracts[1].Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if contracts[1].Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
	if contracts[1].Status != types.ContractStatusActive {
		t.Errorf("Contract Status is not %v", types.ContractStatusActive)
	}
	if contracts[1].PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Contract Payment Status is not %v", types.PaymentStatusPaid)
	}
	if contracts[1].RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if contracts[1].ReviewClientID != 0 {
		t.Errorf("Contract Review Client ID is not 0")
	}
	if contracts[1].ReviewRepetitorID != 0 {
		t.Errorf("Contract Review Repetitor ID is not 0")
	}
	if contracts[1].TransactionID != 0 {
		t.Errorf("Contract Transaction ID is not 0")
	}
	contracts, err = contractRepository.GetContractsByClientID(1, 1, 5, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Failed to get contracts: %v", err)
	}
	if len(contracts) != 1 {
		t.Errorf("Number of contracts is not 1")
	}
	if contracts[0].ID != insertedID {
		t.Errorf("Contract ID is not %d", insertedID)
	}
	if contracts[0].ClientID != 1 {
		t.Errorf("Contract Client ID is not 1")
	}
	if contracts[0].RepetitorID != 2 {
		t.Errorf("Contract Repetitor ID is not 2")
	}
	if contracts[0].Price != 100 {
		t.Errorf("Contract Price is not 100")
	}
	if contracts[0].Commission != 10 {
		t.Errorf("Contract Commission is not 10")
	}
	if contracts[0].Status != types.ContractStatusActive {
		t.Errorf("Contract Status is not %v", types.ContractStatusActive)
	}
	if contracts[0].PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Contract Payment Status is not %v", types.PaymentStatusPaid)
	}
	contracts, err = contractRepository.GetContractsByClientID(5, 0, 5, types.ContractStatusActive)
	if err != nil {
		t.Errorf("Failed to get contracts: %v", err)
	}
	if len(contracts) != 0 {
		t.Errorf("Number of contracts is not 0")
	}
}

func TestUpdateContractStatus(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	err = contractRepository.UpdateContractStatus(insertedID, types.ContractStatusCompleted)
	if err != nil {
		t.Errorf("Failed to update contract status: %v", err)
	}
	resultContract, err := contractRepository.GetContract(insertedID)
	if err != nil {
		t.Errorf("Failed to get contract: %v", err)
	}
	if resultContract.Status != types.ContractStatusCompleted {
		t.Errorf("Contract Status is not %v", types.ContractStatusCompleted)
	}
	err = contractRepository.UpdateContractStatus(insertedID+1, types.ContractStatusActive)
	if err == nil {
		t.Errorf("Contract status was updated for non-existent contract")
	}
}

func TestUpdateContractPaymentStatus(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	err = contractRepository.UpdateContractPaymentStatus(insertedID, types.PaymentStatusPaid)
	if err != nil {
		t.Errorf("Failed to update contract payment status: %v", err)
	}
	resultContract, err := contractRepository.GetContract(insertedID)
	if err != nil {
		t.Errorf("Failed to get contract: %v", err)
	}
	if resultContract.PaymentStatus != types.PaymentStatusPaid {
		t.Errorf("Contract Payment Status is not %v", types.PaymentStatusPaid)
	}
	err = contractRepository.UpdateContractPaymentStatus(insertedID+1, types.PaymentStatusPaid)
	if err == nil {
		t.Errorf("Contract payment status was updated for non-existent contract")
	}
}

func TestUpdateContractReviewClientID(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	err = contractRepository.UpdateContractReviewClientID(insertedID, 1)
	if err != nil {
		t.Errorf("Failed to update contract review client ID: %v", err)
	}
	resultContract, err := contractRepository.GetContract(insertedID)
	if err != nil {
		t.Errorf("Failed to get contract: %v", err)
	}
	if resultContract.ReviewClientID != 1 {
		t.Errorf("Contract Review Client ID is not %d", 1)
	}
	err = contractRepository.UpdateContractReviewClientID(insertedID+1, 1)
	if err == nil {
		t.Errorf("Contract review client ID was updated for non-existent contract")
	}
}

func TestUpdateContractReviewRepetitorID(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_contract_table CASCADE")
	contractRepository := CreateContractRepository(globalDb, "test_contract_table")
	contract := types.Contract{
		ClientID:          1,
		RepetitorID:       2,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Description:       "test",
		Status:            types.ContractStatusActive,
		PaymentStatus:     types.PaymentStatusPaid,
		TransactionID:     0,
		Price:             100,
		Commission:        10,
	}
	insertedID, err := contractRepository.InsertContract(contract)
	if err != nil {
		t.Errorf("Failed to insert contract: %v", err)
	}
	err = contractRepository.UpdateContractReviewRepetitorID(insertedID, 2)
	if err != nil {
		t.Errorf("Failed to update contract review repetitor ID: %v", err)
	}
	resultContract, err := contractRepository.GetContract(insertedID)
	if err != nil {
		t.Errorf("Failed to get contract: %v", err)
	}
	if resultContract.ReviewRepetitorID != 2 {
		t.Errorf("Contract Review Repetitor ID is not %d", 2)
	}
	err = contractRepository.UpdateContractReviewRepetitorID(insertedID+1, 2)
	if err == nil {
		t.Errorf("Contract review repetitor ID was updated for non-existent contract")
	}
}
