package dto

type LoginUser struct {
	ID           uint64
	Email        string
	Username     string
	PasswordHash string
}

type UsersRepository struct {
	Email        string
	Username     string
	PasswordHash string
}
