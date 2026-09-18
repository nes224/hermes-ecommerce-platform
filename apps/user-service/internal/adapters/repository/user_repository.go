package repository

import (
	"context"
	"hermes-ecommerce-platform/apps/user-service/internal/adapters/repository/db"
	"hermes-ecommerce-platform/apps/user-service/internal/core/domain"
	"hermes-ecommerce-platform/apps/user-service/internal/core/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) ports.UserRepository {
	return &userRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: pgtype.Text{String: user.PhoneNumber, Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        userID,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Time,
	}, nil
}
