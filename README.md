# Cauciones API

[![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)](https://golang.org)
[![Coverage](badges/coverage.svg)](coverage.out)
[![Lint](badges/lint.svg)](https://golangci-lint.run)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

REST API for **Cauciones** built with Go 1.26+, designed for AWS Lambda behind API Gateway, following **hexagonal architecture**.

---

## Proposed Architecture

```
┌──────────────────────────────────────────────────────┐
│                  Adapters (Primary)                  │
│          handlers | dto | clients (input)            │
└────────────────────┬─────────────────────────────────┘
                     │ primary ports
┌────────────────────▼─────────────────────────────────┐
│                   Application                        │
│                   services/usecases                  │
└────────────────────┬─────────────────────────────────┘
                     │ secondary ports
┌────────────────────▼─────────────────────────────────┐
│                  Domain + Ports                      │
│              entities and business rules             │
└────────────────────┬─────────────────────────────────┘
                     │
┌────────────────────▼─────────────────────────────────┐
                    [![Coverage](badges/coverage.svg)](coverage.out)
│                repositories (DynamoDB)               │
└──────────────────────────────────────────────────────┘
```

Required layers:

- `domain`
- `application`
- `adapters`
- `infrastructure`

Required component separation:

- `handlers`
- `services`
- `repositories`
- `dto`
- `clients`

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/cauciones` | Create caucion |
| `GET` | `/cauciones` | List cauciones |
| `GET` | `/cauciones/{id}` | Get caucion by ID |
| `PUT` | `/cauciones/{id}` | Update caucion |
| `DELETE` | `/cauciones/{id}` | Delete caucion |
| `PATCH` | `/cauciones/{id}/estado` | Change status |

### Expected endpoint format

Each endpoint must include:

- handler
- request dto
- response dto
- service
- repository
- tests

### States and transitions

```
pending ──► active ──► expired
     │              │
     └──────────────┴──► canceled
```

---

## Prerequisites

- [Go 1.26+](https://golang.org/dl/)
- [Task](https://taskfile.dev) — `go install github.com/go-task/task/v3/cmd/task@latest`
- [AWS CLI](https://docs.aws.amazon.com/cli/)
- [Serverless Framework](https://www.serverless.com/framework/docs/getting-started)
- [golangci-lint](https://golangci-lint.run/usage/install/) (recommended: latest)
- [pre-commit](https://pre-commit.com/#install)
- Docker (optional, for local workflows)



### Install dependencies

```bash
# 1. Install Task runner
go install github.com/go-task/task/v3/cmd/task@latest

# 2. Install pre-commit
pip install pre-commit          # Linux/macOS
py -m pip install pre-commit    # Windows

# 3. Install golangci-lint
# Linux/macOS:
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
# Windows (winget):
winget install -e --id GolangCI.golangci-lint

# 4. Register Git hooks (required once per clone)
pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg

# 5. Validate toolchain and hooks
golangci-lint --version
pre-commit run --all-files

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

### Troubleshooting (new machine)

```bash
# If pre-commit does not run on commit, reinstall hooks in this clone:
pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg

# Verify hook files exist:
ls -la .git/hooks | grep -E "pre-commit|commit-msg"

# If golangci-lint reports
# "unsupported version of the configuration: \"\""
# update golangci-lint to latest and re-run:
golangci-lint --version
pre-commit run --all-files
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

### Security

- OAuth2/JWT
- Scope validation
- Do not hardcode secrets
- Use Secrets Manager for credentials

---

## Current project structure

```
.
├── internal/
│   ├── domain/                                 # Entities, errors and ports
│   │   └── ports/
│   │       ├── usecase/
│   │       └── repository/
│   ├── application/                            # Application services
│   │   └── services/
│   ├── adapters/
│   │   ├── http/                               # Router, handlers, DTOs and middleware
│   │   └── repositories/
│   │       ├── dynamodb/
│   │       └── sqlserver/
│   └── infrastructure/
│       ├── clients/
│       └── lambda/
│           ├── post_cauciones/main.go          # Lambda POST /cauciones
│           └── get_cauciones/main.go           # Lambda GET /cauciones
├── pkg/logger/                                 # Logger wrapper
├── scripts/                                    # Coverage, badge and packaging scripts
│   └── package_lambda.go
├── build/                                      # Separate binaries per Lambda
│   ├── post_cauciones/bootstrap
│   ├── post_cauciones.zip
│   ├── get_cauciones/bootstrap
│   └── get_cauciones.zip
├── serverless.yml                              # Main API stack (uses outputs from external stacks)
├── serverless.policies.yml                     # Independent IAM policies stack
├── serverless.infrastructure.yml               # Independent infrastructure stack
├── policies/serverless.yml                     # Policies stack variant
├── infra/serverless.yml                        # Infrastructure stack variant
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
