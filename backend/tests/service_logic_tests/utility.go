package service_test

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"database/sql"
	"errors"
	"sort"
	"time"
)

type TestTransactionRepository struct {
	transactions map[int64]*types.Transaction
}

func CreateTestTransaction() TestTransactionRepository {
	return TestTransactionRepository{
		transactions: make(map[int64]*types.Transaction),
	}
}

func (r *TestTransactionRepository) ApproveTransaction(transaction_id int64) error {
	return nil
}

func (r *TestTransactionRepository) GetPendingContractPaymentTransaction() (*types.PendingContractPaymentTransaction, error) {
	return nil, nil
}

func (r *TestTransactionRepository) InsertTransaction(transaction types.Transaction) (int64, error) {
	id := int64(len(r.transactions) + 1)
	transaction.ID = id
	r.transactions[id] = &transaction
	return id, nil
}

func (r *TestTransactionRepository) InsertPendingContractPaymentTransaction(
	transactionPendingContractPayment types.PendingContractPaymentTransaction,
	transaction types.Transaction,
) (int64, error) {
	return r.InsertTransaction(transaction)
}

func (r *TestTransactionRepository) UpdateTransactionStatus(transaction_id int64, status types.TransactionStatus) error {
	transaction, ok := r.transactions[transaction_id]
	if !ok {
		return errors.New("transaction not found")
	}
	transaction.Status = status
	r.transactions[transaction_id] = transaction
	return nil
}

func (r *TestTransactionRepository) GetTransaction(transaction_id int64) (*types.Transaction, error) {
	transaction, ok := r.transactions[transaction_id]
	if !ok {
		return nil, errors.New("transaction not found")
	}
	return transaction, nil
}

func (r *TestTransactionRepository) GetTransactionsList(user_id int64, from int64, size int64) ([]types.Transaction, error) {
	transactions := make([]types.Transaction, 0)
	for _, transaction := range r.transactions {
		if transaction.UserID == user_id {
			transactions = append(transactions, *transaction)
		}
	}
	if len(transactions) == 0 {
		return []types.Transaction{}, nil
	}
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.After(transactions[j].CreatedAt)
	})
	return transactions[from:min(from+size, int64(len(transactions)))], nil
}

type TestUserRepository struct {
	users map[int64]*types.UserData
}

func CreateTestUserRepository() *TestUserRepository {
	return &TestUserRepository{
		users: make(map[int64]*types.UserData),
	}
}

func (r *TestUserRepository) InsertUser(user types.UserData) (int64, error) {
	id := int64(len(r.users) + 1)
	user.ID = id
	r.users[id] = &user
	return id, nil
}

func (r *TestUserRepository) InsertUserInSeq(tx *sql.Tx, user types.UserData) (int64, error) {
	return r.InsertUser(user)
}

func (r *TestUserRepository) GetUser(user_id int64) (*types.UserData, error) {
	return r.users[user_id], nil
}

func (r *TestUserRepository) GetUserList(from int64, size int64) ([]types.UserData, error) {
	users := make([]types.UserData, 0)
	for _, user := range r.users {
		users = append(users, *user)
	}
	if len(users) == 0 {
		return []types.UserData{}, nil
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	return users[from:min(from+size, int64(len(users)))], nil
}

type TestMessageRepository struct {
	messages map[int64]*types.Message
}

func CreateTestMessageRepository() *TestMessageRepository {
	return &TestMessageRepository{
		messages: make(map[int64]*types.Message),
	}
}

func (r *TestMessageRepository) InsertMessage(message types.Message) (int64, error) {
	id := int64(len(r.messages) + 1)
	message.ID = id
	r.messages[id] = &message
	return id, nil
}

func (r *TestMessageRepository) GetMessages(chat_id int64, from int64, size int64) ([]types.Message, error) {
	messages := make([]types.Message, 0)
	for _, message := range r.messages {
		if message.ChatID == chat_id {
			messages = append(messages, *message)
		}
	}
	if len(messages) == 0 {
		return []types.Message{}, nil
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	return messages[from:min(from+size, int64(len(messages)))], nil
}

type TestChatRepository struct {
	chats map[int64]*types.Chat
}

func CreateTestChatRepository() *TestChatRepository {
	return &TestChatRepository{
		chats: make(map[int64]*types.Chat),
	}
}

func (r *TestChatRepository) GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error) {
	return 0, nil
}

func (r *TestChatRepository) GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error) {
	return 0, nil
}

func (r *TestChatRepository) GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error) {
	return 0, nil
}

func (r *TestChatRepository) InsertChat(chat types.Chat) (int64, error) {
	id := int64(len(r.chats) + 1)
	chat.ID = id
	r.chats[id] = &chat
	return id, nil
}

func (r *TestChatRepository) GetChat(chat_id int64) (*types.Chat, error) {
	chat, ok := r.chats[chat_id]
	if !ok {
		return nil, errors.New("chat not found")
	}
	return chat, nil
}

func (r *TestChatRepository) GetChatList(from int64, size int64) ([]types.Chat, error) {
	chats := make([]types.Chat, 0)
	for _, chat := range r.chats {
		chats = append(chats, *chat)
	}
	if len(chats) == 0 {
		return []types.Chat{}, nil
	}
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].CreatedAt.After(chats[j].CreatedAt)
	})
	return chats[from:min(from+size, int64(len(chats)))], nil
}

func (r *TestChatRepository) GetChatListByClientID(client_id int64, from int64, size int64) ([]types.Chat, error) {
	chats := make([]types.Chat, 0)
	for _, chat := range r.chats {
		if chat.ClientID == client_id {
			chats = append(chats, *chat)
		}
	}
	if len(chats) == 0 {
		return []types.Chat{}, nil
	}
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].CreatedAt.After(chats[j].CreatedAt)
	})
	return chats[from:min(from+size, int64(len(chats)))], nil
}

func (r *TestChatRepository) GetChatListByRepetitorID(repetitor_id int64, from int64, size int64) ([]types.Chat, error) {
	chats := make([]types.Chat, 0)
	for _, chat := range r.chats {
		if chat.RepetitorID == repetitor_id {
			chats = append(chats, *chat)
		}
	}
	if len(chats) == 0 {
		return []types.Chat{}, nil
	}
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].CreatedAt.After(chats[j].CreatedAt)
	})
	return chats[from:min(from+size, int64(len(chats)))], nil
}

func (r *TestChatRepository) GetChatListByModeratorID(moderator_id int64, from int64, size int64) ([]types.Chat, error) {
	chats := make([]types.Chat, 0)
	for _, chat := range r.chats {
		if chat.ModeratorID == moderator_id {
			chats = append(chats, *chat)
		}
	}
	if len(chats) == 0 {
		return []types.Chat{}, nil
	}
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].CreatedAt.After(chats[j].CreatedAt)
	})
	return chats[from:min(from+size, int64(len(chats)))], nil
}

func ChatCompareWithoutTime(chat1 types.Chat, chat2 types.Chat) bool {
	return chat1.ID == chat2.ID &&
		chat1.ClientID == chat2.ClientID &&
		chat1.RepetitorID == chat2.RepetitorID &&
		chat1.ModeratorID == chat2.ModeratorID
}

func ChatCompare(chat1 types.Chat, chat2 types.Chat) bool {
	return ChatCompareWithoutTime(chat1, chat2) &&
		chat1.CreatedAt == chat2.CreatedAt
}

type TestContractRepository struct {
	contracts        map[int64]*types.Contract
	reviewRepository data_base.IReviewRepository
}

func CreateTestContractRepository() *TestContractRepository {
	return &TestContractRepository{
		contracts:        make(map[int64]*types.Contract),
		reviewRepository: CreateTestReviewRepository(),
	}
}

func (r *TestContractRepository) GetAllContracts(from int64, size int64) ([]types.Contract, error) {
	contracts := make([]types.Contract, 0)
	for _, contract := range r.contracts {
		contracts = append(contracts, *contract)
	}
	if len(contracts) == 0 {
		return []types.Contract{}, nil
	}
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].CreatedAt.After(contracts[j].CreatedAt)
	})
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) UpdateContractReviewClientIDInSeq(tx *sql.Tx, contractID int64, reviewClientID int64) error {
	r.contracts[contractID].ReviewClientID = reviewClientID
	return nil
}

func (r *TestContractRepository) GetContractList(from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	contracts := make([]types.Contract, 0)
	for _, contract := range r.contracts {
		if contract.Status == status {
			contracts = append(contracts, *contract)
		}
	}
	if len(contracts) == 0 {
		return []types.Contract{}, nil
	}
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].CreatedAt.After(contracts[j].CreatedAt)
	})
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) InsertContract(contract types.Contract) (int64, error) {
	id := int64(len(r.contracts) + 1)
	contract.ID = id
	r.contracts[id] = &contract
	return id, nil
}

func (r *TestContractRepository) GetContract(contract_id int64) (*types.Contract, error) {
	contract, ok := r.contracts[contract_id]
	if !ok {
		return nil, errors.New("contract not found")
	}
	return contract, nil
}

func (r *TestContractRepository) GetContractsByRepetitorID(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	contracts := make([]types.Contract, 0)
	for _, contract := range r.contracts {
		if contract.RepetitorID == repetitorID && contract.Status == status {
			contracts = append(contracts, *contract)
		}
	}
	if len(contracts) == 0 {
		return []types.Contract{}, nil
	}
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].CreatedAt.After(contracts[j].CreatedAt)
	})
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) GetContractsByClientID(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.Contract, error) {
	contracts := make([]types.Contract, 0)
	for _, contract := range r.contracts {
		if contract.ClientID == clientID && contract.Status == status {
			contracts = append(contracts, *contract)
		}
	}
	if len(contracts) == 0 {
		return []types.Contract{}, nil
	}
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].CreatedAt.After(contracts[j].CreatedAt)
	})
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) UpdateContractStatus(id int64, status types.ContractStatus) error {
	r.contracts[id].Status = status
	return nil
}

func (r *TestContractRepository) UpdateContractPaymentStatus(id int64, paymentStatus types.PaymentStatus) error {
	r.contracts[id].PaymentStatus = paymentStatus
	return nil
}

func (r *TestContractRepository) UpdateContractReviewClientID(id int64, reviewClientID int64) error {
	r.contracts[id].ReviewClientID = reviewClientID
	return nil
}

func (r *TestContractRepository) UpdateContractReviewRepetitorID(id int64, reviewRepetitorID int64) error {
	r.contracts[id].ReviewRepetitorID = reviewRepetitorID
	return nil
}

func (r *TestContractRepository) UpdateContractRepetitorID(id int64, repetitorID int64) error {
	r.contracts[id].RepetitorID = repetitorID
	r.contracts[id].Status = types.ContractStatusActive
	return nil
}

func (r *TestContractRepository) BeginTx() (*sql.Tx, error) {
	return nil, nil
}

type TestAuthRepository struct {
	authData map[int64]*types.AuthInfo
}

func CreateTestAuthRepository() *TestAuthRepository {
	return &TestAuthRepository{
		authData: make(map[int64]*types.AuthInfo),
	}
}

func (r *TestAuthRepository) Authorize(authData types.AuthData) (types.AuthVerdict, error) {
	for _, auth := range r.authData {
		if auth.Login == authData.Login {
			if auth.Password == authData.Password {
				return types.AuthVerdict{
					UserID:   auth.UserID,
					UserType: auth.UserType,
				}, nil
			} else {
				return types.AuthVerdict{}, errors.New("wrong password, expected: " + auth.Password + ", got: " + authData.Password)
			}
		}
	}
	return types.AuthVerdict{}, errors.New("auth not found")
}

func (r *TestAuthRepository) CheckLogin(login string) (bool, error) {
	for _, auth := range r.authData {
		if auth.Login == login {
			return true, nil
		}
	}
	return false, nil
}

func (r *TestAuthRepository) InsertAuth(auth types.AuthInfo) (int64, error) {
	id := int64(len(r.authData) + 1)
	r.authData[id] = &types.AuthInfo{
		ID:       id,
		UserID:   auth.UserID,
		UserType: auth.UserType,
		Login:    auth.Login,
		Password: auth.Password,
	}
	return id, nil
}

func (r *TestAuthRepository) InsertAuthInSeq(tx *sql.Tx, auth types.AuthInfo) (int64, error) {
	return r.InsertAuth(auth)
}

func (r *TestAuthRepository) ChangePassword(userId int64, authData types.AuthData, newPassword string) error {
	for _, auth := range r.authData {
		if auth.Login == authData.Login {
			if auth.Password == authData.Password {
				auth.Password = newPassword
				return nil
			} else {
				return errors.New("wrong password, expected: " + auth.Password + ", got: " + authData.Password)
			}
		}
	}
	return errors.New("auth not found")
}

type TestPersonalDataRepository struct {
	personalData map[int64]*types.PersonalData
}

func CreateTestPersonalDataRepository() *TestPersonalDataRepository {
	return &TestPersonalDataRepository{
		personalData: make(map[int64]*types.PersonalData),
	}
}

func (r *TestPersonalDataRepository) InsertPersonalData(personalData types.PersonalData) (int64, error) {
	id := int64(len(r.personalData) + 1)
	personalData.ID = id
	r.personalData[id] = &personalData
	return id, nil
}

func (r *TestPersonalDataRepository) GetPersonalData(userID int64) (*types.PersonalData, error) {
	personalData, ok := r.personalData[userID]
	if !ok {
		return nil, errors.New("personal data not found")
	}
	return personalData, nil
}

func (r *TestPersonalDataRepository) UpdatePersonalData(id int64, personalData types.PersonalData) error {
	if _, ok := r.personalData[id]; !ok {
		return errors.New("personal data not found")
	}
	r.personalData[id] = &personalData
	return nil
}

func (r *TestPersonalDataRepository) InsertPersonalDataInSeq(tx *sql.Tx, personalData types.PersonalData) (int64, error) {
	return r.InsertPersonalData(personalData)
}

type TestResumeRepository struct {
	resumes map[int64]*types.Resume
}

func CreateTestResumeRepository() *TestResumeRepository {
	return &TestResumeRepository{
		resumes: make(map[int64]*types.Resume),
	}
}

func (r *TestResumeRepository) InsertResume(resume types.Resume) (int64, error) {
	id := int64(len(r.resumes) + 1)
	resume.ID = id
	r.resumes[id] = &resume
	return id, nil
}

func (r *TestResumeRepository) InsertResumeInSeq(tx *sql.Tx, resume types.Resume) (int64, error) {
	return r.InsertResume(resume)
}

func (r *TestResumeRepository) GetResume(resumeID int64) (*types.Resume, error) {
	resume, ok := r.resumes[resumeID]
	if !ok {
		return nil, errors.New("resume not found")
	}
	return resume, nil
}

func (r *TestResumeRepository) UpdateResumeDescription(resumeID int64, description string, updatedAt time.Time) error {
	if _, ok := r.resumes[resumeID]; !ok {
		return errors.New("resume not found")
	}
	r.resumes[resumeID].Description = description
	r.resumes[resumeID].UpdatedAt = updatedAt
	return nil
}

func (r *TestResumeRepository) UpdateResumeTitle(resumeID int64, title string, updatedAt time.Time) error {
	if _, ok := r.resumes[resumeID]; !ok {
		return errors.New("resume not found")
	}
	r.resumes[resumeID].Title = title
	r.resumes[resumeID].UpdatedAt = updatedAt
	return nil
}

func (r *TestResumeRepository) UpdateResumePrices(resumeID int64, prices map[string]int, updatedAt time.Time) error {
	if _, ok := r.resumes[resumeID]; !ok {
		return errors.New("resume not found")
	}
	r.resumes[resumeID].Prices = prices
	r.resumes[resumeID].UpdatedAt = updatedAt
	return nil
}
func (r *TestResumeRepository) DeleteResume(resumeID int64) error {
	if _, ok := r.resumes[resumeID]; !ok {
		return errors.New("resume not found")
	}
	delete(r.resumes, resumeID)
	return nil
}

type TestReviewRepository struct {
	reviews map[int64]*types.Review
}

func CreateTestReviewRepository() *TestReviewRepository {
	return &TestReviewRepository{
		reviews: make(map[int64]*types.Review),
	}
}

func (r *TestReviewRepository) InsertReview(review types.Review) (int64, error) {
	id := int64(len(r.reviews) + 1)
	review.ID = id
	r.reviews[id] = &review
	return id, nil
}

func (r *TestReviewRepository) GetReview(reviewID int64) (*types.Review, error) {
	review, ok := r.reviews[reviewID]
	if !ok {
		return nil, errors.New("review not found")
	}
	return review, nil
}

func (r *TestReviewRepository) GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.Review, error) {
	reviews := make([]types.Review, 0)
	for _, review := range r.reviews {
		if review.ClientID == clientID {
			reviews = append(reviews, *review)
		}
	}
	if len(reviews) == 0 {
		return []types.Review{}, nil
	}
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})
	return reviews[from:min(from+size, int64(len(reviews)))], nil
}

func (r *TestReviewRepository) GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.Review, error) {
	reviews := make([]types.Review, 0)
	for _, review := range r.reviews {
		if review.RepetitorID == repetitorID {
			reviews = append(reviews, *review)
		}
	}
	if len(reviews) == 0 {
		return []types.Review{}, nil
	}
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})
	return reviews[from:min(from+size, int64(len(reviews)))], nil
}

func (r *TestReviewRepository) UpdateReview(review types.Review) error {
	if _, ok := r.reviews[review.ID]; !ok {
		return errors.New("review not found")
	}
	r.reviews[review.ID] = &review
	return nil
}

func (r *TestReviewRepository) InsertReviewInSeq(tx *sql.Tx, review types.Review) (int64, error) {
	return r.InsertReview(review)
}

type TestAdminRepository struct {
	authRepository         data_base.IAuthRepository
	personalDataRepository data_base.IPersonalDataRepository
	userRepository         data_base.IUserRepository
	admins                 map[int64]*types.AdminData
}

func CreateTestAdminRepository(authRepository data_base.IAuthRepository, personalDataRepository data_base.IPersonalDataRepository, userRepository data_base.IUserRepository) *TestAdminRepository {
	return &TestAdminRepository{
		authRepository:         authRepository,
		personalDataRepository: personalDataRepository,
		userRepository:         userRepository,
		admins:                 make(map[int64]*types.AdminData),
	}
}

func (r *TestAdminRepository) InsertAdmin(admin types.AdminData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	id := int64(len(r.admins) + 1)
	admin.ID = id
	r.admins[id] = &admin
	r.authRepository.InsertAuth(types.AuthInfo{
		UserID:   id,
		Login:    auth.Login,
		Password: auth.Password,
	})

	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	admin.UserData.PersonalDataID = personalDataID
	userID, err := r.userRepository.InsertUser(admin.UserData)
	if err != nil {
		return 0, err
	}
	r.admins[id].UserData.ID = userID
	return id, nil
}

func (r *TestAdminRepository) InsertAdminInSeq(tx *sql.Tx, admin types.AdminData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	return r.InsertAdmin(admin, personalData, auth)
}

func (r *TestAdminRepository) GetAdmin(adminID int64) (*types.AdminData, error) {
	return r.admins[adminID], nil
}

func (r *TestAdminRepository) UpdateAdminPersonalData(adminID int64, personalData types.PersonalData) error {
	if _, ok := r.admins[adminID]; !ok {
		return errors.New("admin not found")
	}
	r.personalDataRepository.UpdatePersonalData(r.admins[adminID].PersonalDataID, personalData)
	return nil
}

func (r *TestAdminRepository) UpdateAdminPassword(adminID int64, authData types.AuthData, newPassword string) error {
	if _, ok := r.admins[adminID]; !ok {
		return errors.New("admin not found")
	}
	return r.authRepository.ChangePassword(r.admins[adminID].UserData.ID, authData, newPassword)
}

func (r *TestAdminRepository) UpdateAdminDepartment(adminID int64, departmentID int64) error {
	if _, ok := r.admins[adminID]; !ok {
		return errors.New("admin not found")
	}
	r.admins[adminID].DepartmentID = departmentID
	return nil
}

func (r *TestAdminRepository) UpdateAdminSalary(adminID int64, salary int64) error {
	if _, ok := r.admins[adminID]; !ok {
		return errors.New("admin not found")
	}
	r.admins[adminID].Salary = salary
	return nil
}

type TestRepetitorRepository struct {
	authRepository         data_base.IAuthRepository
	personalDataRepository data_base.IPersonalDataRepository
	userRepository         data_base.IUserRepository
	repetitors             map[int64]*types.RepetitorData
	reviewRepository       data_base.IReviewRepository
}

func CreateTestRepetitorRepository(
	authRepository data_base.IAuthRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	userRepository data_base.IUserRepository,
	reviewRepository data_base.IReviewRepository,
) *TestRepetitorRepository {
	return &TestRepetitorRepository{
		authRepository:         authRepository,
		personalDataRepository: personalDataRepository,
		userRepository:         userRepository,
		repetitors:             make(map[int64]*types.RepetitorData),
		reviewRepository:       reviewRepository,
	}
}

func (r *TestRepetitorRepository) GetRepetitorsIds(repetitorsOffset int64, repetitorsLimit int64) ([]int64, error) {
	ids := make([]int64, 0)
	for id := range r.repetitors {
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *TestRepetitorRepository) InsertRepetitor(repetitor types.RepetitorData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	id := int64(len(r.repetitors) + 1)
	repetitor.ID = id
	r.repetitors[id] = &repetitor
	r.authRepository.InsertAuth(types.AuthInfo{
		UserID:   id,
		Login:    auth.Login,
		Password: auth.Password,
	})

	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	repetitor.UserData.PersonalDataID = personalDataID
	userID, err := r.userRepository.InsertUser(repetitor.UserData)
	if err != nil {
		return 0, err
	}
	r.repetitors[id].UserData.ID = userID
	return id, nil
}

func (r *TestRepetitorRepository) InsertRepetitorInSeq(tx *sql.Tx, repetitor types.RepetitorData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	return r.InsertRepetitor(repetitor, personalData, auth)
}

func (r *TestRepetitorRepository) GetRepetitor(repetitorID int64) (*types.RepetitorData, error) {
	return r.repetitors[repetitorID], nil
}

func (r *TestRepetitorRepository) UpdateRepetitorPersonalData(repetitorID int64, personalData types.PersonalData) error {
	if _, ok := r.repetitors[repetitorID]; !ok {
		return errors.New("repetitor not found")
	}
	r.personalDataRepository.UpdatePersonalData(r.repetitors[repetitorID].PersonalDataID, personalData)
	return nil
}

func (r *TestRepetitorRepository) GetUserIdInRepetitor(repetitorID int64) (int64, error) {
	return r.repetitors[repetitorID].UserData.ID, nil
}

func (r *TestRepetitorRepository) UpdateRepetitorPassword(repetitorID int64, authData types.AuthData, newPassword string) error {
	if _, ok := r.repetitors[repetitorID]; !ok {
		return errors.New("repetitor not found")
	}
	return r.authRepository.ChangePassword(r.repetitors[repetitorID].UserData.ID, authData, newPassword)
}

type TestClientRepository struct {
	authRepository         data_base.IAuthRepository
	personalDataRepository data_base.IPersonalDataRepository
	userRepository         data_base.IUserRepository
	clients                map[int64]*types.ClientData
}

func CreateTestClientRepository(authRepository data_base.IAuthRepository, personalDataRepository data_base.IPersonalDataRepository, userRepository data_base.IUserRepository) *TestClientRepository {
	return &TestClientRepository{
		authRepository:         authRepository,
		personalDataRepository: personalDataRepository,
		userRepository:         userRepository,
		clients:                make(map[int64]*types.ClientData),
	}
}

func (r *TestClientRepository) InsertClient(client types.ClientData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	id := int64(len(r.clients) + 1)
	client.ID = id
	r.clients[id] = &client
	r.authRepository.InsertAuth(types.AuthInfo{
		UserID:   id,
		Login:    auth.Login,
		Password: auth.Password,
	})

	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	client.UserData.PersonalDataID = personalDataID
	userID, err := r.userRepository.InsertUser(client.UserData)
	if err != nil {
		return 0, err
	}
	r.clients[id].UserData.ID = userID
	return id, nil
}

func (r *TestClientRepository) InsertClientInSeq(tx *sql.Tx, client types.ClientData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	return r.InsertClient(client, personalData, auth)
}

func (r *TestClientRepository) GetClient(clientID int64) (*types.ClientData, error) {
	return r.clients[clientID], nil
}

func (r *TestClientRepository) GetUserIdInClient(clientID int64) (int64, error) {
	return r.clients[clientID].UserData.ID, nil
}

func (r *TestClientRepository) UpdateClientPassword(clientID int64, authData types.AuthData, newPassword string) error {
	if _, ok := r.clients[clientID]; !ok {
		return errors.New("client not found")
	}
	return r.authRepository.ChangePassword(r.clients[clientID].UserData.ID, authData, newPassword)
}

func (r *TestClientRepository) UpdateClientPersonalData(clientID int64, personalData types.PersonalData) error {
	if _, ok := r.clients[clientID]; !ok {
		return errors.New("client not found")
	}
	r.personalDataRepository.UpdatePersonalData(r.clients[clientID].PersonalDataID, personalData)
	return nil
}

type TestModeratorRepository struct {
	authRepository         data_base.IAuthRepository
	personalDataRepository data_base.IPersonalDataRepository
	userRepository         data_base.IUserRepository
	moderators             map[int64]*types.ModeratorData
}

func CreateTestModeratorRepository(authRepository data_base.IAuthRepository, personalDataRepository data_base.IPersonalDataRepository, userRepository data_base.IUserRepository) *TestModeratorRepository {
	return &TestModeratorRepository{
		authRepository:         authRepository,
		personalDataRepository: personalDataRepository,
		userRepository:         userRepository,
		moderators:             make(map[int64]*types.ModeratorData),
	}
}

func (r *TestModeratorRepository) GetModerators() ([]int64, error) {
	ids := make([]int64, 0)
	for id := range r.moderators {
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *TestModeratorRepository) InsertModerator(moderator types.ModeratorData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	id := int64(len(r.moderators) + 1)
	moderator.ID = id
	r.moderators[id] = &moderator
	r.authRepository.InsertAuth(types.AuthInfo{
		UserID:   id,
		Login:    auth.Login,
		Password: auth.Password,
	})

	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	moderator.UserData.PersonalDataID = personalDataID
	userID, err := r.userRepository.InsertUser(moderator.UserData)
	if err != nil {
		return 0, err
	}
	r.moderators[id].UserData.ID = userID
	return id, nil
}

func (r *TestModeratorRepository) InsertModeratorInSeq(tx *sql.Tx, moderator types.ModeratorData, personalData types.PersonalData, auth types.AuthData) (int64, error) {
	return r.InsertModerator(moderator, personalData, auth)
}

func (r *TestModeratorRepository) GetModerator(moderatorID int64) (*types.ModeratorData, error) {
	return r.moderators[moderatorID], nil
}

func (r *TestModeratorRepository) UpdateModeratorPersonalData(moderatorID int64, personalData types.PersonalData) error {
	if _, ok := r.moderators[moderatorID]; !ok {
		return errors.New("moderator not found")
	}
	r.personalDataRepository.UpdatePersonalData(r.moderators[moderatorID].PersonalDataID, personalData)
	return nil
}

func (r *TestModeratorRepository) UpdateModeratorPassword(moderatorID int64, authData types.AuthData, newPassword string) error {
	if _, ok := r.moderators[moderatorID]; !ok {
		return errors.New("moderator not found")
	}
	return r.authRepository.ChangePassword(r.moderators[moderatorID].UserData.ID, authData, newPassword)
}

func (r *TestModeratorRepository) GetUserIdInModerator(moderatorID int64) (int64, error) {
	return r.moderators[moderatorID].UserData.ID, nil
}

func (r *TestModeratorRepository) UpdateModeratorSalary(moderatorID int64, salary int64) error {
	if _, ok := r.moderators[moderatorID]; !ok {
		return errors.New("moderator not found")
	}
	r.moderators[moderatorID].Salary = salary
	return nil
}

type TestDepartmentRepository struct {
	departments map[int64]*types.Department
	hireInfo    map[int64]*types.HireInfo
}

func CreateTestDepartmentRepository() *TestDepartmentRepository {
	return &TestDepartmentRepository{
		departments: make(map[int64]*types.Department),
		hireInfo:    make(map[int64]*types.HireInfo),
	}
}

func (r *TestDepartmentRepository) GetDepartmentsByHeadID(headID int64) ([]types.Department, error) {
	departments := make([]types.Department, 0)
	for _, department := range r.departments {
		if department.HeadID == headID {
			departments = append(departments, *department)
		}
	}
	return departments, nil
}

func (r *TestDepartmentRepository) GetDepartmentIdByName(name string) (int64, error) {
	for _, department := range r.departments {
		if department.Name == name {
			return department.ID, nil
		}
	}
	return 0, errors.New("department not found")
}

func (r *TestDepartmentRepository) InsertDepartment(department types.Department) (int64, error) {
	id := int64(len(r.departments) + 1)
	department.ID = id
	r.departments[id] = &department
	return id, nil
}

func (r *TestDepartmentRepository) GetDepartment(departmentID int64) (*types.Department, error) {
	return r.departments[departmentID], nil
}

func (r *TestDepartmentRepository) ChangeDepartmentHead(departmentID int64, headID int64) error {
	if _, ok := r.departments[departmentID]; !ok {
		return errors.New("department not found")
	}
	r.departments[departmentID].HeadID = headID
	return nil
}

func (r *TestDepartmentRepository) HireInfoInsert(hireInfo types.HireInfo) error {
	id := int64(len(r.hireInfo) + 1)
	hireInfo.ID = id
	r.hireInfo[id] = &hireInfo
	return nil
}

func (r *TestDepartmentRepository) HireInfoDelete(userId int64, departmentId int64) error {
	for _, hireInfo := range r.hireInfo {
		if hireInfo.UserID == userId && hireInfo.DepartmentID == departmentId {
			delete(r.hireInfo, hireInfo.ID)
			return nil
		}
	}
	return errors.New("hire info not found")
}

func (r *TestDepartmentRepository) GetUserDepartmentsIDs(userID int64) ([]int64, error) {
	departmentsIDs := make([]int64, 0)
	for _, hireInfo := range r.hireInfo {
		if hireInfo.UserID == userID {
			departmentsIDs = append(departmentsIDs, hireInfo.DepartmentID)
		}
	}
	return departmentsIDs, nil
}

func (r *TestDepartmentRepository) GetDepartmentUsersIDs(departmentID int64) ([]int64, error) {
	usersIDs := make([]int64, 0)
	for _, hireInfo := range r.hireInfo {
		if hireInfo.DepartmentID == departmentID {
			usersIDs = append(usersIDs, hireInfo.UserID)
		}
	}
	return usersIDs, nil
}
