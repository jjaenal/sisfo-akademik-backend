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

## 🚀 Next Priority Tasks (Sorted)

### High Priority (Critical) 🔴

- [ ] 🔄 📝 Achieve 70% code coverage (40h)
  - [x] File Service — UseCase coverage ≥85% (88.7%)
  - [x] File Service — Handler coverage ≥90% (93.3%)
  - [x] File Service — Repository coverage ≥90% (90.4%)
  - [x] Attendance Service — School client coverage ≥80% (88.9%)
- [ ] 🔄 Complete security audit (16h)
  - [x] Security review (Gosec scan)
  - [x] Fix vulnerabilities (Permissions, Unhandled errors)
  - [ ] Manual review of file inclusion warnings
- [ ] 📝 Stakeholder sign-off (4h)

### Medium Priority (Important) 🟡

- [ ] 📝 Performance testing (16h)
- [ ] 📝 Setup Kubernetes cluster (12h)
- [x] ✅ 📝 Create Helm charts (16h)
  - Chart: `deploy/helm/microservice/`
  - Values per service: `deploy/helm/values/*.yaml`
- [ ] 📝 Configure auto-scaling (8h)
- [ ] 📝 Setup Prometheus (6h)
- [ ] 📝 Create Grafana dashboards (8h)
- [ ] 📝 Configure alerting (6h)
- [ ] 📝 Setup automated backups (8h)
- [ ] 📝 Test disaster recovery (8h)
- [ ] 📝 Write system architecture doc (8h)
- [ ] 📝 Create ADRs (12h)
- [ ] 📝 Write deployment guide (6h)
- [ ] 📝 Write runbooks (16h)
- [ ] 📝 Complete Swagger docs (12h)
- [ ] 📝 Create Postman collections (8h)
- [ ] 📝 Write E2E tests (24h)
- [ ] 📝 Setup monitoring & alerting (12h)
- [ ] 📝 Backup & DR verification (8h)
- [ ] 📝 Documentation review (8h)

---

## Global Infrastructure

### Repository & Project Setup

- [x] ✅ 📝 Create monorepo structure (4h) 🔴

  - Setup services/, shared/, infrastructure/, docs/ folders
  - Configure Go workspace
  - Setup .gitignore for Go projects
  - **AC**: Directory structure matches standard layout

- [x] ✅ 📝 Ensure Go version consistency (1.25.5) 🔴

  - Update all go.mod files
  - Update go.work
  - **AC**: All modules using Go 1.25.5

- [x] ✅ 📝 Setup GitHub/GitLab organization (2h) 🔴

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

- [x] ✅ 📝 Create docker-compose.yml for local dev (6h) 🔴

  - PostgreSQL container
  - Redis container
  - RabbitMQ container
  - Jaeger container
  - Service containers (placeholder)
  - **AC**: `docker-compose up` starts all services

- [x] ✅ 📝 Create Makefile with common commands (3h) 🔴

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

- [x] ✅ 📝 Setup environment files (2h) 🔴

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

- [x] ✅ 📝 Implement config package (4h) 🔴

  - Viper integration
  - Environment variable loading
  - Config validation
  - Hot reload support
  - **AC**: Config loaded correctly with validation

- [x] ✅ 📝 Implement database package (6h) 🔴

  - PostgreSQL connection
  - Connection pooling (pgxpool)
  - Health check
  - Transaction helper
  - **AC**: DB connection stable with pooling

- [x] ✅ 📝 Implement Redis package (4h) 🔴

  - Redis client
  - Connection pooling
  - Health check
  - Helper functions (Get, Set, Delete)
  - **AC**: Redis operations working

- [x] ✅ 📝 Implement logger package (5h) 🔴

  - Structured logging (zap/logrus)
  - Log levels
  - Context-aware logging
  - JSON output
  - **AC**: Logs formatted correctly

- [x] ✅ 📝 Implement middleware package (8h) 🔴

  - Authentication middleware
  - Authorization middleware
  - Logging middleware
  - CORS middleware
  - Rate limiting middleware
  - Error handling middleware
  - **AC**: All middleware functional

- [x] ✅ 📝 Implement errors package (3h) 🟡

  - Custom error types
  - Error codes
  - Error wrapping
  - HTTP error responses
  - **AC**: Consistent error handling

- [x] ✅ 📝 Implement validator package (4h) 🟡

  - go-playground/validator wrapper
  - Custom validators
  - Validation error formatting
  - **AC**: Input validation working

- [x] ✅ 📝 Implement JWT package (5h) 🔴

  - JWT generation
  - JWT validation
  - Token refresh
  - Claims extraction
  - **AC**: JWT operations secure & functional

- [x] ✅ 📝 Implement httputil package (3h) 🟡

  - Response helpers
  - Request parsing
  - Pagination helpers
  - **AC**: HTTP utilities working

- [x] ✅ 📝 Implement testutil package (4h) 🟡
  - Test database helpers
  - Mock helpers
  - Test fixtures
  - **AC**: Tests easier to write

### CI/CD Pipeline

- [x] ✅ 📝 Setup GitHub Actions workflow (6h) 🔴

  - Lint job
  - Test job
  - Security scan job
  - Build job
  - Deploy staging job
  - **AC**: Pipeline runs on push

- [x] ✅ 📝 Configure linting (golangci-lint) (2h) 🔴

  - Install golangci-lint
  - Configure .golangci.yml
  - Add to CI
  - **AC**: Code passes linting

- [ ] 📝 Setup test coverage reporting (3h) 🟡

  - Integrate Codecov
  - Coverage badge
  - Coverage threshold (70%)
  - **AC**: Coverage tracked in CI

- [x] ✅ 📝 Setup security scanning (4h) 🔴

  - Gosec for static analysis
  - Trivy for container scanning
  - Snyk for dependencies
  - **AC**: Security scans in CI

- [x] ✅ 📝 Configure Docker registry (2h) 🔴
  - GitHub Container Registry or Docker Hub
  - Setup credentials
  - Image tagging strategy
  - **AC**: Images pushed to registry

### Observability

- [x] ✅ 📝 Setup ELK Stack (8h) 🟡

  - Elasticsearch container
  - Logstash container
  - Kibana container
  - Configure log shipping
  - **AC**: Logs viewable in Kibana

- [x] ✅ 📝 Setup Prometheus (6h) 🟡

  - Prometheus container
  - Configure scraping
  - Define metrics
  - **AC**: Metrics scraped

- [x] ✅ 📝 Setup Grafana (6h) 🟡

  - Grafana container
  - Connect to Prometheus
  - Create dashboards
  - **AC**: Dashboards showing metrics

- [x] ✅ 📝 Setup Jaeger (4h) 🟡
  - Jaeger container
  - Trace instrumentation
  - Log TraceID Injection
  - **AC**: Traces visible in Jaeger

---

## Auth / Identity Service

### Core Setup

- [x] ✅ 📝 Create auth-service structure (3h) 🔴

  - Initialize Go module
  - Setup directory structure
  - Create Dockerfile
  - Create docker-compose.yml
  - **AC**: Service structure ready

- [x] ✅ 📝 Setup database connection (2h) 🔴

  - Use shared database package
  - Test connection
  - **AC**: Auth service connects to DB

- [x] ✅ 📝 Create database migrations (4h) 🔴
  - Users table
  - Roles table
  - Permissions table
  - Role_permissions table
  - User_roles table
  - Audit_logs table
  - **AC**: Migrations run successfully

Status: ✅ Created migrations for users, roles, permissions, role_permissions, user_roles, and audit_logs with up/down SQL.

### User Management

- [x] ✅ 📝 Implement User entity (2h) 🔴

- Define User struct
- Validation rules
- Methods (BeforeCreate, etc)
- **AC**: User entity complete

- [x] ✅ 📝 Implement User repository (6h) 🔴

- Create()
- GetByID()
- GetByEmail()
- List() with pagination
- Update()
- Delete() (soft delete)
- **AC**: CRUD operations working

- [x] ✅ 📝 Implement User use case (6h) 🔴

- Register user
- Get user
- Update user
- Delete user
- Search users
- **AC**: Business logic implemented

- [x] ✅ 📝 Implement User handlers (8h) 🔴

- POST /api/v1/users
- GET /api/v1/users/:id
- GET /api/v1/users
- PUT /api/v1/users/:id
- DELETE /api/v1/users/:id
- PATCH /api/v1/users/:id/activate
- **AC**: All endpoints working

- [x] ✅ 📝 Add input validation (3h) 🔴

- Email format
- Password strength
- Required fields
- **AC**: Invalid input rejected

- [x] ✅ 📝 Implement password hashing (2h) 🔴

- bcrypt implementation
- Cost factor: 12
- **AC**: Passwords hashed securely

- [x] ✅ 📝 Unit tests for User (8h) 🔴

- Repository tests
- Use case tests
- Handler tests
- Coverage >70%
- **AC**: Tests passing with good coverage

### Authentication

- [x] ✅ 📝 Implement login handler (6h) 🔴

- POST /api/v1/auth/login
- Validate credentials
- Generate tokens
- Return access & refresh tokens
- **AC**: Login working

- [x] ✅ 📝 Implement JWT generation (4h) 🔴

  - Access token (15 min TTL)
  - Refresh token (7 days TTL)
  - Include claims (user_id, tenant_id, roles, permissions)
  - **AC**: JWT generated correctly

- [x] ✅ 📝 Implement logout handler (3h) 🔴

- POST /api/v1/auth/logout
- Invalidate refresh token
- Blacklist access token in Redis
- **AC**: Logout working

- [x] ✅ 📝 Implement token refresh (5h) 🔴

- POST /api/v1/auth/refresh
- Validate refresh token
- Generate new access token
- Rotate refresh token
- **AC**: Token refresh working

- [x] ✅ 📝 Implement forgot password (6h) 🟡

  - POST /api/v1/auth/forgot-password
  - Generate reset token
  - Send reset email
  - **AC**: Reset email sent

- [x] ✅ 📝 Implement reset password (5h) 🟡

  - POST /api/v1/auth/reset-password
  - Validate reset token
  - Update password
  - **AC**: Password reset working

- [x] ✅ 📝 Implement change password (4h) 🟡

  - POST /api/v1/auth/change-password
  - Validate old password
  - Update password
  - **AC**: Password change working

- [x] ✅ 📝 Implement failed login tracking (4h) 🟡

  - Track failed attempts
  - Lock account after 5 failures
  - Auto-unlock after 30 minutes
  - **AC**: Account lockout working

- [x] ✅ 🔄 📝 Unit tests for Auth (10h) 🔴

- Login tests
- Token generation tests
- Token refresh tests
- Logout tests
- Coverage >70%
- **AC**: Tests passing

### RBAC (Role-Based Access Control)

- [x] ✅ 📝 Implement Role entity (2h) 🔴

  - Define Role struct
  - Validation rules
  - **AC**: Role entity complete

- [x] ✅ 📝 Implement Permission entity (2h) 🔴

  - Define Permission struct
  - Resource:Action format
  - **AC**: Permission entity complete

- [x] ✅ 📝 Implement Role repository (6h) 🔴

  - Create()
  - GetByID()
  - List()
  - Update()
  - Delete()
  - AssignPermissions()
  - **AC**: Role CRUD working

- [x] ✅ 📝 Implement Permission repository (4h) 🔴

  - Create()
  - List()
  - GetByRole()
  - **AC**: Permission operations working

- [x] ✅ 📝 Seed default roles & permissions (4h) ✅

  - Super Admin role
  - School Admin role
  - Teacher role
  - Student role
  - Parent role
  - Default permissions
  - **AC**: Default roles created

- [x] ✅ 📝 Implement role assignment (5h) 🔴

  - Assign role to user
  - Remove role from user
  - Get user roles
  - Get user permissions (effective)
  - **AC**: Role assignment working

- [x] ✅ 📝 Implement RBAC middleware (8h) 🔴

  - Check user authentication
  - Check user permissions
  - Context-aware (tenant-based)
  - **AC**: Protected endpoints secured

- [x] ✅ 📝 Implement role handlers (8h) 🟡

  - POST /api/v1/roles
  - GET /api/v1/roles
  - GET /api/v1/roles/:id
  - PUT /api/v1/roles/:id
  - DELETE /api/v1/roles/:id
  - **AC**: Role management working

- [x] ✅ 📝 Unit tests for RBAC (10h) 🔴
  - Role tests
  - Permission tests
  - Middleware tests
  - Coverage >70%
  - **AC**: Tests passing

### Audit Logging

- [x] ✅ 📝 Implement audit log entity (2h) 🔴

  - Define AuditLog struct
  - Fields: user, action, resource, changes
  - **AC**: Audit log entity complete

- [x] ✅ 📝 Implement audit log repository (4h) 🔴

  - Create()
  - List() with filters
  - Search()
  - **AC**: Audit logs stored

- [x] ✅ 📝 Implement audit middleware (6h) 🔴

  - Capture request details
  - Log after response
  - Async logging (don't block)
  - **AC**: All actions logged

- [x] ✅ 📝 Implement audit log handlers (4h) 🟡

  - GET /api/v1/audit-logs
  - GET /api/v1/audit-logs/search
  - Export audit logs
  - **AC**: Audit logs viewable

- [x] ✅ 📝 Setup log retention (2h) 🟡
  - 90-day retention
  - Automated cleanup job
  - **AC**: Old logs cleaned up

### Security Enhancements

- [x] ✅ 📝 Implement rate limiting (6h) 🔴

  - Redis-based rate limiter
  - Different limits per endpoint type
  - Rate limit headers
  - **AC**: Rate limiting active

- [x] ✅ 📝 Implement security headers (3h) 🔴

  - X-Content-Type-Options
  - X-Frame-Options
  - X-XSS-Protection
  - Strict-Transport-Security
  - CSP
  - **AC**: Security headers present

- [x] ✅ 📝 Configure CORS (2h) 🔴

  - Whitelist origins
  - Allowed methods & headers
  - **AC**: CORS working

- [x] ✅ 📝 Implement password validation (3h) 🟡

  - Minimum 8 characters
  - Complexity requirements
  - Common password check
  - **AC**: Weak passwords rejected

- [x] ✅ 📝 Implement password history (3h) 🟡
  - Track last 5 passwords
  - Prevent reuse
  - **AC**: Password reuse prevented

### Integration Tests

- [x] ✅ 📝 Auth service integration tests (12h) 🔴
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

- [x] ✅ 📝 Create Postman collection (3h) 🟡
  - All endpoints
  - Example requests
  - Environment variables
  - **AC**: Postman collection works

---

## Academic Core Service

### Core Setup

- [x] ✅ 📝 Create academic-service structure (3h) 🔴

  - Initialize Go module
  - Setup directory structure
  - Create Dockerfile
  - **AC**: Service structure ready

- [x] ✅ 📝 Setup database connection (2h) 🔴

  - Use shared database package
  - Test connection
  - **AC**: Service connects to DB

- [x] ✅ 📝 Create database migrations (8h) 🔴
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

- [x] ✅ 📝 Implement School entity (2h) 🔴

  - Define School struct
  - Validation rules
  - **AC**: School entity complete

- [x] ✅ 📝 Implement School repository (6h) 🔴

  - CRUD operations
  - GetByTenantID()
  - **AC**: School CRUD working

- [x] ✅ 📝 Implement School use case (4h) 🔴

  - Business logic
  - Validation
  - **AC**: School operations working

- [x] ✅ 📝 Implement School handlers (8h) 🔴

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

- [x] ✅ 📝 Unit tests for School (8h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Academic Year & Semester

- [x] ✅ 📝 Implement AcademicYear entity (2h) 🔴

  - Define AcademicYear struct
  - Validation (dates, active flag)
  - **AC**: AcademicYear entity complete

- [x] ✅ 📝 Implement Semester entity (2h) 🔴

  - Define Semester struct
  - Validation
  - **AC**: Semester entity complete

- [x] ✅ 📝 Implement AcademicYear repository (6h) 🔴

  - CRUD operations
  - GetActive()
  - ValidateNonOverlap()
  - **AC**: AcademicYear CRUD working

- [x] ✅ 📝 Implement Semester repository (5h) 🔴

  - CRUD operations
  - GetBySemester()
  - GetActive()
  - **AC**: Semester CRUD working

- [x] ✅ 📝 Implement academic year handlers (8h) 🔴

  - POST /api/v1/academic-years
  - GET /api/v1/academic-years
  - GET /api/v1/academic-years/:id
  - PUT /api/v1/academic-years/:id
  - PATCH /api/v1/academic-years/:id/activate
  - DELETE /api/v1/academic-years/:id
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement semester handlers (8h) 🔴

  - POST /api/v1/semesters
  - GET /api/v1/semesters
  - GET /api/v1/semesters/:id
  - PUT /api/v1/semesters/:id
  - PATCH /api/v1/semesters/:id/activate
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement active year/semester validation (3h) 🔴

  - Only 1 active year per tenant
  - Only 1 active semester per year
  - **AC**: Validation working

- [x] ✅ 📝 Unit tests for AcademicYear & Semester (10h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Class Management

- [x] ✅ 📝 Implement Class entity (3h) 🔴

  - Define Class struct
  - Validation rules
  - **AC**: Class entity complete

- [x] ✅ 📝 Implement ClassStudent entity (Enrollment) (2h) 🔴

  - Enrollment tracking
  - Status (active, transferred, graduated)
  - **AC**: ClassStudent entity complete

- [x] ✅ 📝 Implement Class repository (8h) 🔴

  - CRUD operations
  - GetByAcademicYear()
  - GetStudents()
  - EnrollStudent()
  - **AC**: Class operations working

- [x] ✅ 📝 Implement class handlers (10h) 🔴

  - POST /api/v1/classes
  - GET /api/v1/classes
  - GET /api/v1/classes/:id
  - PUT /api/v1/classes/:id
  - DELETE /api/v1/classes/:id
  - POST /api/v1/classes/:id/students (Implemented via Enrollment)
  - GET /api/v1/classes/:id/students
  - DELETE /api/v1/classes/:id/students/:student_id
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement bulk enrollment (5h) 🟡

  - POST /api/v1/classes/:id/students/bulk
  - CSV import
  - Validation
  - **AC**: Bulk enrollment working

- [x] ✅ 📝 Implement capacity management (3h) 🟡

  - Check max_students
  - Prevent over-enrollment
  - **AC**: Capacity enforced

- [x] ✅ 📝 Unit tests for Class (10h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Subject Management

- [x] ✅ 📝 Implement Subject entity (2h) 🔴

  - Define Subject struct
  - Categories (Wajib, Peminatan, Mulok)
  - **AC**: Subject entity complete

- [x] ✅ 📝 Implement ClassSubject entity (2h) 🔴

  - Subject-class-teacher mapping
  - **AC**: ClassSubject entity complete

- [x] ✅ 📝 Implement Subject repository (6h) 🔴

  - CRUD operations
  - GetByCategory()
  - AssignToClass()
  - **AC**: Subject operations working

- [x] ✅ 📝 Implement subject handlers (10h) 🔴

  - POST /api/v1/subjects
  - GET /api/v1/subjects
  - GET /api/v1/subjects/:id
  - PUT /api/v1/subjects/:id
  - DELETE /api/v1/subjects/:id
  - POST /api/v1/classes/:id/subjects
  - GET /api/v1/classes/:id/subjects
  - DELETE /api/v1/classes/:id/subjects/:subject_id
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement teacher assignment (4h) 🔴

  - PUT /api/v1/classes/:id/subjects/:subject_id/teacher
  - Validation
  - **AC**: Teacher assignment working

- [x] ✅ 📝 Unit tests for Subject (8h) 🔴
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Curriculum Management

- [x] ✅ 📝 Implement Curriculum entity (3h) 🔴

  - Define Curriculum struct
  - Support multiple curricula per tenant
  - **AC**: Curriculum entity complete

- [x] ✅ 📝 Implement GradingRule entity (3h) 🔴

  - KKM configuration
  - Grade components & weights
  - **AC**: GradingRule entity complete

- [x] ✅ 📝 Implement Curriculum repository (6h) 🔴

  - CRUD operations
  - GetSubjects()
  - GetGradingRules()
  - **AC**: Curriculum operations working

- [x] ✅ 📝 Implement curriculum handlers (10h) 🟡

  - POST /api/v1/curricula
  - GET /api/v1/curricula
  - GET /api/v1/curricula/:id
  - PUT /api/v1/curricula/:id
  - POST /api/v1/curricula/:id/subjects
  - GET /api/v1/curricula/:id/subjects
  - POST /api/v1/curricula/:id/grading-rules
  - **AC**: All endpoints working

- [x] ✅ 📝 Unit tests for Curriculum (8h) 🟡
  - Repository tests
  - Use case tests
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Schedule Management

- [x] ✅ 📝 Implement Schedule entity (3h) 🔴

  - Define Schedule struct
  - Day of week, time slots
  - **AC**: Schedule entity complete

- [x] ✅ 📝 Implement Schedule repository (6h) 🔴

  - CRUD operations
  - GetWeeklySchedule()
  - CheckConflicts()
  - **AC**: Schedule operations working

- [x] ✅ 📝 Implement conflict detection (6h) 🔴

  - Class conflict check
  - Teacher conflict check
  - Room conflict check
  - **AC**: Conflicts detected

- [x] ✅ 📝 Implement schedule handlers (10h) 🔴

  - POST /api/v1/schedules
  - GET /api/v1/schedules
  - PUT /api/v1/schedules/:id
  - DELETE /api/v1/schedules/:id
  - GET /api/v1/schedules/class/:class_id/weekly
  - GET /api/v1/schedules/teacher/:teacher_id/weekly
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement bulk schedule creation (5h) 🟡

  - Template system
  - Batch creation
  - **AC**: Bulk creation working

- [x] ✅ 📝 Unit tests for Schedule (8h) 🔴
  - Repository tests (Mocked)
  - Conflict tests (Mocked)
  - Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [x] ✅ 📝 Academic service integration tests (12h) 🔴
  - School creation flow
  - Academic year setup flow
  - Class & student enrollment flow
  - Schedule creation flow
  - **AC**: Integration tests passing

### Documentation

- [x] ✅ 📝 Write Academic service README (3h) 🟡

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

- [x] 📝 Create attendance-service structure (3h) ✅

  - Initialize Go module
  - Setup directory structure
  - **AC**: Service structure ready

- [x] 📝 Setup database connection (2h) ✅

  - Use shared database package
  - **AC**: Service connects to DB

- [x] ✅ 📝 Create database migrations (4h) 🔴
  - Student_attendance table
  - Teacher_attendance table
  - **AC**: Migrations run

### Student Attendance

- [x] ✅ 📝 Implement StudentAttendance entity (2h) 🔴

  - Define struct
  - Status types (present, absent, late, excused, sick)
  - **AC**: Entity complete

- [x] ✅ 📝 Implement StudentAttendance repository (6h) 🔴

  - Create()
  - Update()
  - GetByStudentAndDate()
  - List() with filters
  - GetSummary()
  - **AC**: CRUD working

- [x] ✅ 📝 Implement attendance handlers (10h) 🔴

  - POST /api/v1/attendance/students
  - POST /api/v1/attendance/students/bulk
  - GET /api/v1/attendance/students
  - PUT /api/v1/attendance/students/:id
  - GET /api/v1/attendance/students/:student_id/summary
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement bulk check-in (5h) 🔴

  - Full class check-in
  - Validation
  - **AC**: Bulk check-in working

- [x] ✅ 📝 Implement GPS validation (4h) 🟡

  - Validate location against school location
  - Distance calculation
  - **AC**: GPS validation working

- [x] ✅ 📝 Unit tests for StudentAttendance (8h) ✅
  - [x] Repository tests
  - [x] Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Teacher Attendance

- [x] ✅ 📝 Implement TeacherAttendance entity (2h) 🔴

  - Define struct
  - Check-in/check-out times
  - **AC**: Entity complete

- [x] ✅ 📝 Implement TeacherAttendance repository (5h) 🔴

  - Create()
  - Update()
  - GetByTeacherAndDate()
  - List()
  - **AC**: CRUD working

- [x] ✅ 📝 Implement teacher attendance handlers (8h) 🔴

  - POST /api/v1/attendance/teachers/check-in
  - POST /api/v1/attendance/teachers/check-out
  - GET /api/v1/attendance/teachers
  - GET /api/v1/attendance/teachers/:teacher_id/summary
  - **AC**: All endpoints working

- [x] ✅ 📝 Unit tests for TeacherAttendance (6h) ✅
  - [x] Repository tests
  - [x] Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Reports

- [x] 📝 Implement attendance reports (8h) ✅
  - [x] `GET /api/v1/attendance/reports/daily`
  - [x] `GET /api/v1/attendance/reports/monthly`
  - [x] `GET /api/v1/attendance/reports/class/:class_id`
  - **AC**: Reports working

### Integration Tests

- [x] ✅ 📝 Attendance service integration tests (8h) 🔴
  - Student attendance flow
  - Bulk check-in flow
  - Teacher attendance flow
  - **AC**: Integration tests passing

### Documentation

- [x] 📝 Write Attendance service README (2h) ✅
- [x] 📝 Generate Swagger docs (3h) ✅
- [x] 📝 Create Postman collection (2h) ✅

---

## Assessment Service

### Core Setup

- [x] 📝 Create assessment-service structure (3h) ✅
- [x] 📝 Setup database connection (2h) ✅
- [x] ✅ 📝 Create database migrations (6h) ✅
  - Grade_categories table
  - Assessments table
  - Grades table
  - Report_cards table
  - Report_card_details table

### Grading System

- [x] ✅ 📝 Implement GradeCategory entity (2h) ✅
- [x] ✅ 📝 Implement Assessment entity (3h) ✅
- [x] ✅ 📝 Implement Grade entity (3h) ✅
- [x] ✅ 📝 Implement ReportCard entity (3h) ✅

- [x] ✅ 📝 Implement grade repositories (8h) ✅

  - GradeCategory CRUD
  - Assessment CRUD
  - Grade CRUD (bulk update)

- [x] ✅ 📝 Implement grade calculation engine (8h) 🔴

  - Calculate weighted scores
  - Final score calculation
  - Grade letter assignment
  - KKM validation
  - **AC**: Grades calculated correctly

- [x] ✅ 📝 Implement grade handlers (12h) 🔴

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

- [x] 📝 Implement grade approval workflow (5h) ✅

  - Draft → Submitted → Approved
  - Audit trail
  - **AC**: Workflow working

- [x] ✅ 📝 Unit tests for Grading (10h) 🔴
  - [x] GradeCategory Repository & Handler tests
  - [x] Calculation tests
  - [x] Repository tests
  - [x] Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Report Card Generation

- [x] ✅ 📝 Implement ReportCard entity (3h) ✅

  - Define struct
  - Status (draft, generated, published)
  - **AC**: Entity complete

- [x] ✅ 📝 Implement report card data aggregation (8h) ✅

  - Collect all grades
  - Calculate final scores
  - Get attendance summary
  - **AC**: Data aggregated correctly

- [x] ✅ 📝 Implement report card generation (12h) ✅

  - POST /api/v1/report-cards/generate/:student_id/:semester_id
  - POST /api/v1/report-cards/generate/class/:class_id/:semester_id
  - Generate report data
  - **AC**: Report cards generated

- [x] ✅ 📝 Implement PDF generation (12h) ✅

  - HTML template
  - Convert to PDF (maroto)
  - Store in object storage (Local storage implemented)
  - **AC**: PDF generated correctly

- [x] ✅ 📝 Implement report card handlers (8h) ✅

  - GET /api/v1/report-cards/:id (Done)
  - GET /api/v1/report-cards/student/:student_id (Done)
  - PATCH /api/v1/report-cards/:id/publish (Done)
  - GET /api/v1/report-cards/:id/pdf (Done)
  - GET /api/v1/report-cards/:id/download (Done)
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement template customization (6h) ✅

  - Template management
  - Variable replacement (Partially implemented in PDF generation)
  - **AC**: Templates customizable

- [x] ✅ 📝 Unit tests for ReportCard (10h) 🔴
  - [x] Generation tests
  - [x] PDF tests
  - [x] Handler tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [x] ✅ 📝 Assessment service integration tests (12h) ✅
  - [x] Grading flow
  - [x] Report card generation flow
  - **AC**: Integration tests passing

### Documentation

- [x] ✅ 📝 Write Assessment service README (3h) ✅
- [x] ✅ 📝 Generate Swagger docs (4h) 🟡
- [x] ✅ 📝 Create Postman collection (3h) ✅

---

## Admission Service (PPDB)

### Core Setup

- [x] ✅ 📝 Create admission-service structure (3h) ✅
- [x] ✅ 📝 Setup database connection (2h) ✅
- [x] ✅ 📝 Create database migrations (5h) 🔴
  - Admission_periods table
  - Applications table
  - Application_documents table

### Admission Management

- [x] ✅ 📝 Implement AdmissionPeriod entity (2h) 🔴
- [x] ✅ 📝 Implement Application entity (3h) 🔴
- [x] ✅ 📝 Implement ApplicationDocument entity (2h) 🔴

- [x] ✅ 📝 Implement admission repositories (8h) 🔴

  - AdmissionPeriod CRUD
  - Application CRUD
  - ApplicationDocument CRUD

- [x] ✅ 📝 Implement admission period handlers (8h) 🔴

  - POST /api/v1/admission/periods
  - GET /api/v1/admission/periods
  - GET /api/v1/admission/periods/:id
  - PUT /api/v1/admission/periods/:id
  - PATCH /api/v1/admission/periods/:id/close
  - **AC**: All endpoints working

- [x] ✅ 📝 Implement public application (10h) 🔴

  - GET /api/v1/admission/public/periods
  - POST /api/v1/admission/applications
  - GET /api/v1/admission/applications/:number/status
  - Application number generation
  - **AC**: Public application working

- [x] ✅ 📝 Implement document upload (8h) ✅

  - POST /api/v1/admission/applications/:id/documents
  - File validation (size, type)
  - Store in object storage (Local storage implemented)
  - **AC**: Upload working

- [x] ✅ 📝 Implement application management (10h) 🔴

  - GET /api/v1/admission/applications
  - GET /api/v1/admission/applications/:id
  - PUT /api/v1/admission/applications/:id
  - PATCH /api/v1/admission/applications/:id/verify
  - PATCH /api/v1/admission/applications/:id/accept
  - PATCH /api/v1/admission/applications/:id/reject
  - **AC**: Management working

- [x] ✅ � Implement selection process (10h) 🟡

  - [x] POST /api/v1/admission/applications/:id/test-score
  - [x] POST /api/v1/admission/applications/:id/interview-score
  - [x] POST /api/v1/admission/periods/:id/calculate-final-scores
  - [x] POST /api/v1/admission/periods/:id/announce
  - [x] Final score calculation
  - **AC**: Selection working

- [x] ✅ 📝 Implement student registration (8h) 🔴

  - POST /api/v1/admission/applications/:id/register
  - Create user account
  - Create student record
  - **AC**: Registration working

- [x] ✅ 📝 Unit tests for Admission (12h) 🔴
  - [x] AdmissionPeriod Repository & Handler tests
  - Repository tests
  - Handler tests
  - Selection logic tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [x] ✅ 📝 Admission service integration tests (10h) ✅
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

- [x] ✅ 📝 Create finance-service structure (3h) 🔴
- [x] ✅ 📝 Setup database connection (2h) 🔴
- [x] ✅ 📝 Create database migrations (4h) 🔴
  - Billing_configurations table
  - Invoices table
  - Payments table

### Finance Management

- [x] ✅ 📝 Implement BillingConfig entity (2h) 🔴
- [x] ✅ 📝 Implement Invoice entity (3h) 🔴
- [x] ✅ 📝 Implement Payment entity (2h) 🔴

- [x] ✅ 📝 Implement finance repositories (8h) 🔴

  - BillingConfig CRUD
  - Invoice CRUD
  - Payment CRUD

- [x] ✅ 📝 Implement billing configuration (8h) 🔴

  - POST /api/v1/finance/billing-configs
  - GET /api/v1/finance/billing-configs
  - PUT /api/v1/finance/billing-configs/:id
  - **AC**: Billing config working

- [x] ✅ 📝 Implement invoice generation (10h) 🔴

  - POST /api/v1/finance/invoices/generate
  - POST /api/v1/finance/invoices/generate/bulk
  - POST /api/v1/finance/invoices/generate/auto
  - Invoice number generation
  - **AC**: Invoice generation working

- [x] ✅ 📝 Implement auto-generation (8h) 🟡

  - Scheduled job (cron)
  - Monthly SPP generation
  - **AC**: Auto-generation working

- [x] ✅ 📝 Implement invoice handlers (8h) 🔴

  - GET /api/v1/finance/invoices
  - GET /api/v1/finance/invoices/:id
  - PUT /api/v1/finance/invoices/:id
  - GET /api/v1/finance/invoices/student/:student_id
  - GET /api/v1/finance/invoices/student/:student_id/outstanding
  - **AC**: Invoice management working

- [x] ✅ 📝 Implement payment recording (8h) 🔴

  - POST /api/v1/finance/payments
  - GET /api/v1/finance/payments
  - GET /api/v1/finance/payments/:id
  - Payment number generation
  - Receipt generation
  - **AC**: Payment recording working

- [x] ✅ 📝 Implement financial reports (10h) 🟡

  - GET /api/v1/finance/reports/revenue/daily
  - GET /api/v1/finance/reports/revenue/monthly
  - GET /api/v1/finance/reports/outstanding
  - GET /api/v1/finance/reports/student/:student_id/history
  - **AC**: Reports working

- [x] ✅ 📝 Implement overdue tracking (5h) 🟡

  - Scheduled job
  - Mark overdue invoices
  - **AC**: Overdue tracking working

- [x] ✅ 📝 Unit tests for Finance (12h) ✅
  - [x] Repository tests
  - [x] Handler tests
  - [x] Calculation tests
  - [x] Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [x] ✅ 📝 Finance service integration tests (10h) ✅
  - Invoice generation flow
  - Payment recording flow
  - **AC**: Integration tests passing

### Documentation

- [x] ✅ 📝 Write Finance service README (3h) ✅
- [x] ✅ 📝 Generate Swagger docs (4h) ✅
- [x] ✅ 📝 Create Postman collection (3h) ✅

---

## Notification Service

### Core Setup

- [x] ✅ 📝 Create notification-service structure (3h) 🔴
- [x] ✅ 📝 Setup database connection (2h) 🔴
- [x] ✅ 📝 Create database migrations (3h) 🔴
  - Notification_templates table
  - Notifications table

### Notification Core

- [x] ✅ 📝 Implement NotificationTemplate entity (2h) 🔴
- [x] ✅ 📝 Implement Notification entity (2h) 🔴

- [x] ✅ 📝 Implement notification repositories (6h) 🔴

  - Template CRUD
  - Notification CRUD

- [x] ✅ 📝 Implement template management (8h) 🟡
  - POST /api/v1/notifications/templates
  - GET /api/v1/notifications/templates
  - PUT /api/v1/notifications/templates/:id
  - Variable replacement logic
  - **AC**: Templates working

### Email Service

- [x] ✅ 📝 Configure SMTP (3h) 🔴

  - SMTP settings
  - Connection testing
  - **AC**: Email connection working

- [x] ✅ 📝 Implement email sending (8h) 🔴

  - HTML templates
  - Send function
  - Error handling
  - **AC**: Emails sent successfully

- [x] ✅ 📝 Implement email queue (6h) 🟡
  - Queue emails
  - Process queue
  - Retry on failure
  - **AC**: Queue working

### WhatsApp Integration

- [x] ✅ 📝 Configure WhatsApp API (4h) 🔴

  - API credentials
  - Connection testing
  - **AC**: WhatsApp connection working

- [x] ✅ 📝 Implement WhatsApp sending (8h) 🔴

  - Text messages
  - Template messages
  - Error handling
  - **AC**: WhatsApp messages sent

- [x] ✅ 📝 Implement webhook handler (5h) 🟡
  - Receive status updates
  - Update notification status
  - **AC**: Webhook working

### Event-Driven Messaging

- [x] ✅ 📝 Setup RabbitMQ (4h) 🔴

  - RabbitMQ container
  - Connection configuration
  - **AC**: RabbitMQ running

- [x] ✅ 📝 Implement event publisher (6h) 🔴

  - Publish function
  - Event schema
  - **AC**: Events published

- [x] ✅ 📝 Implement event consumer (8h) 🔴

  - Subscribe to events
  - Process events
  - Send notifications
  - **AC**: Events consumed

- [x] ✅ 📝 Implement retry mechanism (5h) 🟡
  - Retry failed notifications (3 attempts)
  - Dead letter queue
  - **AC**: Retry working

### Notification Handlers

- [x] ✅ 📝 Implement notification handlers (8h) 🔴

  - POST /api/v1/notifications/send
  - POST /api/v1/notifications/send/bulk
  - GET /api/v1/notifications
  - GET /api/v1/notifications/:id
  - GET /api/v1/notifications/user/:user_id
  - **AC**: All endpoints working

- [x] ✅ 📝 Unit tests for Notification (10h) 🔴
  - Repository tests
  - Email tests
  - WhatsApp tests
  - Event tests
  - Coverage >70%
  - **AC**: Tests passing

### Integration Tests

- [x] ✅ 📝 Notification service integration tests (10h) 🔴
  - Email sending flow
  - WhatsApp sending flow
  - Event-driven flow
  - **AC**: Integration tests passing

### Documentation

- [x] ✅ 📝 Write Notification service README (3h) 🟡
- [x] ✅ 📝 Generate Swagger docs (4h) 🟡
- [x] ✅ 📝 Create Postman collection (3h) ✅

---

## API Gateway

### Core Setup

- [x] ✅ 📝 Create api-gateway structure (4h) 🔴
  - Initialize project
  - Choose gateway (Kong/Traefik/Custom)
  - **AC**: Gateway structure ready

### Gateway Configuration

- [x] ✅ 📝 Configure routing (8h) 🔴

  - Route to auth-service
  - Route to academic-service
  - Route to attendance-service
  - Route to assessment-service
  - Route to admission-service
  - Route to finance-service
  - Route to notification-service
  - **AC**: All routes working

- [x] ✅ 📝 Implement rate limiting (5h) 🔴

  - Global rate limiting
  - Per-service rate limiting
  - **AC**: Rate limiting active

- [x] ✅ 📝 Implement request logging (4h) 🟡

  - Log all requests
  - Request ID generation
  - **AC**: Logging working

- [x] ✅ 📝 Implement authentication (6h) 🔴

  - JWT validation
  - Forward user context
  - **AC**: Auth working

- [x] ✅ 📝 Setup load balancing (5h) 🟡

- Load balancing strategy
- Health check integration
- **AC**: Load balancing working

Status: Round-robin multi-upstream aktif (env `*_URLS`), dengan circuit breaker per upstream dan health aggregator `/api/v1/gateway/health` mendukung multi upstream.

- [x] ✅ Unit tests for Gateway (8h) 🔴
  - Routing tests
  - Rate limiting tests
  - Coverage >70%
  - **AC**: Tests passing

### Documentation

- [x] ✅ 📝 Write Gateway README (3h) ✅

  - Service overview
  - Routing rules
  - Configuration
  - **AC**: README complete

- [x] ✅ 📝 Gateway architecture diagram (2h) ✅

---

## Supporting Services

### File Service (Optional)

- [x] ✅ 📝 Create file-service structure (3h) 🟢
- [x] ✅ 📝 Implement file upload (8h) 🟢

  - POST /api/v1/files/upload
  - File validation
  - Store in object storage (MinIO/S3) (Local storage implemented)
  - **AC**: Upload working

- [x] ✅ 📝 Implement file download (4h) 🟢

  - GET /api/v1/files/:id
  - Signed URLs (Direct download implemented for local)
  - **AC**: Download working

- [x] ✅ 📝 Implement file deletion (3h) 🟢
  - DELETE /api/v1/files/:id
  - Soft delete (Hard delete implemented for local)
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

- [x] 📝 Implement connection pooling (3h) 🔴
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

- [x] 📝 Setup Gosec scanning (3h) 🔴

  - Configure Gosec
  - Add to CI
  - Fix issues
  - **AC**: Security scan passing

- [x] 📝 Setup Trivy scanning (3h) 🔴

  - Configure Trivy
  - Scan containers
  - Fix vulnerabilities
  - **AC**: No critical vulnerabilities

- [x] 📝 Implement security headers (2h) 🔴
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

- [ ] 🔄 📝 Achieve 70% code coverage (40h) 🔴
  - Write missing tests
  - Fix failing tests
  - **AC**: Coverage >70%

### Integration Tests

- [x] ✅ 📝 Write integration tests (32h) 🔴
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

- [ ] 🔄 Complete security audit (16h) 🔴

  - [x] Security review (Gosec scan)
  - [x] Fix vulnerabilities (Permissions, Unhandled errors)
  - [ ] Manual review of file inclusion warnings
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

| Service                    | Estimated Hours |
| -------------------------- | --------------- |
| Global Infrastructure      | 150h            |
| Auth Service               | 180h            |
| Academic Core Service      | 220h            |
| Attendance Service         | 90h             |
| Assessment Service         | 150h            |
| Admission Service          | 120h            |
| Finance Service            | 120h            |
| Notification Service       | 110h            |
| API Gateway                | 42h             |
| Supporting Services        | 50h             |
| Performance & Optimization | 60h             |
| Security Hardening         | 50h             |
| Deployment & DevOps        | 100h            |
| Documentation              | 70h             |
| Testing                    | 112h            |
| Production Readiness       | 64h             |
| **TOTAL**                  | **1,688h**      |

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

- [x] Setup monorepo structure
- [x] Create docker-compose.yml
- [x] Create Makefile
- [x] Implement shared packages (config, database, logger)
- [x] Setup CI/CD pipeline basics
- [x] Create auth-service structure
- [x] Setup first database migrations

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
