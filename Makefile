DB_URL ?= postgres://dev:dev@localhost:55432/devdb?sslmode=disable
AUTH_MIGRATIONS_DIR := services/auth-service/migrations

setup:
	@docker compose pull

run-local:
	@docker compose up -d

down:
	@docker compose down

test:
	@cd shared && env -u GOROOT go test ./...

test-coverage:
	@cd shared && env -u GOROOT go test ./... -covermode=atomic -coverprofile=coverage_shared.out
	@cd services/auth-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=coverage_auth.out

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then (cd shared && golangci-lint run --config .golangci.shared.yml); else echo "golangci-lint not installed"; fi

deps:
	@cd shared && env -u GOROOT go mod tidy

migrate-up:
	@echo "Running auth-service migrations (up)"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(AUTH_MIGRATIONS_DIR) -database "$(DB_URL)" up; \
	else \
		echo "migrate CLI not found; using docker-based psql fallback"; \
		docker cp $(AUTH_MIGRATIONS_DIR) postgres:/migrations; \
		docker exec postgres psql -U postgres -c "CREATE DATABASE devdb TEMPLATE template0;" || true; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/001_users_roles_permissions.up.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/002_audit_logs.up.sql; \
	fi

migrate-down:
	@echo "Running auth-service migrations (down)"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(AUTH_MIGRATIONS_DIR) -database "$(DB_URL)" down; \
	else \
		echo "migrate CLI not found; using docker-based psql fallback"; \
		docker cp $(AUTH_MIGRATIONS_DIR) postgres:/migrations; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/002_audit_logs.down.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/001_users_roles_permissions.down.sql; \
	fi

clean:
	@rm -f coverage.out coverage_shared.out coverage_auth.out
	@find . -name \"*.out\" -delete

docker-build-all:
	@for d in services/*; do if [ -f $$d/Dockerfile ]; then docker build -t sisfo-`basename $$d` $$d; else echo "skip $$d (no Dockerfile)"; fi; done
