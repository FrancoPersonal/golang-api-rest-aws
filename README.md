# Cauciones API

[![CI](https://github.com/FrancoPersonal/golang-api-rest-aws/actions/workflows/ci.yml/badge.svg)](https://github.com/FrancoPersonal/golang-api-rest-aws/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-0%25-red)](coverage.out)
[![Lint](https://img.shields.io/badge/lint-failing-red)](https://golangci-lint.run)
[![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

REST API for **Cauciones** built with Go 1.26+, designed for AWS Lambda behind API Gateway, following **hexagonal architecture**.

---

## Arquitectura propuesta

```
┌──────────────────────────────────────────────────────┐
│                  Adapters (Primary)                  │
│          handlers | dto | clients (entrada)          │
└────────────────────┬─────────────────────────────────┘
                     │ puertos primarios
┌────────────────────▼─────────────────────────────────┐
│                   Application                        │
│                   services/usecases                  │
└────────────────────┬─────────────────────────────────┘
                     │ puertos secundarios
┌────────────────────▼─────────────────────────────────┐
│                  Domain + Ports                      │
│              entidades y reglas de negocio           │
└────────────────────┬─────────────────────────────────┘
                     │
┌────────────────────▼─────────────────────────────────┐
│              Infrastructure / Adapters               │
│                repositories (DynamoDB)               │
└──────────────────────────────────────────────────────┘
```

Capas obligatorias:

- `domain`
- `application`
- `adapters`
- `infrastructure`

Separacion obligatoria de componentes:

- `handlers`
- `services`
- `repositories`
- `dto`
- `clients`

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/cauciones` | Crear caución |
| `GET` | `/cauciones` | Listar cauciones |
| `GET` | `/cauciones/{id}` | Obtener caución por ID |
| `PUT` | `/cauciones/{id}` | Actualizar caución |
| `DELETE` | `/cauciones/{id}` | Eliminar caución |
| `PATCH` | `/cauciones/{id}/estado` | Cambiar estado |

### Formato esperado por endpoint

Cada endpoint debe incluir:

- handler
- request dto
- response dto
- service
- repository
- tests

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
- [AWS CLI](https://docs.aws.amazon.com/cli/)
- [Serverless Framework](https://www.serverless.com/framework/docs/getting-started)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [pre-commit](https://pre-commit.com/#install)
- Docker (optional, for local workflows)



### Install dependencies

```bash
# 1. Install Task runner
go install github.com/go-task/task/v3/cmd/task@latest

# 2. Install pre-commit
pip install pre-commit          # Linux/macOS
py -m pip install pre-commit    # Windows

# 3. Install golangci-lint and register the git hooks
task setup

# Run full pipeline (tidy → build → lint → test + coverage gate ≥ 80%)
task check

# Individual tasks
task test        # run tests with race detector
task lint        # run golangci-lint
task coverage    # open HTML coverage report
task docs        # serve godoc on :6060
```

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

# Build output (separated per Lambda)
# build/post_cauciones/bootstrap
# build/post_cauciones.zip
# build/get_cauciones/bootstrap
# build/get_cauciones.zip

# Start locally (requires Docker)
task local:api
# → API available at http://localhost:3000
```

---

## Deploy

```bash
# Review and adjust stage variables
serverless print --stage dev

# Deploy policies independently
serverless deploy --config serverless.policies.yml --stage dev

# Deploy infrastructure independently
serverless deploy --config serverless.infrastructure.yml --stage dev

# Deploy API (consumes outputs from policies + infrastructure stacks)
serverless deploy --stage dev
serverless deploy --stage prod

# Equivalent task commands
task deploy:policies:dev
task deploy:infrastructure:dev
task deploy:dev
```

### Environment variables (Lambda)

| Variable | Description |
|----------|-------------|
| `CAUCION_TABLE_NAME` | DynamoDB table name (injected by Serverless Framework) |
| `JWT_SECRET` | JWT signing secret used by authorization middleware |
| `LOG_LEVEL` | Log verbosity: `debug`, `info`, `warn`, `error` (default: `info`) |

### Seguridad

- OAuth2/JWT
- Validacion de scopes
- No hardcodear secretos
- Secrets Manager para credenciales

---

## Estructura actual del proyecto

```
.
├── internal/
│   ├── domain/                                 # Entidades, errores y puertos
│   │   └── ports/
│   │       ├── usecase/
│   │       └── repository/
│   ├── application/                            # Servicios de aplicacion
│   │   └── services/
│   ├── adapters/
│   │   ├── http/                               # Router, handlers, DTOs y middleware
│   │   └── repositories/
│   │       ├── dynamodb/
│   │       └── sqlserver/
│   └── infrastructure/
│       ├── clients/
│       └── lambda/
│           ├── post_cauciones/main.go          # Lambda POST /cauciones
│           └── get_cauciones/main.go           # Lambda GET /cauciones
├── pkg/logger/                                 # Wrapper de logger
├── scripts/                                    # Scripts de coverage, badge y packaging
│   └── package_lambda.go
├── build/                                      # Compilados separados por Lambda
│   ├── post_cauciones/bootstrap
│   ├── post_cauciones.zip
│   ├── get_cauciones/bootstrap
│   └── get_cauciones.zip
├── serverless.yml                              # API principal (usa outputs de stacks externos)
├── serverless.policies.yml                     # Stack independiente de IAM policies
├── serverless.infrastructure.yml               # Stack independiente de infraestructura
├── policies/serverless.yml                     # Variante de stack de policies
├── infra/serverless.yml                        # Variante de stack de infraestructura
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
