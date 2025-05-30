package sqlite

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

type Database struct {
	Connect *sql.DB
}

func New() (*Database, error) {
	os.Mkdir("internal/sqlite/db", 0755)

	connect, err := sql.Open("sqlite", "internal/sqlite/db/database.db")
	if err != nil {
		return nil, err
	}

	err = connect.Ping()
	if err != nil {
		return nil, err
	}

	connect.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Username TEXT UNIQUE NOT NULL,
			Password TEXT NOT NULL
		)
	`)

	connect.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER NOT NULL,
			Title TEXT NOT NULL,
			Description TEXT,
			Priorety TEXT,
			Date TEXT
		)
	`)

	return &Database{Connect: connect}, nil
}
