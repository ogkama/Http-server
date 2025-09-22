package usecases

import "context"

type UserAndSession interface {
	Register(username string, password string) (string, error)
	Login(ctx context.Context, username, password string) (*string, error)
	CheckSession(ctx context.Context, session_id string) (*string, error)
}