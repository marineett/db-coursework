package data_base

import (
	"data_base_project/types"
	"testing"
	"time"
)

func TestCreateMessageRepository(t *testing.T) {
	messageRepository := CreateMessageRepository(globalDb, "test_message_table")
	if messageRepository == nil {
		t.Errorf("Failed to create message repository")
	}
}

func TestInsertMessage(t *testing.T) {
	InsertTestChat(1, 1, 2, 0)
	InsertTestUser(1)
	messageRepository := CreateMessageRepository(globalDb, "test_message_table")
	message := types.Message{
		ChatID:    1,
		SenderID:  1,
		Content:   "Hello, world!",
		CreatedAt: time.Now(),
	}
	insertID, err := messageRepository.InsertMessage(message)
	if err != nil {
		t.Errorf("Failed to insert message: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_auth_table, test_chat_table, test_message_table CASCADE")
	result := types.Message{}
	err = globalDb.QueryRow("SELECT * FROM test_message_table WHERE id = $1", insertID).Scan(&result.ID, &result.ChatID, &result.SenderID, &result.Content, &result.CreatedAt)
	if err != nil {
		t.Errorf("Failed to get message: %v", err)
	}
	if result.ID != insertID {
		t.Errorf("Inserted message ID does not match: %d != %d", result.ID, insertID)
	}
	if result.ChatID != message.ChatID {
		t.Errorf("Inserted message ChatID does not match: %d != %d", result.ChatID, message.ChatID)
	}
	if result.SenderID != message.SenderID {
		t.Errorf("Inserted message SenderID does not match: %d != %d", result.SenderID, message.SenderID)
	}
	if result.Content != message.Content {
		t.Errorf("Inserted message Content does not match: %s != %s", result.Content, message.Content)
	}
}

func TestGetMessages(t *testing.T) {
	InsertTestUser(1)
	InsertTestUser(2)
	InsertTestChat(1, 1, 2, 0)
	messageRepository := CreateMessageRepository(globalDb, "test_message_table")
	insertIDs := make([]int64, 0)
	message := types.Message{
		ChatID:    1,
		SenderID:  1,
		Content:   "Hello!",
		CreatedAt: time.Now(),
	}
	insertID, err := messageRepository.InsertMessage(message)
	if err != nil {
		t.Errorf("Failed to insert message: %v", err)
	}
	insertIDs = append(insertIDs, insertID)
	message = types.Message{
		ChatID:    1,
		SenderID:  2,
		Content:   "Hi!",
		CreatedAt: time.Now(),
	}
	insertID, err = messageRepository.InsertMessage(message)
	if err != nil {
		t.Errorf("Failed to insert message: %v", err)
	}
	insertIDs = append(insertIDs, insertID)
	message = types.Message{
		ChatID:    1,
		SenderID:  1,
		Content:   "...",
		CreatedAt: time.Now(),
	}
	insertID, err = messageRepository.InsertMessage(message)
	insertIDs = append(insertIDs, insertID)
	if err != nil {
		t.Errorf("Failed to insert message: %v", err)
	}
	defer globalDb.Exec("TRUNCATE TABLE test_user_table, test_personal_data_table, test_auth_table, test_chat_table, test_message_table CASCADE")

	resultMessages, err := messageRepository.GetMessages(1, 0, 10)
	if err != nil {
		t.Errorf("Failed to get messages: %v", err)
	}
	if len(resultMessages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(resultMessages))
	}
	if resultMessages[0].ID != insertIDs[2] {
		t.Errorf("Message ID does not match: %d != %d", resultMessages[0].ID, insertIDs[2])
	}
	if resultMessages[1].ID != insertIDs[1] {
		t.Errorf("Message ID does not match: %d != %d", resultMessages[1].ID, insertIDs[1])
	}
	if resultMessages[2].ID != insertIDs[0] {
		t.Errorf("Message ID does not match: %d != %d", resultMessages[2].ID, insertIDs[0])
	}
	if resultMessages[0].Content != "..." {
		t.Errorf("Message Content does not match: %s != %s", resultMessages[0].Content, "...")
	}
	if resultMessages[1].Content != "Hi!" {
		t.Errorf("Message Content does not match: %s != %s", resultMessages[1].Content, "Hi!")
	}
	if resultMessages[2].Content != "Hello!" {
		t.Errorf("Message Content does not match: %s != %s", resultMessages[2].Content, "Hello!")
	}
	resultMessages, err = messageRepository.GetMessages(2, 1, 5)
	if err != nil {
		t.Errorf("Failed to get messages: %v", err)
	}
	if len(resultMessages) != 0 {
		t.Errorf("Expected 0 message, got %d", len(resultMessages))
	}
	resultMessages, err = messageRepository.GetMessages(2, 0, 1)
	if err != nil {
		t.Errorf("Failed to get messages: %v", err)
	}
	if len(resultMessages) != 0 {
		t.Errorf("Expected 0 message, got %d", len(resultMessages))
	}
	resultMessages, err = messageRepository.GetMessages(1, 2, 5)
	if err != nil {
		t.Errorf("Failed to get messages: %v", err)
	}
	if len(resultMessages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(resultMessages))
	}
	if resultMessages[0].ID != insertIDs[0] {
		t.Errorf("Message ID does not match: %d != %d", resultMessages[0].ID, insertIDs[2])
	}
	if resultMessages[0].Content != "Hello!" {
		t.Errorf("Message Content does not match: %s != %s", resultMessages[0].Content, "Hello!")
	}
}
