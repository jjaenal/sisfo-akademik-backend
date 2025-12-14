# TODO — Backend Sistem Akademik

Checklist pekerjaan backend yang detail, terorganisir, dan actionable. Setiap task memiliki priority, estimation, dan acceptance criteria.

**Legend**:
- 🔴 High Priority (Critical Path)
- 🟡 Medium Priority (Important)
- 🟢 Low Priority (Nice to Have)
- ⏱️ Estimated Hours
- ✅ Completed
- 🔄 In Progress
- ⏳ Blocked
- 📝 Not Started

---

## Global Infrastructure

### Repository & Project Setup
- [ ] 📝 Create monorepo structure (4h) 🔴
  - Setup services/, shared/, infrastructure/, docs/ folders
  - Configure Go workspace
  - Setup .gitignore for Go projects
  - **AC**: Directory structure matches standard layout

- [ ] 📝 Setup GitHub/GitLab organization (2h) 🔴
  - Create organization/group
  - Setup team permissions
  - Configure branch protection (main, develop)
  - **AC**: Team members have appropriate access

- [ ] 📝 Initialize project management board (2h) 🟡
  - Create Jira project or GitHub Projects
  - Define issue templates
  - Setup workflows (Backlog → In Progress → Review → Done)
  - **AC**: Board accessible and organized

### Development Environment

- [ ] 📝 Create docker-compose.yml for local dev (6h) 🔴
  - PostgreSQL container
  - Redis container
  - RabbitMQ container
  - Jaeger container
  - Service containers (placeholder)
  - **AC**: `docker-compose up` starts all services

- [ ] 📝 Create Makefile with common commands (3h) 🔴
  ```makefile
  Commands needed:
  - make setup (initial setup)
  - make run-local (run all services)
  - make test (run tests)
  - make test-coverage (coverage report)
  - make lint (code linting)
  - make migrate-up (run migrations)
  - make migrate-down (rollback migrations)
  - make docker-build-all (build all images)
  - make clean (cleanup)
  ```
  - **AC**: All commands working correctly

- [ ] 📝 Setup environment files (2h) 🔴
  - Create .env.example
  - Document all required variables
  - Create .env.local for local dev
  - **AC**: Services start with provided env vars

- [ ] 📝 Write local setup documentation (3h) 🟡
  - README.md in root
  - Prerequisites
  - Installation steps
  - Troubleshooting guide
  - **AC**: New developer can setup in <30 minutes

### Shared Packages

- [ ] 📝 Implement config package (4h) 🔴
  - Viper integration
  - Environment variable loading
  - Config validation
  - Hot reload support
  - **AC**: Config loaded correctly with validation

- [ ] 📝 Implement database package (6h) 🔴
  - PostgreSQL connection
  - Connection pooling (pgxpool)
  - Health check
  - Transaction helper
  - **AC**: DB connection stable with pooling

- [ ] 📝 Implement Redis package (4h) 🔴
  - Redis client
  - Connection pooling
  - Health check
  - Helper functions (Get, Set, Delete)
  - **AC**: Redis operations working

- [ ] 📝 Implement logger package (5h) 🔴
  - Structured logging (zap/logrus)
  - Log levels
  - Context-aware logging
  - JSON output
  - **AC**: Logs formatted correctly

- [ ] 📝 Implement middleware package (8h) 🔴
  - Authentication middleware
  - Authorization middleware
  - Logging middleware
  - CORS middleware
  - Rate limiting middleware
  - Error handling middleware
  - **AC**: All middleware functional

- [ ] 📝 Implement errors package (3h) 🟡
  - Custom error types
  - Error codes
  - Error wrapping
  - HTTP error responses
  - **AC**: Consistent error handling

- [ ] 📝 Implement validator package (4h) 🟡
  - go-playground/validator wrapper
  - Custom validators
  - Validation error formatting
  - **AC**: Input validation working

- [ ] 📝 Implement JWT package (5h) 🔴
  - JWT generation
  - JWT validation
  - Token refresh
  - Claims extraction
  - **AC**: JWT operations secure & functional

- [ ] 📝 Implement httputil package (3h) 🟡
  - Response helpers
  - Request parsing
  - Pagination helpers
  - **AC**: HTTP utilities working

- [ ] 📝 Implement testutil package (4h) 🟡
  - Test database helpers
  - Mock helpers
  - Test fixtures
  - **AC**: Tests easier to write

### CI/CD Pipeline

- [ ] 📝 Setup GitHub Actions workflow (6h) 🔴
  - Lint job
  - Test job
  - Security scan job
  - Build job
  - Deploy staging job
  - **AC**: Pipeline runs on push

- [ ] 📝 Configure linting (golangci-lint) (2h) 🔴
  - Install golangci-lint
  - Configure .golangci.yml
  - Add to CI
  - **AC**: Code passes linting

- [ ] 📝 Setup test coverage reporting (3h) 🟡
  - Integrate Codecov
  - Coverage badge
  - Coverage threshold (70%)
  - **AC**: Coverage tracked in CI

- [ ] 📝 Setup security scanning (4h) 🔴
  - Gosec for static analysis
  - Trivy for container scanning
  - Snyk for dependencies
  - **AC**: Security scans in CI

- [ ] 📝 Configure Docker registry (2h) 🔴
  - GitHub Container Registry or Docker Hub
  - Setup credentials
  - Image tagging strategy
  - **AC**: Images pushed to registry

### Observability

- [ ] 📝 Setup ELK Stack (8h) 🟡
  - Elasticsearch container
  - Logstash container
  - Kibana container
  - Configure log shipping
  - **AC**: Logs viewable in Kibana

- [ ] 📝 Setup Prometheus (6h) 🟡
  - Prometheus container
  - Configure scraping
  - Define metrics
  - **AC**: Metrics scraped

- [ ] 📝 Setup Grafana (6h) 🟡
  - Grafana container
  - Connect to Prometheus
  - Create dashboards
  - **AC**: Dashboards showing metrics

- [ ] 📝 Setup Jaeger (4h) 🟡
  - Jaeger container
  - Trace instrumentation
  - **AC**: Traces visible in Jaeger

---

## Auth / Identity Service

### Core Setup

- [ ] 📝 Create auth-service structure (3h) 🔴
  - Initialize Go module
  - Setup directory structure
  - Create Dockerfile
  - Create docker-compose.yml
  - **AC**: Service structure ready

- [ ] 📝 Setup database connection (2h) 🔴
  - Use shared database package
  - Test connection
  - **AC**: Auth service connects to DB

- [ ] 📝 Create database migrations (4h) 🔴
  - Users table
  - Roles table
  - Permissions table
  - Role_permissions table
  - User_roles table
  - Audit_logs table
  - **AC**: Migrations run successfully

### User Management

- [ ] 📝 Implement User entity (2h) 🔴
  - Define User struct
  - Validation rules
  - Methods (BeforeCreate, etc)
  - **AC**: User entity complete

- [ ] 📝 Implement User repository (6h) 🔴
  - Create()
  - GetByID()
  - GetByEmail()
  - List() with pagination
  - Update()
  - Delete() (soft delete)
  - **AC**: CRUD operations working

- [ ] 📝 Implement User use case (6h) 🔴
  - Register user
  - Get user
  - Update user
  - Delete user
  - Search users
  - **AC**: Business logic implemented

- [ ] 📝 Implement User handlers (8h) 🔴
  - POST /api/v1/users
  - GET /api/v1/users/:id
  - GET /api/v1/users
  - PUT /api/v1/users/:id
  - DELETE /api/v1/users/:id
  - PATCH /api/v1/users/:id/activate
  - **AC**: All endpoints working

- [ ] 📝 Add input validation (3h) 🔴
  - Email format
  - Password strength
  - Required fields
  - **AC**: Invalid input rejected

- [ ] 📝 Implement password hashing (2h) 🔴
  - bcrypt implementation
  - Cost factor: 12
  - **AC**: Passwords hashed securely

- [ ] 📝 Unit tests for User (8h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing with good coverage

### Authentication

- [ ] 📝 Implement login handler (6h) 🔴
  - POST /api/v1/auth/login
  - Validate credentials
  - Generate tokens
  - Return access & refresh tokens
  - **AC**: Login working

- [ ] 📝 Implement JWT generation (4h) 🔴
  - Access token (15 min TTL)
  - Refresh token (7 days TTL)
  - Include claims (user_id, tenant_id, roles, permissions)
  - **AC**: JWT generated correctly

- [ ] 📝 Implement logout handler (3h) 🔴
  - POST /api/v1/auth/logout
  - Invalidate refresh token
  - Blacklist access token in Redis
  - **AC**: Logout working

- [ ] 📝 Implement token refresh (5h) 🔴
  - POST /api/v1/auth/refresh
  - Validate refresh token
  - Generate new access token
  - Rotate refresh token
  - **AC**: Token refresh working

- [ ] 📝 Implement forgot password (6h) 🟡
  - POST /api/v1/auth/forgot-password
  - Generate reset token
  - Send reset email
  - **AC**: Reset email sent

- [ ] 📝 Implement reset password (5h) 🟡
  - POST /api/v1/auth/reset-password
  - Validate reset token
  - Update password
  - **AC**: Password reset working

- [ ] 📝 Implement change password (4h) 🟡
  - POST /api/v1/auth/change-password
  - Validate old password
  - Update password
  - **AC**: Password change working

- [ ] 📝 Implement failed login tracking (4h) 🟡
  - Track failed attempts
  - Lock account after 5 failures
  - Auto-unlock after 30 minutes
  - **AC**: Account lockout working

- [ ] 📝 Unit tests for Auth (10h) 🔴
  - Login tests
  - Token generation tests
  - Token refresh tests
  - Logout tests
  - Coverage >70%
  - **AC**: Tests passing

### RBAC (Role-Based Access Control)

- [ ] 📝 Implement Role entity (2h) 🔴
  - Define Role struct
  - Validation rules
  - **AC**: Role entity complete

- [ ] 📝 Implement Permission entity (2h) 🔴
  - Define Permission struct
  - Resource:Action format
  - **AC**: Permission entity complete

- [ ] 📝 Implement Role repository (6h) 🔴
  - Create()
  - GetByID()
  - List()
  - Update()
  - Delete()
  - AssignPermissions()
  - **AC**: Role CRUD working

- [ ] 📝 Implement Permission repository (4h) 🔴
  - Create()
  - List()
  - GetByRole()
  - **AC**: Permission operations working

- [ ] 📝 Seed default roles & permissions (4h) 🔴
  - Super Admin role
  - School Admin role
  - Teacher role
  - Student role
  - Parent role
  - Default permissions
  - **AC**: Default roles created

- [ ] 📝 Implement role assignment (5h) 🔴
  - Assign role to user
  - Remove role from user
  - Get user roles
  - Get user permissions (effective)
  - **AC**: Role assignment working

- [ ] 📝 Implement RBAC middleware (8h) 🔴
  - Check user authentication
  - Check user permissions
  - Context-aware (tenant-based)
  - **AC**: Protected endpoints secured

- [ ] 📝 Implement role handlers (8h) 🟡
  - POST /api/v1/roles
  - GET /api/v1/roles
  - GET /api/v1/roles/:id
  - PUT /api/v1/roles/:id
  - DELETE /api/v1/roles/:id
  - **AC**: Role management working

- [ ] 📝 Unit tests for RBAC (10h) 🔴
  - Role tests
  - Permission tests
  - Middleware tests
  - Coverage >70%
  - **AC**: Tests passing

### Audit Logging

- [ ] 📝 Implement audit log entity (2h) 🔴
  - Define AuditLog struct
  - Fields: user, action, resource, changes
  - **AC**: Audit log entity complete

- [ ] 📝 Implement audit log repository (4h) 🔴
  - Create()
  - List() with filters
  - Search()
  - **AC**: Audit logs stored

- [ ] 📝 Implement audit middleware (6h) 🔴
  - Capture request details
  - Log after response
  - Async logging (don't block)
  - **AC**: All actions logged

- [ ] 📝 Implement audit log handlers (4h) 🟡
  - GET /api/v1/audit-logs
  - GET /api/v1/audit-logs/search
  - Export audit logs
  - **AC**: Audit logs viewable

- [ ] 📝 Setup log retention (2h) 🟡
  - 90-day retention
  - Automated cleanup job
  - **AC**: Old logs cleaned up

### Security Enhancements

- [ ] 📝 Implement rate limiting (6h) 🔴
  - Redis-based rate limiter
  - Different limits per endpoint type
  - Rate limit headers
  - **AC**: Rate limiting active

- [ ] 📝 Implement security headers (3h) 🔴
  - X-Content-Type-Options
  - X-Frame-Options
  - X-XSS-Protection
  - Strict-Transport-Security
  - CSP
  - **AC**: Security headers present

- [ ] 📝 Configure CORS (2h) 🔴
  - Whitelist origins
  - Allowed methods & headers
  - **AC**: CORS working

- [ ] 📝 Implement password validation (3h) 🟡
  - Minimum 8 characters
  - Complexity requirements
  - Common password check
  - **AC**: Weak passwords rejected

- [ ] 📝 Implement password history (3h) 🟡
  - Track last 5 passwords
  - Prevent reuse
  - **AC**: Password reuse prevented

### Integration Tests

- [ ] 📝 Auth service integration tests (12h) 🔴
  - User registration flow
  - Login flow
  - Token refresh flow
  - RBAC flow
  - **AC**: Integration tests passing

### Documentation

- [ ] 📝 Write Auth service README (3h) 🟡
  - Service overview
  - Setup instructions
  - API endpoints
  - **AC**: README complete

- [ ] 📝 Generate Swagger docs (4h) 🟡
  - Annotate handlers
  - Generate OpenAPI spec
  - Setup Swagger UI
  - **AC**: API docs accessible

- [ ] 📝 Create Postman collection (3h) 🟡
  - All endpoints
  - Example requests
  - Environment variables
  - **AC**: Postman collection works

---

## Academic Core Service

### Core Setup

- [ ] 📝 Create academic-service structure (3h) 🔴
  - Initialize Go module
  - Setup directory structure
  - Create Dockerfile
  - **AC**: Service structure ready

- [ ] 📝 Setup database connection (2h) 🔴
  - Use shared database package
  - Test connection
  - **AC**: Service connects to DB

- [ ] 📝 Create database migrations (8h) 🔴
  - Schools table
  - Academic_years table
  - Semesters table
  - Classes table
  - Class_students table
  - Subjects table
  - Class_subjects table
  - Curricula table
  - Curriculum_subjects table
  - Grading_rules table
  - Schedules table
  - **AC**: Migrations run successfully

### School Management

- [ ] 📝 Implement School entity (2h) 🔴
  - Define School struct
  - Validation rules
  - **AC**: School entity complete

- [ ] 📝 Implement School repository (6h) 🔴
  - CRUD operations
  - GetByTenantID()
  - **AC**: School CRUD working

- [ ] 📝 Implement School use case (4h) 🔴
  - Business logic
  - Validation
  - **AC**: School operations working

- [ ] 📝 Implement School handlers (8h) 🔴
  - POST /api/v1/schools
  - GET /api/v1/schools
  - GET /api/v1/schools/:id
  - PUT /api/v1/schools/:id
  - DELETE /api/v1/schools/:id
  - **AC**: All endpoints working

- [ ] 📝 Implement school logo upload (4h) 🟡
  - File upload endpoint
  - Image validation
  - Store in object storage
  - **AC**: Logo upload working

- [ ] 📝 Unit tests for School (8h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Academic Year & Semester

- [ ] 📝 Implement AcademicYear entity (2h) 🔴
  - Define AcademicYear struct
  - Validation (dates, active flag)
  - **AC**: AcademicYear entity complete

- [ ] 📝 Implement Semester entity (2h) 🔴
  - Define Semester struct
  - Validation
  - **AC**: Semester entity complete

- [ ] 📝 Implement AcademicYear repository (6h) 🔴
  - CRUD operations
  - GetActive()
  - ValidateNonOverlap()
  - **AC**: AcademicYear CRUD working

- [ ] 📝 Implement Semester repository (5h) 🔴
  - CRUD operations
  - GetBySemester()
  - GetActive()
  - **AC**: Semester CRUD working

- [ ] 📝 Implement academic year handlers (8h) 🔴
  - POST /api/v1/academic-years
  - GET /api/v1/academic-years
  - GET /api/v1/academic-years/:id
  - PUT /api/v1/academic-years/:id
  - PATCH /api/v1/academic-years/:id/activate
  - DELETE /api/v1/academic-years/:id
  - **AC**: All endpoints working

- [ ] 📝 Implement semester handlers (8h) 🔴
  - POST /api/v1/semesters
  - GET /api/v1/semesters
  - GET /api/v1/semesters/:id
  - PUT /api/v1/semesters/:id
  - PATCH /api/v1/semesters/:id/activate
  - **AC**: All endpoints working

- [ ] 📝 Implement active year/semester validation (3h) 🔴
  - Only 1 active year per tenant
  - Only 1 active semester per year
  - **AC**: Validation working

- [ ] 📝 Unit tests for AcademicYear & Semester (10h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Class Management

- [ ] 📝 Implement Class entity (3h) 🔴
  - Define Class struct
  - Validation rules
  - **AC**: Class entity complete

- [ ] 📝 Implement ClassStudent entity (2h) 🔴
  - Enrollment tracking
  - Status (active, transferred, graduated)
  - **AC**: ClassStudent entity complete

- [ ] 📝 Implement Class repository (8h) 🔴
  - CRUD operations
  - GetByAcademicYear()
  - GetStudents()
  - EnrollStudent()
  - **AC**: Class operations working

- [ ] 📝 Implement class handlers (10h) 🔴
  - POST /api/v1/classes
  - GET /api/v1/classes
  - GET /api/v1/classes/:id
  - PUT /api/v1/classes/:id
  - DELETE /api/v1/classes/:id
  - POST /api/v1/classes/:id/students
  - GET /api/v1/classes/:id/students
  - DELETE /api/v1/classes/:id/students/:student_id
  - **AC**: All endpoints working

- [ ] 📝 Implement bulk enrollment (5h) 🟡
  - POST /api/v1/classes/:id/students/bulk
  - CSV import
  - Validation
  - **AC**: Bulk enrollment working

- [ ] 📝 Implement capacity management (3h) 🟡
  - Check max_students
  - Prevent over-enrollment
  - **AC**: Capacity enforced

- [ ] 📝 Unit tests for Class (10h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Subject Management

- [ ] 📝 Implement Subject entity (2h) 🔴
  - Define Subject struct
  - Categories (Wajib, Peminatan, Mulok)
  - **AC**: Subject entity complete

- [ ] 📝 Implement ClassSubject entity (2h) 🔴
  - Subject-class-teacher mapping
  - **AC**: ClassSubject entity complete

- [ ] 📝 Implement Subject repository (6h) 🔴
  - CRUD operations
  - GetByCategory()
  - AssignToClass()
  - **AC**: Subject operations working

- [ ] 📝 Implement subject handlers (10h) 🔴
  - POST /api/v1/subjects
  - GET /api/v1/subjects
  - GET /api/v1/subjects/:id
  - PUT /api/v1/subjects/:id
  - DELETE /api/v1/subjects/:id
  - POST /api/v1/classes/:id/subjects
  - GET /api/v1/classes/:id/subjects
  - DELETE /api/v1/classes/:id/subjects/:subject_id
  - **AC**: All endpoints working

- [ ] 📝 Implement teacher assignment (4h) 🔴
  - PUT /api/v1/classes/:id/subjects/:subject_id/teacher
  - Validation
  - **AC**: Teacher assignment working

- [ ] 📝 Unit tests for Subject (8h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Curriculum Management

- [ ] 📝 Implement Curriculum entity (3h) 🔴
  - Define Curriculum struct
  - Support multiple curricula per tenant
  - **AC**: Curriculum entity complete

- [ ] 📝 Implement GradingRule entity (3h) 🔴
  - KKM configuration
  - Grade components & weights
  - **AC**: GradingRule entity complete

- [ ] 📝 Implement Curriculum repository (6h) 🔴
  - CRUD operations
  - GetSubjects()
  - GetGradingRules()
  - **AC**: Curriculum operations working

- [ ] 📝 Implement curriculum handlers (10h) 🟡
  - POST /api/v1/curricula
  - GET /api/v1/curricula
  - GET /api/v1/curricula/:id
  - PUT /api/v1/curricula/:id
  - POST /api/v1/curricula/:id/subjects
  - GET /api/v1/curricula/:id/subjects
  - POST /api/v1/curricula/:id/grading-rules
  - **AC**: All endpoints working

- [ ] 📝 Unit tests for Curriculum (8h) 🟡
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Schedule Management

- [ ] 📝 Implement Schedule entity (3h) 🔴
  - Define Schedule struct
  - Day of week, time slots
  - **AC**: Schedule entity complete

- [ ] 📝 Implement Schedule repository (6h) 🔴
  - CRUD operations
  - GetWeeklySchedule()
  - CheckConflicts()
  - **AC**: Schedule operations working

- [ ] 📝 Implement conflict detection (6h) 🔴
  - Class conflict check
  - Teacher conflict check
  - Room conflict check
  - **AC**: Conflicts detected

- [ ] 📝 Implement schedule handlers (10h) 🔴
  - POST /api/v1/schedules
  - GET /api/v1/schedules
  - PUT /api/v1/schedules/:id
  - DELETE /api/v1/schedules/:id
  - GET /api/v1/schedules/class/:class_id/weekly
  - GET /api/v1/schedules/teacher/:teacher_id/weekly
  - **AC**: All endpoints working

- [ ] 📝 Implement bulk schedule creation (5h) 🟡
  - Template system
  - Batch creation
  - **AC**: Bulk creation working

- [ ] 📝 Unit tests for Schedule (8h) 🔴
  - Repository tests
  - Conflict tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [ ] 📝 Academic service integration tests (12h) 🔴
  - School creation flow
  - Academic year setup flow
  - Class & student enrollment flow
  - Schedule creation flow
  - **AC**: Integration tests passing

### Documentation

- [ ] 📝 Write Academic service README (3h) 🟡
  - Service overview
  - Setup instructions
  - API endpoints
  - **AC**: README complete

- [ ] 📝 Generate Swagger docs (4h) 🟡
  - Annotate handlers
  - Generate spec
  - **AC**: API docs accessible

- [ ] 📝 Create Postman collection (3h) 🟡
  - All endpoints
  - Example requests
  - **AC**: Postman collection works

---

## Attendance Service

### Core Setup

- [ ] 📝 Create attendance-service structure (3h) 🔴
  - Initialize Go module
  - Setup directory structure
  - **AC**: Service structure ready

- [ ] 📝 Setup database connection (2h) 🔴
  - Use shared database package
  - **AC**: Service connects to DB

- [ ] 📝 Create database migrations (4h) 🔴
  - Student_attendance table
  - Teacher_attendance table
  - **AC**: Migrations run

### Student Attendance

- [ ] 📝 Implement StudentAttendance entity (2h) 🔴
  - Define struct
  - Status types (present, absent, late, excused, sick)
  - **AC**: Entity complete

- [ ] 📝 Implement StudentAttendance repository (6h) 🔴
  - Create()
  - Update()
  - GetByStudentAndDate()
  - List() with filters
  - GetSummary()
  - **AC**: CRUD working

- [ ] 📝 Implement attendance handlers (10h) 🔴
  - POST /api/v1/attendance/students
  - POST /api/v1/attendance/students/bulk
  - GET /api/v1/attendance/students
  - PUT /api/v1/attendance/students/:id
  - GET /api/v1/attendance/students/:student_id/summary
  - **AC**: All endpoints working

- [ ] 📝 Implement bulk check-in (5h) 🔴
  - Full class check-in
  - Validation
  - **AC**: Bulk check-in working

- [ ] 📝 Implement GPS validation (4h) 🟡
  - Validate location against school location
  - Distance calculation
  - **AC**: GPS validation working

- [ ] 📝 Unit tests for StudentAttendance (8h) 🔴
  - Repository tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Teacher Attendance

- [ ] 📝 Implement TeacherAttendance entity (2h) 🔴
  - Define struct
  - Check-in/check-out times
  - **AC**: Entity complete

- [ ] 📝 Implement TeacherAttendance repository (5h) 🔴
  - Create()
  - Update()
  - GetByTeacherAndDate()
  - List()
  - **AC**: CRUD working

- [ ] 📝 Implement teacher attendance handlers (8h) 🔴
  - POST /api/v1/attendance/teachers/check-in
  - POST /api/v1/attendance/teachers/check-out
  - GET /api/v1/attendance/teachers
  - GET /api/v1/attendance/teachers/:teacher_id/summary
  - **AC**: All endpoints working

- [ ] 📝 Unit tests for TeacherAttendance (6h) 🔴
  - Repository tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Reports

- [ ] 📝 Implement attendance reports (8h) 🟡
  - GET /api/v1/attendance/reports/daily
  - GET /api/v1/attendance/reports/monthly
  - GET /api/v1/attendance/reports/class/:class_id
  - **AC**: Reports working

### Integration Tests

- [ ] 📝 Attendance service integration tests (8h) 🔴
  - Student attendance flow
  - Bulk check-in flow
  - Teacher attendance flow
  - **AC**: Integration tests passing

### Documentation

- [ ] 📝 Write Attendance service README (2h) 🟡
- [ ] 📝 Generate Swagger docs (3h) 🟡
- [ ] 📝 Create Postman collection (2h) 🟡

---

## Assessment Service

### Core Setup

- [ ] 📝 Create assessment-service structure (3h) 🔴
- [ ] 📝 Setup database connection (2h) 🔴
- [ ] 📝 Create database migrations (6h) 🔴
  - Grade_categories table
  - Assessments table
  - Grades table
  - Report_cards table
  - Report_card_details table

### Grading System

- [ ] 📝 Implement GradeCategory entity (2h) 🔴
- [ ] 📝 Implement Assessment entity (3h) 🔴
- [ ] 📝 Implement Grade entity (3h) 🔴

- [ ] 📝 Implement grade repositories (8h) 🔴
  - GradeCategory CRUD
  - Assessment CRUD
  - Grade CRUD

- [ ] 📝 Implement grade calculation engine (8h) 🔴
  - Calculate weighted scores
  - Final score calculation
  - Grade letter assignment
  - KKM validation
  - **AC**: Grades calculated correctly

- [ ] 📝 Implement grade handlers (12h) 🔴
  - POST /api/v1/grade-categories
  - GET /api/v1/grade-categories
  - POST /api/v1/assessments
  - GET /api/v1/assessments
  - POST /api/v1/grades
  - POST /api/v1/grades/bulk
  - GET /api/v1/grades/assessment/:assessment_id
  - PUT /api/v1/grades/:id
  - PATCH /api/v1/grades/:id/approve
  - GET /api/v1/grades/student/:student_id/semester/:semester_id
  - **AC**: All endpoints working

- [ ] 📝 Implement grade approval workflow (5h) 🟡
  - Draft → Submitted → Approved
  - Audit trail
  - **AC**: Workflow working

- [ ] 📝 Unit tests for Grading (10h) 🔴
  - Calculation tests
  - Repository tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Report Card Generation

- [ ] 📝 Implement ReportCard entity (3h) 🔴
  - Define struct
  - Status (draft, generated, published)
  - **AC**: Entity complete

- [ ] 📝 Implement report card data aggregation (8h) 🔴
  - Collect all grades
  - Calculate final scores
  - Get attendance summary
  - **AC**: Data aggregated correctly

- [ ] 📝 Implement report card generation (12h) 🔴
  - POST /api/v1/report-cards/generate/:student_id/:semester_id
  - POST /api/v1/report-cards/generate/class/:class_id/:semester_id
  - Generate report data
  - **AC**: Report cards generated

- [ ] 📝 Implement PDF generation (12h) 🔴
  - HTML template
  - Convert to PDF (chromedp/gotenberg)
  - Store in object storage
  - **AC**: PDF generated correctly

- [ ] 📝 Implement report card handlers (8h) 🔴
  - GET /api/v1/report-cards/:id
  - GET /api/v1/report-cards/student/:student_id
  - PATCH /api/v1/report-cards/:id/publish
  - GET /api/v1/report-cards/:id/pdf
  - GET /api/v1/report-cards/:id/download
  - **AC**: All endpoints working

- [ ] 📝 Implement template customization (6h) 🟡
  - Template management
  - Variable replacement
  - **AC**: Templates customizable

- [ ] 📝 Unit tests for ReportCard (10h) 🔴
  - Generation tests
  - PDF tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [ ] 📝 Assessment service integration tests (12h) 🔴
  - Grading flow
  - Report card generation flow
  - **AC**: Integration tests passing

### Documentation

- [ ] 📝 Write Assessment service README (3h) 🟡
- [ ] 📝 Generate Swagger docs (4h) 🟡
- [ ] 📝 Create Postman collection (3h) 🟡

---

## Admission Service (PPDB)

### Core Setup

- [ ] 📝 Create admission-service structure (3h) 🔴
- [ ] 📝 Setup database connection (2h) 🔴
- [ ] 📝 Create database migrations (5h) 🔴
  - Admission_periods table
  - Applications table
  - Application_documents table

### Admission Management

- [ ] 📝 Implement AdmissionPeriod entity (2h) 🔴
- [ ] 📝 Implement Application entity (3h) 🔴
- [ ] 📝 Implement ApplicationDocument entity (2h) 🔴

- [ ] 📝 Implement admission repositories (8h) 🔴
  - AdmissionPeriod CRUD
  - Application CRUD
  - ApplicationDocument CRUD

- [ ] 📝 Implement admission period handlers (8h) 🔴
  - POST /api/v1/admission/periods
  - GET /api/v1/admission/periods
  - GET /api/v1/admission/periods/:id
  - PUT /api/v1/admission/periods/:id
  - PATCH /api/v1/admission/periods/:id/close
  - **AC**: All endpoints working

- [ ] 📝 Implement public application (10h) 🔴
  - GET /api/v1/admission/public/periods
  - POST /api/v1/admission/applications
  - GET /api/v1/admission/applications/:number/status
  - Application number generation
  - **AC**: Public application working

- [ ] 📝 Implement document upload (8h) 🔴
  - POST /api/v1/admission/applications/:id/documents
  - File validation (size, type)
  - Store in object storage
  - **AC**: Upload working

- [ ] 📝 Implement application management (10h) 🔴
  - GET /api/v1/admission/applications
  - GET /api/v1/admission/applications/:id
  - PUT /api/v1/admission/applications/:id
  - PATCH /api/v1/admission/applications/:id/verify
  - PATCH /api/v1/admission/applications/:id/accept
  - PATCH /api/v1/admission/applications/:id/reject
  - **AC**: Management working

- [ ] 📝 Implement selection process (10h) 🟡
  - POST /api/v1/admission/applications/:id/test-score
  - POST /api/v1/admission/applications/:id/interview-score
  - POST /api/v1/admission/periods/:id/calculate-final-scores
  - POST /api/v1/admission/periods/:id/announce
  - Final score calculation
  - **AC**: Selection working

- [ ] 📝 Implement student registration (8h) 🔴
  - POST /api/v1/admission/applications/:id/register
  - Create user account
  - Create student record
  - **AC**: Registration working

- [ ] 📝 Unit tests for Admission (12h) 🔴
  - Repository tests
  - Handler tests
  - Selection logic tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [ ] 📝 Admission service integration tests (10h) 🔴
  - Application submission flow
  - Document upload flow
  - Selection flow
  - **AC**: Integration tests passing

### Documentation

- [ ] 📝 Write Admission service README (3h) 🟡
- [ ] 📝 Generate Swagger docs (4h) 🟡
- [ ] 📝 Create Postman collection (3h) 🟡

---

## Finance Service

### Core Setup

- [ ] 📝 Create finance-service structure (3h) 🔴
- [ ] 📝 Setup database connection (2h) 🔴
- [ ] 📝 Create database migrations (4h) 🔴
  - Billing_configurations table
  - Invoices table
  - Payments table

### Finance Management

- [ ] 📝 Implement BillingConfig entity (2h) 🔴
- [ ] 📝 Implement Invoice entity (3h) 🔴
- [ ] 📝 Implement Payment entity (2h) 🔴

- [ ] 📝 Implement finance repositories (8h) 🔴
  - BillingConfig CRUD
  - Invoice CRUD
  - Payment CRUD

- [ ] 📝 Implement billing configuration (8h) 🔴
  - POST /api/v1/finance/billing-configs
  - GET /api/v1/finance/billing-configs
  - PUT /api/v1/finance/billing-configs/:id
  - **AC**: Billing config working

- [ ] 📝 Implement invoice generation (10h) 🔴
  - POST /api/v1/finance/invoices/generate
  - POST /api/v1/finance/invoices/generate/bulk
  - POST /api/v1/finance/invoices/generate/auto
  - Invoice number generation
  - **AC**: Invoice generation working

- [ ] 📝 Implement auto-generation (8h) 🟡
  - Scheduled job (cron)
  - Monthly SPP generation
  - **AC**: Auto-generation working

- [ ] 📝 Implement invoice handlers (8h) 🔴
  - GET /api/v1/finance/invoices
  - GET /api/v1/finance/invoices/:id
  - PUT /api/v1/finance/invoices/:id
  - GET /api/v1/finance/invoices/student/:student_id
  - GET /api/v1/finance/invoices/student/:student_id/outstanding
  - **AC**: Invoice management working

- [ ] 📝 Implement payment recording (8h) 🔴
  - POST /api/v1/finance/payments
  - GET /api/v1/finance/payments
  - GET /api/v1/finance/payments/:id
  - Payment number generation
  - Receipt generation
  - **AC**: Payment recording working

- [ ] 📝 Implement financial reports (10h) 🟡
  - GET /api/v1/finance/reports/revenue/daily
  - GET /api/v1/finance/reports/revenue/monthly
  - GET /api/v1/finance/reports/outstanding
  - GET /api/v1/finance/reports/student/:student_id/history
  - **AC**: Reports working

- [ ] 📝 Implement overdue tracking (5h) 🟡
  - Scheduled job
  - Mark overdue invoices
  - **AC**: Overdue tracking working

- [ ] 📝 Unit tests for Finance (12h) 🔴
  - Repository tests
  - Handler tests
  - Calculation tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [ ] 📝 Finance service integration tests (10h) 🔴
  - Invoice generation flow
  - Payment recording flow
  - **AC**: Integration tests passing

### Documentation

- [ ] 📝 Write Finance service README (3h) 🟡
- [ ] 📝 Generate Swagger docs (4h) 🟡
- [ ] 📝 Create Postman collection (3h) 🟡

---

## Notification Service

### Core Setup

- [ ] 📝 Create notification-service structure (3h) 🔴
- [ ] 📝 Setup database connection (2h) 🔴
- [ ] 📝 Create database migrations (3h) 🔴
  - Notification_templates table
  - Notifications table

### Notification Core

- [ ] 📝 Implement NotificationTemplate entity (2h) 🔴
- [ ] 📝 Implement Notification entity (2h) 🔴

- [ ] 📝 Implement notification repositories (6h) 🔴
  - Template CRUD
  - Notification CRUD

- [ ] 📝 Implement template management (8h) 🟡
  - POST /api/v1/notifications/templates
  - GET /api/v1/notifications/templates
  - PUT /api/v1/notifications/templates/:id
  - Variable replacement logic
  - **AC**: Templates working

### Email Service

- [ ] 📝 Configure SMTP (3h) 🔴
  - SMTP settings
  - Connection testing
  - **AC**: Email connection working

- [ ] 📝 Implement email sending (8h) 🔴
  - HTML templates
  - Send function
  - Error handling
  - **AC**: Emails sent successfully

- [ ] 📝 Implement email queue (6h) 🟡
  - Queue emails
  - Process queue
  - Retry on failure
  - **AC**: Queue working

### WhatsApp Integration

- [ ] 📝 Configure WhatsApp API (4h) 🔴
  - API credentials
  - Connection testing
  - **AC**: WhatsApp connection working

- [ ] 📝 Implement WhatsApp sending (8h) 🔴
  - Text messages
  - Template messages
  - Error handling
  - **AC**: WhatsApp messages sent

- [ ] 📝 Implement webhook handler (5h) 🟡
  - Receive status updates
  - Update notification status
  - **AC**: Webhook working

### Event-Driven Messaging

- [ ] 📝 Setup RabbitMQ (4h) 🔴
  - RabbitMQ container
  - Connection configuration
  - **AC**: RabbitMQ running

- [ ] 📝 Implement event publisher (6h) 🔴
  - Publish function
  - Event schema
  - **AC**: Events published

- [ ] 📝 Implement event consumer (8h) 🔴
  - Subscribe to events
  - Process events
  - Send notifications
  - **AC**: Events consumed

- [ ] 📝 Implement retry mechanism (5h) 🟡
  - Retry failed notifications (3 attempts)
  - Dead letter queue
  - **AC**: Retry working

### Notification Handlers

- [ ] 📝 Implement notification handlers (8h) 🔴
  - POST /api/v1/notifications/send
  - POST /api/v1/notifications/send/bulk
  - GET /api/v1/notifications
  - GET /api/v1/notifications/:id
  - GET /api/v1/notifications/user/:user_id
  - **AC**: All endpoints working

- [ ] 📝 Unit tests for Notification (10h) 🔴
  - Repository tests
  - Email tests
  - WhatsApp tests
  - Event tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [ ] 📝 Notification service integration tests (10h) 🔴
  - Email sending flow
  - WhatsApp sending flow
  - Event-driven flow
  - **AC**: Integration tests passing

### Documentation

- [ ] 📝 Write Notification service README (3h) 🟡
- [ ] 📝 Generate Swagger docs (4h) 🟡
- [ ] 📝 Create Postman collection (3h) 🟡

---

## API Gateway

### Core Setup

- [ ] 📝 Create api-gateway structure (4h) 🔴
  - Initialize project
  - Choose gateway (Kong/Traefik/Custom)
  - **AC**: Gateway structure ready

### Gateway Configuration

- [ ] 📝 Configure routing (8h) 🔴
  - Route to auth-service
  - Route to academic-service
  - Route to attendance-service
  - Route to assessment-service
  - Route to admission-service
  - Route to finance-service
  - Route to notification-service
  - **AC**: All routes working

- [ ] 📝 Implement rate limiting (5h) 🔴
  - Global rate limiting
  - Per-service rate limiting
  - **AC**: Rate limiting active

- [ ] 📝 Implement request logging (4h) 🟡
  - Log all requests
  - Request ID generation
  - **AC**: Logging working

- [ ] 📝 Implement authentication (6h) 🔴
  - JWT validation
  - Forward user context
  - **AC**: Auth working

- [ ] 📝 Setup load balancing (5h) 🟡
  - Load balancing strategy
  - Health check integration
  - **AC**: Load balancing working

- [ ] 📝 Unit tests for Gateway (8h) 🔴
  - Routing tests
  - Rate limiting tests
  - Coverage >70%
  - **AC**: Tests passing

### Documentation

- [ ] 📝 Write Gateway README (2h) 🟡
- [ ] 📝 Gateway architecture diagram (2h) 🟡

---

## Supporting Services

### File Service (Optional)

- [ ] 📝 Create file-service structure (3h) 🟢
- [ ] 📝 Implement file upload (8h) 🟢
  - POST /api/v1/files/upload
  - File validation
  - Store in object storage (MinIO/S3)
  - **AC**: Upload working

- [ ] 📝 Implement file download (4h) 🟢
  - GET /api/v1/files/:id
  - Signed URLs
  - **AC**: Download working

- [ ] 📝 Implement file deletion (3h) 🟢
  - DELETE /api/v1/files/:id
  - Soft delete
  - **AC**: Deletion working

### Report Service (Optional)

- [ ] 📝 Create report-service structure (3h) 🟢
- [ ] 📝 Implement custom reports (12h) 🟢
  - Report builder
  - Data aggregation
  - Export (PDF, Excel)
  - **AC**: Reports working

---

## Performance & Optimization

### Caching

- [ ] 📝 Implement Redis caching (8h) 🟡
  - Cache user sessions
  - Cache frequently accessed data
  - Cache invalidation strategy
  - **AC**: Caching working

### Database Optimization

- [ ] 📝 Add database indexes (6h) 🟡
  - Analyze slow queries
  - Add indexes
  - Test performance
  - **AC**: Queries faster

- [ ] 📝 Implement connection pooling (3h) 🔴
  - Configure pool size
  - Monitor connections
  - **AC**: Pooling working

### Load Testing

- [ ] 📝 Write k6 load test scripts (10h) 🟡
  - Baseline test
  - Stress test
  - Spike test
  - **AC**: Scripts ready

- [ ] 📝 Run load tests (8h) 🟡
  - Execute tests
  - Analyze results
  - Optimize bottlenecks
  - **AC**: Performance targets met

---

## Security Hardening

### Security Measures

- [ ] 📝 Implement field-level encryption (8h) 🟡
  - Encrypt PII fields
  - Key management
  - **AC**: Sensitive data encrypted

- [ ] 📝 Setup Gosec scanning (3h) 🔴
  - Configure Gosec
  - Add to CI
  - Fix issues
  - **AC**: Security scan passing

- [ ] 📝 Setup Trivy scanning (3h) 🔴
  - Configure Trivy
  - Scan containers
  - Fix vulnerabilities
  - **AC**: No critical vulnerabilities

- [ ] 📝 Implement security headers (2h) 🔴
  - Add security headers
  - Test headers
  - **AC**: Headers present

### Penetration Testing

- [ ] 📝 Conduct penetration testing (16h) 🟡
  - Engage security team
  - Test vulnerabilities
  - Fix issues
  - **AC**: No critical issues

---

## Deployment & DevOps

### Kubernetes Setup

- [ ] 📝 Setup Kubernetes cluster (12h) 🟡
  - Provision cluster
  - Configure networking
  - Setup ingress
  - **AC**: Cluster operational

- [ ] 📝 Create Helm charts (16h) 🟡
  - Chart per service
  - ConfigMaps & Secrets
  - Deployments & Services
  - **AC**: Helm charts working

- [ ] 📝 Configure auto-scaling (8h) 🟡
  - HPA configuration
  - Resource limits
  - **AC**: Auto-scaling working

### Monitoring

- [ ] 📝 Setup Prometheus (6h) 🟡
  - Prometheus deployment
  - Service monitors
  - **AC**: Metrics collected

- [ ] 📝 Create Grafana dashboards (8h) 🟡
  - Service dashboards
  - System dashboards
  - **AC**: Dashboards working

- [ ] 📝 Configure alerting (6h) 🟡
  - Alert rules
  - Notification channels
  - **AC**: Alerts working

### Backup & Recovery

- [ ] 📝 Setup automated backups (8h) 🟡
  - Database backups
  - File backups
  - **AC**: Backups running

- [ ] 📝 Test disaster recovery (8h) 🟡
  - DR procedures
  - Restore testing
  - **AC**: DR working

---

## Documentation

### Technical Documentation

- [ ] 📝 Write system architecture doc (8h) 🟡
  - Architecture diagrams
  - Service interactions
  - **AC**: Architecture documented

- [ ] 📝 Create ADRs (12h) 🟡
  - Document key decisions
  - Rationale & consequences
  - **AC**: ADRs complete

- [ ] 📝 Write deployment guide (6h) 🟡
  - Deployment steps
  - Rollback procedures
  - **AC**: Guide complete

- [ ] 📝 Write runbooks (16h) 🟡
  - Incident response
  - Troubleshooting
  - Common issues
  - **AC**: Runbooks complete

### API Documentation

- [ ] 📝 Complete Swagger docs (12h) 🟡
  - All services documented
  - Examples included
  - **AC**: API docs complete

- [ ] 📝 Create Postman collections (8h) 🟡
  - All services
  - Environment setup
  - **AC**: Collections working

---

## Testing

### Unit Tests

- [ ] 📝 Achieve 70% code coverage (40h) 🔴
  - Write missing tests
  - Fix failing tests
  - **AC**: Coverage >70%

### Integration Tests

- [ ] 📝 Write integration tests (32h) 🔴
  - Test service interactions
  - Test database operations
  - **AC**: Integration tests passing

### E2E Tests

- [ ] 📝 Write E2E tests (24h) 🟡
  - Critical user flows
  - Full system tests
  - **AC**: E2E tests passing

---

## Production Readiness

### Pre-Production Checklist

- [ ] 📝 Complete security audit (16h) 🔴
  - Security review
  - Fix vulnerabilities
  - **AC**: Audit passed

- [ ] 📝 Performance testing (16h) 🟡
  - Load testing
  - Stress testing
  - **AC**: Performance targets met

- [ ] 📝 Setup monitoring & alerting (12h) 🟡
  - Monitoring operational
  - Alerts configured
  - **AC**: Monitoring working

- [ ] 📝 Backup & DR verification (8h) 🟡
  - Test backups
  - Test restore
  - **AC**: Backup/restore working

- [ ] 📝 Documentation review (8h) 🟡
  - Review all docs
  - Update as needed
  - **AC**: Docs complete

- [ ] 📝 Stakeholder sign-off (4h) 🔴
  - Demo to stakeholders
  - Get approval
  - **AC**: Sign-off obtained

---

## Summary Statistics

### Total Estimated Hours by Service

| Service | Estimated Hours |
|---------|----------------|
| Global Infrastructure | 150h |
| Auth Service | 180h |
| Academic Core Service | 220h |
| Attendance Service | 90h |
| Assessment Service | 150h |
| Admission Service | 120h |
| Finance Service | 120h |
| Notification Service | 110h |
| API Gateway | 42h |
| Supporting Services | 50h |
| Performance & Optimization | 60h |
| Security Hardening | 50h |
| Deployment & DevOps | 100h |
| Documentation | 70h |
| Testing | 112h |
| Production Readiness | 64h |
| **TOTAL** | **1,688h** |

### Priority Breakdown

- 🔴 High Priority: ~950h (56%)
- 🟡 Medium Priority: ~590h (35%)
- 🟢 Low Priority: ~148h (9%)

### Team Size Estimation

Assuming:
- 1 developer = 160h/month (40h/week × 4 weeks)
- Total hours = 1,688h

**Options**:
1. **4 developers × 3 months** = 1,920h (buffer: 232h)
2. **5 developers × 2 months** = 1,600h (tight schedule)
3. **3 developers × 4 months** = 1,920h (comfortable pace)

**Recommended**: 4 developers × 3 months

---

## Quick Start Checklist

### Week 1 Priority Tasks (Must Complete)
- [ ] Setup monorepo structure
- [ ] Create docker-compose.yml
- [ ] Create Makefile
- [ ] Implement shared packages (config, database, logger)
- [ ] Setup CI/CD pipeline basics
- [ ] Create auth-service structure
- [ ] Setup first database migrations

### Critical Path Items (Blocking Others)
1. ✅ Shared packages (blocks all services)
2. ✅ Auth service (blocks all protected endpoints)
3. ✅ Academic core (blocks attendance, assessment)
4. ✅ Database schema (blocks all data operations)
5. ✅ API Gateway (blocks external access)

---

**Last Updated**: 2025-01-15  
**Version**: 2.0  
**Owner**: Engineering Team  
**Status**: Active Task List

---

## Notes

### Development Best Practices
- Always write tests before marking task complete
- Update documentation as you code
- Create small, focused PRs
- Get code reviews before merging
- Run linter before committing
- Keep task list up-to-date

### When Task is Blocked
1. Update task status to ⏳
2. Document blocker in task notes
3. Notify team lead
4. Work on non-blocked tasks
5. Regularly check blocker status

### Definition of Done
A task is complete when:
- [ ] Code written & tested
- [ ] Unit tests passing (>70% coverage)
- [ ] Code reviewed & approved
- [ ] Documentation updated
- [ ] Changes merged to develop
- [ ] Task marked as ✅ in this list
