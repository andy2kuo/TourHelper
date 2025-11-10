# Getting Started with TourHelper Development

This guide will help you set up your development environment and start contributing to TourHelper.

## Prerequisites

### Required
- **Go**: Version 1.24 or higher
  - Download from [golang.org](https://golang.org/dl/)
  - Verify: `go version`

### Optional
- **Docker**: For containerized deployment
  - Download from [docker.com](https://www.docker.com/get-started)
  - Verify: `docker --version`

- **Docker Compose**: For multi-container setup
  - Usually included with Docker Desktop
  - Verify: `docker-compose --version`

- **PostgreSQL**: Version 16 or higher (if not using Docker)
  - Download from [postgresql.org](https://www.postgresql.org/download/)
  - Verify: `psql --version`

- **Make**: For build automation
  - Usually pre-installed on macOS/Linux
  - Windows: Install via chocolatey or use WSL

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://github.com/andy2kuo/TourHelper.git
cd TourHelper
```

### 2. Install Go Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Environment

Create a `.env` file from the example:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```bash
# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Database Configuration (optional for development)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tourhelper
DB_SSLMODE=disable

# Logger Configuration
LOG_LEVEL=debug
LOG_FORMAT=console
```

**Note**: The application can run without a database in degraded mode for development.

### 4. Choose Your Development Path

#### Option A: Run Without Database (Quickest)

Perfect for working on API structure, middleware, or testing without data persistence.

```bash
# Run directly
go run cmd/api/main.go

# Or use Make
make run
```

The server will start on `http://localhost:8080` with limited functionality.

#### Option B: Run With Local PostgreSQL

If you have PostgreSQL installed locally:

1. Create the database:
```bash
psql -U postgres
CREATE DATABASE tourhelper;
\q
```

2. Initialize the schema:
```bash
psql -U postgres -d tourhelper -f scripts/init.sql
```

3. Run the application:
```bash
make run
```

#### Option C: Run With Docker Compose (Recommended)

This starts both the API server and PostgreSQL database:

```bash
docker-compose up -d
```

Check the logs:
```bash
docker-compose logs -f api
```

Stop the services:
```bash
docker-compose down
```

## Verify Installation

### 1. Check Health

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "database": "healthy",
  "version": "1.0.0"
}
```

### 2. Test API Endpoints

List tours:
```bash
curl http://localhost:8080/api/v1/tours
```

### 3. Run Tests

```bash
make test

# Or directly
go test -v ./...
```

## Development Workflow

### Project Structure

```
TourHelper/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── config/          # Configuration
│   ├── database/        # Database connection
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── model/           # Domain models
│   ├── repository/      # Data access
│   └── service/         # Business logic
├── pkg/                 # Public libraries
├── docs/                # Documentation
├── scripts/             # Database scripts
└── Makefile            # Build automation
```

### Common Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Start Docker containers
make docker-up

# Stop Docker containers
make docker-down
```

### Making Changes

1. **Create a feature branch**
```bash
git checkout -b feature/my-new-feature
```

2. **Make your changes**
   - Follow the existing code style
   - Update tests as needed
   - Add comments for complex logic

3. **Format your code**
```bash
make fmt
```

4. **Run tests**
```bash
make test
```

5. **Build and test locally**
```bash
make build
./bin/api
```

6. **Commit your changes**
```bash
git add .
git commit -m "Add: description of changes"
```

7. **Push and create a Pull Request**
```bash
git push origin feature/my-new-feature
```

## Adding New Features

### Adding a New API Endpoint

1. **Define the model** in `internal/model/`
```go
type MyModel struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}
```

2. **Create repository methods** in `internal/repository/`
```go
type MyRepository interface {
    Create(ctx context.Context, model *MyModel) error
    GetByID(ctx context.Context, id int64) (*MyModel, error)
}
```

3. **Implement the repository**
```go
func (r *myRepository) Create(ctx context.Context, model *MyModel) error {
    // Implementation
}
```

4. **Add service layer** in `internal/service/`
```go
type MyService interface {
    CreateItem(ctx context.Context, req *CreateRequest) (*MyModel, error)
}
```

5. **Create handler** in `internal/handler/`
```go
func (h *MyHandler) CreateItem(c *gin.Context) {
    // Handler implementation
}
```

6. **Register routes** in `cmd/api/main.go`
```go
v1.POST("/items", myHandler.CreateItem)
```

7. **Add tests** for each layer

8. **Update API documentation** in `docs/API.md`

### Database Migrations

For schema changes:

1. Create a new SQL file in `scripts/`
2. Apply manually or update `scripts/init.sql`
3. Test with a fresh database

## Testing

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/config/

# With coverage
go test -cover ./...

# With verbose output
go test -v ./...
```

### Writing Tests

Create test files with `_test.go` suffix:

```go
package mypackage

import "testing"

func TestMyFunction(t *testing.T) {
    result := MyFunction()
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

## Debugging

### Using Logs

Set log level to debug:
```bash
LOG_LEVEL=debug go run cmd/api/main.go
```

### Using Delve Debugger

Install Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Debug the application:
```bash
dlv debug cmd/api/main.go
```

### Common Issues

#### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### Database Connection Failed
- Check PostgreSQL is running
- Verify credentials in `.env`
- Check database exists

#### Build Errors
```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

## Code Style Guidelines

### Go Best Practices

1. **Use `gofmt`**: Always format code
```bash
gofmt -w .
```

2. **Follow naming conventions**:
   - Exported names start with uppercase
   - Unexported names start with lowercase
   - Use camelCase for multi-word names

3. **Error handling**:
```go
if err != nil {
    return fmt.Errorf("descriptive message: %w", err)
}
```

4. **Comments**:
   - Public functions need comments
   - Start with function name
```go
// CreateTour creates a new tour in the system.
func CreateTour(...) {...}
```

5. **Keep functions small**: Max 50 lines when possible

6. **Use meaningful names**: `getUserByID` not `get`

### Project-Specific Guidelines

1. **Layer separation**: Don't import handler in service
2. **Use interfaces**: Define in the package that uses them
3. **Context propagation**: Always pass context
4. **Structured logging**: Use zap logger with fields
5. **Error wrapping**: Use `fmt.Errorf` with `%w`

## Resources

### Documentation
- [API Documentation](API.md)
- [Architecture Guide](ARCHITECTURE.md)
- [Go Documentation](https://golang.org/doc/)

### Tools
- [Gin Framework](https://gin-gonic.com/)
- [Zap Logger](https://github.com/uber-go/zap)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)

### Community
- [GitHub Issues](https://github.com/andy2kuo/TourHelper/issues)
- [Go Forum](https://forum.golangbridge.org/)

## Next Steps

1. ✅ Complete this setup guide
2. 📖 Read the [Architecture Guide](ARCHITECTURE.md)
3. 📖 Review the [API Documentation](API.md)
4. 🔍 Explore the codebase
5. 🐛 Pick an issue from GitHub
6. 💻 Start coding!

## Getting Help

- **Documentation**: Check docs/ directory
- **Issues**: Open a GitHub issue
- **Code**: Add comments and ask in PRs

Happy coding! 🚀
