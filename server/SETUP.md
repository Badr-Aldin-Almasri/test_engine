# Backend Setup Instructions

## Quick Start

### 1. Start Database and Backend with Docker Compose

```bash
cd server
docker-compose up -d
```

This will:
- Start PostgreSQL on port **5433** (custom port to avoid conflicts)
- Start the Go backend on port **8080**
- Automatically run database migrations from `database/init.sql`

### 2. Check Logs

```bash
# View backend logs
docker-compose logs -f backend

# View database logs
docker-compose logs -f postgres
```

### 3. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Database Connection

- **Host**: `localhost`
- **Port**: `5433` (custom port)
- **User**: `visual_testing`
- **Password**: `visual_testing_pass`
- **Database**: `visual_testing_db`

## API Base URL

- Backend: `http://localhost:8080`
- WebSocket: `ws://localhost:8080/api/ws?testRunId=<uuid>`

## Running Locally (without Docker)

If you want to run the backend locally without Docker:

```bash
# 1. Start only the database
docker-compose up -d postgres

# 2. Install Go dependencies
cd server
go mod download

# 3. Create .env file (copy from .env.example)
cp .env.example .env

# 4. Run the server
go run cmd/main.go
```

## Troubleshooting

### Database Connection Issues

If you see connection errors:
1. Check PostgreSQL is running: `docker-compose ps`
2. Verify port 5433 is available
3. Check database logs: `docker-compose logs postgres`

### Backend Won't Start

1. Check Go is installed: `go version`
2. Install dependencies: `go mod download`
3. Check logs: `docker-compose logs backend`

## Next Steps

1. Connect frontend to `http://localhost:8080`
2. Update frontend API base URL in your React app
3. Add WebSocket connection for real-time updates

