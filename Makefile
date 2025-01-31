ifndef CI
  ifneq ("$(wildcard .env)","")
    include .env
  endif
endif

# Build auth-service Docker image
build-auth-service:
	docker build -t auth_records/auth-service -f deploy/Dockerfile.auth-service .

# Build records-service Docker image
build-records-service:
	docker build -t auth_records/records-service -f deploy/Dockerfile.records-service .

# Build both services
build-services: build-auth-service build-records-service

# Run migrations for Records DB
migrate-records-server-db:
	DATABASE_URL="postgres://${POSTGRES_RECORS_USER}:${POSTGRES_RECORS_PASSWORD}@localhost:5434/${POSTGRES_RECORS_DB}?sslmode=disable" \
	MIGRATE_PROJECT=records \
	go run scripts/db/migrations.go $(cmd)

# Run migrations for Users DB
migrate-users-server-db:
	DATABASE_URL="postgres://${POSTGRES_USERS_USER}:${POSTGRES_USERS_PASSWORD}@localhost:5434/${POSTGRES_USERS_DB}?sslmode=disable" \
	MIGRATE_PROJECT=auth \
	go run scripts/db/migrations.go $(cmd)

# Seed the Users service DB
seed-users-service-db:
	DATABASE_URL="postgres://${POSTGRES_USERS_USER}:${POSTGRES_USERS_PASSWORD}@localhost:5432/${POSTGRES_USERS_DB}?sslmode=disable" \
	go run internal/auth/db/seed/seed.go

# Run Go linter to check code quality
run-go-linter:
	golangci-lint run -v --timeout=10m

# Run Go tests with race condition checks
run-go-tests:
	CGO_ENABLED=1 go test -v -race ./...

# Initialize .env file if it doesn't exist
init-env:
	@test -f .env || cp .env_example .env

