GOOSE_DRIVER = mysql
GOOSE_DBSTRING = root:shopdev@tcp(127.0.0.1:3300)/shopdevgo
GOOSE_MIGRATION_DIR = sql/schema

APP_NAME = go-ecommerce-backend-api

dev:
	go run ./cmd/server/

docker_build:
	docker compose up -d --build
	docker compose ps

docker_down:
	docker-compose down

docker_up:
	docker-compose up -d

upgoose:
	powershell -Command "$$env:GOOSE_DRIVER='$(GOOSE_DRIVER)'; $$env:GOOSE_DBSTRING='$(GOOSE_DBSTRING)'; goose -dir=$(GOOSE_MIGRATION_DIR) up"

downgoose:
	powershell -Command "$$env:GOOSE_DRIVER='$(GOOSE_DRIVER)'; $$env:GOOSE_DBSTRING='$(GOOSE_DBSTRING)'; goose -dir=$(GOOSE_MIGRATION_DIR) down"

resetgoose:
	powershell -Command "$$env:GOOSE_DRIVER='$(GOOSE_DRIVER)'; $$env:GOOSE_DBSTRING='$(GOOSE_DBSTRING)'; goose -dir=$(GOOSE_MIGRATION_DIR) reset"

sqlgen:
	sqlc generate

up_by_one:
	powershell -Command "$$env:GOOSE_DRIVER='$(GOOSE_DRIVER)'; $$env:GOOSE_DBSTRING='$(GOOSE_DBSTRING)'; goose -dir=$(GOOSE_MIGRATION_DIR) up-by-one"

#create new migration
create_migration:
	powershell -Command "goose -dir=$(GOOSE_MIGRATION_DIR) create $(name) sql" 

swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs

.PHONY: dev docker_build docker_down docker docker_up upgoose downgoose resetgoose sqlgen swag

# .PHONY: air
