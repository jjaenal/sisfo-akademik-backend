# Academic Service

Core service for managing academic operations in the Sisfo Akademik system. Handles curriculum, classes, subjects, teachers, students, and schedules.

## Features

- **Curriculum Management**: Create and manage curriculums.
- **Class Management**: Manage classes and assign students/teachers.
- **Subject Management**: Manage subjects offered.
- **Student & Teacher Management**: Basic profiles and assignments.
- **Schedule Management**: 
  - Create, update, delete schedules.
  - Conflict detection.
  - Bulk creation.
  - Template-based schedule generation.

## Architecture

Follows Clean Architecture principles:
- **Domain**: Entities and repository interfaces.
- **UseCase**: Business logic.
- **Handler**: HTTP handlers (Gin).
- **Repository**: PostgreSQL implementation.

## Tech Stack

- Go 1.21+
- Gin Web Framework
- PostgreSQL (pgx)
- Golang-Migrate

## Setup

1. **Environment Variables**:
   Copy `.env.example` to `.env` and configure:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=academic_db
   ```

2. **Run Migrations**:
   ```bash
   make migrate-up
   ```

3. **Run Service**:
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

### Schedules
- `POST /api/v1/schedules`: Create a schedule
- `POST /api/v1/schedules/bulk`: Bulk create schedules
- `POST /api/v1/schedules/from-template`: Create schedules from a template
- `GET /api/v1/schedules`: List schedules (filter by class_id, etc.)
- `GET /api/v1/schedules/:id`: Get schedule details
- `PUT /api/v1/schedules/:id`: Update a schedule
- `DELETE /api/v1/schedules/:id`: Delete a schedule

## API Documentation

This service uses Swagger for API documentation.

1. **Generate Docs**:
   ```bash
   swag init -g cmd/server/main.go --parseDependency --parseInternal
   ```

2. **Access Docs**:
   Run the service and visit:
   `http://localhost:8081/swagger/index.html`

## Testing

Run all tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test -cover ./...
```
