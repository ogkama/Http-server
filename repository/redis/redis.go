package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}
func NewRedisClient(addr string, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
  		Addr:     addr,
  		Password: password,
  		DB:       db,
	})

 	if err := client.Ping(context.Background()).Err(); err != nil {
  		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

	func (rs *RedisClient) GetUser(ctx context.Context, session_id string) (*string, error) {
		user_id, err := rs.client.Get(ctx, session_id).Result()
		if err != nil {
			return nil, err
		}

		return &user_id, err
	}

	func (rs *RedisClient) SetSession(ctx context.Context, session_id string, user_id string, expiration time.Duration) error {
		return rs.client.Set(ctx, session_id, user_id, expiration).Err()
	}

	func (rs *RedisClient) DeleteSession(session_id string) error {
		return rs.client.Close()
	}