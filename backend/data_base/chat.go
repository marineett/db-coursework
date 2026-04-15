package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IChatRepository interface {
	InsertChat(chat types.Chat) (int64, error)
	GetChat(id int64) (*types.Chat, error)
	GetChatListByClientID(clientID int64, from int64, size int64) ([]types.Chat, error)
	GetChatListByRepetitorID(repetitorID int64, from int64, size int64) ([]types.Chat, error)
	GetChatListByModeratorID(moderatorID int64, from int64, size int64) ([]types.Chat, error)
	GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error)
	GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error)
	GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error)
}

func CreateChatTable(db *sql.DB, chatTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + chatTableName + ` (
		id SERIAL PRIMARY KEY,
		client_id INTEGER NOT NULL,
		repetitor_id INTEGER NOT NULL,
		moderator_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", chatTableName, err)
	}
	return nil
}

type ChatRepository struct {
	db        *sql.DB
	chatTable string
}

func CreateChatRepository(db *sql.DB, chatTable string) *ChatRepository {
	return &ChatRepository{
		db:        db,
		chatTable: chatTable,
	}
}

func (r *ChatRepository) InsertChat(chat types.Chat) (int64, error) {
	query := `
	INSERT INTO ` + r.chatTable + ` (client_id, repetitor_id, moderator_id, created_at) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	var lastInsertId int64
	err := r.db.QueryRow(query, chat.ClientID, chat.RepetitorID, chat.ModeratorID, chat.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (r *ChatRepository) GetChat(id int64) (*types.Chat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var chat types.Chat
	err := row.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *ChatRepository) GetChatListByClientID(clientID int64, from int64, size int64) ([]types.Chat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, clientID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]types.Chat, 0)
	for rows.Next() {
		var chat types.Chat
		err := rows.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepository) GetChatListByRepetitorID(repetitorID int64, from int64, size int64) ([]types.Chat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE repetitor_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, repetitorID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]types.Chat, 0)
	for rows.Next() {
		var chat types.Chat
		err := rows.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepository) GetChatListByModeratorID(moderatorID int64, from int64, size int64) ([]types.Chat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE moderator_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, moderatorID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]types.Chat, 0)
	for rows.Next() {
		var chat types.Chat
		err := rows.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepository) GetChatListByUserID(userID int64, from int64, size int64) ([]types.Chat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE client_id = $1 OR repetitor_id = $1 OR moderator_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]types.Chat, 0)
	for rows.Next() {
		var chat types.Chat
		err := rows.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepository) GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error) {
	query := `
	SELECT id FROM ` + r.chatTable + ` WHERE client_id = $1 AND moderator_id = $2
	`
	row := r.db.QueryRow(query, clientID, moderatorID)
	var chatID int64
	err := row.Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatID, nil
}

func (r *ChatRepository) GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error) {
	query := `
	SELECT id FROM ` + r.chatTable + ` WHERE client_id = $1 AND repetitor_id = $2
	`
	row := r.db.QueryRow(query, clientID, repetitorID)
	var chatID int64
	err := row.Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatID, nil
}

func (r *ChatRepository) GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error) {
	query := `
	SELECT id FROM ` + r.chatTable + ` WHERE moderator_id = $1 AND repetitor_id = $2
	`
	row := r.db.QueryRow(query, moderatorID, repetitorID)
	var chatID int64
	err := row.Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatID, nil
}
