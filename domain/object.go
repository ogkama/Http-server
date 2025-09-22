package domain

type Task struct {
	Task_id string `json:"task_id"`
	Status  string `json:"status"`
	Result  string `json:"result"`
	Data    string `json:"data"`
}