# AGENT.md – Agentic AI System Instructions & Project Context

> **Project:** Geospatial Entity Management Backend Service  
> **Repository:** `be-entity-info`  
> **Pair Programming Context:** Human Developer + Agentic AI Assistant (Antigravity)  
> **Target Framework:** Go (Gin) + PostgreSQL 16 + Docker  

---

## 1. Role & Purpose

This repository was developed using an **Agentic AI Pair Programming workflow**. The AI agent acted as a senior software developer assisting in:
- Architectural planning and separation of concerns.
- Idiomatic Go backend implementation (REST API).
- Relational database schema design and migrations.
- Input validation and standardized error handling.
- Containerization with multi-stage Docker and Docker Compose.
- API documentation, Frontend Integration PRD, and verification testing.

---

## 2. Core Architectural Principles & Guardrails

When inspecting, extending, or maintaining this codebase, the agent must adhere to the following principles:

### 2.1 Clean Architecture & Separation of Concerns (SoC)
The project strictly isolates responsibilities across layers:
```
cmd/server/main.go       --> Application entry point, dependency injection & graceful shutdown
internal/
  ├── config/            --> Environment configuration loading (.env & system env)
  ├── database/          --> PostgreSQL connection pooling (pgxpool) & schema migrations
  ├── handler/           --> HTTP transport layer (Gin handlers, request binding, error mapping)
  ├── model/             --> Domain models & DTOs (Data Transfer Objects) with validation tags
  ├── repository/        --> Data access layer (SQL queries using squirrel & pgx)
  ├── response/          --> Standardized JSON response envelope
  ├── router/            --> Gin routing, CORS middleware, and custom validator registrations
  └── service/           --> Business logic layer (ID generation, timestamps, validation)
migrations/              --> SQL migration files (golang-migrate up/down scripts)
docs/                    --> System & integration documentation (PRD, API specs)
```

### 2.2 Lean Engineering (No Over-Engineering)
- Avoid unnecessarily heavy ORM frameworks (e.g. GORM); prefer composable SQL query builders (`squirrel`) with raw connection pools (`pgx/v5`).
- Keep abstractions purposeful—do not introduce unused interfaces or excessive boilerplate.
- Avoid mock or in-memory persistence; all entities are persisted directly to PostgreSQL.

---

## 3. Technology Stack & Dependencies

| Component | Technology | Rationale |
| :--- | :--- | :--- |
| **Language** | Go (`1.24+` / `1.25`) | High performance, type safety, low latency |
| **HTTP Framework** | `github.com/gin-gonic/gin` | Fast routing, middleware support, built-in validation |
| **Database** | PostgreSQL 16 Alpine | ACID compliance, relational integrity, geospatial coordinate storage |
| **DB Driver / Pool** | `github.com/jackc/pgx/v5` | High-performance PostgreSQL driver and connection pooling |
| **Query Builder** | `github.com/Masterminds/squirrel` | Type-safe, composable SQL query builder |
| **Migrations** | `github.com/golang-migrate/migrate/v4` | Automated, versioned database schema migrations |
| **Validation** | `github.com/go-playground/validator/v10` | Struct and custom validator tags |
| **Orchestration** | Docker & Docker Compose | Single-command deployment for DB and API |

---

## 4. Domain Model & Business Logic Standards

### 4.1 Entity Identification
- Entity IDs must follow the format **`ENT-XXXXXX`** (where `XXXXXX` is a 6-character uppercase hex string, generated via `crypto/rand`).

### 4.2 Enumerations
- **Kind**: `'Vehicle'`, `'IoT Device'`, `'Facility'`, `'Asset'`
- **Status**: `'Active'`, `'Idle'`, `'Maintenance'`, `'Offline'`

### 4.3 Coordinate Constraints
- **Latitude**: Decimal between `-90.000000` and `90.000000`.
- **Longitude**: Decimal between `-180.000000` and `180.000000`.
- *Note:* Request DTOs use pointer types (`*float64`) to properly distinguish coordinate `0.0` from absent fields.

### 4.4 Dynamic Attributes
- Stored relationally in table `entity_attributes` with `ON DELETE CASCADE` referencing `entities(id)`.
- Structured as `{ "label": string, "value": string }`.

---

## 5. API Response Contracts

All API endpoints must return a standardized JSON envelope.

### 5.1 Success Envelope
```json
{
  "success": true,
  "data": {},
  "message": "Operation description"
}
```

### 5.2 Error Envelope
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR | NOT_FOUND | INTERNAL_ERROR",
    "message": "Human readable summary",
    "details": [
      {
        "field": "latitude",
        "message": "latitude is required and must be between -90 and 90"
      }
    ]
  }
}
```

---

## 6. Development & Operational Commands

### 6.1 Docker Orchestration (Recommended)
```bash
# Start both PostgreSQL and Go API backend
docker compose up --build -d

# View live API logs
docker compose logs -f api

# Stop all containers (data is preserved in pgdata volume)
docker compose down
```

### 6.2 Local Development Run
```bash
# 1. Start PostgreSQL only
docker compose up -d postgres

# 2. Run backend locally
go run ./cmd/server/main.go
```

### 6.3 Build & Validation
```bash
# Compile server binary
go build ./cmd/server/...

# Format code
go fmt ./...
```

---

## 7. Verification & Smoke Test Matrix

| Action | HTTP Method | Endpoint | Expected Status |
| :--- | :--- | :--- | :--- |
| **List Entities** | `GET` | `/api/v1/entities` | `200 OK` |
| **Filtered List** | `GET` | `/api/v1/entities?kind=Vehicle&status=Active` | `200 OK` |
| **Get Entity by ID** | `GET` | `/api/v1/entities/:id` | `200 OK` / `404 Not Found` |
| **Create Entity** | `POST` | `/api/v1/entities` | `201 Created` / `400 Bad Request` |
| **Update Entity** | `PUT` | `/api/v1/entities/:id` | `200 OK` / `404 Not Found` |
| **Delete Entity** | `DELETE` | `/api/v1/entities/:id` | `200 OK` / `404 Not Found` |
| **Aggregate Metrics** | `GET` | `/api/v1/entities/metrics` | `200 OK` |
| **CORS Preflight** | `OPTIONS` | `/api/v1/entities` | `204 No Content` / `200 OK` |
