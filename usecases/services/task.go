package services

import (
	"encoding/base64"
	"encoding/json"
	"http_server/domain"
	"http_server/repository"
	"log"

	"github.com/google/uuid"
)

type Task struct {
	repo repository.Task
	sender repository.TaskSender
}

func NewTask(repo repository.Task, sender repository.TaskSender) *Task {
	return &Task{
		repo: repo,
		sender: sender,
	}
}

func (rs *Task) GetStatus(task_id string) (*string, error) {
	return rs.repo.GetStatus(task_id)
}

type Response struct {
	Result string `json:"result"`
}

func (rs *Task) GetResult(task_id string) ([]byte, error) {
	jsonBytes, err := rs.repo.GetResult(task_id)
	if err != nil {
		return nil, err
	}

	var img Response
	json.Unmarshal(jsonBytes, &img)

	decoded, err := base64.StdEncoding.DecodeString(img.Result)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (rs *Task) Post(data string, user_id string) (string, error) {
	task_id := uuid.New().String()	
	log.Printf("Generated task_id: %s", task_id)

	err := rs.sender.Send(domain.Task{Task_id: task_id, Status: "in progress", Result: "", Data: data})
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(data)
    if err != nil {
        log.Fatal(err)
    }
	
	return task_id, rs.repo.Post(task_id, jsonBytes, user_id)
}