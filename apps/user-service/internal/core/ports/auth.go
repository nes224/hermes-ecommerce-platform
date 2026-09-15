package ports

import (
	"context"
	"hermes-ecommerce-platform/apps/user-service/internal/adapters/repository/db"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RedisRepository interface {
	CreateSession(ctx context.Context, session *Session, duration time.Duration) error
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
}

type RegisterInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	UserAgent   string `json:"user_agent"`
	ClientIP    string `json:"client_ip"`
}

type LoginInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
}

type AuthResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type AuthService interface {
	SignUp(ctx context.Context, input *RegisterInput) (*db.User, error)
	SignIn(ctx context.Context, input *LoginInput) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
}

type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	Role      string    `json:"role"`
}

type TokenService interface {
	GenerateAccessToken(userID uuid.UUID, sessionID uuid.UUID, role string) (string, time.Time, error)
	ValidateAccessToken(tokenStr string) (*TokenClaims, error)
}