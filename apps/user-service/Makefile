DB_URL=postgresql://root:secret@localhost:5432/hermes_user?sslmode=disable

.PHONY: sqlc migrate-up migrate-down seed-admin run

# Generate Go code from SQL queries
sqlc:
	sqlc generate

# Create new migration file (usage: make migrate-create name=init_schema)
migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

# Database Migrations
migrate-up:
	golang-migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	golang-migrate -path db/migrations -database "$(DB_URL)" -verbose down

# Run App with Admin Seeding Flag
seed-admin:
	go run apps/user-service/cmd/main.go --seed-admin

# Run User Service
run:
	go run apps/user-service/cmd/main.g