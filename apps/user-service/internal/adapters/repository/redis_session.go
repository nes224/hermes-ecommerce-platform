package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"hermes-ecommerce-platform/apps/user-service/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) redisRepository {
	return redisRepository{
		client: client,
	}
}

func formatSessionKey(sessionID uuid.UUID) string {
	return fmt.Sprintf("session:%s", sessionID.String())
}

func (r *redisRepository) CreateSession(ctx context.Context, session *ports.Session, duration time.Duration) error {
	key := formatSessionKey(session.ID)
	
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	err = r.client.Set(ctx, key, data, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to save session to redis: %w", err)
	}
	return nil
}

func (r *redisRepository) GetSession(ctx context.Context, id uuid.UUID) (*ports.Session, error) {
	key := formatSessionKey(id)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}

	var session ports.Session
	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (r *redisRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	key := formatSessionKey(id)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}
	return nil
}