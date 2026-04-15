package main

import (
	"data_base_project/data_base"
	"data_base_project/service_logic"
	"data_base_project/utility_module"
	"log"
	"os"

	"console/console_module"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	defer utility_module.UnsetEnv()
	db, err := data_base.CreateConnection(data_base.GetConnectionString())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	err = data_base.DropTables(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("CHAT_TABLE_NAME"),
		os.Getenv("MESSAGE_TABLE_NAME"),
		os.Getenv("DEPARTMENT_TABLE_NAME"),
		os.Getenv("HIRE_INFO_TABLE_NAME"),
		os.Getenv("CLIENT_TABLE_NAME"),
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("REVIEW_TABLE_NAME"),
		os.Getenv("REPEATITOR_TABLE_NAME"),
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("ADMIN_TABLE_NAME"),
		os.Getenv("MODERATOR_TABLE_NAME"),
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("PENDING_CONTRACT_PAYMENT_TRANSACTIONS"),
		os.Getenv("LESSON_TABLE_NAME"),
	)
	if err != nil {
		log.Fatalf("Error dropping tables: %v", err)
		return
	}
	err = data_base.CreateTables(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("CHAT_TABLE_NAME"),
		os.Getenv("MESSAGE_TABLE_NAME"),
		os.Getenv("DEPARTMENT_TABLE_NAME"),
		os.Getenv("HIRE_INFO_TABLE_NAME"),
		os.Getenv("CLIENT_TABLE_NAME"),
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("REVIEW_TABLE_NAME"),
		os.Getenv("REPEATITOR_TABLE_NAME"),
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("ADMIN_TABLE_NAME"),
		os.Getenv("MODERATOR_TABLE_NAME"),
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("PENDING_CONTRACT_PAYMENT_TRANSACTIONS"),
		os.Getenv("LESSON_TABLE_NAME"),
	)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	adminRepository := data_base.CreateAdminRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("ADMIN_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
	)
	moderatorRepository := data_base.CreateModeratorRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("MODERATOR_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
	)
	clientRepository := data_base.CreateClientRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("CLIENT_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
	)
	repetitorRepository := data_base.CreateRepetitorRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("REPEATITOR_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("REVIEW_TABLE_NAME"),
	)
	contractRepository := data_base.CreateContractRepository(
		db,
		os.Getenv("CONTRACT_TABLE_NAME"),
	)
	chatRepository := data_base.CreateChatRepository(
		db,
		os.Getenv("CHAT_TABLE_NAME"),
	)
	messageRepository := data_base.CreateMessageRepository(
		db,
		os.Getenv("MESSAGE_TABLE_NAME"),
	)
	resumeRepository := data_base.CreateResumeRepository(
		db,
		os.Getenv("RESUME_TABLE_NAME"),
	)
	transactionRepository := data_base.CreateTransactionRepository(
		db,
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("PENDING_CONTRACT_PAYMENT_TRANSACTIONS"),
	)
	reviewRepository := data_base.CreateReviewRepository(
		db,
		os.Getenv("REVIEW_TABLE_NAME"),
	)
	authRepository := data_base.CreateAuthRepository(
		db,
		os.Getenv("AUTH_TABLE_NAME"),
	)
	departmentRepository := data_base.CreateDepartmentRepository(
		db,
		os.Getenv("DEPARTMENT_TABLE_NAME"),
		os.Getenv("HIRE_INFO_TABLE_NAME"),
	)
	personalDataRepository := data_base.CreatePersonalDataRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
	)
	lessonRepository := data_base.CreateLessonRepository(db,
		os.Getenv("LESSON_TABLE_NAME"),
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("TRANSACTION_TABLE_NAME"),
	)
	serviceModule := service_logic.CreateServiceModule(
		authRepository,
		adminRepository,
		moderatorRepository,
		clientRepository,
		repetitorRepository,
		contractRepository,
		reviewRepository,
		chatRepository,
		messageRepository,
		resumeRepository,
		transactionRepository,
		departmentRepository,
		personalDataRepository,
		lessonRepository,
	)
	console_module.MainMenu(serviceModule, db)
}
