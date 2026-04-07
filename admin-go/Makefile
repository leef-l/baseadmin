include hack/hack-cli.mk

MIGRATE_DIR ?= database/migrations
MIGRATE_NAME ?=
MIGRATE_CONFIG ?=
MIGRATE_ARGS ?=

.PHONY: migrate.up
migrate.up:
	@go run ./cmd/migrate --dir $(MIGRATE_DIR) $(if $(MIGRATE_CONFIG),--config $(MIGRATE_CONFIG),) up $(MIGRATE_ARGS)

.PHONY: migrate.down
migrate.down:
	@go run ./cmd/migrate --dir $(MIGRATE_DIR) $(if $(MIGRATE_CONFIG),--config $(MIGRATE_CONFIG),) down $(MIGRATE_ARGS)

.PHONY: migrate.version
migrate.version:
	@go run ./cmd/migrate --dir $(MIGRATE_DIR) $(if $(MIGRATE_CONFIG),--config $(MIGRATE_CONFIG),) version

.PHONY: migrate.create
migrate.create:
	@if [ "$(MIGRATE_NAME)" = "" ]; then echo "usage: make migrate.create MIGRATE_NAME=add_users_index"; exit 1; fi
	@go run ./cmd/migrate --dir $(MIGRATE_DIR) create $(MIGRATE_NAME)
