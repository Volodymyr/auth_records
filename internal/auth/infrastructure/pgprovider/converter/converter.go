package converter

import (
	"auth_records/internal/auth/adapter/storage/dto"
	"auth_records/internal/auth/infrastructure/pgprovider/model"
)

func ToLoginUser(user *model.User) *dto.LoginUser {
	u := &dto.LoginUser{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}

	return u
}
