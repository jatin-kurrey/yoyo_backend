# YOYO Go API

## Architecture

The server follows:

`routes -> controllers -> services -> repositories -> PostgreSQL`

- `cmd/api`: API entrypoint, also supports `go run cmd/api/main.go seed`
- `cmd/seed`: dedicated seed entrypoint
- `internal/config`: env-based config
- `internal/database`: GORM PostgreSQL connection
- `internal/models`: UUID GORM models
- `internal/routes`: route definitions only
- `internal/controllers`: request/response handling
- `internal/services`: business logic and transactions
- `internal/repositories`: database queries
- `internal/middleware`: auth, role checks, security headers, rate limits
- `uploads` or configured `UPLOAD_DIR`: local upload fallback
- `migrations`: golang-migrate style SQL migrations

## Run Locally

```bash
cp .env.example .env
createdb yoyo_booking
go mod tidy
go run cmd/api/main.go seed
go run cmd/api/main.go
```

Health check:

```bash
curl http://localhost:8080/api/health
```

## Migrations

Development can use `AUTO_MIGRATE=true`. In production, set `AUTO_MIGRATE=false` and run:

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

## Seed Data

The seed command creates:

- initial super admin from env
- default YOYO ticket catalog from the old frontend hardcoded data
- default public site settings

```bash
go run cmd/seed/main.go
```

## Production Notes

- Use a strong `JWT_SECRET`.
- Restrict `CORS_ALLOWED_ORIGINS`.
- Keep Razorpay secret and webhook secret only on the server.
- Run behind TLS and Nginx.
- Use SQL migrations, not AutoMigrate.
- Use `APP_ENV=production`.
