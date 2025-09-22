package postgres

import (
	"database/sql"
	"errors"
	"http_server/domain"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func (ps *PostgresStorage) GetID(login string) (*string, error) {
	var user_id string
	err := ps.db.QueryRow("SELECT user_id FROM users WHERE user_login = $1", login).Scan(&user_id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user_id, nil
}

func (ps *PostgresStorage) GetPassword(login string) (*string, error) {
	var password string
	err := ps.db.QueryRow("SELECT user_password FROM users WHERE user_login = $1", login).Scan(&password)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &password, nil
}

func (ps *PostgresStorage) AddUser(login string, id string, password string) (string, error) {
	_, err := ps.db.Exec("INSERT INTO users (user_id, user_login, user_password) VALUES ($1, $2, $3)", id, login, password)
	if err != nil && err.(*pq.Error).Code == "23505"{
		return "", domain.ErrUserAlreadyExists
	}

	return "New user added to database", nil
}

func (ps *PostgresStorage) DeleteUser(id string) error {
	result, err := ps.db.Exec("DELETE FROM users where user_id = $1", id)
	if err != nil {
		return err
	}

	rowsAfected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAfected == 0 {
		return errors.New("key not found")
	}

	return nil
}