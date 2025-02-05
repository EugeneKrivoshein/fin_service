DCOMPOSE = docker-compose
MIGRATE_CMD = docker exec app_container sh -c "cd /app/internal/postgres/migrations && go run migrate.go"

.PHONY: run up migrate stop clean

run: up migrate

up:
	$(DCOMPOSE) up -d --build

migrate:
	$(MIGRATE_CMD)

stop:
	$(DCOMPOSE) down

clean: stop
	docker volume rm $$(docker volume ls -q)