package domain

type User struct {
	ID           uint64
	Email        string
	Username     string
	PasswordHash string
}
