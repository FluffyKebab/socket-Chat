package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID          int    `json:"id"`
	SessionID   string `json:"sessionID"`
	DisplayName string `json:"displayName"`
	Chats       string `json:"chats"`
}

type Storer struct {
	DataBase *sql.DB
}

func New(path string) (*Storer, error) {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Klarte ikke å åpne databasen %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, sessionID string, displayName string, chats string);")
	if err != nil {
		return nil, fmt.Errorf("Klart ikke å lage users table %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS chats (id INTEGER PRIMARY KEY AUTOINCREMENT, name string, key string);")
	if err != nil {
		return nil, fmt.Errorf("klarte ikke å lage chats tatble %v", err)
	}

	return &Storer{
		DataBase: db,
	}, nil
}
