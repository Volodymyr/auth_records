package main

import (
	"auth_records/pkg/utils"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	devDBURL := "postgres://user:password@localhost:5432/database?sslmode=disable"

	dbURL := utils.GetEnv("DATABASE_URL", devDBURL)
	project := utils.GetEnv("MIGRATE_PROJECT", "")

	m, err := migrate.New(
		"file://internal/"+project+"/db/migrations",
		dbURL)
	if err != nil {
		log.Fatal(err)
	}

	arg := os.Args[1]

	switch arg {
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "force":
		step := os.Args[2]

		i, err := strconv.Atoi(step)
		if err != nil {
			log.Fatal(err)
		}

		if err := m.Force(i); err != nil {
			log.Fatal(err)
		}
	case "step":
		step := os.Args[2]

		i, err := strconv.Atoi(step)
		if err != nil {
			log.Fatal(err)
		}

		if err := m.Steps(i); err != nil {
			log.Fatal(err)
		}
	case "to":
		step := os.Args[2]

		i, err := strconv.Atoi(step)
		if err != nil {
			log.Fatal(err)
		}

		if err := m.Migrate(uint(i)); err != nil {
			log.Fatal(err)
		}
	}
}
