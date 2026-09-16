package ports

import (
	"context"
	"hermes-ecommerce-platform/apps/user-service/internal/core/domain"

	"github.com/google/uuid"
)

type UpdateProfileInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input *UpdateProfileInput) (*domain.User, error) 
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}