# Development

## Quick start

```bash
# Install pre-commit hooks (one-time per clone)
pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg

task tidy        # tidy and verify Go modules
task test        # run tests with race detector
task lint        # run golangci-lint
task coverage    # enforce ≥80% coverage and update badge
task build       # build Lambda binaries and zip artifacts
task local:api   # start API locally via Serverless offline → http://localhost:3000
```

---

## Available Tasks

Run `task` (no arguments) to list every task.

### Development

| Task | Description |
|---|---|
| `task lint` | Run golangci-lint across all packages |
| `task test` | Run unit tests with race detector |
| `task coverage` | Run tests, enforce ≥80% coverage and update badge |
| `task build` | Build Lambda binaries and zip artifacts under `build/` |
| `task tidy` | Tidy and verify Go modules |
| `task local:api` | Start the API locally via Serverless offline (requires Docker) |

### Documentation

| Task | Description |
|---|---|
| `task docs` | Generate Swagger spec + Postman collection |
| `task docs:swagger` | Generate `docs/openapi.json` and `docs/openapi.yaml` from code annotations |
| `task docs:swagger:open` | Regenerate spec and open Swagger UI in the browser |
| `task docs:postman` | Generate `docs/Cauciones-API.postman_collection.json` |

### Badges

| Task | Description |
|---|---|
| `task badge` | Update both lint and coverage badges in README |
| `task badge:lint` | Run golangci-lint and update lint badge |
| `task badge:coverage` | Run tests and update coverage badge |

### Deploy

| Task | Description |
|---|---|
| `task deploy:dev` | Build + deploy policies, infrastructure and API to dev |
| `task deploy:prod` | Deploy policies, infrastructure and API to prod |
| `task deploy:policies:dev` | Deploy only the IAM policies stack to dev |
| `task deploy:infrastructure:dev` | Deploy only the infrastructure stack to dev |

---

## Build artifacts

```
build/
├── post_cauciones/bootstrap   ← Linux/ARM64 binary
├── post_cauciones.zip         ← Lambda artifact for POST /cauciones
├── get_cauciones/bootstrap
└── get_cauciones.zip          ← Lambda artifact for GET /cauciones
```

---

## API documentation

```bash
task docs:swagger       # regenerate docs/openapi.json + docs/openapi.yaml
task docs:swagger:open  # regenerate and open Swagger UI in browser
task docs:postman       # regenerate Postman collection from openapi.json
```

Swagger annotations live in:

| File | Purpose |
|---|---|
| `internal/adapters/http/swagger_docs.go` | Global API metadata (title, host, auth) |
| `internal/adapters/http/handlers/caucion_handler.go` | `@Router`, `@Param`, `@Success`, `@Failure` |
| `internal/adapters/http/dto/` | Request/response schemas with `example` tags |
