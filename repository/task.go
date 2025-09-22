package repository

type Task interface {
	Post(task_id string, data []byte, user_id string) error
	GetStatus(task_id string) (*string, error)
	GetResult(task_id string) ([]byte, error)
	SetStatus(task_id string,status string) error
	SetResult(task_id string, result []byte) error
}

