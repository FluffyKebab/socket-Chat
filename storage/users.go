package storage

import (
	"database/sql"
	"fmt"

	"github.com/beevik/guid"
	_ "github.com/mattn/go-sqlite3"
)

func (s *Storer) GetUserData(sessionID string) (*User, error) {
	row := s.DataBase.QueryRow("SELECT * FROM users WHERE sessionID = ?", sessionID)

	var user User
	err := row.Scan(&user.ID, &user.SessionID, &user.DisplayName, &user.Chats)
	if err != nil {
		return nil, fmt.Errorf("Noe gikk galt under scaningen %v", err)
	}

	return &user, nil
}

func (s *Storer) NewGuestUser() (string, error) {
	var sessionID string

	for {
		sessionID = guid.New().String()

		row := s.DataBase.QueryRow("SELECT * FROM users WHERE sessionID = ?", sessionID)
		err := row.Scan()

		if err != nil {
			if err == sql.ErrNoRows {
				break
			}

			return "", fmt.Errorf("klarte ikke å scane %v", err)
		}
	}

	runes := []rune(sessionID)
	displayName := "NO-USER-" + string(runes[0:4])

	_, err := s.DataBase.Exec("INSERT INTO users (sessionID, displayName, chats) VALUES (?, ?, ?)", sessionID, displayName, "")
	if err != nil {
		return "", fmt.Errorf("klarte ikke å sette in den nye brukeren %v", err)
	}

	return sessionID, nil
}
