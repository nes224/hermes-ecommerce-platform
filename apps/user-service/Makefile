-include app.env
export

MIGRATION_PATH = db/migrations

.PHONY: migrate-up migrate-down migrate-up1 migrate-down1 migrate-create sqlc
# Generate Go code from SQL queries
sqlc:
	sqlc generate

migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose down

migrate-up1:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose up 1

migrate-down1:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose down 1

# แถม: คำสั่งสร้างไฟล์ Migration ใหม่ (ใช้: make migrate-create name=create_users_table)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

# Run App with Admin Seeding Flag
seed-admin:
	go run apps/user-service/cmd/main.go --seed-admin

# Run User Service
run:
	go run apps/user-service/cmd/main.g