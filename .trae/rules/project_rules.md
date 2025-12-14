# PROJECT RULES — Backend Sistem Akademik

Aturan wajib untuk menjaga kualitas, keamanan, dan skalabilitas sistem enterprise-grade.

---

## 1. Architecture Rules

### 1.1 Microservices Design

- **Wajib microservices pattern** dengan single responsibility per service
- **Database per service** (Database per Service pattern) - tidak boleh shared database
- **Event-driven architecture** untuk komunikasi antar service yang loosely coupled
- **API Gateway** sebagai single entry point
- **Service mesh consideration** untuk production (Istio/Linkerd)
- **Circuit breaker pattern** untuk fault tolerance (menggunakan library seperti gobreaker)
- **Saga pattern** untuk distributed transactions

### 1.2 Domain Boundaries

```
Core Services:
├── auth-service          # Identity & Access Management
├── academic-service      # Core academic operations
├── attendance-service    # Presensi siswa & guru
├── assessment-service    # Penilaian & raport
├── admission-service     # PPDB management
├── finance-service       # Billing & payment
└── notification-service  # Email & WhatsApp

Supporting Services:
├── api-gateway          # Entry point & routing
├── file-service         # Document & media management
└── report-service       # Report generation & export
```

### 1.3 Communication Patterns

- **Synchronous**: REST API dengan HTTP/gRPC untuk read operations
- **Asynchronous**: Message broker (RabbitMQ/Kafka) untuk events & background jobs
- **Service discovery**: Consul atau built-in Kubernetes service discovery
- **Load balancing**: Reverse proxy (Nginx/Traefik) atau cloud load balancer

---

## 2. Coding Standards

### 2.1 Project Structure

```
service-name/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/          # Entities & business logic
│   ├── usecase/         # Application logic
│   ├── repository/      # Data access layer
│   ├── handler/         # HTTP/gRPC handlers
│   └── middleware/      # Middleware components
├── pkg/                 # Reusable packages
├── config/              # Configuration files
├── migrations/          # Database migrations
├── docs/                # Documentation
├── tests/               # Integration & E2E tests
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

### 2.2 Clean Architecture Layers

- **Domain Layer**: Entities, value objects, domain errors (zero dependencies)
- **Use Case Layer**: Business logic & orchestration
- **Repository Layer**: Database operations & external APIs
- **Handler Layer**: HTTP/gRPC request handling & validation
- **Dependency injection** wajib menggunakan wire atau manual DI

### 2.3 Code Quality

- **Conventional Commits** format: `type(scope): subject`
  - Types: feat, fix, docs, style, refactor, test, chore
  - Example: `feat(auth): add refresh token endpoint`
- **Go standard library preference** over third-party packages
- **Error handling**: Custom error types dengan context & stack trace
- **Linting**: golangci-lint dengan minimal 90% compliance
- **Code review** mandatory untuk setiap PR
- **No magic numbers**: Gunakan constants
- **Comment**: Godoc format untuk exported functions

### 2.4 Testing Requirements

- **Unit tests**: Minimal 70% code coverage
- **Integration tests**: Testing dengan real database (testcontainers)
- **E2E tests**: Critical user flows
- **Mock interfaces** menggunakan mockgen atau testify/mock
- **Table-driven tests** untuk multiple scenarios
- **Test naming**: `TestFunctionName_Scenario_ExpectedBehavior`

#### Coverage Threshold Policy

- **Global targets**: ≥90% statement coverage and ≥85% branch coverage across critical paths
- **Commit gating**: Jika total coverage < 80%, commit perubahan yang berorientasi coverage dengan pesan komit yang menyertakan persentase saat ini, contoh: `chore(coverage): raise coverage to 78.4%`
- **Quality gate**: Semua komit coverage harus lulus lint dan tes sebelum merge
- **Backward compatibility**: Perluas test suite tanpa mengubah perilaku publik; gunakan mocking untuk dependensi eksternal
- **Integration tests**: Gunakan Testcontainers untuk Postgres/Redis saat memungkinkan; test akan di-skip aman jika lingkungan tidak mendukung

---

## 3. API Rules

### 3.1 API Design

- **OpenAPI First Design**: Spec dulu baru implementasi
- **API versioning**: Mandatory `/api/v1/` prefix
- **RESTful conventions**:
  - GET: Read resources
  - POST: Create resources
  - PUT: Full update
  - PATCH: Partial update
  - DELETE: Remove resources
- **Resource naming**: Plural nouns (`/users`, `/classes`, `/subjects`)
- **Filtering & sorting**: Query parameters (`?status=active&sort=created_at:desc`)
- **Pagination**: Cursor-based atau offset-based dengan limit default

### 3.2 Request/Response Format

```json
// Success Response
{
  "success": true,
  "data": { ... },
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "request_id": "uuid"
  }
}

// Error Response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Email format is invalid"
      }
    ]
  },
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "request_id": "uuid"
  }
}
```

### 3.3 Error Codes Standard

```
1xxx - General errors (1000: Unknown, 1001: Internal Server Error)
2xxx - Authentication errors (2001: Unauthorized, 2002: Token Expired)
3xxx - Authorization errors (3001: Forbidden, 3002: Insufficient Permission)
4xxx - Validation errors (4001: Invalid Input, 4002: Required Field)
5xxx - Business logic errors (5001: Duplicate Entry, 5002: Resource Not Found)
6xxx - External service errors (6001: Third Party Error, 6002: Timeout)
```

### 3.4 API Documentation

- **Swagger/OpenAPI 3.0** specification
- **Postman collection** untuk testing
- **Example requests & responses** untuk setiap endpoint
- **Auto-generate** docs dari code annotations

---

## 4. Security Rules

### 4.1 Authentication

- **JWT access token**: Short-lived (15 minutes)
- **JWT refresh token**: Long-lived (7 days), stored dengan secure httpOnly cookie
- **Token rotation**: Refresh token harus di-rotate setiap digunakan
- **Revocation**: Token blacklist menggunakan Redis
- **Multi-factor authentication** untuk admin users
- **Password policy**: Minimal 8 karakter, kombinasi huruf, angka, simbol

### 4.2 Authorization

- **RBAC (Role-Based Access Control)** dinamis, tidak hardcode
- **Permission granularity**: Resource:Action format (`user:create`, `report:read`)
- **Role hierarchy**: Support untuk role inheritance
- **Policy engine**: Open Policy Agent (OPA) atau custom implementation
- **Context-aware authorization**: Tenant, organization, location-based

### 4.3 Data Security

- **Encryption at rest**: Sensitive data (passwords, API keys, PII)
- **Encryption in transit**: TLS 1.3 mandatory
- **Field-level encryption**: AES-256 untuk data sensitif
- **Key management**: Vault (HashiCorp) atau cloud KMS
- **Data masking**: Untuk logs & debugging
- **PII handling**: GDPR-compliant (right to be forgotten)

### 4.4 Security Best Practices

- **Rate limiting**: Per user, per IP, per endpoint
  - Authentication endpoints: 5 req/min
  - Read endpoints: 100 req/min
  - Write endpoints: 30 req/min
- **CORS policy**: Whitelist domains only
- **Input validation**: Sanitize semua input (SQL injection, XSS prevention)
- **SQL injection prevention**: Parameterized queries wajib
- **Audit logging**: WHO did WHAT, WHEN, WHERE, WHY
- **Security headers**: HSTS, CSP, X-Frame-Options, X-Content-Type-Options
- **Dependency scanning**: Automated vulnerability scanning (Dependabot, Snyk)
- **Secrets management**: Never commit secrets, use environment variables atau vault

---

## 5. Database Rules

### 5.1 Database Strategy

- **PostgreSQL** sebagai primary database
- **Isolated database** per microservice
- **Connection pooling**: pgx/pgxpool dengan proper configuration
- **Read replicas** untuk read-heavy operations
- **Sharding strategy** untuk horizontal scaling (jika diperlukan)

### 5.2 Schema Management

- **Migration tools**: golang-migrate atau goose
- **Versioned migrations**: Sequential numbering (001_initial_schema.up.sql)
- **Rollback support**: Setiap migration harus punya down migration
- **Idempotent migrations**: Safe untuk re-run
- **Zero-downtime deployments**: Backward compatible migrations

### 5.3 Database Design

- **Soft delete**: `deleted_at` timestamp instead of hard delete
- **Audit fields**: `created_at`, `updated_at`, `created_by`, `updated_by`
- **UUID primary keys**: Untuk distributed systems
- **Indexes**: Proper indexing untuk query performance
- **Foreign keys**: Enforce referential integrity
- **Constraints**: Check constraints untuk business rules
- **Naming conventions**:
  - Tables: snake_case, plural (`users`, `student_grades`)
  - Columns: snake_case (`first_name`, `email_address`)
  - Indexes: `idx_tablename_columnname`
  - Foreign keys: `fk_tablename_columnname`

### 5.4 Query Optimization

- **N+1 query prevention**: Eager loading atau batch queries
- **Query timeout**: Maximum 30 seconds
- **Prepared statements**: Untuk repeated queries
- **EXPLAIN ANALYZE**: Untuk performance tuning
- **Database monitoring**: Query performance tracking

---

## 6. DevOps Rules

### 6.1 Containerization

- **Docker** untuk semua services
- **Multi-stage builds**: Minimize image size
- **Distroless/Alpine** base images
- **Non-root user**: Security best practice
- **Health checks**: Mandatory dalam Dockerfile
- **Resource limits**: CPU & memory constraints
- **.dockerignore**: Exclude unnecessary files

### 6.2 Configuration Management

- **12-factor app principles**
- **Environment variables** untuk configuration
- **Config validation**: Startup validation untuk required configs
- **Secret management**: Never in environment variables, use vault
- **Multiple environments**: dev, staging, production
- **Feature flags**: LaunchDarkly atau custom implementation

### 6.3 Observability

#### Logging

- **Structured logging**: JSON format dengan contextual information
- **Log levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Centralized logging**: ELK stack atau Loki
- **Request ID**: Trace requests across services
- **No sensitive data**: Mask passwords, tokens, PII

#### Monitoring

- **Metrics collection**: Prometheus + Grafana
- **Key metrics**:
  - Request rate & latency
  - Error rate
  - CPU & memory usage
  - Database connections
  - Queue depth
- **Alerting**: PagerDuty, Slack integration
- **Uptime monitoring**: External monitoring (Pingdom, UptimeRobot)

#### Tracing

- **Distributed tracing**: Jaeger atau Zipkin
- **Span context propagation**: Across service boundaries
- **Performance profiling**: pprof endpoints

### 6.4 CI/CD Pipeline

```
Pipeline Stages:
1. Lint (golangci-lint)
2. Unit Tests (go test -cover)
3. Security Scan (gosec, trivy)
4. Build Docker Image
5. Integration Tests
6. Push to Registry
7. Deploy to Staging
8. E2E Tests
9. Deploy to Production (with approval)
```

- **GitOps workflow**: Infrastructure as code (Terraform, Helm)
- **Blue-green deployment** atau canary releases
- **Automated rollback**: Jika health check fails
- **Zero-downtime deployment**: Mandatory untuk production

### 6.5 Infrastructure

- **Kubernetes** untuk orchestration (production)
- **Docker Compose** untuk local development
- **Service mesh**: Istio untuk advanced traffic management
- **API Gateway**: Kong, Traefik, atau custom gateway
- **Load balancer**: Nginx, HAProxy, atau cloud LB
- **Cache layer**: Redis/Memcached
- **Message broker**: RabbitMQ atau Kafka
- **Object storage**: MinIO atau cloud storage (S3)

---

## 7. Documentation Rules

### 7.1 Code Documentation

- **README per service**: Purpose, setup, API endpoints, dependencies
- **Godoc comments**: Untuk exported functions & types
- **Inline comments**: Untuk complex logic
- **Architecture diagrams**: C4 model atau diagram.net
- **Sequence diagrams**: Untuk complex flows

### 7.2 API Documentation

- **OpenAPI/Swagger spec**: Auto-generated atau hand-written
- **Postman collection**: Import-ready
- **Code examples**: cURL, Go, JavaScript
- **Authentication guide**: How to obtain & use tokens
- **Rate limit documentation**: Limits per endpoint

### 7.3 Architecture Documentation

- **Architecture Decision Records (ADR)**:
  - Format: MADR (Markdown Any Decision Records)
  - Sections: Context, Decision, Consequences
  - Version controlled in `/docs/adr/`
- **System design docs**: High-level architecture
- **Data flow diagrams**: Service interactions
- **Database schema docs**: ER diagrams, table descriptions
- **Runbooks**: Operational procedures
- **Incident response**: Troubleshooting guides

### 7.4 Documentation Standards

- **Up-to-date**: Documentation drift adalah violation
- **Versioned**: Sync dengan API versions
- **Searchable**: Markdown format in repository
- **Change log**: CHANGELOG.md per service
- **Migration guides**: Major version upgrades

---

## 8. Performance Rules

### 8.1 Response Time Targets

- **API response time**:
  - P50 < 100ms
  - P95 < 500ms
  - P99 < 1000ms
- **Database queries**: < 50ms average
- **Background jobs**: Complete within SLA

### 8.2 Optimization Strategies

- **Caching**: Redis untuk frequently accessed data
- **Database indexes**: Query optimization
- **Connection pooling**: Reuse connections
- **Async processing**: Background jobs untuk heavy operations
- **CDN**: Static assets
- **Compression**: Gzip/Brotli untuk responses
- **Pagination**: Limit result sets
- **Lazy loading**: Load data on-demand

### 8.3 Scalability

- **Horizontal scaling**: Stateless services
- **Database read replicas**: Scale reads
- **Queue-based processing**: Handle spikes
- **Auto-scaling**: Based on metrics
- **Load testing**: Regular performance tests (k6, JMeter)

---

## 9. Multi-Tenancy Rules

### 9.1 Tenant Isolation

- **Tenant identification**: Header-based atau subdomain
- **Data isolation**: Row-level security atau schema per tenant
- **Resource quotas**: CPU, memory, storage limits per tenant
- **Rate limiting**: Per tenant basis

### 9.2 Configuration

- **Tenant-specific config**: Override default configs
- **Feature flags**: Per tenant features
- **Branding**: Custom themes per tenant (jika applicable)

---

## 10. Compliance & Standards

### 10.1 Code Review Requirements

- **Mandatory PR reviews**: Minimal 1 reviewer
- **Review checklist**:
  - Tests included & passing
  - Documentation updated
  - Security considerations
  - Performance impact
  - Breaking changes documented
- **No direct commits** to main/master branch

### 10.2 Definition of Done

- [ ] Code written & peer reviewed
- [ ] Unit tests written (>70% coverage)
- [ ] Integration tests (if applicable)
- [ ] Documentation updated
- [ ] API spec updated
- [ ] Security review passed
- [ ] Performance acceptable
- [ ] Deployed to staging
- [ ] QA approved

### 10.3 Violation Handling

- **Code review rejection**: Must fix before merge
- **Post-merge issues**: Immediate hotfix required
- **Repeated violations**: Team discussion & training
- **Security violations**: Immediate rollback & incident report

---

## Compliance Matrix

| Rule Category    | Priority | Auto-Enforced   | Manual Review |
| ---------------- | -------- | --------------- | ------------- |
| Architecture     | HIGH     | Partial         | Required      |
| Coding Standards | HIGH     | Yes (linter)    | Required      |
| API Design       | HIGH     | Yes (OpenAPI)   | Required      |
| Security         | CRITICAL | Partial         | Required      |
| Database         | HIGH     | Yes (migration) | Recommended   |
| DevOps           | HIGH     | Yes (CI/CD)     | Recommended   |
| Documentation    | MEDIUM   | Partial         | Required      |
| Performance      | HIGH     | Yes (tests)     | Required      |
| Testing          | HIGH     | Yes (coverage)  | Required      |

---

**Last Updated**: 2025-01-15  
**Version**: 2.0  
**Status**: Mandatory Compliance  
**Owner**: Engineering Team

---

## Quick Reference

### Common Commands

```bash
# Run tests with coverage
make test-coverage

# Lint code
make lint

# Run locally
make run-local

# Build docker image
make docker-build

# Run migrations
make migrate-up

# Generate mocks
make generate-mocks

# View API docs
make docs-serve
```

### Emergency Contacts

- **Security Issues**: security@company.com
- **Production Issues**: oncall@company.com
- **Architecture Questions**: #architecture-team (Slack)
