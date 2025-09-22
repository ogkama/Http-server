package types

import (
	"encoding/json"
	"http_server/domain"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err != nil {
		switch err{
			case domain.ErrUserNotFound:
				w.WriteHeader(http.StatusNotFound) 
				json.NewEncoder(w).Encode(ErrorResponse{Message: "User not found"})
				return
			case domain.ErrUserAlreadyExists:
				w.WriteHeader(http.StatusConflict) 
				json.NewEncoder(w).Encode(ErrorResponse{Message: "User already exists"})
				return
			case domain.ErrSessionNotFound:
				w.WriteHeader(http.StatusNotFound) 
				json.NewEncoder(w).Encode(ErrorResponse{Message: "Session not found"})
				return
			case domain.ErrTaskNotFound:
				w.WriteHeader(http.StatusNotFound) 
				json.NewEncoder(w).Encode(ErrorResponse{Message: "Task not found"})
				return
			case domain.ErrWrongPassword:
				w.WriteHeader(http.StatusUnauthorized) 
				json.NewEncoder(w).Encode(ErrorResponse{Message: "Wrong password"})
				return
			case domain.ErrMissingToken:
				w.WriteHeader(http.StatusUnauthorized) 
				json.NewEncoder(w).Encode(ErrorResponse{Message: "Missing token"})
				return	
			case domain.ErrUnauthorized:
				w.WriteHeader(http.StatusUnauthorized) 
				json.NewEncoder(w).Encode(ErrorResponse{Message: "Unathorized"})
				return
			default:
				w.WriteHeader(http.StatusInternalServerError) 						
				json.NewEncoder(w).Encode(ErrorResponse{Message: "Iternal Error"})
				w.Write([]byte(err.Error()))
				return
			//http.StatusUnsupportedMediaType
		}
	}

	if resp != nil {
		switch v := resp.(type) {
    		case []byte:
            	contentType := http.DetectContentType(v)
        		w.Header().Set("Content-Type", contentType)
        		w.WriteHeader(http.StatusOK)
        		w.Write(v)
    		default:
				w.Header().Set("Content-Type", "application/json")
        		if err := json.NewEncoder(w).Encode(resp); err != nil {
				//http.Error(w, "Internal Error", http.StatusInternalServerError)
					w.WriteHeader(http.StatusInternalServerError) 
					json.NewEncoder(w).Encode(ErrorResponse{Message: "Iternal Error"})
					return
				}
		}
    }

}