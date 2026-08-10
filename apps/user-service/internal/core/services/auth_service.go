package services

import (
	"hermes-ecommerce-platform/apps/user-service/internal/adapters/repository/db"
	"hermes-ecommerce-platform/apps/user-service/internal/core/ports"
	"hermes-ecommerce-platform/pkg/token"
	"time"
)

type AuthService struct {
	userRepo             db.Queries
	sessionRepo          ports.RedisRepository
	tokenMaker           token.Maker
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}
