package usecases


type Task interface {
	Post(data string, user_id string) (string, error)
	GetStatus(task_id string) (*string, error)
	GetResult(task_id string) ([]byte, error)
}