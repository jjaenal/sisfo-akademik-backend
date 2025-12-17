# Academic Service

Core service for managing academic operations in the Sisfo Akademik system. This service is the heart of the platform, handling the structure of schools, academic calendars, curriculum, classes, enrollments, and scheduling.

## Features

- **School Management**: Multi-tenant support for managing school profiles.
- **Academic Calendar**: Manage academic years and semesters.
- **Curriculum Management**: Define curricula, subjects, and grading rules.
- **Class Management**: Create classes, assign homeroom teachers.
- **Enrollment**: Manage student enrollments in classes (including bulk enrollment).
- **Schedule Management**: 
  - Create, update, delete schedules.
  - Conflict detection (Room, Teacher, Class).
  - Bulk creation.
  - Template-based schedule generation.
- **People Management**: Basic management for Students and Teachers (linked to Auth service).

## Architecture

Follows Clean Architecture principles:
- **Domain**: Entities and repository interfaces (Pure Go, no external deps).
- **UseCase**: Business logic and orchestration.
- **Handler**: HTTP handlers (Gin).
- **Repository**: Data access implementation (PostgreSQL).

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL (pgx driver)
- **Caching**: Redis
- **Messaging**: RabbitMQ (for event publication)
- **Documentation**: Swagger/OpenAPI

## Setup

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (for dependencies)

### Configuration

Copy `.env.example` (from root or service dir) and configure environment variables:

```env
APP_ENV=development
HTTP_PORT=8081
POSTGRES_URL=postgres://user:password@localhost:5432/academic_db?sslmode=disable
REDIS_ADDR=localhost:6379
RABBIT_URL=amqp://guest:guest@localhost:5672/
```

### Running Locally

1. **Start Dependencies**:
   ```bash
   make run-local # from project root
   ```

2. **Run Migrations**:
   ```bash
   migrate -path migrations -database "postgres://user:password@localhost:5432/academic_db?sslmode=disable" up
   ```

3. **Run Service**:
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

### Schools
- `POST /api/v1/schools` - Create a new school
- `GET /api/v1/schools/:id` - Get school details
- `GET /api/v1/schools/tenant/:tenant_id` - Get school by tenant ID
- `PUT /api/v1/schools/:id` - Update school details
- `DELETE /api/v1/schools/:id` - Delete a school

### Academic Years
- `POST /api/v1/academic-years` - Create academic year
- `GET /api/v1/academic-years/:id` - Get details
- `GET /api/v1/academic-years/tenant/:tenant_id` - List by tenant
- `PUT /api/v1/academic-years/:id` - Update
- `DELETE /api/v1/academic-years/:id` - Delete

### Semesters
- `POST /api/v1/semesters` - Create semester
- `GET /api/v1/semesters/:id` - Get details
- `GET /api/v1/semesters` - List semesters
- `PUT /api/v1/semesters/:id` - Update
- `PATCH /api/v1/semesters/:id/activate` - Set as active semester
- `DELETE /api/v1/semesters/:id` - Delete

### Classes
- `POST /api/v1/classes` - Create class
- `GET /api/v1/classes/:id` - Get details
- `GET /api/v1/classes` - List classes
- `PUT /api/v1/classes/:id` - Update
- `DELETE /api/v1/classes/:id` - Delete
- `POST /api/v1/classes/:id/subjects` - Add subject to class
- `GET /api/v1/classes/:id/subjects` - List class subjects
- `DELETE /api/v1/classes/:id/subjects/:subject_id` - Remove subject
- `POST /api/v1/classes/:id/subjects/:subject_id/teacher` - Assign teacher to subject

### Enrollments
- `POST /api/v1/enrollments` - Enroll a student
- `GET /api/v1/enrollments/:id` - Get enrollment details
- `PUT /api/v1/enrollments/:id/status` - Update status
- `DELETE /api/v1/enrollments/:id` - Unenroll
- `GET /api/v1/classes/:id/students` - List students in a class
- `POST /api/v1/classes/:id/students/bulk` - Bulk enroll students (CSV)
- `GET /api/v1/students/:id/classes` - List classes for a student

### Schedules
- `POST /api/v1/schedules` - Create schedule
- `POST /api/v1/schedules/bulk` - Bulk create
- `POST /api/v1/schedules/from-template` - Generate from template
- `GET /api/v1/schedules` - List schedules
- `GET /api/v1/schedules/:id` - Get details
- `GET /api/v1/schedules/class/:class_id` - Get class schedule
- `PUT /api/v1/schedules/:id` - Update
- `DELETE /api/v1/schedules/:id` - Delete

### Schedule Templates
- `POST /api/v1/schedule-templates` - Create template
- `GET /api/v1/schedule-templates` - List templates
- `GET /api/v1/schedule-templates/:id` - Get template details
- `POST /api/v1/schedule-templates/:id/items` - Add item to template
- `DELETE /api/v1/schedule-templates/items/:item_id` - Remove item

### Curricula
- `POST /api/v1/curricula` - Create curriculum
- `GET /api/v1/curricula` - List curricula
- `GET /api/v1/curricula/:id` - Get details
- `PUT /api/v1/curricula/:id` - Update
- `DELETE /api/v1/curricula/:id` - Delete
- `POST /api/v1/curricula/:id/subjects` - Add subject to curriculum
- `POST /api/v1/curricula/:id/grading-rules` - Add grading rule

### People (Students & Teachers)
- `POST /api/v1/students` - Create student profile
- `GET /api/v1/students` - List students
- `POST /api/v1/teachers` - Create teacher profile
- `GET /api/v1/teachers` - List teachers

## API Documentation

This service uses Swagger for API documentation.

1. **Generate Docs**:
   ```bash
   swag init -g cmd/server/main.go --parseDependency --parseInternal
   ```

2. **Access Docs**:
   Run the service and visit: `http://localhost:8081/swagger/index.html`

## Testing

Run integration tests:
```bash
go test -v ./tests/integration/...
```

Run unit tests:
```bash
go test -v ./internal/usecase/... ./internal/handler/...
```
