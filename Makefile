include .envrc

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database "postgres://postgres:12345@localhost/go_social?sslmode=disable" up

.PHONY: migrate-down
migrate-down:	
	@migrate -path $(MIGRATIONS_PATH) -database "postgres://postgres:12345@localhost/go_social?sslmode=disable" down $(filter-out $@,$(MAKECMDGOALS))