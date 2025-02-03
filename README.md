# auth_records

## To Run

Initialize .env

```bash
make init-env
```
Build services

```bash
make build-services
```

Run services

```bash
docker compose up
```

Run migrations for Records DB

```bash
make migrate-records-server-db
```

Run migrations for Users DB

```bash
make migrate-users-server-db
```

Seed the Users auth DB

```bash
make seed-users-service-db
```

Generate records protos

```bash
make generate-records-protos
```

## Check the services

```bash
curl -X POST http://localhost:8081/api/v1/login \
     -H "Content-Type: application/json" \
     -d '{"email": "admin@admin.com", "password": "admin"}'
```