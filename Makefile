DB_URL ?= postgres://dev:dev@localhost:55432/devdb?sslmode=disable
AUTH_MIGRATIONS_DIR := services/auth-service/migrations

setup:
	@docker compose pull

run-local:
	@docker compose up -d

down:
	@docker compose down

test:
	@echo "Running tests for all services..."
	@cd shared && env -u GOROOT go test ./...
	@cd services/auth-service && env -u GOROOT go test ./...
	@cd services/academic-service && env -u GOROOT go test ./...
	@cd services/admission-service && env -u GOROOT go test ./...
	@cd services/assessment-service && env -u GOROOT go test ./...
	@cd services/attendance-service && env -u GOROOT go test ./...
	@cd services/finance-service && env -u GOROOT go test ./...
	@cd services/notification-service && env -u GOROOT go test ./...
	@cd services/file-service && env -u GOROOT go test ./...
	@cd services/api-gateway && env -u GOROOT go test ./...

test-coverage:
	@echo "Running coverage for all services..."
	@mkdir -p coverage
	@cd shared && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/shared.out
	@cd services/auth-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/auth.out
	@cd services/academic-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/academic.out
	@cd services/admission-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/admission.out
	@cd services/assessment-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/assessment.out
	@cd services/attendance-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/attendance.out
	@cd services/finance-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/finance.out
	@cd services/notification-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/notification.out
	@cd services/file-service && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/file.out
	@cd services/api-gateway && env -u GOROOT go test ./... -covermode=atomic -coverprofile=$(shell pwd)/coverage/gateway.out

load-test-smoke:
	@echo "Running smoke test..."
	@docker run --rm -i --add-host=host.docker.internal:host-gateway -v $(shell pwd)/performance-tests:/scripts --network host grafana/k6 run /scripts/smoke.js

load-test-ratelimit:
	@echo "Running rate limit test..."
	@docker run --rm -i --add-host=host.docker.internal:host-gateway -v $(shell pwd)/performance-tests:/scripts --network host grafana/k6 run /scripts/rate_limit.js

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		(golangci-lint run --config .golangci.yml) || \
		(cd shared && golangci-lint run --config .golangci.shared.yml) || \
		(cd shared && golangci-lint run) || \
		(cd shared && env -u GOROOT go vet ./...); \
	else echo "golangci-lint not installed"; fi

deps:
	@cd shared && env -u GOROOT go mod tidy

migrate-up:
	@echo "Running auth-service migrations (up)"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(AUTH_MIGRATIONS_DIR) -database "$(DB_URL)" up; \
	else \
		echo "migrate CLI not found; using docker-based psql fallback"; \
		docker cp $(AUTH_MIGRATIONS_DIR)/. postgres:/migrations/; \
		docker exec postgres psql -U postgres -c "CREATE DATABASE devdb TEMPLATE template0;" || true; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/001_users_roles_permissions.up.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/002_audit_logs.up.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/003_password_resets.up.sql || docker exec postgres psql -U postgres -d devdb -f /migrations/migrations/003_password_resets.up.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/004_password_history.up.sql || docker exec postgres psql -U postgres -d devdb -f /migrations/migrations/004_password_history.up.sql; \
	fi

migrate-down:
	@echo "Running auth-service migrations (down)"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(AUTH_MIGRATIONS_DIR) -database "$(DB_URL)" down; \
	else \
		echo "migrate CLI not found; using docker-based psql fallback"; \
		docker cp $(AUTH_MIGRATIONS_DIR)/. postgres:/migrations/; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/004_password_history.down.sql || docker exec postgres psql -U postgres -d devdb -f /migrations/migrations/004_password_history.down.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/003_password_resets.down.sql || docker exec postgres psql -U postgres -d devdb -f /migrations/migrations/003_password_resets.down.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/002_audit_logs.down.sql; \
		docker exec postgres psql -U postgres -d devdb -f /migrations/001_users_roles_permissions.down.sql; \
	fi

clean:
	@rm -f coverage.out coverage_shared.out coverage_auth.out
	@find . -name \"*.out\" -delete

docker-build-all:
	@for d in services/*; do \
		if [ -f $$d/Dockerfile ]; then \
			docker build -f $$d/Dockerfile -t sisfo-`basename $$d` .; \
		else \
			echo "skip $$d (no Dockerfile)"; \
		fi; \
	done
