package main

import (
	"auth_records/pkg/password"
	"auth_records/pkg/utils"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var userData = []struct {
	UserName string `db:"username"`
	Password string `db:"password_hash"`
	Email    string `db:"email"`
}{
	{"admin", "admin", "admin@admin.com"},
	{"bob", "bobsneg", "bob@mail.com"},
	{"igor", "igor", "igor@mail.com"},
}

func main() {
	// Replace with your database URL
	dbURL := utils.GetEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable")

	db, err := sqlx.Connect("postgres", dbURL)

	if err != nil {
		log.Fatalln(err)
	}

	tx := db.MustBegin()

	for _, user := range userData {
		heshPassword, err := password.HashPassword([]byte(user.Password))
		if err != nil {
			log.Fatalln("Error hashing password", err)
		}

		_, err = tx.Exec(`INSERT INTO users (username, password_hash, email) VALUES($1, $2, $3)`, user.UserName, heshPassword, user.Email)
		if err != nil {
			log.Fatalln("Error inserting data:", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln("Failed to commit transaction:", err)
	}

	log.Println("User seeded successfully!")
}
