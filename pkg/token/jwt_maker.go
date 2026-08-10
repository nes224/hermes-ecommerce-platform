package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least 32 chracters")
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

type JWTCustomerClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func (maker *JWTMaker) CreateToken(userID uuid.UUID, role string, duration time.Duration) (string, *Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	claims := JWTCustomerClaims{
		UserID: payload.UserID,
		Role:   payload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        payload.ID.String(),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}

	return tokenStr, payload, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTCustomerClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*JWTCustomerClaims)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	tokenID, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        tokenID,
		UserID:    claims.UserID,
		Role:      claims.Role,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}, nil

}
