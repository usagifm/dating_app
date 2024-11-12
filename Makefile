fmt:
	go fmt ./...

build:
	go build cmd/main.go

run:
	go run cmd/main.go

migrate.up:
	go run migration/main/main.go up

migrate.rollback:
	go run migration/main/main.go rollback


stop:
	docker compose stop

docs-update:
	swag fmt
	swag init -d ./cmd/,./  -o swagger/v1 
