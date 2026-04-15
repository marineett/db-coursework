package service_logic

import (
	service_test "data_base_project/tests/service_logic_tests"
	"data_base_project/types"
	"testing"
	"time"
)

func TestGetChat(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	first_chat := types.Chat{
		ID:          1,
		ClientID:    0,
		RepetitorID: 1,
		ModeratorID: 2,
		CreatedAt:   time.Now(),
	}

	second_chat := types.Chat{
		ID:          2,
		ClientID:    3,
		RepetitorID: 4,
		ModeratorID: 0,
		CreatedAt:   time.Now(),
	}

	third_chat := types.Chat{
		ID:          3,
		ClientID:    5,
		RepetitorID: 0,
		ModeratorID: 2,
		CreatedAt:   time.Now(),
	}
	chatRepository.InsertChat(first_chat)
	chatRepository.InsertChat(second_chat)
	chatRepository.InsertChat(third_chat)
	chat, err := chatService.GetChat(1)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}

	if chat.ID != 1 {
		t.Errorf("Expected chat ID 1, got %d", chat.ID)
	}
	if !service_test.ChatCompare(*chat, first_chat) {
		t.Errorf("Expected chat %v, got %v", first_chat, chat)
	}
	chat, err = chatService.GetChat(2)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompare(*chat, second_chat) {
		t.Errorf("Expected chat %v, got %v", second_chat, chat)
	}
	chat, err = chatService.GetChat(3)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompare(*chat, third_chat) {
		t.Errorf("Expected chat %v, got %v", third_chat, chat)
	}
	_, err = chatService.GetChat(4)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreateCRChat(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	chatFirstID, err := chatService.CreateCRChat(1, 2)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chatSecondID, err := chatService.CreateCRChat(1, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chatThirdID, err := chatService.CreateCRChat(2, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chat, err := chatService.GetChat(chatFirstID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatFirstID,
		ClientID:    1,
		RepetitorID: 2,
		ModeratorID: 0,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatFirstID,
			ClientID:    1,
			RepetitorID: 2,
			ModeratorID: 0,
		}, *chat)
	}
	chat, err = chatService.GetChat(chatSecondID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatSecondID,
		ClientID:    1,
		RepetitorID: 3,
		ModeratorID: 0,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatSecondID,
			ClientID:    1,
			RepetitorID: 3,
			ModeratorID: 0,
		}, *chat)
	}
	chat, err = chatService.GetChat(chatThirdID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatThirdID,
		ClientID:    2,
		RepetitorID: 3,
		ModeratorID: 0,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatThirdID,
			ClientID:    2,
			RepetitorID: 3,
			ModeratorID: 0,
		}, *chat)
	}
	var invalid_id int64
	for i := range 5 {
		if int64(i) != chatFirstID && int64(i) != chatSecondID && int64(i) != chatThirdID {
			invalid_id = int64(i)
			break
		}
	}
	_, err = chatService.GetChat(invalid_id)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
func TestCreateRMChat(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	chatFirstID, err := chatService.CreateRMChat(1, 2)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chatSecondID, err := chatService.CreateRMChat(1, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chatThirdID, err := chatService.CreateRMChat(2, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chat, err := chatService.GetChat(chatFirstID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatFirstID,
		ClientID:    0,
		RepetitorID: 1,
		ModeratorID: 2,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatFirstID,
			ClientID:    0,
			RepetitorID: 1,
			ModeratorID: 2,
		}, *chat)
	}
	chat, err = chatService.GetChat(chatSecondID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatSecondID,
		ClientID:    0,
		RepetitorID: 1,
		ModeratorID: 3,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatSecondID,
			ClientID:    0,
			RepetitorID: 1,
			ModeratorID: 3,
		}, *chat)
	}
	chat, err = chatService.GetChat(chatThirdID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatThirdID,
		ClientID:    0,
		RepetitorID: 2,
		ModeratorID: 3,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatThirdID,
			ClientID:    0,
			RepetitorID: 2,
			ModeratorID: 3,
		}, *chat)
	}
	var invalid_id int64
	for i := range 5 {
		if int64(i) != chatFirstID && int64(i) != chatSecondID && int64(i) != chatThirdID {
			invalid_id = int64(i)
			break
		}
	}
	_, err = chatService.GetChat(invalid_id)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreateCMChat(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	chatFirstID, err := chatService.CreateCMChat(1, 2)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chatSecondID, err := chatService.CreateCMChat(1, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chatThirdID, err := chatService.CreateCMChat(2, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	chat, err := chatService.GetChat(chatFirstID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatFirstID,
		ClientID:    1,
		RepetitorID: 0,
		ModeratorID: 2,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatFirstID,
			ClientID:    1,
			RepetitorID: 0,
			ModeratorID: 2,
		}, *chat)
	}
	chat, err = chatService.GetChat(chatSecondID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatSecondID,
		ClientID:    1,
		RepetitorID: 0,
		ModeratorID: 3,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatSecondID,
			ClientID:    1,
			RepetitorID: 0,
			ModeratorID: 3,
		}, *chat)
	}
	chat, err = chatService.GetChat(chatThirdID)
	if err != nil {
		t.Errorf("Error getting chat: %v", err)
	}
	if !service_test.ChatCompareWithoutTime(*chat, types.Chat{
		ID:          chatThirdID,
		ClientID:    2,
		RepetitorID: 0,
		ModeratorID: 3,
	}) {
		t.Errorf("Expected chat %v, got %v", types.Chat{
			ID:          chatThirdID,
			ClientID:    2,
			RepetitorID: 0,
			ModeratorID: 3,
		}, *chat)
	}
	var invalid_id int64
	for i := range 5 {
		if int64(i) != chatFirstID && int64(i) != chatSecondID && int64(i) != chatThirdID {
			invalid_id = int64(i)
			break
		}
	}
	_, err = chatService.GetChat(invalid_id)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetChatListByClientID(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	chatService.CreateCRChat(1, 2)
	chatService.CreateCRChat(1, 3)
	chatService.CreateCRChat(2, 3)
	chatService.CreateCRChat(1, 4)
	chatService.CreateCRChat(2, 4)
	chatService.CreateCRChat(3, 4)
	chatService.CreateCMChat(4, 5)
	chatService.CreateCMChat(1, 6)
	chatService.CreateCMChat(1, 7)
	chatService.CreateCMChat(7, 8)
	chatService.CreateCMChat(8, 9)
	chatService.CreateCMChat(1, 10)
	chatService.CreateRMChat(1, 11)
	chatService.CreateRMChat(1, 12)
	chatService.CreateRMChat(1, 13)
	chatService.CreateRMChat(1, 14)
	chatService.CreateRMChat(1, 15)
	chatService.CreateRMChat(1, 16)

	firstChats, err := chatService.GetChatListByClientID(1, 0, 10)
	if err != nil {
		t.Errorf("Error getting chat list by client ID: %v", err)
	}

	if len(firstChats) != 6 {
		t.Errorf("Expected 6 chats, got %d", len(firstChats))
	}
	for _, chat := range firstChats {
		if chat.ClientID != 1 {
			t.Errorf("Expected chat client ID %d, got %d", 1, chat.ClientID)
		}
	}
	for i := 1; i < len(firstChats); i++ {
		if firstChats[i].CreatedAt.After(firstChats[i-1].CreatedAt) {
			t.Errorf("Expected chat %d to be after chat %d", firstChats[i].ID, firstChats[i-1].ID)
		}
	}
	secondChats, err := chatService.GetChatListByClientID(1, 2, 3)
	if err != nil {
		t.Errorf("Error getting chat list by client ID: %v", err)
	}
	if len(secondChats) != 3 {
		t.Errorf("Expected 3 chats, got %d", len(secondChats))
	}
	for _, chat := range secondChats {
		if chat.ClientID != 1 {
			t.Errorf("Expected chat client ID %d, got %d", 1, chat.ClientID)
		}
	}
	for i := range secondChats {
		if !service_test.ChatCompareWithoutTime(secondChats[i], firstChats[i+2]) {
			t.Errorf("Order between windows error, Expected chat %d to be after chat %d", secondChats[i].ID, firstChats[i+2].ID)
		}
	}
}

func TestGetChatListByRepetitorID(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	chatService.CreateCRChat(2, 1)
	chatService.CreateCRChat(3, 1)
	chatService.CreateCRChat(3, 2)
	chatService.CreateCRChat(4, 1)
	chatService.CreateCRChat(4, 2)
	chatService.CreateCRChat(4, 3)
	chatService.CreateRMChat(4, 5)
	chatService.CreateRMChat(1, 6)
	chatService.CreateRMChat(1, 7)
	chatService.CreateRMChat(1, 8)
	chatService.CreateRMChat(8, 9)
	chatService.CreateRMChat(9, 10)
	chatService.CreateCMChat(1, 10)
	chatService.CreateCMChat(1, 11)
	chatService.CreateCMChat(1, 12)
	chatService.CreateCMChat(1, 13)
	chatService.CreateCMChat(1, 14)
	chatService.CreateCMChat(1, 15)

	firstChats, err := chatService.GetChatListByRepetitorID(1, 0, 10)
	if err != nil {
		t.Errorf("Error getting chat list by repetitor ID: %v", err)
	}
	if len(firstChats) != 6 {
		t.Errorf("Expected 6 chats, got %d", len(firstChats))
	}
	for _, chat := range firstChats {
		if chat.RepetitorID != 1 {
			t.Errorf("Expected chat repetitor ID %d, got %d", 1, chat.RepetitorID)
		}
	}
	for i := 1; i < len(firstChats); i++ {
		if firstChats[i].CreatedAt.After(firstChats[i-1].CreatedAt) {
			t.Errorf("Expected chat %d to be after chat %d", firstChats[i].ID, firstChats[i-1].ID)
		}
	}
	secondChats, err := chatService.GetChatListByRepetitorID(1, 2, 3)
	if err != nil {
		t.Errorf("Error getting chat list by repetitor ID: %v", err)
	}
	if len(secondChats) != 3 {
		t.Errorf("Expected 3 chats, got %d", len(secondChats))
	}
	for _, chat := range secondChats {
		if chat.RepetitorID != 1 {
			t.Errorf("Expected chat repetitor ID %d, got %d", 1, chat.RepetitorID)
		}
	}
	for i := range secondChats {
		if !service_test.ChatCompareWithoutTime(secondChats[i], firstChats[i+2]) {
			t.Errorf("Order between windows error, Expected chat %d to be after chat %d", secondChats[i].ID, firstChats[i+2].ID)
		}
	}
}

func TestGetChatListByModeratorID(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	chatService.CreateCMChat(2, 1)
	chatService.CreateCMChat(3, 1)
	chatService.CreateCMChat(3, 2)
	chatService.CreateCMChat(4, 1)
	chatService.CreateCMChat(4, 2)
	chatService.CreateCMChat(4, 3)
	chatService.CreateRMChat(4, 5)
	chatService.CreateRMChat(6, 1)
	chatService.CreateRMChat(7, 1)
	chatService.CreateRMChat(8, 7)
	chatService.CreateRMChat(9, 6)
	chatService.CreateRMChat(10, 1)
	chatService.CreateCRChat(1, 10)
	chatService.CreateCRChat(1, 11)
	chatService.CreateCRChat(1, 12)
	chatService.CreateCRChat(1, 13)
	chatService.CreateCRChat(1, 14)
	chatService.CreateCRChat(1, 15)

	firstChats, err := chatService.GetChatListByModeratorID(1, 0, 10)
	if err != nil {
		t.Errorf("Error getting chat list by moderator ID: %v", err)
	}
	if len(firstChats) != 6 {
		t.Errorf("Expected 6 chats, got %d", len(firstChats))
	}
	for _, chat := range firstChats {
		if chat.ModeratorID != 1 {
			t.Errorf("Expected chat moderator ID %d, got %d", 1, chat.ModeratorID)
		}
	}
	for i := 1; i < len(firstChats); i++ {
		if firstChats[i].CreatedAt.After(firstChats[i-1].CreatedAt) {
			t.Errorf("Expected chat %d to be after chat %d", firstChats[i].ID, firstChats[i-1].ID)
		}
	}
	secondChats, err := chatService.GetChatListByModeratorID(1, 2, 3)
	if err != nil {
		t.Errorf("Error getting chat list by moderator ID: %v", err)
	}
	if len(secondChats) != 3 {
		t.Errorf("Expected 3 chats, got %d", len(secondChats))
	}
	for _, chat := range secondChats {
		if chat.ModeratorID != 1 {
			t.Errorf("Expected chat moderator ID %d, got %d", 1, chat.ModeratorID)
		}
	}
	for i := range secondChats {
		if !service_test.ChatCompareWithoutTime(secondChats[i], firstChats[i+2]) {
			t.Errorf("Order between windows error, Expected chat %d to be after chat %d", secondChats[i].ID, firstChats[i+2].ID)
		}
	}
}

func TestSendMessage(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	firstChatID, err := chatService.CreateCRChat(1, 2)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	_, err = chatService.CreateCRChat(1, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	_, err = chatService.CreateCRChat(2, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	firstMessage := "Hello!"
	secondMessage := "Hello again!"
	chatService.SendMessage(firstChatID, 1, firstMessage)
	message, err := messageRepository.GetMessages(firstChatID, 0, 1)
	if err != nil {
		t.Errorf("Error getting messages: %v", err)
	}
	if len(message) != 1 {
		t.Errorf("Expected 1 message, got %d", len(message))
	}
	if message[0].Content != firstMessage {
		t.Errorf("Expected message content %s, got %s", firstMessage, message[0].Content)
	}
	chatService.SendMessage(firstChatID, 2, secondMessage)
	message, err = messageRepository.GetMessages(firstChatID, 0, 2)
	if err != nil {
		t.Errorf("Error getting messages: %v", err)
	}
	if len(message) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(message))
	}
	if message[0].Content != secondMessage {
		t.Errorf("Expected message content %s, got %s", secondMessage, message[0].Content)
	}
	if message[1].Content != firstMessage {
		t.Errorf("Expected message content %s, got %s", firstMessage, message[1].Content)
	}
}

func TestGetMessages(t *testing.T) {
	chatRepository := service_test.CreateTestChatRepository()
	messageRepository := service_test.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)

	firstChatID, err := chatService.CreateCRChat(1, 2)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	_, err = chatService.CreateCRChat(1, 3)
	if err != nil {
		t.Errorf("Error creating chat: %v", err)
	}
	firstMessage := "Hello!"
	secondMessage := "Hello again!"
	chatService.SendMessage(firstChatID, 1, firstMessage)
	chatService.SendMessage(firstChatID, 2, secondMessage)

	message, err := chatService.GetMessages(firstChatID, 0, 10)
	if err != nil {
		t.Errorf("Error getting messages: %v", err)
	}
	if len(message) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(message))
	}
	if message[0].Content != secondMessage {
		t.Errorf("Expected message content %s, got %s", secondMessage, message[0].Content)
	}
	if message[1].Content != firstMessage {
		t.Errorf("Expected message content %s, got %s", firstMessage, message[1].Content)
	}
}
