# TourHelper Backend Architecture

## System Overview

TourHelper is a RESTful API service designed to help users discover and manage tour destinations. The system follows clean architecture principles with clear separation of concerns.

## Architecture Layers

### 1. Presentation Layer (Handler)
**Location**: `internal/handler/`

Responsibilities:
- HTTP request/response handling
- Input validation
- Response formatting
- Error handling

Components:
- `health_handler.go`: Health check and readiness endpoints
- `tour_handler.go`: Tour management endpoints

### 2. Business Logic Layer (Service)
**Location**: `internal/service/`

Responsibilities:
- Business rules enforcement
- Data transformation
- Orchestrating repository calls
- Complex business logic

Components:
- `tour_service.go`: Tour service interface
- `tour_service_impl.go`: Tour service implementation

### 3. Data Access Layer (Repository)
**Location**: `internal/repository/`

Responsibilities:
- Database queries
- Data persistence
- CRUD operations
- Query optimization

Components:
- `tour_repository.go`: Repository interface
- `tour_repository_impl.go`: PostgreSQL implementation

### 4. Domain Layer (Model)
**Location**: `internal/model/`

Responsibilities:
- Domain entities
- Data transfer objects (DTOs)
- Request/response models

Components:
- `tour.go`: Tour entity and related models
- `response.go`: Standard API response models

### 5. Infrastructure Layer
**Location**: `internal/database/`, `internal/config/`, `internal/middleware/`

Responsibilities:
- Database connections
- Configuration management
- Cross-cutting concerns (logging, CORS, recovery)

Components:
- `database/database.go`: Database connection management
- `config/config.go`: Configuration loading
- `middleware/`: HTTP middleware components

### 6. Utilities Layer
**Location**: `pkg/utils/`

Responsibilities:
- Shared utilities
- Helper functions
- Reusable components

Components:
- `logger.go`: Logging utilities

## Data Flow

```
Client Request
    ↓
HTTP Router (Gin)
    ↓
Middleware (CORS, Logger, Recovery)
    ↓
Handler (HTTP handling, validation)
    ↓
Service (Business logic)
    ↓
Repository (Data access)
    ↓
Database (PostgreSQL)
```

## Component Diagram

```
┌─────────────────────────────────────────────────┐
│                   Client                        │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│            HTTP Server (Gin)                    │
│  ┌──────────────────────────────────────────┐   │
│  │    Middleware Layer                      │   │
│  │  - CORS                                  │   │
│  │  - Logger                                │   │
│  │  - Recovery                              │   │
│  └──────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│            Handler Layer                        │
│  - HealthHandler                                │
│  - TourHandler                                  │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│            Service Layer                        │
│  - TourService (interface)                      │
│  - tourService (implementation)                 │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│          Repository Layer                       │
│  - TourRepository (interface)                   │
│  - tourRepository (implementation)              │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│            Database Layer                       │
│  - PostgreSQL Connection Pool                   │
│  - Connection Health Check                      │
└─────────────────────────────────────────────────┘
```

## API Endpoints

### Health Endpoints
- `GET /health` - Overall system health
- `GET /ready` - Readiness check

### Tour Management
- `POST /api/v1/tours` - Create tour
- `GET /api/v1/tours` - List tours (with filters)
- `GET /api/v1/tours/:id` - Get tour by ID
- `PUT /api/v1/tours/:id` - Update tour
- `DELETE /api/v1/tours/:id` - Delete tour
- `POST /api/v1/tours/suggest` - Get tour suggestions

## Database Schema

### Tours Table
```sql
CREATE TABLE tours (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    duration INTEGER NOT NULL,
    season VARCHAR(50),
    budget VARCHAR(20),
    image_url TEXT,
    rating DECIMAL(3, 2) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Indexes
- `idx_tours_category` on category
- `idx_tours_country` on country
- `idx_tours_budget` on budget
- `idx_tours_season` on season
- `idx_tours_rating` on rating

## Configuration

Configuration is managed through environment variables:

| Category | Variable | Description | Default |
|----------|----------|-------------|---------|
| Server | `SERVER_PORT` | HTTP port | `8080` |
| Server | `GIN_MODE` | Gin mode | `debug` |
| Database | `DB_HOST` | Database host | `localhost` |
| Database | `DB_PORT` | Database port | `5432` |
| Database | `DB_USER` | Database user | `postgres` |
| Database | `DB_PASSWORD` | Database password | `postgres` |
| Database | `DB_NAME` | Database name | `tourhelper` |
| Database | `DB_SSLMODE` | SSL mode | `disable` |
| Logging | `LOG_LEVEL` | Log level | `info` |
| Logging | `LOG_FORMAT` | Log format | `json` |

## Dependency Injection

The application uses constructor-based dependency injection:

```go
// Initialize layers from bottom up
database → repository → service → handler

// Example flow:
db := database.New(config)
repo := repository.NewTourRepository(db, logger)
service := service.NewTourService(repo, logger)
handler := handler.NewTourHandler(service, logger)
```

## Error Handling

### Levels of Error Handling

1. **Repository Layer**: Database errors, not found errors
2. **Service Layer**: Business logic errors, validation errors
3. **Handler Layer**: HTTP status codes, user-facing error messages
4. **Middleware Layer**: Panic recovery, logging

### Error Response Format

```json
{
  "success": false,
  "error": "User-friendly error message"
}
```

## Logging

### Structured Logging with Zap

All logs are structured JSON for easy parsing:

```json
{
  "level": "info",
  "timestamp": "2024-01-01T00:00:00Z",
  "caller": "handler/tour_handler.go:45",
  "message": "Tour created successfully",
  "tour_id": 123,
  "name": "Tokyo Tour"
}
```

### Log Levels
- **DEBUG**: Detailed debugging information
- **INFO**: General informational messages
- **WARN**: Warning messages
- **ERROR**: Error messages

## Middleware Pipeline

Request flows through middleware in order:

1. **Recovery**: Catches panics and returns 500 error
2. **Logger**: Logs request details and response
3. **CORS**: Handles cross-origin requests

## Testing Strategy

### Unit Tests
- Configuration loading
- Business logic validation
- Data transformation

### Integration Tests (Future)
- API endpoint testing
- Database operations
- Service interactions

### End-to-End Tests (Future)
- Full request/response cycle
- Multiple component interactions

## Deployment Options

### 1. Standalone Binary
```bash
go build -o api ./cmd/api
./api
```

### 2. Docker Container
```bash
docker build -t tourhelper:latest .
docker run -p 8080:8080 tourhelper:latest
```

### 3. Docker Compose (with PostgreSQL)
```bash
docker-compose up -d
```

## Security Considerations

1. **Database**: Uses parameterized queries to prevent SQL injection
2. **CORS**: Configurable CORS policy
3. **Logging**: Sensitive data not logged
4. **Graceful Shutdown**: Proper resource cleanup
5. **Health Checks**: Readiness/liveness probes for orchestrators

## Performance Features

1. **Database Connection Pooling**: Max 25 connections, 5 idle
2. **Timeout Configuration**: Read/Write/Idle timeouts
3. **Structured Logging**: Minimal performance overhead
4. **Pagination**: Limit/offset for large datasets

## Monitoring & Observability

### Health Checks
- `/health`: Overall system health
- `/ready`: Readiness for traffic

### Logs
- Structured JSON logs
- Request/response logging
- Error tracking
- Performance metrics

## Future Enhancements

1. **Authentication & Authorization**: JWT tokens, role-based access
2. **Caching**: Redis for frequently accessed data
3. **Rate Limiting**: Per-user or per-IP rate limits
4. **API Versioning**: Support for multiple API versions
5. **Metrics**: Prometheus metrics endpoint
6. **Tracing**: Distributed tracing with OpenTelemetry
7. **Documentation**: OpenAPI/Swagger specification
8. **Database Migrations**: Automated schema migrations
9. **Message Queue**: Async processing for heavy tasks
10. **Search**: Full-text search with Elasticsearch

## Development Workflow

### Adding a New Feature

1. **Define Model**: Add/update models in `internal/model/`
2. **Create Repository**: Add repository methods
3. **Implement Service**: Add business logic
4. **Create Handler**: Add HTTP endpoints
5. **Update Router**: Register routes in `main.go`
6. **Write Tests**: Add unit and integration tests
7. **Document API**: Update API documentation

### Best Practices

- Follow clean architecture principles
- Keep layers independent
- Use dependency injection
- Write comprehensive tests
- Log important events
- Handle errors gracefully
- Document public APIs
- Use meaningful variable names
- Keep functions small and focused

## Conclusion

This architecture provides a solid foundation for a scalable, maintainable tour management system. The clean separation of concerns allows for easy testing, modification, and extension of the system.
