# Finance Service

Finance Service handles all financial operations within the Academic System, including billing configuration, invoice generation, payment processing, and financial reporting.

## Features

- **Billing Configuration**: Manage tuition fees and other charges.
- **Invoice Management**: Generate invoices automatically or manually.
- **Payment Processing**: Record payments and track transaction history.
- **Financial Reporting**: Generate revenue reports and track outstanding invoices.
- **Student Integration**: Sync student data for billing purposes.

## Architecture

This service follows Clean Architecture principles:

- `cmd/`: Entry point and server configuration.
- `internal/domain/`: Entities and interfaces (business core).
- `internal/usecase/`: Application business logic.
- `internal/handler/`: HTTP handlers (Gin).
- `internal/repository/`: Data access layer (PostgreSQL).

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Documentation**: Swagger/OpenAPI

## Setup

1. **Prerequisites**:
   - Go 1.21+
   - PostgreSQL
   - Docker (optional)

2. **Environment Variables**:
   Copy `.env.example` to `.env`:
   ```env
   APP_ENV=development
   APP_HTTP_PORT=9096
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=finance_service
   ```

3. **Run Service**:
   ```bash
   go run cmd/server/main.go
   ```

4. **Run Tests**:
   ```bash
   go test ./...
   ```

## API Documentation

Swagger documentation is available at:
`http://localhost:9096/swagger/index.html`

### Key Endpoints

- `POST /api/v1/finance/billing-configs`: Create billing configuration
- `POST /api/v1/finance/invoices/generate`: Generate invoice
- `POST /api/v1/finance/payments`: Record payment
- `GET /api/v1/finance/reports/revenue/monthly`: Get monthly revenue

## Database Schema

- `billing_configs`: Stores fee configurations.
- `invoices`: Records generated invoices for students.
- `payments`: Tracks payments made against invoices.
- `students`: Read-only replica/cache of student data for integrity.
