package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"http_server/domain"
	"http_server/repository"
	"time"

	"github.com/google/uuid"
)

type UserAndSession struct {
	user repository.User
	session repository.Session
}

func NewUser(user repository.User, session repository.Session) *UserAndSession {
	return &UserAndSession{
		user: user,
		session: session,
	}
}

func (s *UserAndSession) Register(login string, password string) (string, error){
	id := uuid.New().String()
	
	//hash := argon2.IDKey()
	hash := sha256.Sum256([]byte(password))
	hashStr := hex.EncodeToString(hash[:]) 

	return s.user.AddUser(login, id, hashStr)
}

func (s *UserAndSession) sessionId() string {
    /*b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
	}
        return base64.URLEncoding.EncodeToString(b)*/
	return uuid.New().String()
}

func (s *UserAndSession) Login(ctx context.Context, login string, password string) (*string, error){

	pass_db, err := s.user.GetPassword(login)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	hash := sha256.Sum256([]byte(password))
	hashStr := hex.EncodeToString(hash[:]) 
	if hashStr != *pass_db {
		return nil, domain.ErrWrongPassword
	}
	
	user_id, _ := s.user.GetID(login)
	session_id := s.sessionId()
	s.session.SetSession(ctx, session_id, *user_id, 15 * time.Minute)
	return &session_id, nil

}

func (s *UserAndSession) CheckSession(ctx context.Context, session_id string) (*string, error) {
	user_id, err := s.session.GetUser(ctx, session_id)
	if err != nil {
		return nil, err
	}

	return user_id, nil
}