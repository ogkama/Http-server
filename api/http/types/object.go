package types

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetTaskStatusHandlerResponse struct {
	Status  *string `json:"status"`
}

type GetTaskResultHandlerResponse struct {
	Result  *string `json:"result"`
}

type GetTaskHandlerRequest struct {
	Task_id string `json:"task_id"`
}

type GetTaskIdPostHandlerResponse struct {
	Task_id string `json:"task_id"`
}

func CreateGetTaskHandlerRequest(r *http.Request) (*GetTaskHandlerRequest, error) {
	//task_id := r.URL.Query().Get("task_id")
	task_id := chi.URLParam(r, "task_id")
	if task_id == "" {
		return nil, fmt.Errorf("missing task_id")
	}
	return &GetTaskHandlerRequest{Task_id: task_id}, nil
}

type PostTaskHandlerRequest struct {
	Data string `json:"data"`
}

func CreatePostTaskHandlerRequest(r *http.Request) (*PostTaskHandlerRequest, error) {
	img, header, err := r.FormFile("file")
	log.Printf("Got file: %s, size: %d byteb", header.Filename, header.Size)

	if err != nil {
		log.Printf("error while decoding json: %v", err)
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}

	defer img.Close()

	_, err = img.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("error while resetting file reader: %v", err)
		return nil, fmt.Errorf("error while resetting file reader: %v", err)
	}

    fileBytes, err  := io.ReadAll(img)
	if err != nil {
		log.Printf("error while reading file: %v", err)
		return nil, fmt.Errorf("error while reading file: %v", err)
	}
	
    bufferReader := bytes.NewReader(fileBytes)
	_, _, err = image.Decode(bufferReader)
	if err != nil {
		log.Printf("file is not a valid image: %v", err)
		return nil, fmt.Errorf("file is not a valid image: %v", err)
	}
	
	encoded := base64.StdEncoding.EncodeToString(fileBytes)

    return &PostTaskHandlerRequest{
		Data: encoded,
	}, nil
}

type PostUserHandlerRequest struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type Message struct {
	Message string `json:"message"`
}

func CreatPostUserHandlerRequest(r *http.Request) (*PostUserHandlerRequest, error){
	var req PostUserHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	if req.Login == "" || req.Password == "" {
		return nil, fmt.Errorf("missing login or password")
	}
	return &req, nil	
}