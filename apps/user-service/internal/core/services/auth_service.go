package services

import (
	"context"
	"errors"
	"fmt"
	"hermes-ecommerce-platform/apps/user-service/internal/adapters/repository/db"
	"hermes-ecommerce-platform/apps/user-service/internal/core/ports"
	"hermes-ecommerce-platform/pkg/token"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo             db.Querier
	sessionRepo          ports.RedisRepository
	tokenMaker           token.Maker
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(userRepo db.Querier, sessionRepo ports.RedisRepository, tokenMaker token.Maker, accessDuration time.Duration, refreshDuration time.Duration) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		tokenMaker:           tokenMaker,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input *ports.RegisterInput) (*db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	arg := db.CreateUserParams{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PhoneNumber:  pgtype.Text{String: input.PhoneNumber, Valid: true},
		Role:         db.UserRoleCUSTOMER,
	}

	row, err := s.userRepo.CreateUser(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &db.User{
		ID:          row.ID,
		Email:       row.Email,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		PhoneNumber: row.PhoneNumber,
		Role:        row.Role,
		IsActive:    row.IsActive,
		IsVerified:  row.IsVerified,
		LastLoginAt: row.LastLoginAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, err
}

func (s *AuthService) SignIn(ctx context.Context, input *ports.LoginInput) (*ports.AuthResponse, error) {
	user, err := s.userRepo.GetUserForAuth(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password: %w", err)
	}

	userID, err := uuid.FromBytes(user.ID[:])
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(userID, string(user.Role), s.accessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(userID, string(user.Role), s.refreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	session := &ports.Session{
		ID:           refreshPayload.ID,
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    input.UserAgent,
		ClientIP:     input.ClientIP,
		IsBlocked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	if err := s.sessionRepo.CreateSession(ctx, session, s.refreshTokenDuration); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &ports.AuthResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*ports.AuthResponse, error) {
	refreshPayload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token: %w", err)
	}

	session, err := s.sessionRepo.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}

	if session.IsBlocked {
		return nil, errors.New("blocked session")
	}

	if session.RefreshToken != refreshToken {
		return nil, errors.New("mismatched session token")
	}

	user, err := s.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	userID, err := uuid.FromBytes(user.ID[:])
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(userID, string(user.Role), s.accessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	return &ports.AuthResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.sessionRepo.DeleteSession(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
