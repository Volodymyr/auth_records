ifndef CI
  ifneq ("$(wildcard .env)","")
    include .env
  endif
endif

build-auth-service:
	docker build -t auth_records/auth-service -f deploy/Dockerfile.auth-service .

build-records-service:
	docker build -t auth_records/records-service -f deploy/Dockerfile.records-service .

build-services: build-auth-service build-records-service

seed-auth-service-db:
	DATABASE_URL="postgres://${POSTGRES_USERS_USER}:${POSTGRES_USERS_PASSWORD}@localhost:5432/${POSTGRES_USERS_DB}?sslmode=disable" \
	go run internal/auth/db/seed/seed.go

run-go-linter:
	golangci-lint run -v --timeout=10m

run-go-tests:
	CGO_ENABLED=1 go test -v -race ./...

init-env:
	@test -f .env || cp .env_example .env

