setup:
	@docker compose --env-file .env.example pull

run-local:
	@docker compose --env-file .env.example up -d

down:
	@docker compose --env-file .env.example down

test:
	@cd shared && go test ./...

test-coverage:
	@cd shared && go test ./... -coverprofile=coverage.out && mv coverage.out ../coverage.out

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then (cd shared && golangci-lint run); else echo "golangci-lint not installed"; fi

deps:
	@cd shared && go mod tidy

migrate-up:
	@echo "migrate-up not implemented yet"

migrate-down:
	@echo "migrate-down not implemented yet"

docker-build-all:
	@for d in services/*; do if [ -f $$d/Dockerfile ]; then docker build -t sisfo-`basename $$d` $$d; else echo "skip $$d (no Dockerfile)"; fi; done
