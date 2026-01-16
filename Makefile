.PHONY: migrate-up migrate-down migrate-status migrate-create

# Переменные окружения
DB_HOST ?= localhost
DB_PORT ?= 8083
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= postgres

DATABASE_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Команды миграций
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-down-1:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

migrate-status:
	migrate -path migrations -database "$(DATABASE_URL)" version

migrate-create:
	@read -p "Введите название миграции: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Полный цикл пересоздания БД
reset-db:
	psql -h $(DB_HOST) -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	psql -h $(DB_HOST) -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"
	$(MAKE) migrate-up

# Запустить миграции с выводом
run-migrations:
	@echo "🔧 Применяем миграции..."
	$(MAKE) migrate-up
	@echo "✅ Миграции успешно применены!"