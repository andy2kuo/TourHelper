# TourHelper
自動幫我想可以去哪旅遊

A modern Golang backend service for managing and suggesting tour destinations.

## Features

- 🏗️ **Clean Architecture**: Separation of concerns with handlers, services, repositories
- 🚀 **RESTful API**: Complete CRUD operations for tour management
- 🤖 **Tour Suggestions**: AI-ready architecture for suggesting destinations based on preferences
- 🔍 **Filtering & Pagination**: Advanced query capabilities
- 🐳 **Docker Support**: Ready for containerized deployment
- 📝 **Structured Logging**: Using Uber's Zap logger
- 🔄 **Graceful Shutdown**: Proper resource cleanup
- 💾 **PostgreSQL**: Production-ready database support
- ⚡ **High Performance**: Built with Gin web framework

## Architecture

```
TourHelper/
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middleware (CORS, logging, recovery)
│   ├── model/            # Domain models and DTOs
│   ├── repository/       # Data access layer
│   └── service/          # Business logic layer
├── pkg/
│   └── utils/            # Shared utilities
├── api/
│   └── v1/               # API versioning structure
├── docs/                 # Documentation
├── scripts/              # Database scripts and utilities
├── Dockerfile            # Container definition
├── docker-compose.yml    # Multi-container setup
└── Makefile             # Build automation

```

## Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Logging**: Uber Zap
- **Configuration**: Environment variables with godotenv
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 16 (optional - can run without database)
- Docker & Docker Compose (for containerized deployment)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/andy2kuo/TourHelper.git
   cd TourHelper
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   # Without database (degraded mode)
   make run

   # Or directly
   go run cmd/api/main.go
   ```

   The server will start on `http://localhost:8080`

### Docker Deployment

1. **Start all services** (API + PostgreSQL)
   ```bash
   docker-compose up -d
   ```

2. **Check service health**
   ```bash
   curl http://localhost:8080/health
   ```

3. **Stop services**
   ```bash
   docker-compose down
   ```

## API Documentation

See [API Documentation](docs/API.md) for detailed endpoint information.

### Quick Examples

**Health Check**
```bash
curl http://localhost:8080/health
```

**List Tours**
```bash
curl http://localhost:8080/api/v1/tours
```

**Create Tour**
```bash
curl -X POST http://localhost:8080/api/v1/tours \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Amazing Tour",
    "description": "A wonderful experience",
    "location": "Paris",
    "country": "France",
    "category": "cultural",
    "duration": 5,
    "budget": "medium",
    "rating": 4.5
  }'
```

**Get Tour Suggestions**
```bash
curl -X POST http://localhost:8080/api/v1/tours/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "category": "beach",
    "budget": "medium",
    "min_rating": 4.0
  }'
```

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Lint Code

```bash
make lint
```

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `GIN_MODE` | Gin mode (debug/release) | `debug` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `tourhelper` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `LOG_LEVEL` | Log level | `info` |
| `LOG_FORMAT` | Log format (json/console) | `json` |

## Database Schema

The application uses a PostgreSQL database with the following main table:

**tours**
- `id`: Serial primary key
- `name`: Tour name
- `description`: Tour description
- `location`: Location name
- `country`: Country name
- `category`: Tour category (beach, mountain, cultural, city)
- `duration`: Duration in days
- `season`: Best season to visit
- `budget`: Budget level (low, medium, high)
- `image_url`: Tour image URL
- `rating`: Rating (0-5)
- `created_at`: Creation timestamp
- `updated_at`: Update timestamp

## Project Status

✅ Core backend architecture implemented
✅ RESTful API endpoints
✅ Database integration
✅ Docker support
✅ Basic tour management
✅ Tour suggestion system (basic)

### Future Enhancements

- [ ] User authentication and authorization
- [ ] Advanced AI-powered tour recommendations
- [ ] Image upload support
- [ ] Review and rating system
- [ ] Booking integration
- [ ] Multi-language support
- [ ] Caching layer (Redis)
- [ ] API rate limiting
- [ ] Comprehensive test coverage
- [ ] CI/CD pipeline

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For questions or feedback, please open an issue on GitHub.
