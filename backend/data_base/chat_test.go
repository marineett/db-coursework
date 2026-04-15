package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateChatRepository(t *testing.T) {
	chatRepository := CreateChatRepository(globalDb, "test_chat_table")
	if chatRepository == nil {
		t.Errorf("Failed to create chat repository")
	}
}

func TestInsertChat(t *testing.T) {
	chatRepository := CreateChatRepository(globalDb, "test_chat_table")
	createdAt := time.Now()
	chat := types.Chat{
		ClientID:    1,
		RepetitorID: 2,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}

	_, err := chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_message_table, test_chat_table CASCADE")
	result := globalDb.QueryRow("SELECT client_id, repetitor_id, moderator_id, created_at FROM test_chat_table WHERE client_id = 1")
	if result == nil {
		t.Errorf("Expected result, got nil")
	}
	var clientID int64
	var repetitorID int64
	var moderatorID int64
	var resultCreatedAt time.Time
	err = result.Scan(&clientID, &repetitorID, &moderatorID, &resultCreatedAt)
	if err != nil {
		t.Errorf("Failed to scan result: %v", err)
	}
	if clientID != 1 {
		t.Errorf("Expected client_id %d, got %d", 1, clientID)
	}
	if repetitorID != 2 {
		t.Errorf("Expected repetitor_id %d, got %d", 2, repetitorID)
	}
	if moderatorID != 0 {
		t.Errorf("Expected moderator_id %d, got %d", 0, moderatorID)
	}
}

func TestGetChat(t *testing.T) {
	chatRepository := CreateChatRepository(globalDb, "test_chat_table")
	createdAt := time.Now()
	chat := types.Chat{
		ClientID:    1,
		RepetitorID: 2,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}
	chatID, err := chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_message_table, test_chat_table CASCADE")
	result, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Errorf("Failed to get chat: %v", err)
	}
	if result.ID != chatID {
		t.Errorf("Expected user id %d, got %d", chatID, result.ID)
	}
	if result.ClientID != 1 {
		t.Errorf("Expected client_id %d, got %d", 1, result.ClientID)
	}
	if result.RepetitorID != 2 {
		t.Errorf("Expected repetitor_id %d, got %d", 2, result.RepetitorID)
	}
	if result.ModeratorID != 0 {
		t.Errorf("Expected moderator_id %d, got %d", 0, result.ModeratorID)
	}

}

func TestGetChatsByClientID(t *testing.T) {
	chatRepository := CreateChatRepository(globalDb, "test_chat_table")
	createdAt := time.Now()
	chatsId := make([]int64, 0)
	chat := types.Chat{
		ClientID:    1,
		RepetitorID: 2,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}
	id, err := chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)
	createdAt = time.Now()
	chat = types.Chat{
		ClientID:    1,
		RepetitorID: 3,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}
	id, err = chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)
	createdAt = time.Now()
	chat = types.Chat{
		ClientID:    2,
		RepetitorID: 4,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}
	id, err = chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)

	defer globalDb.Exec("TRUNCATE TABLE test_message_table, test_chat_table CASCADE")
	result, err := chatRepository.GetChatListByClientID(1, 0, 10)
	if err != nil {
		t.Errorf("Failed to get chats: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 chats, got %d", len(result))
	}
	if result[0].ID != chatsId[1] {
		t.Errorf("Expected chat id %d, got %d", chatsId[1], result[0].ID)
	}
	if result[1].ID != chatsId[0] {
		t.Errorf("Expected chat id %d, got %d", chatsId[0], result[1].ID)
	}
}

func TestGetChatsByRepetitorID(t *testing.T) {
	chatRepository := CreateChatRepository(globalDb, "test_chat_table")
	createdAt := time.Now()
	chatsId := make([]int64, 0)
	chat := types.Chat{
		ClientID:    1,
		RepetitorID: 2,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}
	id, err := chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)
	createdAt = time.Now()
	chat = types.Chat{
		ClientID:    5,
		RepetitorID: 2,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}
	id, err = chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)
	createdAt = time.Now()
	chat = types.Chat{
		ClientID:    6,
		RepetitorID: 3,
		ModeratorID: 0,
		CreatedAt:   createdAt,
	}
	id, err = chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)

	defer globalDb.Exec("TRUNCATE TABLE test_message_table, test_chat_table CASCADE")
	result, err := chatRepository.GetChatListByRepetitorID(2, 0, 10)
	if err != nil {
		t.Errorf("Failed to get chats: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 chats, got %d", len(result))
	}
	if result[0].ID != chatsId[1] {
		t.Errorf("Expected chat id %d, got %d", chatsId[1], result[0].ID)
	}
	if result[1].ID != chatsId[0] {
		t.Errorf("Expected chat id %d, got %d", chatsId[0], result[1].ID)
	}
}

func TestGetChatsByModeratorID(t *testing.T) {
	chatRepository := CreateChatRepository(globalDb, "test_chat_table")

	createdAt := time.Now()
	chatsId := make([]int64, 0)
	chat := types.Chat{
		ClientID:    3,
		RepetitorID: 4,
		ModeratorID: 1,
		CreatedAt:   createdAt,
	}
	id, err := chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)
	createdAt = time.Now()
	chat = types.Chat{
		ClientID:    7,
		RepetitorID: 4,
		ModeratorID: 1,
		CreatedAt:   createdAt,
	}
	id, err = chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)
	createdAt = time.Now()
	chat = types.Chat{
		ClientID:    8,
		RepetitorID: 4,
		ModeratorID: 2,
		CreatedAt:   createdAt,
	}
	id, err = chatRepository.InsertChat(chat)
	if err != nil {
		t.Errorf("Failed to insert chat: %v", err)
	}
	chatsId = append(chatsId, id)

	defer globalDb.Exec("TRUNCATE TABLE test_message_table, test_chat_table CASCADE")
	result, err := chatRepository.GetChatListByModeratorID(1, 0, 10)
	if err != nil {
		t.Errorf("Failed to get chats: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 chats, got %d", len(result))
	}
	if result[0].ID != chatsId[1] {
		t.Errorf("Expected chat id %d, got %d", chatsId[1], result[0].ID)
	}
	if result[1].ID != chatsId[0] {
		t.Errorf("Expected chat id %d, got %d", chatsId[0], result[1].ID)
	}
}
