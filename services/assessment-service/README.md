# Assessment Service

Assessment Service is a microservice responsible for managing academic assessments, grading, and report card generation in the Sisfo Akademik system.

## Features

- **Grade Categories**: Manage categories for grades (e.g., Assignment, Quiz, Midterm, Final).
- **Assessments**: Create and manage assessments for subjects/classes.
- **Grading**: Input, update, and approve grades for students.
- **Report Cards**: Generate report cards (PDF) based on student grades.
- **Templates**: Customize report card templates.

## Architecture

This service follows Clean Architecture principles:
- **Domain**: Entities and business logic interfaces.
- **Usecase**: Application business rules.
- **Repository**: Data access layer (PostgreSQL).
- **Handler**: HTTP handlers (Gin).

## API Endpoints

### Grade Categories
- `POST /api/v1/grade-categories`: Create a new grade category.
- `GET /api/v1/grade-categories`: List grade categories.
- `PUT /api/v1/grade-categories/:id`: Update a grade category.
- `DELETE /api/v1/grade-categories/:id`: Delete a grade category.

### Assessments
- `POST /api/v1/assessments`: Create a new assessment.
- `GET /api/v1/assessments`: List assessments (filters supported).

### Grades
- `POST /api/v1/grades`: Input a grade for a student.
- `GET /api/v1/grades/student/:student_id`: Get grades for a student.
- `PUT /api/v1/grades/:id/approve`: Approve a grade.
- `GET /api/v1/grades/calculate`: Calculate final score.

### Report Cards
- `POST /api/v1/report-cards/generate`: Generate a report card.
- `GET /api/v1/report-cards/student/:student_id`: Get report card for a student.
- `GET /api/v1/report-cards/:id/download`: Download report card PDF.

### Templates
- `POST /api/v1/templates`: Create a report card template.
- `GET /api/v1/templates`: List templates.
- `GET /api/v1/templates/:id`: Get a template.
- `PUT /api/v1/templates/:id`: Update a template.
- `DELETE /api/v1/templates/:id`: Delete a template.

## Setup

### Prerequisites
- Go 1.20+
- PostgreSQL
- Docker (optional)

### Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8084` |
| `DB_URL` | Database connection string | `postgres://dev:dev@localhost:55432/devdb?sslmode=disable` |
| `STORAGE_PATH` | Local storage path for PDFs | `./storage` |

### Running Locally
```bash
# Install dependencies
go mod tidy

# Run the service
go run cmd/server/main.go
```

### Running Tests
```bash
# Run unit tests
go test ./...

# Run integration tests
go test ./tests/integration/...
```

## Database Migrations
Migrations are managed using `golang-migrate` or the `Makefile` utilities.

```bash
# Run migrations (up)
make migrate-up
```
