package postgres

import (
	"database/sql"
	"http_server/domain"

	"github.com/lib/pq"
)

func (ps *PostgresStorage) GetStatus(task_id string) (*string, error) {
	var task_status string
	err := ps.db.QueryRow("SELECT task_status FROM tasks WHERE task_id = $1", task_id).Scan(&task_status)
	if err == sql.ErrNoRows {
		return nil, domain.ErrTaskNotFound
	} else if err != nil {
		return nil, err
	}

	return &task_status, nil
}

func (ps *PostgresStorage) GetResult(task_id string) ([]byte, error) {
	var task_result sql.NullString
	err := ps.db.QueryRow("SELECT task_result FROM tasks WHERE task_id = $1", task_id).Scan(&task_result)
	if err == sql.ErrNoRows {
		return nil, domain.ErrTaskNotFound
	} else if err != nil {
		return nil, err
	}
	if !task_result.Valid {
		return nil, nil
	}

	return []byte(task_result.String), nil
}

func (ps *PostgresStorage) Post(task_id string, task_data []byte, user_id string) error {
	_, err := ps.db.Exec("INSERT INTO tasks (task_id, user_id, task_data, task_status) VAlUES ($1, $2, $3, $4)", task_id, user_id, task_data, "in progress")
	if err != nil && err.(*pq.Error).Code == "23505"{
		return domain.ErrTaskAlreadyExists
	}

	return nil
}

func (ps *PostgresStorage) SetStatus(task_id string, task_status string) error {
	_, err := ps.db.Exec("UPDATE tasks SET task_status = $1, updated_at = now() WHERE task_id = $2", task_status, task_id)
	if err != nil{
		return err
	}

	return nil
}

func (ps *PostgresStorage) SetResult(task_id string, task_result []byte) error {
	_, err := ps.db.Exec("UPDATE tasks SET task_result = $1, updated_at = now() WHERE task_id = $2", task_result, task_id)
	if err != nil{
		return err
	}

	return nil
}