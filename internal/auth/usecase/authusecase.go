package usecase

import (
	"auth_records/internal/auth/apperrors"
	"auth_records/internal/auth/domain"
	"auth_records/pkg/password"
	"context"

	"go.uber.org/zap"
)

type userRepository interface {
	UserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TokenService interface {
	GenerateTokenWithID(userID uint64) (string, error)
	CheckJWTAndGetUserID(tokenString string) (uint64, error)
}

type usecase struct {
	repo         userRepository
	tokenService TokenService
	log          *zap.Logger
}

func New(log *zap.Logger, repo userRepository, tokenService TokenService) *usecase {
	return &usecase{
		repo,
		tokenService,
		log,
	}
}

func (u *usecase) Login(ctx context.Context, email string, pswd string) (string, error) {
	user, err := u.repo.UserByEmail(ctx, email)
	if err != nil {
		u.log.Error("User with the specified email not found", zap.Error(err))

		return "", apperrors.ErrUserNotFound
	}

	if !password.IsMatch(pswd, user.PasswordHash) {
		u.log.Warn("Invalid login credentials", zap.String("email", email))

		return "", apperrors.ErrInvalidLoginCredantials
	}

	token, err := u.tokenService.GenerateTokenWithID(user.ID)
	if err != nil {
		u.log.Error("Token generation failed", zap.Error(err))

		return "", apperrors.ErrInternalServer
	}

	return token, nil
}
