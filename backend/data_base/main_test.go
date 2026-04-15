package data_base

import (
	"data_base_project/types"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"
)

var globalDb *sql.DB

func InsertTestUser(userId int64) {
	personalData := types.PersonalData{
		TelephoneNumber: "+88005553535",
		Email:           "test@example.com",
		PassportData: types.PassportData{
			PassportNumber:   "1234567890",
			PassportDate:     time.Now(),
			PassportSeries:   "1024",
			PassportIssuedBy: "test",
		},
		FirstName:  "Jhon",
		LastName:   "Doe",
		MiddleName: "Jhonovich",
	}
	var lastInsertedID int64
	err := globalDb.QueryRow(`INSERT INTO test_personal_data_table
	 (telephone_number, email, passport_number, passport_series, passport_date, 
	 passport_issued_by, first_name, last_name, middle_name) 
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		personalData.TelephoneNumber,
		personalData.Email,
		personalData.PassportData.PassportNumber,
		personalData.PassportData.PassportSeries,
		personalData.PassportData.PassportDate,
		personalData.PassportData.PassportIssuedBy,
		personalData.FirstName,
		personalData.LastName,
		personalData.MiddleName).Scan(&lastInsertedID)
	if err != nil {
		log.Fatalf("Error inserting personal data: %v", err)
	}
	globalDb.Exec(`INSERT INTO test_user_table 
	(id, registration_date, last_login_date, personal_data_id) 
	VALUES ($1, $2, $3, $4)`,
		userId,
		time.Now(),
		time.Now(),
		lastInsertedID)
}

func InsertTestChat(chatId int64, clientId int64, repetitorId int64, moderatorId int64) {
	globalDb.Exec(`INSERT INTO test_chat_table 
	(id, client_id, repetitor_id, moderator_id, created_at) 
	VALUES ($1, $2, $3, $4, $5)`,
		chatId,
		clientId,
		repetitorId,
		moderatorId, time.Now())
}

func InsertTestPesonalData(personalDataId int64) {
	globalDb.Exec(`INSERT INTO test_personal_data_table 
	(id, telephone_number, email, passport_number, passport_series, passport_date, passport_issued_by, first_name, last_name, middle_name) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		personalDataId,
		"+88005553535",
		"test@example.com",
		"1234567890", "1024", time.Now(), "test", "Jhon", "Doe", "Jhonovich")
}

func InsertTestContract(contractId int64, clientId int64, repetitorId int64) {
	contract := types.Contract{
		ID:                contractId,
		ClientID:          clientId,
		RepetitorID:       repetitorId,
		Status:            types.ContractStatusActive,
		TransactionID:     0,
		Price:             1000,
		Commission:        10,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		StartDate:         time.Now(),
		EndDate:           time.Now().AddDate(0, 0, 30),
		CreatedAt:         time.Now(),
		Description:       "test",
		PaymentStatus:     types.PaymentStatusNull,
	}
	globalDb.Exec(`INSERT INTO test_contract_table 
	(id, client_id, repetitor_id, transaction_id, created_at, description, status, payment_status, price, commission, start_date, end_date, review_client_id, review_repetitor_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		contract.ID,
		contract.ClientID,
		contract.RepetitorID,
		contract.TransactionID,
		contract.CreatedAt,
		contract.Description,
		contract.Status,
		contract.PaymentStatus,
		contract.Price,
		contract.Commission,
		contract.StartDate,
		contract.EndDate,
		contract.ReviewClientID,
		contract.ReviewRepetitorID)
}

func InsertTestDepartment(departmentId int64) {
	department := types.Department{
		ID:     departmentId,
		Name:   "Test Department",
		HeadID: 0,
	}
	globalDb.Exec(`INSERT INTO test_department_table 
	(id, name, head_id) 
	VALUES ($1, $2, $3)`,
		department.ID,
		department.Name,
		department.HeadID)
}

func SetupTestFrameWork(m *testing.M) {
	var err error
	globalDb, err = CreateConnection("host=db user=testuser password=testpass dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer globalDb.Close()
	DropTables(globalDb,
		"test_personal_data_table",
		"test_user_table",
		"test_auth_table",
		"test_chat_table",
		"test_message_table",
		"test_department_table",
		"test_hire_info_table",
		"test_client_table",
		"test_resume_table",
		"test_review_table",
		"test_repetitor_table",
		"test_contract_table",
		"test_admin_table",
		"test_moderator_table",
		"test_transaction_table",
		"test_pending_contract_payment_transactions_table",
		"test_lesson_table")
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	defer DropTables(globalDb,
		"test_personal_data_table",
		"test_user_table",
		"test_auth_table",
		"test_chat_table",
		"test_message_table",
		"test_department_table",
		"test_hire_info_table",
		"test_client_table",
		"test_resume_table",
		"test_review_table",
		"test_repetitor_table",
		"test_contract_table",
		"test_admin_table",
		"test_moderator_table",
		"test_transaction_table",
		"test_pending_contract_payment_transactions_table",
		"test_lesson_table")
	if err != nil {
		log.Fatalf("Error dropping tables: %v", err)
	}
	CreateTables(globalDb,
		"test_personal_data_table",
		"test_user_table",
		"test_auth_table",
		"test_chat_table",
		"test_message_table",
		"test_department_table",
		"test_hire_info_table",
		"test_client_table",
		"test_resume_table",
		"test_review_table",
		"test_repetitor_table",
		"test_contract_table",
		"test_admin_table",
		"test_moderator_table",
		"test_transaction_table",
		"test_pending_contract_payment_transactions_table",
		"test_lesson_table")
	os.Exit(m.Run())
}

func TestMain(m *testing.M) {
	SetupTestFrameWork(m)
}
