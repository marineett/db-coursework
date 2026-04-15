package server

import (
	"data_base_project/service_logic"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func SetupChatRouter(
	chatService service_logic.IChatService,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(CHAT_GET_CLIENT_CHATS, ChatGetClientChatsHandler(chatService))
	router.HandleFunc(CHAT_GET_REPETITOR_CHATS, ChatGetRepetitorChatsHandler(chatService))
	router.HandleFunc(CHAT_GET_MODERATOR_CHATS, ChatGetModeratorChatsHandler(chatService))
	router.HandleFunc(CHAT_START_CM_CHAT, ChatStartCMHandler(chatService))
	router.HandleFunc(CHAT_START_RM_CHAT, ChatStartRMHandler(chatService))
	router.HandleFunc(CHAT_START_CR_CHAT, ChatStartCRHandler(chatService))
	router.HandleFunc(CHAT_GET_CHAT, ChatGetChatHandler(chatService))
	router.HandleFunc(CHAT_SEND_MESSAGE, ChatSendMessageHandler(chatService))
	router.HandleFunc(CHAT_GET_MESSAGES, ChatGetChatMessagesHandler(chatService))
	return router
}

func ChatGetClientChatsHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("id")
		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		chatsOffsetStr := r.URL.Query().Get("chats_offset")
		chatsOffset, err := strconv.ParseInt(chatsOffsetStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatsOffset to int: %v", err)
			http.Error(w, "Invalid chatsOffset", http.StatusBadRequest)
			return
		}
		chatsLimitStr := r.URL.Query().Get("chats_limit")
		chatsLimit, err := strconv.ParseInt(chatsLimitStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatsLimit to int: %v", err)
			http.Error(w, "Invalid chatsLimit", http.StatusBadRequest)
			return
		}
		chats, err := chatService.GetChatListByClientID(clientID, chatsOffset, chatsLimit)
		if err != nil {
			log.Printf("Error getting chats: %v", err)
			http.Error(w, "Error getting chats", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chats)
		log.Printf("Chats retrieved: %v", chats)
	}
}

func ChatGetRepetitorChatsHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		chatsOffsetStr := r.URL.Query().Get("chats_offset")
		chatsOffset, err := strconv.ParseInt(chatsOffsetStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatsOffset to int: %v", err)
			http.Error(w, "Invalid chatsOffset", http.StatusBadRequest)
			return
		}
		chatsLimitStr := r.URL.Query().Get("chats_limit")
		chatsLimit, err := strconv.ParseInt(chatsLimitStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatsLimit to int: %v", err)
			http.Error(w, "Invalid chatsLimit", http.StatusBadRequest)
			return
		}
		chats, err := chatService.GetChatListByRepetitorID(repetitorID, chatsOffset, chatsLimit)
		if err != nil {
			log.Printf("Error getting chats: %v", err)
			http.Error(w, "Error getting chats", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chats)
		log.Printf("Chats retrieved: %v", chats)
	}
}

func ChatGetModeratorChatsHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		moderatorIDStr := r.URL.Query().Get("id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting moderatorID to int: %v", err)
			http.Error(w, "Invalid moderatorID", http.StatusBadRequest)
			return
		}
		chatsOffsetStr := r.URL.Query().Get("chats_offset")
		chatsOffset, err := strconv.ParseInt(chatsOffsetStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatsOffset to int: %v", err)
			http.Error(w, "Invalid chatsOffset", http.StatusBadRequest)
			return
		}
		chatsLimitStr := r.URL.Query().Get("chats_limit")
		chatsLimit, err := strconv.ParseInt(chatsLimitStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatsLimit to int: %v", err)
			http.Error(w, "Invalid chatsLimit", http.StatusBadRequest)
			return
		}
		chats, err := chatService.GetChatListByModeratorID(moderatorID, chatsOffset, chatsLimit)
		if err != nil {
			log.Printf("Error getting chats: %v", err)
			http.Error(w, "Error getting chats", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chats)
		log.Printf("Chats retrieved: %v", chats)
	}
}

func ChatStartCMHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("c_id")
		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting moderatorID to int: %v", err)
			http.Error(w, "Invalid moderatorID", http.StatusBadRequest)
			return
		}
		_, err = chatService.CreateCMChat(clientID, moderatorID)
		if err != nil {
			log.Printf("Error creating chat: %v", err)
			http.Error(w, "Error creating chat", http.StatusInternalServerError)
			return
		}
		log.Printf("Chat created")
		chatID, err := chatService.GetChatIdByCIDAndMID(clientID, moderatorID)
		if err != nil {
			log.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(chatID)
	}
}

func ChatStartRMHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("r_id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting moderatorID to int: %v", err)
			http.Error(w, "Invalid moderatorID", http.StatusBadRequest)
			return
		}
		_, err = chatService.CreateRMChat(repetitorID, moderatorID)
		if err != nil {
			log.Printf("Error creating chat: %v", err)
			http.Error(w, "Error creating chat", http.StatusInternalServerError)
			return
		}
		chatID, err := chatService.GetChatIdByMIDAndRID(moderatorID, repetitorID)
		if err != nil {
			log.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusInternalServerError)
			return
		}
		log.Printf("Chat created")
		json.NewEncoder(w).Encode(chatID)
	}
}

func ChatStartCRHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("c_id")
		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		repetitorIDStr := r.URL.Query().Get("r_id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		_, err = chatService.CreateCRChat(clientID, repetitorID)
		if err != nil {
			log.Printf("Error creating chat: %v", err)
			http.Error(w, "Error creating chat", http.StatusInternalServerError)
			return
		}
		chatID, err := chatService.GetChatIdByCIDAndRID(clientID, repetitorID)
		if err != nil {
			log.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusInternalServerError)
			return
		}
		log.Printf("Chat created")
		json.NewEncoder(w).Encode(chatID)
	}
}

func ChatGetChatHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		chatIDStr := r.URL.Query().Get("id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		chat, err := chatService.GetChat(chatID)
		if err != nil {
			log.Printf("Error getting chat: %v", err)
			http.Error(w, "Error getting chat", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(*chat)
		log.Printf("Chat retrieved: %v", chat)
	}
}

func ChatSendMessageHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		message := ""
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			http.Error(w, "Error decoding message", http.StatusBadRequest)
			return
		}
		senderIDStr := r.URL.Query().Get("sender_id")
		senderID, err := strconv.ParseInt(senderIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting senderID to int: %v", err)
			http.Error(w, "Invalid senderID", http.StatusBadRequest)
			return
		}
		chatIDStr := r.URL.Query().Get("chat_id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		err = chatService.SendMessage(chatID, senderID, message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			http.Error(w, "Error sending message", http.StatusInternalServerError)
			return
		}
		log.Printf("Message sent")
		w.WriteHeader(http.StatusOK)
	}
}

func ChatGetChatMessagesHandler(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		chatIDStr := r.URL.Query().Get("id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		messagesOffsetStr := r.URL.Query().Get("messages_offset")
		messagesOffset, err := strconv.ParseInt(messagesOffsetStr, 10, 64)
		if err != nil {
			log.Printf("Error converting messagesOffset to int: %v", err)
			http.Error(w, "Invalid messagesOffset", http.StatusBadRequest)
			return
		}
		messagesLimitStr := r.URL.Query().Get("messages_limit")
		messagesLimit, err := strconv.ParseInt(messagesLimitStr, 10, 64)
		if err != nil {
			log.Printf("Error converting messagesLimit to int: %v", err)
			http.Error(w, "Invalid messagesLimit", http.StatusBadRequest)
			return
		}
		messages, err := chatService.GetMessages(chatID, messagesOffset, messagesLimit)
		if err != nil {
			log.Printf("Error getting messages: %v", err)
			http.Error(w, "Error getting messages", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(messages)
		log.Printf("Messages retrieved: %v", messages)
	}
}
