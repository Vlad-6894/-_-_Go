include .env
export

export PROJECT_ROOT=$(shell pwd)

start-app:
	@docker compose up -d --build security-golang-app

finish-app:
	@docker compose down security-golang-app

start-postgres:
	@docker compose up -d security-postgres-db

finish-postgres:
	@docker compose down security-postgres-db

cleanup-postgres:
	@read -p "Очистить pg_data? ВНИМАНИЕ! Опасность утери данных! [y/N]: " choice; \
	if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
		docker compose down security-postgres-db security-postgres-port-forwarder && \
		sudo rm -rf ${PROJECT_ROOT}/out/pg_data && \
		echo "Очищено"; \
	else \
		echo "Операция отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Вы не передали название миграции!" \
		exit 1; \
	fi;
	@docker compose run --rm security-postgres-db-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Нет параметра action"; \
		exit 1; \
	fi;
	@docker compose run --rm security-postgres-db-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@security-postgres-db:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

start-port-forwarder:
	@docker compose up -d security-postgres-port-forwarder

finish-port-forwarder:
	@docker compose down security-postgres-port-forwarder

generate-mocks:
	@go generate ./...

start-unit-tests:
	@go test -v ./internal/features/service

start-integration-tests:
	@go test -v -tags=integration ./internal/features/repository/postgres

generate-https-certc:
	@openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt -subj "/CN=localhost"