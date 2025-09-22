package http

import (
	"context"
	"http_server/api/http/types"
	"http_server/domain"
	"http_server/usecases"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Server represents an HTTP handler for managing users.
type User struct {
	user usecases.UserAndSession
}

// NewServerHandler creates a new instance of User.
func NewUserHandler(user usecases.UserAndSession) *User {
	return &User{user: user}
	
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags user
// @Accept  json
// @Produce json
// @Param request body types.PostUserHandlerRequest true "User credentials"
// @Success 200 {object} types.Message "User registered successfully"
// @Failure 400 {string} string "Bad request"
// @Router /register [post]
func (u *User) postUserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatPostUserHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	message, err := u.user.Register(req.Login, req.Password)
	types.ProcessError(w, err, &types.Message{Message: message})	
}

// @Summary User login
// @Description Authenticate user and start a session
// @Tags user
// @Accept  json
// @Produce json
// @Param request body types.PostUserHandlerRequest true "User credentials"
// @Success 200 {object} types.Message "Successful login"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func (u *User) postUserLoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatPostUserHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second) //
	defer cancel()
	
	session_id, err := u.user.Login(ctx, req.Login, req.Password)
	if err != nil {
		types.ProcessError(w, err, nil)	
		return
	}
	w.Header().Set("Authorization", *session_id)
	types.ProcessError(w, err, &types.Message{Message: "Succesfull login"})	
}

func Authorization(u *User, r *http.Request) (*string, error){
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second) //
	defer cancel()

	token :=r.Header.Get("Authorization")
	if token == ""{
		return nil, domain.ErrMissingToken
	} 
	user_id, err := u.user.CheckSession(ctx, token)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return user_id, nil
}

// WithUserHandlers registers user-related HTTP handlers.
func (u *User) WithUserHandlers(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/register", u.postUserRegisterHandler)
		r.Post("/login", u.postUserLoginHandler)
	})
}