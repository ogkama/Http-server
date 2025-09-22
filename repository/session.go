package repository

import (
	"context"
	"time"
)

type Session interface {
	GetUser(ctx context.Context, session_id string) (*string, error)
	SetSession(ctx context.Context, session_id string, user_id string, expiration time.Duration) error
	DeleteSession(session_id string) error
}