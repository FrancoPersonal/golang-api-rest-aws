# Architecture

## Hexagonal Architecture

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
│                repositories (DynamoDB)               │
└──────────────────────────────────────────────────────┘
```

### Layers

| Layer | Package | Responsibility |
|---|---|---|
| Domain | `internal/domain` | Entities, business rules, port interfaces |
| Application | `internal/application` | Use case orchestration |
| Adapters | `internal/adapters` | HTTP handlers, DTOs, repository implementations |
| Infrastructure | `internal/infrastructure` | Lambda bootstrap, AWS clients wiring |

### Component separation

- `handlers` — HTTP entry points
- `services` — application use cases
- `repositories` — data access (DynamoDB / SQL Server)
- `dto` — request and response shapes
- `clients` — external service clients

---

## Project Structure

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
├── scripts/
│   ├── go/                                     # Go scripts (package_lambda, generate_postman)
│   └── sh/                                     # Shell scripts (lint, coverage, swagger, postman)
├── build/                                      # Separate binaries per Lambda
│   ├── post_cauciones/bootstrap
│   ├── post_cauciones.zip
│   ├── get_cauciones/bootstrap
│   └── get_cauciones.zip
├── docs/                                       # Generated API docs (OpenAPI, Postman, Swagger UI)
├── badges/                                     # Local SVG badges
├── serverless.yml                              # Main API stack
├── serverless.policies.yml                     # Independent IAM policies stack
├── serverless.infrastructure.yml               # Independent infrastructure stack
├── Taskfile.yml                                # Task runner
└── .pre-commit-config.yaml                     # Pre-commit hooks
```
