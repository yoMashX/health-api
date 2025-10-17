# Healthcare Management System

A full-stack healthcare system with Go API backend, PostgreSQL database, and React frontend, supporting patients, physicians, and prescription drugs with role-based access control.

## Features

### Backend API
- **Entities**: Patients, Physicians, Drugs, Prescriptions
- **RBAC**: Role-based access control via `X-Role` header (physician|patient|admin)
- **Endpoints**:
  - `POST /prescriptions`: Create new prescriptions (validation + authorization)
  - `GET /analytics/top-drugs`: Get top N drugs by quantity in date range

### Frontend UI
- **Authentication**: Role-based login system (fake auth with localStorage)
- **Dashboard**: Role-aware interface with different views per user type
- **Top Drugs Report**: Interactive charts and tables with data filtering

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start the full application stack
docker compose up

# Access the application:
# - Frontend UI: http://localhost:3000
# - Backend API: http://localhost:8080
```

This will:
- Start a PostgreSQL database with the schema and seed data
- Build and run the Go API server with CORS support
- Build and run the React frontend with hot reload
- Set up proper networking between all services

### Manual Setup

1. **Start PostgreSQL**:
```bash
# Using Docker
docker run --name postgres -e POSTGRES_DB=health_api -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15

# Or use your local PostgreSQL installation
```

2. **Set up the database**:
```bash
# Connect to PostgreSQL and run:
psql -h localhost -U postgres -d health_api -f sql/schema.sql
psql -h localhost -U postgres -d health_api -f sql/seed.sql
```

3. **Run the API**:
```bash
go mod tidy
go run main.go
```

## API Usage

### Authentication

All requests require the `X-Role` header with one of: `physician`, `patient`, `admin`

Non-admin roles also require `X-User-ID` header with the user's ID.

### Create Prescription

```bash
curl -X POST http://localhost:8080/prescriptions \
  -H "Content-Type: application/json" \
  -H "X-Role: physician" \
  -H "X-User-ID: 1" \
  -d '{
    "patient_id": 1,
    "physician_id": 1,
    "drug_id": 1,
    "quantity": 30,
    "sig": "Take one tablet daily with food"
  }'
```

### Get Top Drugs

```bash
# Admin view (all prescriptions)
curl "http://localhost:8080/analytics/top-drugs?from=2024-01-01&to=2024-12-31&limit=5" \
  -H "X-Role: admin"

# Patient view (scoped to patient's prescriptions)
curl "http://localhost:8080/analytics/top-drugs?from=2024-01-01&to=2024-12-31&limit=5" \
  -H "X-Role: patient" \
  -H "X-User-ID: 1"

# Physician view (all prescriptions, but could be scoped in future)
curl "http://localhost:8080/analytics/top-drugs?from=2024-01-01&to=2024-12-31&limit=5" \
  -H "X-Role: physician" \
  -H "X-User-ID: 1"
```

## Frontend Usage

### Login System

Access the web interface at `http://localhost:3000` and select your role:

**Test Credentials:**
- **Admin**: No ID required - full system access
- **Physician**: Use IDs 1-5 (e.g., login as Physician with ID 1)
- **Patient**: Use IDs 1-8 (e.g., login as Patient with ID 1)

### Features by Role

**Admin Dashboard:**
- View all patients' prescription data in Top Drugs Report
- See system-wide analytics and statistics

**Physician Dashboard:**
- View Top Drugs Report filtered to your assigned patients only
- Dashboard shows patient management (placeholder)

**Patient Dashboard:**
- View your personal prescription history in Top Drugs Report
- Data filtered to show only your own prescriptions

### Top Drugs Report

Interactive analytics screen featuring:
- **Bar Charts**: Visual comparison of drug quantities and prescription counts
- **Data Table**: Detailed breakdown with rankings and averages
- **Date Filters**: Adjust reporting period (from/to dates)
- **Result Limits**: Show top 5, 10, 15, or 20 drugs
- **Responsive Design**: Works on desktop and mobile devices
- **Accessibility**: Screen reader support and keyboard navigation

## Authorization Rules

- **Physicians**: Can only create prescriptions for patients they are linked to
- **Patients**: Cannot create prescriptions; analytics scoped to their own prescriptions
- **Admin**: Unrestricted access to all endpoints and data

## Database Schema

- **physicians**: Doctor information and license details
- **patients**: Patient demographics
- **drugs**: Medication catalog with generic names and dosage info
- **patient_physicians**: Many-to-many relationship (patients can see multiple doctors)
- **prescriptions**: Links patients, physicians, and drugs with quantities and instructions

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test -v ./internal/api -run TestGetTopDrugsQuery
```

Note: Tests require a test database. Create `health_api_test` database for integration tests.

## Performance

The application includes optimized indexes for common queries:
- Prescription lookups by patient, physician, drug, and date
- Composite index on (prescribed_date, drug_id) for analytics queries
- Patient-physician relationship lookups

## Development

```bash
# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run

# Build binary
go build -o health-api
```