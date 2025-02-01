package converter

import (
	"auth_records/internal/auth/adapter/storage/dto"
	"auth_records/internal/auth/domain"
)

func ToServiceLoginUser(user *dto.LoginUser) *domain.User {
	return &domain.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

}
