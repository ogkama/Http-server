package http

import (
	"http_server/api/http/types"
	"http_server/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server represents an HTTP handler for managing tasks.
type Task struct {
	service usecases.Task
}

// NewServerHandler creates a new instance of Task.
func NewTaskHandler(service usecases.Task) *Task {
	return &Task{service: service}
	
}

// @Summary Get task status
// @Description Retrieve the status of a task by its ID
// @Tags task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Session ID"
// @Param task_id path string true "Task ID"
// @Success 200 {object} types.GetTaskStatusHandlerResponse "Task status received successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /task/{task_id}/status [get]
func (t *Task) getTaskStatusHandler(w http.ResponseWriter, r *http.Request, u *User) {
	_, err := Authorization(u, r)
	if err != nil {
		types.ProcessError(w, err, nil)
		return
	}

	req, err := types.CreateGetTaskHandlerRequest(r)	
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	status, err := t.service.GetStatus(req.Task_id)
	types.ProcessError(w, err, &types.GetTaskStatusHandlerResponse{Status: status})
}

// @Summary Get task result
// @Description Retrieve the result of a completed task by its ID
// @Tags task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Session ID"
// @Param task_id path string true "Task ID"
// @Success 200 {object} types.GetTaskResultHandlerResponse "Task result received successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /task/{task_id}/result [get]
func (t *Task) getTaskResultHandler(w http.ResponseWriter, r *http.Request, u *User) {
	_, err := Authorization(u, r)
	if err != nil {
		types.ProcessError(w, err, nil)
		return
	}

	req, err := types.CreateGetTaskHandlerRequest(r)	
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	result, err := t.service.GetResult(req.Task_id)
	types.ProcessError(w, err, result)
}

// @Summary Create a new task with generated task_id
// @Description Submit a new task for processing
// @Tags task
// @Accept  json
// @Produce json
// @Param Authorization header string true "Session ID"
// @Param file formData file true "Image file to upload"
// @Success 200 {string} string "Task ID returned successfully"
// @Failure 400 {string} string "Bad request"
// @Router /task/upload [post]
func (t *Task) postTaskHandler(w http.ResponseWriter, r *http.Request, u *User) {
	user_id, err := Authorization(u, r)
	if err != nil {
		types.ProcessError(w, err, nil)
		return
	}
	req, err := types.CreatePostTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	task_id, err := t.service.Post(req.Data, *user_id)
	types.ProcessError(w, err, &types.GetTaskIdPostHandlerResponse{Task_id: task_id})
}

// WithServerHandlers registers task-related HTTP handlers.
func (t *Task) WithTaskHandlers(r chi.Router, user *User) {
	r.Route("/task", func(r chi.Router) {
		r.Post("/upload", func(w http.ResponseWriter, r *http.Request) { 
			t.postTaskHandler(w, r, user) 
		})
		r.Get("/{task_id}/status", func(w http.ResponseWriter, r *http.Request) { 
			t.getTaskStatusHandler(w, r, user) 
		})
		r.Get("/{task_id}/result", func(w http.ResponseWriter, r *http.Request) { 
			t.getTaskResultHandler(w, r, user) 
		})
	})
}