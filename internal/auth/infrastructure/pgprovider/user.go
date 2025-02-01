package pgprovider

import (
	"auth_records/internal/auth/adapter/storage/dto"
	"auth_records/internal/auth/infrastructure/pgprovider/converter"
	"auth_records/internal/auth/infrastructure/pgprovider/model"
	"context"
)

const queryUserByEmail = "SELECT id, username, email, password_hash FROM users WHERE email = $1 LIMIT 1"

func (p *pgProvider) UserByEmail(ctx context.Context, email string) (*dto.LoginUser, error) {
	var user model.User

	err := p.db.QueryRowContext(ctx, queryUserByEmail, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return converter.ToLoginUser(&user), nil
}
