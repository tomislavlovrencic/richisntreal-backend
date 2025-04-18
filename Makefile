# Customize as needed:
DB_USER?=root
DB_PASSWORD?=
DB_HOST?=localhost
DB_PORT?=3306
DB_NAME?=richisntreal

MIGRATE = migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?multiStatements=true"

.PHONY: migrate-up migrate-down migrate-version

migrate-up:
	@$(MIGRATE) up

migrate-down:
	@$(MIGRATE) down

migrate-version:
	@$(MIGRATE) version
