# Copilot Instructions for Go Inventory

## Overview

Go Inventory is a warehouse/pallet management system built with Go 1.24, using the Gin web framework, GORM ORM with MySQL, and Docker for containerization. The API manages pallets, pallet racks, and palletized products with JWT authentication.

**Size:** ~37 Go files, ~2,700 lines of code  
**Type:** REST API backend  
**Port:** 3000

---

## Build & Run Commands

### Prerequisites
- Go 1.24+
- Docker (for local development with MySQL)

### Build (Verified Working)
```bash
# Always run these commands in sequence:
go mod download
go build -buildvcs=false -o main .
```

### Run Tests
```bash
go test ./...
```
> **Note:** There are currently no test files in this repository.

### Linting (Verified Working)
```bash
go vet ./...      # Static analysis
gofmt -l .        # Format check (lists unformatted files)
gofmt -w .        # Auto-format all files
```

### Swagger Generation (Required when modifying Controllers)
```bash
# Install swag if not present:
go install github.com/swaggo/swag/cmd/swag@latest

# Generate/update swagger docs:
swag init -g main.go
```
> **Always run `swag init -g main.go` after modifying any Controller file** that has swagger annotations (`@Summary`, `@Tags`, etc.).

### Docker Development
```bash
make up           # Start containers (app + MySQL)
docker compose up # Alternative
```

---

## Project Architecture

```
go_inventory/
├── main.go                         # Application entry point
├── Container/container.go          # Dependency injection (uber/dig)
├── SupplyInventory/
│   ├── Application/
│   │   ├── Controllers/            # HTTP handlers with swagger annotations
│   │   ├── Middlewares/auth.go     # JWT authentication middleware
│   │   ├── Requests/               # Request DTOs with validation tags
│   │   ├── Responses/              # Response DTOs
│   │   ├── Routes/Api.go           # Route registration
│   │   ├── Services/               # Business logic interfaces & implementations
│   │   └── ApiContracts/           # API contract definitions
│   ├── Domain/                     # Entity definitions (GORM models)
│   └── Infrastructure/
│       ├── Db/                     # Database connection & migrations
│       └── *Repository.go          # Data access layer
├── Helpers/
│   ├── Errors/errors.go            # AppError type for error handling
│   ├── RequestsHelper/             # Request parameter utilities
│   └── Development/                # Debug utilities
├── docs/                           # Generated swagger files (docs.go, swagger.json, swagger.yaml)
├── Makefile                        # Build automation
├── docker-compose.yml              # Docker services (app + MySQL)
├── Dockerfile.dev                  # Development container
├── .air.toml                       # Hot reload config
└── .env.example                    # Environment template
```

---

## Key Patterns & Conventions

### Layer Architecture
1. **Controllers** → Handle HTTP, call Services
2. **Services** → Business logic, interfaces defined here
3. **Repositories** → Data access via GORM
4. **Domain** → Entity structs with GORM/JSON tags

### Dependency Injection
All dependencies are wired via `Container/container.go` using `go.uber.org/dig`.

### Error Handling
Use `errors.NewAppError(message, httpCode)` from `Helpers/Errors/errors.go`.

### Request Validation
Use `binding:"required"` tags on request structs. Format errors with `requestsHelper.FormatValidationErrors(err)`.

### Swagger Annotations
Controllers use swag annotations for API documentation:
```go
// @Summary List pallets
// @Tags Pallets
// @Accept json
// @Produce json
// @Success 200 {array} domain.PalletEntity
// @Router /api/v1/pallets [get]
```

---

## Environment Configuration

Copy `.env.example` to `.env` and configure:
```
BASE_URL=http://localhost
PORT=3000
DB_HOST=db
DB_PORT=3306
DB_DB_NAME=inventory
DB_USER=root
DB_PASSWORD=secret
JWT_SECRET=your-secret-key
```

> **Note:** The app falls back to default Docker MySQL settings if `.env` is missing.

---

## API Routes

All API routes are prefixed with `/api/v1/`:
- `/auth/login`, `/auth/refreshToken` - Authentication (no auth required)
- `/users/create` - User registration (no auth required)
- `/pallets/*` - Pallet CRUD (requires JWT)
- `/racks/*` - Rack CRUD (requires JWT)
- `/pallet/products/*` - Palletized products (requires JWT)
- `/swagger/*` - Swagger UI

---

## Adding New Features

### New Entity
1. Create entity in `SupplyInventory/Domain/`
2. Add to migration in `SupplyInventory/Infrastructure/Db/Migrate.go`
3. Create repository in `SupplyInventory/Infrastructure/`
4. Create service in `SupplyInventory/Application/Services/`
5. Create controller in `SupplyInventory/Application/Controllers/`
6. Create request DTO in `SupplyInventory/Application/Requests/`
7. Register in `Container/container.go`
8. Add routes in `SupplyInventory/Application/Routes/Api.go`
9. Run `swag init -g main.go`

### New Endpoint
1. Add handler method to controller with swagger annotations
2. Register route in controller's `Register*` method
3. Run `swag init -g main.go`

---

## Validation Checklist

Before completing a task, verify:
1. `go build -buildvcs=false -o main .` succeeds
2. `go vet ./...` reports no issues
3. `swag init -g main.go` succeeds (if controllers modified)
4. `go test ./...` passes (when tests exist)

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Build fails with VCS error | Use `-buildvcs=false` flag |
| Import not found | Run `go mod download` first |
| Swagger not updating | Run `swag init -g main.go` |
| Database connection fails | Ensure Docker MySQL is running or check `.env` settings |

---

Trust these instructions. Only search the codebase if information here is incomplete or incorrect.
