package services

import (
	"context"
	"hermes-ecommerce-platform/apps/user-service/internal/core/domain"
	"hermes-ecommerce-platform/apps/user-service/internal/core/ports"

	"github.com/google/uuid"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, input *ports.UpdateProfileInput) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.PhoneNumber = input.PhoneNumber

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
