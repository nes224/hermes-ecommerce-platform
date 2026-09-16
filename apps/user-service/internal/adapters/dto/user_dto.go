package dto

import (
	"hermes-ecommerce-platform/apps/user-service/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" binding:"required,max=100"`
	LastName    string `json:"last_name" binding:"required,max=100"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UserProfileResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToUserProfileResponse(user *domain.User) *UserProfileResponse {
	return &UserProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
