setup:
	@docker compose --env-file .env.example pull

run-local:
	@docker compose --env-file .env.example up -d

down:
	@docker compose --env-file .env.example down

test:
	@go test ./...

test-coverage:
	@go test ./... -coverprofile=coverage.out

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else echo "golangci-lint not installed"; fi

migrate-up:
	@echo "migrate-up not implemented yet"

migrate-down:
	@echo "migrate-down not implemented yet"

docker-build-all:
	@for d in services/*; do if [ -f $$d/Dockerfile ]; then docker build -t sisfo-`basename $$d` $$d; else echo "skip $$d (no Dockerfile)"; fi; done
