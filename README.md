# 🔐 Go Auth Service

A production-ready, scalable authentication microservice built with Go, featuring JWT tokens, Argon2id password hashing, role-based access control (RBAC), and comprehensive security measures.

## ✨ Features

### 🔒 Security First
- **Argon2id Password Hashing** - Industry-standard password security
- **RS256 JWT Tokens** - RSA-based JSON Web Tokens for stateless authentication
- **Comprehensive Security Headers** - CORS, CSP, XSS protection, and more
- **Rate Limiting** - Built-in protection against brute force attacks
- **CSRF Protection** - Cross-site request forgery prevention

### 🏗️ Architecture & Design
- **Clean Architecture** - Domain-driven design with clear separation of concerns
- **Dependency Injection** - Using Uber FX for clean dependency management
- **Repository Pattern** - Abstracted data access layer
- **Domain Entities** - Rich domain models with business logic
- **CQRS Ready** - Separated command and query responsibilities

### 🚀 Production Ready
- **Docker & Docker Compose** - Containerized deployment
- **PostgreSQL Integration** - Robust relational database with migrations
- **Redis Support** - Caching and session management (configurable)
- **Prometheus Metrics** - Comprehensive application monitoring
- **Health Checks** - Kubernetes-ready liveness and readiness probes
- **Hot Reload** - Air for development productivity

### 🛠️ Developer Experience
- **CLI Commands** - Built-in migration and server management
- **Environment Configuration** - YAML + Environment variable support
- **Database Migrations** - Goose-powered schema management
- **Structured Logging** - Zap for production-grade logging
- **Comprehensive Testing** - Unit tests with benchmarks

## 🏃‍♂️ Quick Start

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 15+

### 1. Clone and Setup
```bash
git clone <repository-url>
cd auth-service
cp .env.example .env
```

### 2. Start Services
```bash
# Start all services (PostgreSQL, Redis, Adminer, Auth Service)
docker-compose up --build

# Or run locally with air for development
go mod download
air
```

### 3. Run Migrations
```bash
# Apply all migrations
go run cmd/auth-service/main.go migrate up

# Or using Docker
docker-compose exec auth ./tmp/main migrate up
```

### 4. Access Services
- **Auth Service**: http://localhost:8080
- **Database Admin (Adminer)**: http://localhost:8081
- **Health Check**: http://localhost:8080/health
- **Metrics**: http://localhost:8080/metrics

## 📁 Project Structure

```
auth-service/
├── cmd/
│   ├── auth-service/           # Application entry point
│   └── commands/              # CLI commands (server, migrate)
├── internal/
│   ├── application/           # Application layer
│   │   ├── dto/              # Data transfer objects
│   │   └── services/         # Business logic services
│   ├── domain/               # Domain layer (entities, repositories)
│   │   ├── auth/            # Authentication domain
│   │   ├── shared/          # Shared domain utilities
│   │   ├── token/           # Token management
│   │   └── user/            # User domain
│   ├── infrastructure/       # Infrastructure layer
│   │   ├── config/          # Configuration management
│   │   ├── crypto/          # Password hashing utilities
│   │   ├── logger/          # Structured logging
│   │   ├── observability/   # Metrics and monitoring
│   │   ├── persistence/     # Database connection & migrations
│   │   └── repository/      # Data access implementations
│   └── interfaces/          # Interface layer
│       ├── handlers/        # HTTP handlers
│       ├── middleware/      # HTTP middleware
│       └── routes/          # Route definitions
├── pkg/                     # Shared packages
│   ├── response/           # HTTP response utilities
│   └── router/             # Router configuration
├── migrations/             # Database schema migrations
├── docker-compose.yml      # Multi-service setup
├── Dockerfile             # Container definition
└── .air.toml              # Hot reload configuration
```

## 🔧 Configuration

The service supports both YAML configuration files and environment variables. Environment variables take precedence.

### Key Configuration Sections

#### Database
```yaml
postgres:
  host: localhost
  port: 5432
  user: your_user
  password: your_password
  dbname: authdb
  sslmode: disable
```

#### JWT Tokens
```yaml
jwt:
  access_token_expiry: 15m
  refresh_token_expiry: 7d
  issuer: "auth-service"
  audience: "go-micro"
  signing_method: "RS256"
  private_key_path: "infrastructure/keys/private.pem"
  public_key_path: "infrastructure/keys/public.pem"
```

#### Security
```yaml
argon2id:
  time: 1
  memory: 65536      # 64MB
  threads: 4
  key_length: 32
  salt_length: 16
```

## 🚀 API Endpoints

### Authentication
```http
POST /v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure_password"
}
```

### Health & Monitoring
```http
GET /health              # Health check
GET /ready               # Readiness probe
GET /live                # Liveness probe
GET /metrics             # Prometheus metrics
```

## 🛠️ CLI Commands

### Server Management
```bash
# Start the server
go run cmd/auth-service/main.go server

# Or using the built binary
./auth-service server
```

### Database Migrations
```bash
# Apply all pending migrations
go run cmd/auth-service/main.go migrate up

# Rollback last migration
go run cmd/auth-service/main.go migrate down

# Check migration status
go run cmd/auth-service/main.go migrate status

# Create new migration
go run cmd/auth-service/main.go migrate create add_user_table sql
```

## 🐳 Docker Deployment

### Development Setup
```bash
# Start all services
docker-compose up --build

# Start only database
docker-compose up db

# View logs
docker-compose logs -f auth
```

### Production Considerations
- Use environment-specific `.env` files
- Configure proper secrets management
- Set up SSL/TLS certificates
- Configure resource limits
- Set up log aggregation

## 📊 Monitoring & Observability

### Prometheus Metrics
The service exposes comprehensive metrics:
- HTTP request metrics (latency, status codes, request count)
- Authentication metrics (login attempts, token generation)
- Database metrics (query performance, connection pool)
- Business metrics (user registrations, active users)

### Health Checks
- `/health` - General health status
- `/ready` - Kubernetes readiness probe
- `/live` - Kubernetes liveness probe

### Logging
Structured JSON logging with configurable levels:
```yaml
logging:
  level: info           # debug, info, warn, error
  encoding: json        # json or console
  output: stdout
  error_output: stderr
```

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./internal/infrastructure/crypto/
```

### Test Categories
- **Unit Tests** - Individual component testing
- **Integration Tests** - Database and service integration
- **Benchmark Tests** - Performance testing for crypto operations

## 🔒 Security Features

### Password Security
- **Argon2id Hashing** - Resistant to GPU and side-channel attacks
- **Configurable Parameters** - Memory, time, and thread costs
- **Salt Generation** - Unique salt per password

### Token Security
- **RS256 Algorithm** - RSA public/private key pair
- **Short-lived Access Tokens** - 15-minute default expiry
- **Refresh Token Rotation** - Long-lived refresh tokens (7 days)

### HTTP Security
- **CORS Configuration** - Cross-origin resource sharing
- **Security Headers** - XSS protection, content sniffing prevention
- **Rate Limiting** - Configurable request limits per IP
- **Request Size Limits** - Prevent large payload attacks

## 🚀 Performance Features

### Database Optimization
- **Connection Pooling** - Optimized database connections
- **Query Monitoring** - Prometheus metrics for query performance
- **Migration Management** - Version-controlled schema changes

### Caching Strategy
- **Redis Integration** - Session and token caching
- **Configurable TTL** - Time-based cache expiration
- **Cache Metrics** - Hit/miss ratio monitoring

### Resource Management
- **Graceful Shutdown** - Proper cleanup on service termination
- **Connection Limits** - Prevents resource exhaustion
- **Timeout Configuration** - Request and database timeouts

## 🔄 Development Workflow

### Local Development
1. **Setup Environment**: Copy `.env.example` to `.env`
2. **Start Dependencies**: `docker-compose up db redis`
3. **Run Migrations**: `go run cmd/auth-service/main.go migrate up`
4. **Start Server**: `air` (hot reload) or `go run cmd/auth-service/main.go server`

### Adding New Features
1. **Domain Layer** - Define entities and repositories
2. **Application Layer** - Implement business logic services
3. **Infrastructure Layer** - Add repository implementations
4. **Interface Layer** - Create handlers and routes
5. **Update Modules** - Wire dependencies with Uber FX

## 📈 Scalability Considerations

### Horizontal Scaling
- **Stateless Design** - No server-side session storage
- **Load Balancer Ready** - Standard HTTP service
- **Database Separation** - Separate read/write operations

### Performance Tuning
- **Connection Pool Sizing** - Based on CPU cores
- **JWT Token Optimization** - Minimal payload size
- **Database Indexing** - Optimized query performance

### Monitoring Integration
- **Prometheus Metrics** - Comprehensive performance monitoring
- **Distributed Tracing** - Request correlation across services
- **Error Tracking** - Structured error reporting

## 🤝 Contributing

### Code Style
- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comprehensive tests for new features
- Update documentation for API changes

### Pull Request Process
1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Update documentation
5. Submit pull request with clear description

---

**Built with ❤️ using Go, PostgreSQL, Redis, and modern DevOps practices.**
