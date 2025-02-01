package storage

import (
	"auth_records/internal/auth/adapter/storage/converter"
	"auth_records/internal/auth/adapter/storage/dto"
	"auth_records/internal/auth/domain"
	"context"
)

type storageProvider interface {
	UserByEmail(ctx context.Context, email string) (*dto.LoginUser, error)
}

type userRepository struct {
	client storageProvider
}

func NewUserRepository(client storageProvider) *userRepository {
	return &userRepository{
		client: client,
	}
}

func (r *userRepository) UserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.client.UserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return converter.ToServiceLoginUser(user), nil
}
