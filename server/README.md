# Visual API Testing Platform - Backend

Go backend server for the Visual API Testing Platform with PostgreSQL database.

## Features

- 🚀 RESTful API for flows, nodes, and test runs
- 🔌 WebSocket support for real-time execution updates
- 🔐 JWT-based authentication
- 📊 PostgreSQL database with migrations
- 🐳 Docker Compose setup
- ⚡ Concurrent flow execution engine

## Architecture

```
server/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── models/              # Data models
│   ├── repository/          # Database access layer
│   ├── handlers/            # HTTP handlers
│   ├── node/                # Node implementations
│   └── engine/              # Flow execution engine
├── database/
│   └── init.sql             # Database schema
└── docker-compose.yml       # Docker setup
```

## Setup

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 16 (via Docker)

### Environment Variables

Copy `.env.example` to `.env` and configure:

```env
DATABASE_URL=postgres://visual_testing:visual_testing_pass@localhost:5433/visual_testing_db?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
GIN_MODE=debug
```

### Running with Docker Compose

```bash
# Start database and backend
docker-compose up -d

# View logs
docker-compose logs -f backend
```

### Running Locally (without Docker)

```bash
# Start PostgreSQL database
docker-compose up -d postgres

# Install dependencies
go mod download

# Run migrations (database/init.sql will auto-run)
# Or manually create database and run init.sql

# Start server
go run cmd/main.go
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Flows (Protected)

- `GET /api/flows` - List all flows for user
- `POST /api/flows` - Create new flow
- `GET /api/flows/:id` - Get flow by ID
- `PUT /api/flows/:id` - Update flow
- `DELETE /api/flows/:id` - Delete flow
- `POST /api/flows/:id/run` - Execute flow
- `GET /api/flows/:id/test-runs` - Get test runs for flow

### Test Runs (Protected)

- `GET /api/test-runs/:id` - Get test run by ID

### WebSocket (Protected)

- `GET /api/ws?testRunId=<uuid>` - WebSocket connection for real-time updates

## Node Types

1. **API Node** - Executes HTTP requests
2. **Mock Node** - Returns predefined responses
3. **Verification Node** - Performs assertions
4. **Report Node** - Generates test reports
5. **Event Trigger Node** - Triggers flows on events

## Database

PostgreSQL runs on **port 5433** (custom port to avoid conflicts).

Database credentials:
- Host: `localhost`
- Port: `5433`
- User: `visual_testing`
- Password: `visual_testing_pass`
- Database: `visual_testing_db`

## Development

```bash
# Run tests
go test ./...

# Build
go build -o bin/server cmd/main.go

# Run
./bin/server
```

## Integration with Frontend

The backend expects the frontend to:
1. Connect to `http://localhost:8080`
2. Send JWT token in `Authorization: Bearer <token>` header
3. Connect to WebSocket at `ws://localhost:8080/api/ws?testRunId=<uuid>`

## Notes

- The flow execution engine runs nodes concurrently where dependencies allow
- WebSocket connections provide real-time updates during flow execution
- All API endpoints (except auth) require authentication
- Test runs are automatically saved to the database after execution

