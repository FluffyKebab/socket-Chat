package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/beevik/guid"
)

type Chat struct {
	ID   int    `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (s *Storer) GetChatDataFromKey(chatKey string) (Chat, error) {
	var chat Chat

	dataRow := s.DataBase.QueryRow("SELECT * FROM chats WHERE key = ?", chatKey)
	err := dataRow.Scan(&chat.ID, &chat.Key, &chat.Name)

	return chat, err
}

func (s *Storer) AddChat(name string) error {
	var key string

	for {
		key = guid.New().String()

		_, err := s.GetChatDataFromKey(key)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}

			return fmt.Errorf("Klarte ikke å scane databasen: %v", err)
		}
	}

	_, err := s.DataBase.Exec("INSERT INTO chats (name, key) VALUES (?, ?)", name, key)
	if err != nil {
		return fmt.Errorf("Klarte ikke å sette inn %v", err)
	}

	return nil
}

func (s *Storer) HasChat(chatKey string) bool {
	row := s.DataBase.QueryRow("SELECT * FROM chats WHERE key = ?", chatKey)
	err := row.Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		fmt.Println("Kunne ikke scane databsen", err)
		return false
	}

	return true
}

func (s *Storer) GetUserChats(sessionID string) ([]Chat, error) {
	data, err := s.GetUserData(sessionID)
	if err != nil {
		return nil, fmt.Errorf("Klarte ikke å få brukerens data: %v", err)
	}

	if data.Chats == "" {
		return make([]Chat, 0), nil
	}

	var chatKeys []string
	err = json.Unmarshal([]byte(data.Chats), &chatKeys)
	if err != nil {
		return nil, fmt.Errorf("Klarte ikke å unmarhle json fra databasen: %v", err)
	}

	chats := make([]Chat, 0)
	for i := 0; i < len(chatKeys); i++ {
		chatData, err := s.GetChatDataFromKey(chatKeys[i])
		if err != nil {
			return nil, err
		}

		chats = append(chats, chatData)
	}

	return chats, nil
}

func (s *Storer) GetNewestChats(min, max int) ([]Chat, error) {
	rows, err := s.DataBase.Query("SELECT * FROM (SELECT * FROM chats ORDER BY id DESC LIMIT ?, ?)s ORDER BY id ASC ", min, max)
	if err != nil {
		return nil, fmt.Errorf("Kan ikke quarye databasen: %v", err)
	}

	chats := make([]Chat, 0)
	for rows.Next() {
		var chat Chat

		err := rows.Scan(chat.Key)
		if err != nil {
			return nil, fmt.Errorf("Klarte ikke å scane: %v", err)
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (s *Storer) GetBiggestChats(min, max int) ([]Chat, error) {
	rows, err := s.DataBase.Query("SELECT * FROM (SELECT * FROM chats ORDER BY id DESC LIMIT ?, ?)s ORDER BY id ASC ", min, max)
	if err != nil {
		return nil, fmt.Errorf("Kan ikke quarye databasen: %v", err)
	}

	chats := make([]Chat, 0)
	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.ID, chat.Key, chat.Name)
		if err != nil {
			return nil, fmt.Errorf("Klarte ikke å scane: %v", err)
		}

		chats = append(chats, chat)
	}

	return chats, nil
}
