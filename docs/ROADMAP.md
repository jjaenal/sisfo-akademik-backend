# ROADMAP — Backend Sistem Akademik (Microservices)

Roadmap komprehensif untuk pengembangan backend sistem akademik enterprise-grade dengan pendekatan microservices, scalable, secure, dan production-ready.

---

## 🎯 Vision & Goals

### Vision
Membangun platform akademik microservices yang:
- **Scalable**: Handle 100K+ concurrent users
- **Secure**: Enterprise-grade security standards
- **Flexible**: Multi-tenant dengan konfigurasi dinamis
- **Reliable**: 99.9% uptime SLA
- **Fast**: Sub-100ms average API response time

### Success Metrics
- [ ] Deployment time < 15 minutes
- [ ] API response P95 < 500ms
- [ ] Code coverage > 70%
- [ ] Zero-downtime deployments
- [ ] < 5 critical bugs per sprint

---

## 📅 Timeline Overview

| Phase | Duration | Status | Target Date |
|-------|----------|--------|-------------|
| Phase 0: Foundation | 2 weeks | 🔄 In Progress | Week 1-2 |
| Phase 1: Infrastructure | 3 weeks | ⏳ Planned | Week 3-5 |
| Phase 2: Identity & Access | 3 weeks | ⏳ Planned | Week 6-8 |
| Phase 3: Academic Core | 4 weeks | ⏳ Planned | Week 9-12 |
| Phase 4: Operations | 4 weeks | ⏳ Planned | Week 13-16 |
| Phase 5: Admission & Finance | 4 weeks | ⏳ Planned | Week 17-20 |
| Phase 6: Integration & Notification | 3 weeks | ⏳ Planned | Week 21-23 |
| Phase 7: Hardening & Production | 4 weeks | ⏳ Planned | Week 24-27 |

**Total Estimated Duration**: 27 weeks (6-7 bulan)

---

## Phase 0 — Foundation & Planning (Week 1-2)

**Goal**: Establish project foundation, define boundaries, and prepare technical specifications.

### 0.1 Requirements Analysis

#### Stakeholder Mapping
- **Admin Sekolah**: User management, konfigurasi sistem, reporting
- **Kepala Sekolah**: Dashboard analytics, approval workflows
- **Guru**: Presensi, penilaian, jadwal mengajar
- **Siswa**: Lihat jadwal, nilai, presensi
- **Orang Tua**: Monitor perkembangan anak, notifikasi
- **Staff Keuangan**: Billing, payment tracking, rekonsiliasi
- **Staff PPDB**: Manage pendaftaran, seleksi, verifikasi

#### Functional Requirements
```
Core Features:
├── Multi-tenant architecture (per sekolah)
├── Dynamic role & permission management
├── School hierarchy (SD, SMP, SMA/SMK)
├── Academic year & semester management
├── Class & subject management
├── Teacher & student management
├── Attendance tracking (real-time)
├── Grading system (flexible KKM)
├── Report card generation
├── Admission management (PPDB)
├── Finance management (SPP, billing)
└── Notification system (email, WhatsApp)

Advanced Features:
├── Parent portal
├── Dashboard & analytics
├── Document management
├── Integration APIs (third-party)
├── Mobile app support
└── Offline-first capability
```

#### Non-Functional Requirements
- **Performance**: 100ms avg response time, 1000 req/s per service
- **Scalability**: Horizontal scaling ready
- **Availability**: 99.9% uptime (8.76 hours downtime/year)
- **Security**: OWASP Top 10 compliance
- **Compliance**: GDPR ready, data residency support

### 0.2 Domain Modeling & Boundaries

#### Service Decomposition
```
┌─────────────────────────────────────────────────────────────┐
│                         API Gateway                          │
│                    (Kong / Traefik / Custom)                 │
└─────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         ▼                    ▼                    ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│   Auth Service   │  │ Academic Service │  │ Attendance Svc   │
│                  │  │                  │  │                  │
│ - User CRUD      │  │ - School CRUD    │  │ - Student attend │
│ - Authentication │  │ - Academic Year  │  │ - Teacher attend │
│ - Authorization  │  │ - Class Mgmt     │  │ - Location check │
│ - Audit Logging  │  │ - Subject Mgmt   │  │ - Validation     │
└──────────────────┘  └──────────────────┘  └──────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│Assessment Service│  │ Admission Svc    │  │  Finance Service │
│                  │  │                  │  │                  │
│ - Grade input    │  │ - Registration   │  │ - Billing setup  │
│ - Grade rules    │  │ - Doc upload     │  │ - Payment track  │
│ - Report card    │  │ - Selection      │  │ - Invoice gen    │
│ - PDF export     │  │ - Verification   │  │ - Reconciliation │
└──────────────────┘  └──────────────────┘  └──────────────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              ▼
                   ┌──────────────────────┐
                   │  Notification Service │
                   │                       │
                   │  - Email gateway      │
                   │  - WhatsApp gateway   │
                   │  - Push notification  │
                   │  - Event processing   │
                   └──────────────────────┘

Supporting Services:
├── File Service (document & media storage)
├── Report Service (advanced reporting & analytics)
└── Integration Service (third-party APIs)
```

#### Domain Events
```
User Domain:
- UserCreated
- UserUpdated
- UserDeleted
- PasswordChanged
- RoleAssigned

Academic Domain:
- ClassCreated
- StudentEnrolled
- TeacherAssigned
- ScheduleUpdated
- AcademicYearStarted

Attendance Domain:
- AttendanceRecorded
- AttendanceValidated
- LateArrival
- AbsenceNotified

Assessment Domain:
- GradeSubmitted
- GradeApproved
- ReportCardGenerated
- ReportCardPublished

Finance Domain:
- InvoiceCreated
- PaymentReceived
- PaymentOverdue
- PaymentReminder

Admission Domain:
- ApplicationSubmitted
- ApplicationApproved
- ApplicationRejected
- DocumentVerified
```

### 0.3 API Contract Design (OpenAPI First)

#### API Gateway Routes
```yaml
/api/v1/auth/*          → Auth Service
/api/v1/schools/*       → Academic Service
/api/v1/classes/*       → Academic Service
/api/v1/subjects/*      → Academic Service
/api/v1/attendance/*    → Attendance Service
/api/v1/grades/*        → Assessment Service
/api/v1/reports/*       → Assessment Service
/api/v1/admissions/*    → Admission Service
/api/v1/finance/*       → Finance Service
/api/v1/notifications/* → Notification Service
/api/v1/files/*         → File Service
```

#### Sample OpenAPI Specs (to be created)
```
docs/
├── openapi/
│   ├── auth-service.yaml
│   ├── academic-service.yaml
│   ├── attendance-service.yaml
│   ├── assessment-service.yaml
│   ├── admission-service.yaml
│   ├── finance-service.yaml
│   └── notification-service.yaml
```

### 0.4 Multi-Tenancy Strategy

#### Tenant Isolation Model: **Shared Database with Row-Level Security**

**Rationale**: 
- Cost-effective untuk small-medium schools
- Easier maintenance
- Good performance dengan proper indexing
- Tenant ID dalam setiap query

**Schema Pattern**:
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,  -- Tenant isolation
    email VARCHAR(255) NOT NULL,
    ...
    CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
```

**Tenant Configuration**:
```json
{
  "tenant_id": "uuid",
  "school_name": "SMA Negeri 1",
  "level": "SMA",
  "features": {
    "attendance": true,
    "grading": true,
    "finance": true,
    "admission": false
  },
  "quotas": {
    "max_students": 1000,
    "max_teachers": 100,
    "storage_gb": 50
  }
}
```

### 0.5 Technology Stack Finalization

```
Backend Framework:
├── Language: Go 1.21+
├── Web Framework: Gin / Echo / Chi
├── Validation: go-playground/validator
└── Config: Viper

Databases:
├── Primary: PostgreSQL 15+ (ACID compliance)
├── Cache: Redis 7+ (session, rate limiting)
└── Message Queue: RabbitMQ / Kafka

Infrastructure:
├── Container: Docker
├── Orchestration: Kubernetes (production)
├── Service Mesh: Istio (optional, phase 7)
├── API Gateway: Kong / Traefik
└── Reverse Proxy: Nginx

Observability:
├── Logging: ELK Stack (Elasticsearch, Logstash, Kibana)
├── Metrics: Prometheus + Grafana
├── Tracing: Jaeger
└── APM: (optional) New Relic / DataDog

Development Tools:
├── Version Control: Git + GitHub/GitLab
├── CI/CD: GitHub Actions / GitLab CI
├── Code Quality: SonarQube
├── Dependency Management: Go Modules
└── Documentation: Swagger UI, Postman
```

### 0.6 Development Standards

**Repository Structure**: Monorepo with multiple services
```
academic-backend/
├── services/
│   ├── auth-service/
│   ├── academic-service/
│   ├── attendance-service/
│   └── ...
├── shared/
│   ├── pkg/           # Shared packages
│   ├── proto/         # gRPC definitions
│   └── events/        # Event schemas
├── infrastructure/
│   ├── docker/
│   ├── kubernetes/
│   └── terraform/
├── docs/
│   ├── adr/           # Architecture decisions
│   ├── openapi/       # API specifications
│   └── diagrams/      # Architecture diagrams
└── scripts/           # Automation scripts
```

### Deliverables
- [x] Requirements document
- [x] Domain model & service boundaries
- [x] OpenAPI specifications (draft)
- [x] Multi-tenancy strategy document
- [x] Technology stack decision (ADR)
- [x] Repository structure
- [ ] Project setup (repos, boards, access)

---

## Phase 1 — Core Infrastructure (Week 3-5)

**Goal**: Setup development environment, CI/CD pipeline, dan foundational infrastructure.

### 1.1 Repository & Project Setup

- [ ] Create monorepo structure
- [ ] Setup GitHub/GitLab organization
- [ ] Configure branch protection rules (main, develop)
- [ ] Setup project management board (Jira/GitHub Projects)
- [ ] Configure team access & permissions

### 1.2 Local Development Environment

```yaml
# docker-compose.yml for local development
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: academic_db
      POSTGRES_USER: dev_user
      POSTGRES_PASSWORD: dev_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # Management UI

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686" # UI
      - "14268:14268" # Collector

  api-gateway:
    build: ./services/api-gateway
    ports:
      - "8080:8080"
    depends_on:
      - auth-service
      - academic-service
```

- [ ] Create docker-compose.yml
- [ ] Setup Makefile dengan common commands
- [ ] Configure environment files (.env.example)
- [ ] Documentation untuk setup lokal

### 1.3 Centralized Configuration Management

**Config Structure**:
```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    Logging  LoggingConfig
}

type ServerConfig struct {
    Host         string        `env:"SERVER_HOST" default:"0.0.0.0"`
    Port         int           `env:"SERVER_PORT" default:"8080"`
    ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" default:"10s"`
    WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" default:"10s"`
}
```

- [ ] Config loading dengan Viper
- [ ] Environment variable validation
- [ ] Config hot-reload support
- [ ] Secrets management strategy (Vault integration)

### 1.4 Database Setup

#### PostgreSQL Configuration
```sql
-- High-availability setup
-- Master-slave replication (future phase)
-- Connection pooling: pgbouncer

-- Initial databases:
CREATE DATABASE auth_db;
CREATE DATABASE academic_db;
CREATE DATABASE attendance_db;
CREATE DATABASE assessment_db;
CREATE DATABASE admission_db;
CREATE DATABASE finance_db;
```

- [ ] PostgreSQL HA setup documentation
- [ ] Connection pooling configuration
- [ ] Database backup strategy
- [ ] Migration tools setup (golang-migrate)

#### Redis Configuration
```yaml
# Redis for:
- Session storage
- Token blacklist
- Rate limiting
- Cache layer
- Pub/sub messaging
```

- [ ] Redis cluster setup (future)
- [ ] Redis persistence configuration
- [ ] Cache invalidation strategy

### 1.5 CI/CD Pipeline Setup

**GitHub Actions Workflow**:
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3

  test:
    runs-on: ubuntu-latest
    steps:
      - name: Run tests
        run: make test-coverage
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  security-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Run Gosec
        run: gosec ./...
      - name: Run Trivy
        run: trivy fs .

  build:
    runs-on: ubuntu-latest
    steps:
      - name: Build Docker images
        run: make docker-build-all
      - name: Push to registry
        run: make docker-push

  deploy-staging:
    needs: [lint, test, security-scan, build]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    steps:
      - name: Deploy to staging
        run: make deploy-staging
```

- [ ] Setup CI pipeline
- [ ] Configure test coverage reporting (Codecov)
- [ ] Setup security scanning (Gosec, Trivy)
- [ ] Configure Docker registry
- [ ] Setup CD to staging environment

### 1.6 Shared Packages

**Common packages** untuk semua services:
```
shared/
├── pkg/
│   ├── config/         # Config loading
│   ├── database/       # DB connection & pooling
│   ├── redis/          # Redis client
│   ├── logger/         # Structured logging
│   ├── middleware/     # Common middleware
│   ├── errors/         # Error handling
│   ├── validator/      # Input validation
│   ├── jwt/            # JWT utilities
│   ├── httputil/       # HTTP helpers
│   └── testutil/       # Testing utilities
```

- [ ] Implement shared packages
- [ ] Unit tests untuk shared packages
- [ ] Documentation per package

### Deliverables
- [x] Monorepo structure created
- [ ] Docker Compose environment
- [ ] CI/CD pipeline configured
- [ ] Shared packages implemented
- [ ] Development documentation

### Success Criteria
- [ ] Developers dapat run full stack locally dengan 1 command
- [ ] CI pipeline passing untuk semua checks
- [ ] Database migrations running successfully
- [ ] All services dapat communicate via docker network

---

## Phase 2 — Identity & Access Management (Week 6-8)

**Goal**: Implement authentication, authorization, dan audit logging system.

### 2.1 Auth Service - Core Features

#### User Management
```go
// User CRUD operations
POST   /api/v1/users              # Create user
GET    /api/v1/users/:id          # Get user by ID
GET    /api/v1/users              # List users (paginated)
PUT    /api/v1/users/:id          # Update user
DELETE /api/v1/users/:id          # Soft delete user
PATCH  /api/v1/users/:id/activate # Activate/deactivate
```

**Features**:
- [ ] User registration
- [ ] Email verification
- [ ] Password hashing (bcrypt, cost=12)
- [ ] Profile management
- [ ] User search & filtering
- [ ] Bulk user import (CSV)

#### Authentication Flow
```go
POST /api/v1/auth/login           # Login (email + password)
POST /api/v1/auth/logout          # Logout (invalidate tokens)
POST /api/v1/auth/refresh         # Refresh access token
POST /api/v1/auth/forgot-password # Password reset request
POST /api/v1/auth/reset-password  # Reset password with token
POST /api/v1/auth/change-password # Change password (authenticated)
```

**Token Strategy**:
```
Access Token:  JWT, 15 minutes, stateless
Refresh Token: JWT, 7 days, stored in Redis (revocable)
```

**JWT Claims**:
```json
{
  "sub": "user_id",
  "tenant_id": "school_id",
  "email": "user@example.com",
  "roles": ["admin", "teacher"],
  "permissions": ["user:create", "grade:update"],
  "exp": 1234567890,
  "iat": 1234567890,
  "jti": "token_unique_id"
}
```

- [ ] Implement login flow
- [ ] JWT access & refresh token generation
- [ ] Token refresh endpoint
- [ ] Token revocation (blacklist in Redis)
- [ ] Password reset flow (email verification)
- [ ] Session management
- [ ] Failed login attempt tracking (account lockout)
- [ ] Multi-device session support

### 2.2 Role-Based Access Control (RBAC)

#### Dynamic RBAC Schema
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    resource VARCHAR(100) NOT NULL,  -- e.g., 'user', 'grade'
    action VARCHAR(50) NOT NULL,     -- e.g., 'create', 'read'
    description TEXT,
    UNIQUE(resource, action)
);

CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id),
    permission_id UUID REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id),
    role_id UUID REFERENCES roles(id),
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id),
    PRIMARY KEY (user_id, role_id)
);
```

**Pre-defined System Roles**:
```
Super Admin:     Full system access (cannot be deleted)
School Admin:    School-level management
Teacher:         Class & grade management
Student:         View own data only
Parent:          View children data only
Finance Staff:   Finance operations
Admission Staff: PPDB operations
```

**Permission Format**: `resource:action`
```
Examples:
- user:create, user:read, user:update, user:delete
- class:create, class:assign_teacher
- grade:create, grade:approve
- report:generate, report:export
```

**API Endpoints**:
```go
// Role Management
POST   /api/v1/roles
GET    /api/v1/roles
GET    /api/v1/roles/:id
PUT    /api/v1/roles/:id
DELETE /api/v1/roles/:id

// Permission Management
GET    /api/v1/permissions
POST   /api/v1/roles/:id/permissions  # Assign permissions
DELETE /api/v1/roles/:id/permissions/:permission_id

// User Role Assignment
POST   /api/v1/users/:id/roles
DELETE /api/v1/users/:id/roles/:role_id
GET    /api/v1/users/:id/permissions  # Get effective permissions
```

- [ ] RBAC database schema
- [ ] Role CRUD operations
- [ ] Permission management
- [ ] Role assignment to users
- [ ] Permission checking middleware
- [ ] Role hierarchy support (future)
- [ ] Context-aware authorization (tenant-based)

### 2.3 Audit Logging

**Audit Log Schema**:
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    user_id UUID,
    action VARCHAR(100) NOT NULL,      -- e.g., 'user.created'
    resource_type VARCHAR(50),         -- e.g., 'user'
    resource_id VARCHAR(255),
    old_values JSONB,                  -- Previous state
    new_values JSONB,                  -- New state
    ip_address INET,
    user_agent TEXT,
    status VARCHAR(20),                -- success, failed
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_audit_tenant (tenant_id),
    INDEX idx_audit_user (user_id),
    INDEX idx_audit_created (created_at)
);
```

**Audit Events**:
```
Authentication:
- auth.login_success
- auth.login_failed
- auth.logout
- auth.token_refreshed

User Management:
- user.created
- user.updated
- user.deleted
- user.password_changed
- user.role_assigned

Data Access:
- resource.read
- resource.created
- resource.updated
- resource.deleted
```

**Audit Middleware**:
```go
func AuditMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Log after request
        logAudit(AuditLog{
            UserID:    getUserID(c),
            Action:    c.Request.Method + " " + c.Request.URL.Path,
            IPAddress: c.ClientIP(),
            Status:    c.Writer.Status(),
            Duration:  time.Since(start),
        })
    }
}
```

- [ ] Audit log schema
- [ ] Audit middleware implementation
- [ ] Async audit logging (don't block requests)
- [ ] Audit log retention policy (90 days)
- [ ] Audit log search & filtering API
- [ ] Audit log export functionality

### 2.4 Security Enhancements

#### Rate Limiting
```go
// Rate limits per endpoint type
Authentication: 5 requests/minute per IP
Read Operations: 100 requests/minute per user
Write Operations: 30 requests/minute per user
```

- [ ] Rate limiting middleware (Redis-based)
- [ ] Rate limit headers (X-RateLimit-*)
- [ ] Rate limit exceeded response (429 Too Many Requests)
- [ ] Whitelist IPs (internal services)

#### Security Headers
```go
// Helmet.js equivalent for Go
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
Content-Security-Policy: default-src 'self'
```

- [ ] Security headers middleware
- [ ] CORS configuration
- [ ] Request size limiting
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (input sanitization)

#### Password Security
```go
Password Requirements:
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character
- Not in common password list
- Not similar to username/email
```

- [ ] Password validation rules
- [ ] Password strength meter (API endpoint)
- [ ] Password history (prevent reuse of last 5 passwords)
- [ ] Force password change on first login (optional)

### Deliverables
- [ ] Auth service fully implemented
- [ ] RBAC system operational
- [ ] Audit logging active
- [ ] Security measures in place
- [ ] Unit tests (>70% coverage)
- [ ] Integration tests
- [ ] API documentation (Swagger)
- [ ] Postman collection

### Success Criteria
- [ ] Users can register, login, logout
- [ ] JWT tokens working correctly
- [ ] RBAC protecting endpoints properly
- [ ] Audit logs capturing all actions
- [ ] Rate limiting preventing abuse
- [ ] All security tests passing

---

## Phase 3 — Academic Core (Week 9-12)

**Goal**: Implement core academic management features.

### 3.1 School Management

**School Entity**:
```sql
CREATE TABLE schools (
    id UUID PRIMARY KEY,
    tenant_id UUID UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    level VARCHAR(20) NOT NULL,  -- SD, SMP, SMA, SMK
    npsn VARCHAR(20) UNIQUE,     -- Nomor Pokok Sekolah Nasional
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    principal_name VARCHAR(255),
    logo_url VARCHAR(500),
    settings JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**API Endpoints**:
```go
POST   /api/v1/schools        # Create school
GET    /api/v1/schools        # List schools (super admin)
GET    /api/v1/schools/:id    # Get school details
PUT    /api/v1/schools/:id    # Update school
DELETE /api/v1/schools/:id    # Deactivate school
```

- [ ] School CRUD operations
- [ ] School settings management (JSONB config)
- [ ] School logo upload
- [ ] Multi-level support (SD/SMP/SMA/SMK)
- [ ] School activation/deactivation

### 3.2 Academic Year & Semester

**Schema**:
```sql
CREATE TABLE academic_years (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,        -- e.g., "2024/2025"
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,   -- Only one active per tenant
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

CREATE TABLE semesters (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    academic_year_id UUID REFERENCES academic_years(id),
    name VARCHAR(50) NOT NULL,         -- "Semester 1", "Semester 2"
    number INT NOT NULL,               -- 1 or 2
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(academic_year_id, number)
);
```

**Business Rules**:
- Only 1 active academic year per tenant
- Only 1 active semester per academic year
- Academic year cannot overlap
- Semester must be within academic year boundaries

**API Endpoints**:
```go
// Academic Year
POST   /api/v1/academic-years
GET    /api/v1/academic-years
GET    /api/v1/academic-years/:id
PUT    /api/v1/academic-years/:id
PATCH  /api/v1/academic-years/:id/activate
DELETE /api/v1/academic-years/:id

// Semester
POST   /api/v1/semesters
GET    /api/v1/semesters
GET    /api/v1/semesters/:id
PUT    /api/v1/semesters/:id
PATCH  /api/v1/semesters/:id/activate
```

- [ ] Academic year management
- [ ] Semester management
- [ ] Active year/semester validation
- [ ] Date overlap validation
- [ ] Automatic semester creation (when year is created)

### 3.3 Class Management

**Schema**:
```sql
CREATE TABLE classes (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    academic_year_id UUID REFERENCES academic_years(id),
    name VARCHAR(100) NOT NULL,        -- e.g., "X IPA 1"
    level VARCHAR(20) NOT NULL,        -- "10", "11", "12" or "7", "8", "9"
    major VARCHAR(50),                 -- IPA, IPS, etc (for SMA)
    homeroom_teacher_id UUID REFERENCES users(id),
    max_students INT DEFAULT 40,
    room_number VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, academic_year_id, name)
);

CREATE TABLE class_students (
    id UUID PRIMARY KEY,
    class_id UUID REFERENCES classes(id),
    student_id UUID REFERENCES users(id),
    enrolled_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',  -- active, transferred, graduated
    UNIQUE(class_id, student_id)
);
```

**API Endpoints**:
```go
// Class Management
POST   /api/v1/classes
GET    /api/v1/classes
GET    /api/v1/classes/:id
PUT    /api/v1/classes/:id
DELETE /api/v1/classes/:id

// Student Enrollment
POST   /api/v1/classes/:id/students        # Enroll student
DELETE /api/v1/classes/:id/students/:student_id
GET    /api/v1/classes/:id/students        # List class students
POST   /api/v1/classes/:id/students/bulk   # Bulk enrollment
```

- [ ] Class CRUD operations
- [ ] Homeroom teacher assignment
- [ ] Student enrollment
- [ ] Bulk student enrollment
- [ ] Student transfer between classes
- [ ] Class capacity management
- [ ] Student graduation process

### 3.4 Subject Management

**Schema**:
```sql
CREATE TABLE subjects (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50),              -- Wajib, Peminatan, Mulok
    credit_hours INT,                  -- Jam per minggu
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE TABLE class_subjects (
    id UUID PRIMARY KEY,
    class_id UUID REFERENCES classes(id),
    subject_id UUID REFERENCES subjects(id),
    teacher_id UUID REFERENCES users(id),
    academic_year_id UUID REFERENCES academic_years(id),
    semester_id UUID REFERENCES semesters(id),
    UNIQUE(class_id, subject_id, semester_id)
);
```

**API Endpoints**:
```go
// Subject Management
POST   /api/v1/subjects
GET    /api/v1/subjects
GET    /api/v1/subjects/:id
PUT    /api/v1/subjects/:id
DELETE /api/v1/subjects/:id

// Class Subject Assignment
POST   /api/v1/classes/:id/subjects        # Assign subject to class
PUT    /api/v1/classes/:id/subjects/:subject_id/teacher  # Assign teacher
GET    /api/v1/classes/:id/subjects
DELETE /api/v1/classes/:id/subjects/:subject_id
```

- [ ] Subject CRUD operations
- [ ] Subject categorization
- [ ] Subject assignment to class
- [ ] Teacher assignment to subject
- [ ] Subject prerequisite management

### 3.5 Curriculum Management (Dynamic)

**Schema**:
```sql
CREATE TABLE curricula (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,        -- e.g., "Kurikulum Merdeka 2024"
    level VARCHAR(20) NOT NULL,        -- SD, SMP, SMA, SMK
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE curriculum_subjects (
    id UUID PRIMARY KEY,
    curriculum_id UUID REFERENCES curricula(id),
    subject_id UUID REFERENCES subjects(id),
    grade_level VARCHAR(20),           -- "10", "11", "12"
    is_mandatory BOOLEAN DEFAULT TRUE,
    credit_hours INT,
    UNIQUE(curriculum_id, subject_id, grade_level)
);

CREATE TABLE grading_rules (
    id UUID PRIMARY KEY,
    curriculum_id UUID REFERENCES curricula(id),
    subject_id UUID REFERENCES subjects(id),
    kkm INT NOT NULL,                  -- Kriteria Ketuntasan Minimal
    grade_components JSONB NOT NULL,   -- Tugas, UH, UTS, UAS weights
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Grading Components Example**:
```json
{
  "components": [
    { "name": "Tugas", "weight": 20 },
    { "name": "Ulangan Harian", "weight": 30 },
    { "name": "UTS", "weight": 20 },
    { "name": "UAS", "weight": 30 }
  ]
}
```

**API Endpoints**:
```go
// Curriculum
POST   /api/v1/curricula
GET    /api/v1/curricula
GET    /api/v1/curricula/:id
PUT    /api/v1/curricula/:id

// Curriculum Subjects
POST   /api/v1/curricula/:id/subjects
GET    /api/v1/curricula/:id/subjects

// Grading Rules
POST   /api/v1/curricula/:id/grading-rules
GET    /api/v1/curricula/:id/grading-rules
PUT    /api/v1/grading-rules/:id
```

- [ ] Curriculum CRUD operations
- [ ] Subject mapping to curriculum
- [ ] Grading rules configuration (KKM, weights)
- [ ] Multi-curriculum support per tenant
- [ ] Curriculum versioning

### 3.6 Schedule Management

**Schema**:
```sql
CREATE TABLE schedules (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    class_subject_id UUID REFERENCES class_subjects(id),
    day_of_week INT NOT NULL,          -- 1=Monday, 7=Sunday
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    room_number VARCHAR(20),
    academic_year_id UUID REFERENCES academic_years(id),
    semester_id UUID REFERENCES semesters(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, class_subject_id, day_of_week, start_time)
);
```

**API Endpoints**:
```go
POST   /api/v1/schedules
GET    /api/v1/schedules              # Filter by class, teacher, day
PUT    /api/v1/schedules/:id
DELETE /api/v1/schedules/:id
GET    /api/v1/schedules/class/:class_id/weekly
GET    /api/v1/schedules/teacher/:teacher_id/weekly
```

**Validations**:
- No overlapping schedules for same class
- No overlapping schedules for same teacher
- No overlapping room usage
- Schedule within semester dates

- [ ] Schedule CRUD operations
- [ ] Weekly schedule view (class, teacher)
- [ ] Schedule conflict detection
- [ ] Bulk schedule creation
- [ ] Schedule template system
- [ ] Room conflict management

### Deliverables
- [ ] Academic service fully implemented
- [ ] All CRUD operations functional
- [ ] Data validation rules enforced
- [ ] Unit tests (>70% coverage)
- [ ] Integration tests
- [ ] API documentation
- [ ] Data migration scripts

### Success Criteria
- [ ] Schools can be created and configured
- [ ] Academic years and semesters managed
- [ ] Classes with students enrolled
- [ ] Subjects assigned to classes
- [ ] Curriculum configured dynamically
- [ ] Weekly schedules generated
- [ ] No schedule conflicts

---

## Phase 4 — Academic Operations (Week 13-16)

**Goal**: Implement attendance tracking, grading, and report card generation.

### 4.1 Attendance Service

**Schema**:
```sql
CREATE TABLE student_attendance (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    student_id UUID REFERENCES users(id),
    class_id UUID REFERENCES classes(id),
    schedule_id UUID REFERENCES schedules(id),
    date DATE NOT NULL,
    status VARCHAR(20) NOT NULL,       -- present, absent, late, excused, sick
    check_in_time TIMESTAMP,
    notes TEXT,
    recorded_by UUID REFERENCES users(id),
    location_lat DECIMAL(10, 8),       -- GPS validation
    location_lng DECIMAL(11, 8),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, schedule_id, date)
);

CREATE TABLE teacher_attendance (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    teacher_id UUID REFERENCES users(id),
    date DATE NOT NULL,
    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    status VARCHAR(20) NOT NULL,
    notes TEXT,
    location_lat DECIMAL(10, 8),
    location_lng DECIMAL(11, 8),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(teacher_id, date)
);
```

**API Endpoints**:
```go
// Student Attendance
POST   /api/v1/attendance/students
POST   /api/v1/attendance/students/bulk     # Bulk check-in
GET    /api/v1/attendance/students          # Filter by class, date
PUT    /api/v1/attendance/students/:id
GET    /api/v1/attendance/students/:student_id/summary

// Teacher Attendance
POST   /api/v1/attendance/teachers/check-in
POST   /api/v1/attendance/teachers/check-out
GET    /api/v1/attendance/teachers
GET    /api/v1/attendance/teachers/:teacher_id/summary

// Reports
GET    /api/v1/attendance/reports/daily
GET    /api/v1/attendance/reports/monthly
GET    /api/v1/attendance/reports/class/:class_id
```

**Features**:
- [ ] Student attendance recording
- [ ] Bulk attendance input (full class)
- [ ] Teacher attendance tracking
- [ ] GPS location validation
- [ ] Attendance status management (present, absent, late, excused)
- [ ] Attendance modification with audit trail
- [ ] Attendance summary per student
- [ ] Attendance reports (daily, monthly)
- [ ] Absent notification trigger

### 4.2 Assessment Service - Grading

**Schema**:
```sql
CREATE TABLE grade_categories (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,        -- Tugas, UH, UTS, UAS
    weight DECIMAL(5, 2) NOT NULL,     -- Percentage
    class_subject_id UUID REFERENCES class_subjects(id),
    UNIQUE(class_subject_id, name)
);

CREATE TABLE assessments (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    class_subject_id UUID REFERENCES class_subjects(id),
    grade_category_id UUID REFERENCES grade_categories(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    max_score INT NOT NULL DEFAULT 100,
    assessment_date DATE NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE grades (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    assessment_id UUID REFERENCES assessments(id),
    student_id UUID REFERENCES users(id),
    score DECIMAL(5, 2) NOT NULL,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'draft',  -- draft, submitted, approved
    graded_by UUID REFERENCES users(id),
    graded_at TIMESTAMP,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(assessment_id, student_id)
);
```

**Grade Calculation**:
```
Final Score = Σ (Category Score × Category Weight)

Example:
Tugas (20%):    85 × 0.20 = 17
UH (30%):       90 × 0.30 = 27
UTS (20%):      80 × 0.20 = 16
UAS (30%):      88 × 0.30 = 26.4
Final Score:    86.4
```

**API Endpoints**:
```go
// Grade Categories
POST   /api/v1/grade-categories
GET    /api/v1/grade-categories
PUT    /api/v1/grade-categories/:id

// Assessments
POST   /api/v1/assessments
GET    /api/v1/assessments
GET    /api/v1/assessments/:id
PUT    /api/v1/assessments/:id
DELETE /api/v1/assessments/:id

// Grading
POST   /api/v1/grades
POST   /api/v1/grades/bulk                # Bulk grade input
GET    /api/v1/grades/assessment/:assessment_id
PUT    /api/v1/grades/:id
PATCH  /api/v1/grades/:id/approve

// Student Grade Summary
GET    /api/v1/grades/student/:student_id/semester/:semester_id
GET    /api/v1/grades/student/:student_id/subject/:subject_id
```

**Validation Rules**:
- Score must be between 0 and max_score
- Only teacher of subject can grade
- Grades can be modified before approval
- Grade modification after approval requires audit log
- KKM validation per subject

- [ ] Grade category setup per subject
- [ ] Assessment creation
- [ ] Grade input (single & bulk)
- [ ] Grade calculation engine
- [ ] Grade approval workflow
- [ ] Grade modification with audit
- [ ] Student grade summary
- [ ] Grade statistics (class average, highest, lowest)

### 4.3 Report Card Generation

**Schema**:
```sql
CREATE TABLE report_cards (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    student_id UUID REFERENCES users(id),
    class_id UUID REFERENCES classes(id),
    semester_id UUID REFERENCES semesters(id),
    academic_year_id UUID REFERENCES academic_years(id),
    status VARCHAR(20) DEFAULT 'draft',     -- draft, generated, published
    generated_at TIMESTAMP,
    generated_by UUID REFERENCES users(id),
    published_at TIMESTAMP,
    file_url VARCHAR(500),                  -- PDF URL
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, semester_id)
);

CREATE TABLE report_card_details (
    id UUID PRIMARY KEY,
    report_card_id UUID REFERENCES report_cards(id),
    subject_id UUID REFERENCES subjects(id),
    final_score DECIMAL(5, 2),
    grade_letter VARCHAR(2),               -- A, B, C, D, E
    predicate VARCHAR(50),                 -- Sangat Baik, Baik, Cukup, Kurang
    kkm INT,
    attendance_summary JSONB,              -- present, absent, late counts
    notes TEXT
);
```

**Report Card Components**:
```
Header:
- School info
- Student info
- Academic year & semester

Academic Performance:
- Subject grades
- Grade letter & predicate
- KKM comparison
- Class rank (optional)

Attendance Summary:
- Present days
- Absent days
- Late days
- Sick days
- Permission days

Teacher Notes:
- Homeroom teacher comment
- Principal signature

Footer:
- Print date
- Barcode/QR code for verification
```

**API Endpoints**:
```go
// Report Card Generation
POST   /api/v1/report-cards/generate/:student_id/:semester_id
POST   /api/v1/report-cards/generate/class/:class_id/:semester_id  # Bulk
GET    /api/v1/report-cards/:id
GET    /api/v1/report-cards/student/:student_id
PATCH  /api/v1/report-cards/:id/publish

// PDF Export
GET    /api/v1/report-cards/:id/pdf
GET    /api/v1/report-cards/:id/download
```

- [ ] Report card data aggregation
- [ ] Grade letter & predicate calculation
- [ ] Attendance summary calculation
- [ ] Class ranking calculation (optional)
- [ ] Report card template design
- [ ] PDF generation (using library: wkhtmltopdf, gotenberg, or chromedp)
- [ ] Bulk report card generation
- [ ] Report card publishing workflow
- [ ] Report card revision tracking

### 4.4 Report Templates & Customization

**Template Variables**:
```
{{school_name}}
{{school_address}}
{{student_name}}
{{student_nis}}
{{class_name}}
{{semester_name}}
{{subjects}}           # Loop
{{attendance_summary}}
{{homeroom_comment}}
{{principal_name}}
{{print_date}}
```

- [ ] HTML template system
- [ ] Template customization per tenant
- [ ] Preview before generation
- [ ] Template versioning
- [ ] Support for school logo
- [ ] Support for signatures

### Deliverables
- [ ] Attendance service operational
- [ ] Grading system functional
- [ ] Report card generation working
- [ ] PDF export capability
- [ ] Unit tests (>70% coverage)
- [ ] Integration tests
- [ ] API documentation

### Success Criteria
- [ ] Teachers can record attendance
- [ ] Bulk attendance input works
- [ ] Grades can be entered and calculated
- [ ] Report cards generate correctly
- [ ] PDFs render properly
- [ ] Bulk operations perform well (< 5 seconds for 40 students)

---

## Phase 5 — Admission & Finance (Week 17-20)

**Goal**: Implement PPDB (admission) and finance management systems.

### 5.1 PPDB (Admission) Service

**Admission Flow**:
```
1. Registration (by prospective student/parent)
2. Document Upload
3. Payment (registration fee)
4. Selection Process
5. Interview/Test (optional)
6. Announcement
7. Re-registration (if accepted)
8. Student Account Creation
```

**Schema**:
```sql
CREATE TABLE admission_periods (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    academic_year_id UUID REFERENCES academic_years(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    registration_fee DECIMAL(10, 2),
    max_applicants INT,
    status VARCHAR(20) DEFAULT 'open',  -- open, closed, selection, announced
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE applications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    admission_period_id UUID REFERENCES admission_periods(id),
    application_number VARCHAR(50) UNIQUE NOT NULL,
    
    -- Personal Info
    full_name VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    gender VARCHAR(10),
    birth_place VARCHAR(100),
    birth_date DATE,
    religion VARCHAR(50),
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    
    -- Parent Info
    parent_name VARCHAR(255),
    parent_phone VARCHAR(20),
    parent_email VARCHAR(255),
    parent_occupation VARCHAR(100),
    
    -- Previous School
    previous_school VARCHAR(255),
    previous_school_address TEXT,
    graduation_year INT,
    
    -- Selection Data
    test_score DECIMAL(5, 2),
    interview_score DECIMAL(5, 2),
    final_score DECIMAL(5, 2),
    
    -- Status
    status VARCHAR(20) DEFAULT 'submitted',  -- submitted, verified, accepted, rejected, registered
    payment_status VARCHAR(20) DEFAULT 'unpaid',
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE application_documents (
    id UUID PRIMARY KEY,
    application_id UUID REFERENCES applications(id),
    document_type VARCHAR(50) NOT NULL,  -- birth_cert, family_card, photo, etc
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255),
    file_size INT,
    verified BOOLEAN DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**API Endpoints**:
```go
// Admission Period Management (Admin)
POST   /api/v1/admission/periods
GET    /api/v1/admission/periods
GET    /api/v1/admission/periods/:id
PUT    /api/v1/admission/periods/:id
PATCH  /api/v1/admission/periods/:id/close

// Public Application
GET    /api/v1/admission/public/periods      # List open periods
POST   /api/v1/admission/applications        # Submit application
GET    /api/v1/admission/applications/:number/status

// Application Management (Staff)
GET    /api/v1/admission/applications
GET    /api/v1/admission/applications/:id
PUT    /api/v1/admission/applications/:id
PATCH  /api/v1/admission/applications/:id/verify
PATCH  /api/v1/admission/applications/:id/accept
PATCH  /api/v1/admission/applications/:id/reject

// Document Management
POST   /api/v1/admission/applications/:id/documents
GET    /api/v1/admission/applications/:id/documents
DELETE /api/v1/admission/documents/:id
PATCH  /api/v1/admission/documents/:id/verify

// Selection Process
POST   /api/v1/admission/applications/:id/test-score
POST   /api/v1/admission/applications/:id/interview-score
POST   /api/v1/admission/periods/:id/calculate-final-scores
POST   /api/v1/admission/periods/:id/announce

// Student Registration
POST   /api/v1/admission/applications/:id/register  # Convert to student account
```

**Features**:
- [ ] Admission period management
- [ ] Public registration form
- [ ] Document upload & verification
- [ ] Application status tracking
- [ ] Test score input
- [ ] Interview score input
- [ ] Final score calculation
- [ ] Acceptance/rejection process
- [ ] Student account creation from application
- [ ] Applicant notification (email/WhatsApp)
- [ ] Application reports & statistics

### 5.2 Finance Service

**Finance Components**:
```
1. Billing Setup (SPP configuration)
2. Invoice Generation
3. Payment Recording
4. Payment Reminder
5. Financial Reports
```

**Schema**:
```sql
CREATE TABLE billing_configurations (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    billing_type VARCHAR(50) NOT NULL,  -- spp, admission, exam, etc
    frequency VARCHAR(20),               -- monthly, yearly, one-time
    class_level VARCHAR(20),            -- "10", "11", "12" or "all"
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE invoices (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    student_id UUID REFERENCES users(id),
    billing_config_id UUID REFERENCES billing_configurations(id),
    amount DECIMAL(10, 2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'unpaid',  -- unpaid, paid, overdue, cancelled
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE payments (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    invoice_id UUID REFERENCES invoices(id),
    payment_number VARCHAR(50) UNIQUE NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50),          -- cash, bank_transfer, credit_card
    payment_date DATE NOT NULL,
    reference_number VARCHAR(100),       -- Bank ref, receipt number
    received_by UUID REFERENCES users(id),
    notes TEXT,
    receipt_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW()
);
```

**API Endpoints**:
```go
// Billing Configuration
POST   /api/v1/finance/billing-configs
GET    /api/v1/finance/billing-configs
GET    /api/v1/finance/billing-configs/:id
PUT    /api/v1/finance/billing-configs/:id
DELETE /api/v1/finance/billing-configs/:id

// Invoice Management
POST   /api/v1/finance/invoices/generate           # Manual generation
POST   /api/v1/finance/invoices/generate/bulk      # Bulk for class
POST   /api/v1/finance/invoices/generate/auto      # Auto-generate monthly
GET    /api/v1/finance/invoices
GET    /api/v1/finance/invoices/:id
PUT    /api/v1/finance/invoices/:id
DELETE /api/v1/finance/invoices/:id

// Student Invoice View
GET    /api/v1/finance/invoices/student/:student_id
GET    /api/v1/finance/invoices/student/:student_id/outstanding

// Payment Recording
POST   /api/v1/finance/payments
GET    /api/v1/finance/payments
GET    /api/v1/finance/payments/:id
GET    /api/v1/finance/payments/invoice/:invoice_id

// Reports
GET    /api/v1/finance/reports/revenue/daily
GET    /api/v1/finance/reports/revenue/monthly
GET    /api/v1/finance/reports/outstanding
GET    /api/v1/finance/reports/student/:student_id/history
```

**Features**:
- [ ] Billing configuration (SPP, fees)
- [ ] Invoice auto-generation (monthly SPP)
- [ ] Manual invoice generation
- [ ] Bulk invoice generation per class
- [ ] Payment recording
- [ ] Payment receipt generation
- [ ] Overdue invoice tracking
- [ ] Payment reminder (automated)
- [ ] Student payment history
- [ ] Financial reports (revenue, outstanding)
- [ ] Reconciliation interface

### 5.3 Payment Integration (Optional - Future)

**Payment Gateway Options**:
- Midtrans
- Xendit
- DOKU
- Manual bank transfer

**Payment Flow**:
```
1. Student views invoice
2. Click "Pay Now"
3. Redirect to payment gateway
4. Complete payment
5. Payment notification callback
6. Update invoice status
7. Send receipt email
```

- [ ] Payment gateway integration preparation
- [ ] Webhook handler for payment notifications
- [ ] Payment verification
- [ ] Refund handling (if needed)

### Deliverables
- [ ] Admission service operational
- [ ] Finance service operational
- [ ] Document upload working
- [ ] Invoice generation automated
- [ ] Payment recording functional
- [ ] Unit tests (>70% coverage)
- [ ] Integration tests
- [ ] API documentation

### Success Criteria
- [ ] Applicants can register online
- [ ] Documents can be uploaded
- [ ] Staff can manage applications
- [ ] Invoices generate correctly
- [ ] Payments can be recorded
- [ ] Financial reports accurate
- [ ] System handles 1000+ applications

---

## Phase 6 — Notification & Integration (Week 21-23)

**Goal**: Implement notification system dan external integrations.

### 6.1 Notification Service

**Notification Channels**:
- Email (SMTP)
- WhatsApp (API - Fonnte, Wablas, Twilio)
- Push Notification (Firebase, future)
- SMS (optional)

**Schema**:
```sql
CREATE TABLE notification_templates (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    channel VARCHAR(20) NOT NULL,        -- email, whatsapp, push, sms
    event_type VARCHAR(100) NOT NULL,    -- user.created, grade.published
    subject VARCHAR(255),                -- For email
    body TEXT NOT NULL,
    variables JSONB,                     -- Template variables
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, event_type, channel)
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    recipient_id UUID REFERENCES users(id),
    channel VARCHAR(20) NOT NULL,
    event_type VARCHAR(100),
    subject VARCHAR(255),
    body TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',  -- pending, sent, failed, retry
    sent_at TIMESTAMP,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    metadata JSONB,                        -- Additional data
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Notification Events**:
```
User Management:
- user.created
- user.password_reset
- user.activated

Academic:
- grade.published
- report_card.available
- class.schedule_changed

Attendance:
- attendance.absent_notification
- attendance.late_notification

Finance:
- invoice.created
- payment.received
- invoice.overdue_reminder

Admission:
- application.submitted
- application.accepted
- application.rejected
```

**API Endpoints**:
```go
// Template Management
POST   /api/v1/notifications/templates
GET    /api/v1/notifications/templates
GET    /api/v1/notifications/templates/:id
PUT    /api/v1/notifications/templates/:id
DELETE /api/v1/notifications/templates/:id

// Send Notification
POST   /api/v1/notifications/send          # Direct send
POST   /api/v1/notifications/send/bulk     # Bulk send
GET    /api/v1/notifications
GET    /api/v1/notifications/:id

// Notification History
GET    /api/v1/notifications/user/:user_id
GET    /api/v1/notifications/stats
```

**Features**:
- [ ] Email notification (SMTP)
- [ ] WhatsApp notification (API)
- [ ] Template management
- [ ] Template variable replacement
- [ ] Notification queuing
- [ ] Retry mechanism (3 attempts)
- [ ] Failed notification tracking
- [ ] Notification history
- [ ] Unsubscribe functionality
- [ ] Notification preferences per user

### 6.2 Email Service

**Email Provider**: SMTP (Gmail, SendGrid, AWS SES)

**Email Templates**:
```html
<!DOCTYPE html>
<html>
<head>
    <style>/* Email CSS */</style>
</head>
<body>
    <div class="container">
        <h1>{{title}}</h1>
        <p>{{message}}</p>
        <a href="{{action_url}}" class="button">{{action_text}}</a>
    </div>
</body>
</html>
```

- [ ] SMTP configuration
- [ ] HTML email templates
- [ ] Email sending service
- [ ] Attachment support
- [ ] Email tracking (opened, clicked)
- [ ] Bulk email sending (batch processing)

### 6.3 WhatsApp Integration

**WhatsApp API Options**:
- Fonnte.com
- Wablas.com
- Twilio WhatsApp Business API
- WhatsApp Cloud API (Official)

**Message Types**:
```
Text Message:
"Halo {{parent_name}}, anak Anda {{student_name}} tidak hadir pada {{date}}."

Template Message (with buttons):
"Invoice SPP {{month}} telah terbit. Total: Rp {{amount}}. 
Jatuh tempo: {{due_date}}."
[Button: Lihat Invoice]
```

- [ ] WhatsApp API integration
- [ ] Text message sending
- [ ] Template message support
- [ ] Message status tracking
- [ ] Webhook for message status

### 6.4 Event-Driven Messaging

**Message Broker**: RabbitMQ

**Event Publishing**:
```go
// Publish event when grade is published
func PublishGradePublished(studentID, subjectID string, grade float64) {
    event := Event{
        Type: "grade.published",
        Timestamp: time.Now(),
        Data: map[string]interface{}{
            "student_id": studentID,
            "subject_id": subjectID,
            "grade": grade,
        },
    }
    messageQueue.Publish("academic.events", event)
}
```

**Event Consuming**:
```go
// Notification service listens to events
func (s *NotificationService) ConsumeEvents() {
    messageQueue.Subscribe("academic.events", func(event Event) {
        switch event.Type {
        case "grade.published":
            s.SendGradeNotification(event.Data)
        case "attendance.absent":
            s.SendAbsentNotification(event.Data)
        }
    })
}
```

- [ ] RabbitMQ setup
- [ ] Event schema definition
- [ ] Event publisher per service
- [ ] Event consumer in notification service
- [ ] Dead letter queue handling
- [ ] Event replay mechanism

### 6.5 Integration APIs (Third-Party)

**Planned Integrations**:
```
Authentication:
- Google OAuth
- Microsoft OAuth
- SSO (LDAP/Active Directory)

Payment:
- Midtrans
- Xendit
- Manual bank transfer verification

Storage:
- AWS S3 / MinIO (file storage)
- Google Drive (backup)

Analytics:
- Google Analytics
- Custom dashboard
```

- [ ] OAuth implementation
- [ ] Payment gateway integration
- [ ] File storage service
- [ ] Third-party API client library
- [ ] API rate limiting handling
- [ ] Webhook security (signature verification)

### Deliverables
- [ ] Notification service operational
- [ ] Email sending functional
- [ ] WhatsApp integration working
- [ ] Event-driven architecture implemented
- [ ] Unit tests (>70% coverage)
- [ ] Integration tests
- [ ] API documentation

### Success Criteria
- [ ] Emails sent successfully
- [ ] WhatsApp messages delivered
- [ ] Events published and consumed
- [ ] Notifications retry on failure
- [ ] System handles 10K+ notifications/day

---

## Phase 7 — Hardening & Production Readiness (Week 24-27)

**Goal**: Optimize, secure, and prepare system for production deployment.

### 7.1 Performance Optimization

#### Caching Strategy
```
Redis Caching:
- User sessions (TTL: 7 days)
- JWT blacklist (TTL: token expiry)
- Frequently accessed data (school info, roles)
- API response caching (read-heavy endpoints)
- Rate limiting counters

Cache Invalidation:
- On data update/delete
- Manual cache clear endpoint
- TTL-based expiration
```

- [ ] Implement Redis caching
- [ ] Cache warming strategy
- [ ] Cache invalidation logic
- [ ] Cache hit rate monitoring

#### Database Optimization
```sql
-- Index optimization
CREATE INDEX CONCURRENTLY idx_users_tenant_email ON users(tenant_id, email);
CREATE INDEX idx_classes_tenant_year ON classes(tenant_id, academic_year_id);
CREATE INDEX idx_grades_student ON grades(student_id, assessment_id);
CREATE INDEX idx_attendance_date ON student_attendance(date, student_id);

-- Query optimization
EXPLAIN ANALYZE SELECT ...;

-- Partitioning (for large tables)
CREATE TABLE audit_logs_2025 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
```

- [ ] Database indexing
- [ ] Query optimization
- [ ] Connection pooling tuning
- [ ] Table partitioning (audit logs, attendance)
- [ ] Read replica setup (for production)

#### API Optimization
```
Techniques:
- Response pagination (default 20, max 100)
- Field selection (?fields=id,name,email)
- Compression (gzip, brotli)
- ETags for caching
- Batch endpoints (reduce round trips)
- GraphQL consideration (future)
```

- [ ] Pagination implementation
- [ ] Response compression
- [ ] ETag support
- [ ] Batch API endpoints
- [ ] API response time optimization (<100ms P50)

### 7.2 Security Hardening

#### Field-Level Encryption
```go
// Encrypt sensitive fields
type User struct {
    Email           string `encrypted:"true"`
    Phone           string `encrypted:"true"`
    TaxID           string `encrypted:"true"`
    BankAccount     string `encrypted:"true"`
}

// Encryption: AES-256-GCM
```

- [ ] Field-level encryption for PII
- [ ] Key rotation strategy
- [ ] Encryption key management (Vault)
- [ ] Data masking in logs

#### Advanced Rate Limiting
```go
// Multi-tier rate limiting
Auth endpoints:      5 req/min per IP
Read operations:     100 req/min per user
Write operations:    30 req/min per user
Admin operations:    50 req/min per admin
Global rate limit:   10K req/min per tenant
```

- [ ] Tiered rate limiting
- [ ] Distributed rate limiting (Redis)
- [ ] Rate limit bypass for internal services
- [ ] Rate limit monitoring & alerting

#### Security Scanning
```bash
# Automated security scans
gosec ./...                    # Static analysis
trivy image myimage:tag        # Container scanning
snyk test                      # Dependency scanning
owasp-zap baseline             # DAST scanning
```

- [ ] Setup Gosec in CI
- [ ] Trivy container scanning
- [ ] Dependency vulnerability scanning
- [ ] Regular penetration testing

#### Security Headers & CORS
```go
// Security headers
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin

// CORS configuration
Allowed Origins: [https://app.example.com, https://admin.example.com]
Allowed Methods: [GET, POST, PUT, PATCH, DELETE]
Allowed Headers: [Content-Type, Authorization]
Max Age: 86400
```

- [ ] Security headers middleware
- [ ] CORS configuration
- [ ] CSP policy
- [ ] HTTPS enforcement

### 7.3 Observability & Monitoring

#### Logging
```go
// Structured logging
logger.Info("User login",
    zap.String("user_id", userID),
    zap.String("tenant_id", tenantID),
    zap.String("ip", ip),
    zap.Duration("duration", duration),
)

// Log levels
DEBUG: Development details
INFO:  Normal operations
WARN:  Warning conditions
ERROR: Error conditions
FATAL: Critical failures
```

- [ ] Centralized logging (ELK Stack)
- [ ] Log aggregation
- [ ] Log retention policy (90 days)
- [ ] Log search & filtering
- [ ] Log alerting (error spikes)

#### Metrics
```go
// Prometheus metrics
http_requests_total{service="auth", method="POST", status="200"}
http_request_duration_seconds{service="auth", endpoint="/login"}
database_connections_active{service="academic"}
cache_hit_ratio{service="auth"}
queue_messages_pending{service="notification"}
```

- [ ] Prometheus setup
- [ ] Grafana dashboards
- [ ] Service-level metrics
- [ ] Business metrics (user growth, revenue)
- [ ] Alert rules (latency, error rate)

#### Distributed Tracing
```go
// Jaeger tracing
span := tracer.StartSpan("GetUserByID")
defer span.Finish()

span.SetTag("user_id", userID)
span.SetTag("tenant_id", tenantID)
```

- [ ] Jaeger setup
- [ ] Trace context propagation
- [ ] Service dependency mapping
- [ ] Latency analysis

#### Health Checks
```go
GET /health              # Basic health
GET /health/ready        # Readiness probe
GET /health/live         # Liveness probe

Response:
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:30:00Z",
  "checks": {
    "database": "healthy",
    "redis": "healthy",
    "message_queue": "healthy"
  },
  "version": "1.2.3"
}
```

- [ ] Health check endpoints
- [ ] Dependency health checks
- [ ] Kubernetes probes configuration
- [ ] Uptime monitoring (external)

### 7.4 Load Testing

**Tools**: k6, JMeter, Gatling

**Test Scenarios**:
```javascript
// k6 load test
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },   // Ramp up
    { duration: '5m', target: 100 },   // Stay at 100 users
    { duration: '2m', target: 200 },   // Spike
    { duration: '5m', target: 200 },
    { duration: '2m', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% < 500ms
    http_req_failed: ['rate<0.01'],    // Error rate < 1%
  },
};

export default function () {
  let response = http.get('https://api.example.com/api/v1/classes');
  check(response, {
    'status is 200': (r) => r.status === 200,
  });
  sleep(1);
}
```

**Load Test Scenarios**:
```
Baseline Test:        10 concurrent users, 5 minutes
Stress Test:          100-500 users, identify breaking point
Spike Test:           Sudden 10x traffic increase
Endurance Test:       Sustained load, 2-4 hours
Scalability Test:     Gradual increase to 1000 users
```

- [ ] k6 load test scripts
- [ ] Baseline performance metrics
- [ ] Stress test (find breaking point)
- [ ] Endurance test (memory leaks)
- [ ] Scalability test (horizontal scaling)
- [ ] Load test reporting

### 7.5 Horizontal Scaling

**Scalability Strategy**:
```
Application Layer:
- Stateless services (scale horizontally)
- Load balancer (Nginx, HAProxy, ALB)
- Auto-scaling based on CPU/memory/requests

Database Layer:
- Master-slave replication
- Read replicas for read-heavy operations
- Connection pooling
- Database sharding (if needed)

Cache Layer:
- Redis cluster
- Cache per service instance

Message Queue:
- RabbitMQ cluster
- Queue partitioning
```

- [ ] Load balancer setup
- [ ] Horizontal scaling configuration
- [ ] Auto-scaling policies
- [ ] Database replication
- [ ] Cache clustering
- [ ] Message queue clustering

### 7.6 Deployment Automation

**Kubernetes Deployment**:
```yaml
# Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: registry.example.com/auth-service:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secrets
              key: url
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
# Service
apiVersion: v1
kind: Service
metadata:
  name: auth-service
spec:
  selector:
    app: auth-service
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
---
# HorizontalPodAutoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

- [ ] Kubernetes cluster setup
- [ ] Helm charts for services
- [ ] ConfigMaps & Secrets management
- [ ] Ingress controller setup
- [ ] Auto-scaling configuration
- [ ] Rolling updates strategy
- [ ] Blue-green deployment
- [ ] Canary deployment strategy

### 7.7 Disaster Recovery & Backup

**Backup Strategy**:
```
Database Backups:
- Daily full backup (automated)
- Hourly incremental backup
- Retention: 30 days
- Off-site backup (cross-region)
- Point-in-time recovery (PITR)

File Backups:
- Daily backup of uploaded files
- Versioning enabled
- Cross-region replication

Disaster Recovery:
- RTO (Recovery Time Objective): 1 hour
- RPO (Recovery Point Objective): 15 minutes
- DR site setup (different region)
- DR drill (quarterly)
```

- [ ] Automated database backup
- [ ] Backup verification
- [ ] Restore procedure documentation
- [ ] DR site setup
- [ ] DR runbook
- [ ] Regular DR drills

### 7.8 Documentation & Runbooks

**Production Runbooks**:
```
Runbooks/
├── 01-deployment.md           # Deployment procedures
├── 02-rollback.md             # Rollback procedures
├── 03-scaling.md              # Scaling procedures
├── 04-incident-response.md    # Incident handling
├── 05-backup-restore.md       # Backup & restore
├── 06-monitoring.md           # Monitoring setup
├── 07-troubleshooting.md      # Common issues
└── 08-disaster-recovery.md    # DR procedures
```

- [ ] Deployment runbook
- [ ] Rollback procedure
- [ ] Incident response playbook
- [ ] Monitoring guide
- [ ] Troubleshooting guide
- [ ] On-call rotation setup

### 7.9 Production Checklist

**Pre-Production Checklist**:
- [ ] All services tested (unit, integration, E2E)
- [ ] Load testing passed
- [ ] Security audit completed
- [ ] Documentation complete
- [ ] Monitoring configured
- [ ] Alerting rules set
- [ ] Backup verified
- [ ] DR plan tested
- [ ] SSL certificates configured
- [ ] DNS configured
- [ ] CDN configured (if applicable)
- [ ] Rate limiting active
- [ ] CORS configured
- [ ] Database migrations tested
- [ ] Secrets rotated
- [ ] Access controls reviewed
- [ ] Compliance requirements met
- [ ] Stakeholder sign-off

### Deliverables
- [ ] Performance optimized (P95 < 500ms)
- [ ] Security hardened
- [ ] Monitoring operational
- [ ] Load testing completed
- [ ] Kubernetes deployment ready
- [ ] Backup & DR procedures in place
- [ ] Production runbooks complete
- [ ] System production-ready

### Success Criteria
- [ ] System handles 1000+ concurrent users
- [ ] API response time P95 < 500ms
- [ ] Zero critical security vulnerabilities
- [ ] 99.9% uptime achieved
- [ ] Automated deployment working
- [ ] Monitoring & alerting operational
- [ ] DR tested successfully

---

## Post-Launch (Continuous)

### Continuous Improvement
- [ ] Monitor production metrics
- [ ] Collect user feedback
- [ ] Analyze performance bottlenecks
- [ ] Security patches & updates
- [ ] Feature enhancements
- [ ] Technical debt reduction

### Maintenance
- [ ] Weekly health checks
- [ ] Monthly security reviews
- [ ] Quarterly DR drills
- [ ] Bi-annual load testing
- [ ] Continuous optimization

---

## Appendix

### Tech Stack Summary
```
Backend:         Go 1.21+
Framework:       Gin / Echo
Database:        PostgreSQL 15+
Cache:           Redis 7+
Message Queue:   RabbitMQ
Container:       Docker
Orchestration:   Kubernetes
API Gateway:     Kong / Traefik
Monitoring:      Prometheus + Grafana
Logging:         ELK Stack
Tracing:         Jaeger
CI/CD:           GitHub Actions
```

### Key Metrics
```
Performance:
- API response P50: < 100ms
- API response P95: < 500ms
- API response P99: < 1000ms

Availability:
- Uptime SLA: 99.9% (8.76 hours downtime/year)
- Error rate: < 0.1%

Scalability:
- Concurrent users: 1000+
- Requests per second: 10,000+
- Data volume: 10GB+ per tenant
```

---

**Last Updated**: 2025-01-15  
**Version**: 2.0  
**Owner**: Engineering Team  
**Status**: Active Roadmap
