# ChatBoltz Service Setup Guide

## Prerequisites
- Go 1.21+
- PostgreSQL
- Redis (optional, for caching/queues)

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd chatboltz-service
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Configure your `.env` file with your database credentials and other settings.

4. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Service

### Development Mode (with Air)
If you have `air` installed for live reloading:
```bash
air
```

### Standard Run
```bash
go run cmd/server/main.go
```

## Database Migration
Ensure your database is running and reachable. The application uses GORM for auto-migration in development mode.

## API Documentation
See [API Reference](api-reference.md) for detailed endpoint documentation.
