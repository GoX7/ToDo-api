package sqlite

import (
	"errors"
)

func (d *Database) Register(username string, password string) (int64, error) {
	ex, err := d.Connect.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return 0, err
	}

	id, err := ex.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) Login(username string, password string) (int, error) {
	var id int
	err := d.Connect.QueryRow("SELECT id FROM users WHERE username=? AND password=?", username, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("User not found")
	}

	return id, nil
}
