# Cauciones API

[![CI](https://github.com/FrancoPersonal/golang-api-rest-aws/actions/workflows/ci.yml/badge.svg)](https://github.com/FrancoPersonal/golang-api-rest-aws/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-0%25-red)](coverage.out)
[![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

REST API for **Cauciones** built with Go 1.26, deployed on AWS Lambda + HTTP API Gateway v2, backed by DynamoDB, following **hexagonal architecture**.

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                  Primary Adapters                    │
│           HTTP Handler (net/http + httpadapter/v2)   │
└────────────────────┬─────────────────────────────────┘
                     │ primary port (CaucionUseCase)
┌────────────────────▼─────────────────────────────────┐
│               Application Layer                      │
│                caucion_usecase.go                    │
└────────────────────┬─────────────────────────────────┘
                     │ secondary port (CaucionRepository)
┌────────────────────▼─────────────────────────────────┐
│              Secondary Adapters                      │
│           DynamoDB Repository                        │
└──────────────────────────────────────────────────────┘
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/cauciones` | Crear caución |
| `GET` | `/cauciones` | Listar cauciones |
| `GET` | `/cauciones/{id}` | Obtener caución por ID |
| `PUT` | `/cauciones/{id}` | Actualizar caución |
| `DELETE` | `/cauciones/{id}` | Eliminar caución |
| `PATCH` | `/cauciones/{id}/estado` | Cambiar estado |

### Estados y transiciones

```
pendiente ──► vigente ──► vencida
     │              │
     └──────────────┴──► cancelada
```

---

## Prerequisites

- [Go 1.26+](https://golang.org/dl/)
- [Task](https://taskfile.dev) — `go install github.com/go-task/task/v3/cmd/task@latest`
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [pre-commit](https://pre-commit.com/#install)
- Docker (for `task local:api`)

### Private module setup

```bash
# Required to access github.com/FrancoPersonal/golang-wappers
export GOPRIVATE=github.com/FrancoPersonal/*
export GONOSUMCHECK=github.com/FrancoPersonal/*
go mod download
```

---

## Development

```bash
# Install pre-commit hooks (one-time)
pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg

# Run all checks (lint + coverage + badge update)
task pre-commit

# Run tests
task test

# Run linter
task lint

# Check coverage (enforces ≥80% and updates badge)
task coverage

# Build Lambda binary
task build

# Start locally (requires Docker)
task local:api
# → API available at http://localhost:3000
```

---

## Deploy

```bash
# First-time interactive deploy
task deploy:guided

# Subsequent deploys
task deploy:dev   # dev stack
task deploy:prod  # prod stack
```

### Environment variables (Lambda)

| Variable | Description |
|----------|-------------|
| `CAUCION_TABLE_NAME` | DynamoDB table name (injected by SAM) |
| `LOG_LEVEL` | Log verbosity: `debug`, `info`, `warn`, `error` (default: `info`) |

---

## Project structure

```
.
├── cmd/lambda/main.go                          # Lambda bootstrap & DI wiring
├── internal/
│   ├── domain/                                 # Entities, errors, state machine
│   │   └── ports/
│   │       ├── primary/caucion_usecase.go      # Input port interface
│   │       └── secondary/caucion_repository.go # Output port interface
│   ├── application/                            # Business logic (use cases)
│   └── adapters/
│       ├── primary/http/                       # HTTP handlers + router
│       └── secondary/dynamodb/                 # DynamoDB repository
├── pkg/logger/                                 # Logger interface + wappers adapter
├── template.yaml                               # AWS SAM infrastructure
├── samconfig.toml                              # SAM deploy config (dev/prod)
├── Taskfile.yml                                # Task runner
└── .pre-commit-config.yaml                     # Pre-commit hooks
```

---

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on every push and PR to `main`:

1. `golangci-lint` — linting
2. `go test -race` — unit tests with race detector
3. Coverage enforcement ≥ **80%**
4. Badge auto-update in `README.md` on merge to `main`
