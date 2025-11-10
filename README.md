# TourHelper
自動幫我想可以去哪旅遊

A tour planning assistant that helps you discover travel destinations with real-time collaboration features.

## Features

- 🌐 **Multi-Platform Support**: Web, Line, and Telegram clients
- 🔐 **Authentication**: JWT-based login and token verification
- 🔄 **Real-time Communication**: WebSocket support for live updates
- 🗺️ **Google Maps Integration**: Interactive map with location search
- 💬 **Real-time Chat**: Collaborate with other users in real-time
- 📍 **Tour Suggestions**: Get recommendations for travel destinations

## Tech Stack

- **Backend**: Go 1.24+
- **WebSocket**: Gorilla WebSocket
- **Authentication**: JWT (JSON Web Tokens)
- **Frontend**: HTML, CSS, JavaScript
- **Maps**: Google Maps JavaScript API

## Prerequisites

- Go 1.24 or higher
- Google Maps API Key (get it from [Google Cloud Console](https://console.cloud.google.com/google/maps-apis))

## Installation

1. Clone the repository:
```bash
git clone https://github.com/andy2kuo/TourHelper.git
cd TourHelper
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file from the example:
```bash
cp .env.example .env
```

4. Edit `.env` and add your Google Maps API Key:
```bash
GOOGLE_MAPS_API_KEY=YOUR_API_KEY_HERE
```

## Running the Application

1. Start the server:
```bash
go run cmd/server/main.go
```

2. Open your browser and navigate to:
```
http://localhost:8080
```

## Default Test Accounts

The application comes with pre-configured test accounts:

| Platform | Username | Password |
|----------|----------|----------|
| Web | webuser | password123 |
| Line | lineuser | password123 |
| Telegram | telegramuser | password123 |

## API Endpoints

### Authentication
- `POST /api/login` - User login
- `GET /api/config` - Get client configuration (requires authentication)

### WebSocket
- `GET /ws?token={jwt_token}` - WebSocket connection (requires authentication)

### Health Check
- `GET /api/health` - Server health check

## Project Structure

```
TourHelper/
├── cmd/
│   └── server/          # Main server application
├── internal/
│   ├── auth/            # Authentication logic
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── websocket/       # WebSocket management
├── web/
│   ├── static/
│   │   ├── css/         # Stylesheets
│   │   └── js/          # JavaScript files
│   └── index.html       # Main HTML page
├── config/              # Configuration management
├── .env.example         # Environment variables example
├── go.mod               # Go module definition
└── README.md            # This file
```

## Development

### Building
```bash
go build -o tourhelper cmd/server/main.go
```

### Running Tests
```bash
go test ./...
```

## Configuration

Configuration is managed through environment variables:

- `SERVER_PORT`: Server port (default: 8080)
- `JWT_SECRET`: Secret key for JWT tokens (change in production!)
- `GOOGLE_MAPS_API_KEY`: Your Google Maps API key

## WebSocket Message Format

Messages sent via WebSocket should follow this format:

```json
{
  "type": "chat|suggestion",
  "payload": {
    // Type-specific data
  }
}
```

### Chat Message Example:
```json
{
  "type": "chat",
  "payload": {
    "sender": "username",
    "message": "Hello!",
    "platform": "web"
  }
}
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
