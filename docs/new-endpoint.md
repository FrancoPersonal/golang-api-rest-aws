# New Endpoint Guide

Step-by-step reference for adding an endpoint following the hexagonal architecture of this project.
Each layer has a single responsibility — do not skip layers or place logic in the wrong place.

> **Example used throughout this guide:** `GET /cauciones/{id}` — get a single surety bond by ID.

---

## Layer map

```
HTTP Request
     │
     ▼
┌─────────────────────────────────────────────┐
│  1. DTO (request.go / response.go)          │  Parse & validate shape
├─────────────────────────────────────────────┤
│  2. Handler (caucion_handler.go)            │  HTTP → domain, domain → HTTP
├─────────────────────────────────────────────┤
│  3. Router (router.go)                      │  Route dispatch + middleware chain
├─────────────────────────────────────────────┤
│  4. Use Case port (ports/usecase/)          │  Contract between adapter and app
├─────────────────────────────────────────────┤
│  5. Service (application/services/)         │  Business logic and orchestration
├─────────────────────────────────────────────┤
│  6. Repository port (ports/repository/)     │  Contract between app and data layer
├─────────────────────────────────────────────┤
│  7. Repository impl (adapters/repositories/)│  DynamoDB query
└─────────────────────────────────────────────┘
```

---

## Step 1 — Domain input/output (`internal/domain/caucion.go`)

Add an input struct for the operation if it needs parameters beyond the entity itself.
Simple lookups by ID may not need a new input struct.

**Rules:**
- No framework types — plain Go structs only.
- No JSON tags — those belong in the DTO.
- Export only what the application layer needs.

```go
// Only add if the operation requires non-trivial input parameters.
type GetSuretyBondInput struct {
    ID string
}
```

---

## Step 2 — Repository port (`internal/domain/ports/repository/caucion_repository.go`)

Add the new method to the interface if the operation requires data access.

**Rules:**
- Accepts `context.Context` as the first parameter — always.
- Returns domain types (`domain.SuretyBond`), never DTO or infrastructure types.
- Returns `error` as the last return value — use domain errors (`domain.ErrNotFound`, etc.).

```go
type SuretyBondRepository interface {
    Create(ctx context.Context, suretyBond domain.SuretyBond) error
    List(ctx context.Context) ([]domain.SuretyBond, error)
    // ✅ New:
    GetByID(ctx context.Context, id string) (domain.SuretyBond, error)
}
```

---

## Step 3 — Use case port (`internal/domain/ports/usecase/caucion_usecase.go`)

Add the new method to the use case interface.

**Rules:**
- One method per operation — no fat interfaces.
- Same signature conventions as the repository port.
- The application service must implement this interface.

```go
type SuretyBondUseCase interface {
    Create(ctx context.Context, input domain.CreateSuretyBondInput) (domain.SuretyBond, error)
    List(ctx context.Context) ([]domain.SuretyBond, error)
    // ✅ New:
    GetByID(ctx context.Context, id string) (domain.SuretyBond, error)
}
```

---

## Step 4 — Application service (`internal/application/services/caucion_service.go`)

Implement the new use case method. This is where **business logic and validation** live.

**Rules:**
- Validate input here — not in the handler, not in the repository.
- Call the repository port — never a concrete implementation.
- Map repository errors to domain errors if needed.
- No HTTP concepts, no AWS types.

```go
func (s *SuretyBondService) GetByID(ctx context.Context, id string) (domain.SuretyBond, error) {
    if strings.TrimSpace(id) == "" {
        return domain.SuretyBond{}, domain.ErrInvalidInput
    }

    bond, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return domain.SuretyBond{}, err
    }

    return bond, nil
}
```

**Test requirements** (`caucion_service_test.go`):
- Table-driven tests covering: valid input, empty ID, not found, repository error.
- Use a mock that implements `SuretyBondRepository`.

---

## Step 5 — DTO (`internal/adapters/http/dto/request.go` / `response.go`)

Define request and response shapes. Only add a request DTO if the operation has a body.
Read-only operations with path/query params do not need a request DTO.

**Rules:**
- Use `json` struct tags matching the API contract.
- Add `example` tags for every field — required for Swagger generation.
- No domain types in DTOs.
- Response shape is already handled by `dto.Success` / `dto.Fail`.

```go
// Only needed for operations with a request body (POST, PUT, PATCH).
// For GET /cauciones/{id} no request DTO is needed — the ID comes from the path.
```

---

## Step 6 — Handler (`internal/adapters/http/handlers/caucion_handler.go`)

Add a new method on `SuretyBondHandler`. This method is the bridge between HTTP and the use case.

**Rules:**
- Parse path/query parameters from `req.PathParameters` or `req.QueryStringParameters`.
- Unmarshal request body only for mutating operations.
- Call the use case — never the repository or service directly.
- Map errors with `mapError(err)`.
- Add complete swaggo annotations (see below).

```go
// GetSuretyBond returns a surety bond by ID.
//
// @Summary      Get a surety bond by ID
// @Tags         Cauciones
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Surety bond ID"
// @Success      200  {object}  dto.StandardResponse
// @Failure      400  {object}  dto.StandardResponse
// @Failure      401  {object}  dto.StandardResponse
// @Failure      404  {object}  dto.StandardResponse
// @Failure      500  {object}  dto.StandardResponse
// @Router       /cauciones/{id} [get]
func (h *SuretyBondHandler) GetSuretyBond(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    id := req.PathParameters["id"]
    if id == "" {
        return dto.Fail(400, "missing_param", "id is required"), nil
    }

    bond, err := h.useCase.GetByID(ctx, id)
    if err != nil {
        return mapError(err), nil
    }

    return dto.Success(200, bond), nil
}
```

**Swaggo annotation requirements:**

| Annotation | Required | Notes |
|---|---|---|
| `@Summary` | ✅ | One-line description |
| `@Tags` | ✅ | Must match existing tag (`Cauciones`) |
| `@Accept` | Mutating only | `json` for POST/PUT/PATCH |
| `@Produce` | ✅ | `json` |
| `@Security` | ✅ | `BearerAuth` — all endpoints require auth |
| `@Param` | ✅ | One per path param, query param, and body |
| `@Success` | ✅ | All expected 2xx codes |
| `@Failure` | ✅ | At minimum: 400, 401, 500. Add 404 if applicable |
| `@Router` | ✅ | Exact path and HTTP method |

**Test requirements** (`caucion_handler_test.go`):
- Table-driven tests covering: valid request, missing ID, not found, service error.
- Mock the `SuretyBondUseCase` interface.

---

## Step 7 — Repository implementation (`internal/adapters/repositories/dynamodb/caucion_repository.go`)

Implement the new method from the repository port.

**Rules:**
- Use parameterized DynamoDB expressions — never string interpolation.
- Return domain errors (`domain.ErrNotFound`, `domain.ErrInternal`), not AWS errors.
- Guard against empty `tableName` or nil `client`.

```go
func (r *SuretyBondRepository) GetByID(ctx context.Context, id string) (domain.SuretyBond, error) {
    if r.tableName == "" || r.client == nil {
        return domain.SuretyBond{}, domain.ErrInternal
    }

    out, err := r.client.GetItem(ctx, &awsdynamodb.GetItemInput{
        TableName: aws.String(r.tableName),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: id},
        },
    })
    if err != nil {
        return domain.SuretyBond{}, domain.ErrInternal
    }
    if out.Item == nil {
        return domain.SuretyBond{}, domain.ErrNotFound
    }

    var bond domain.SuretyBond
    if err := attributevalue.UnmarshalMap(out.Item, &bond); err != nil {
        return domain.SuretyBond{}, domain.ErrInternal
    }

    return bond, nil
}
```

**Test requirements** (`caucion_repository_test.go`):
- Table-driven tests covering: found, not found, DynamoDB error, nil client.
- Use a mock that implements the DynamoDB client interface (`PutItemClient` or equivalent).

---

## Step 8 — Router (`internal/adapters/http/router.go`)

Register the route and wire the middleware chain.

**Rules:**
- Add a field to `Router` for the new handler.
- Wrap with `middleware.Chain(handler, logging, jwtAuth)` in `NewRouter`.
- Add the dispatch case in `Handle` using `req.PathParameters` for path params.

```go
// Add field:
type Router struct {
    createSuretyBond middleware.Handler
    listSuretyBonds  middleware.Handler
    getSuretyBond    middleware.Handler  // ✅ New
}

// Wire in NewRouter:
get := middleware.Chain(
    handler.GetSuretyBond,
    middleware.LoggingMiddleware(logger),
    middleware.JWTAuthMiddleware(jwtSecret),
)
return &Router{..., getSuretyBond: get}

// Dispatch in Handle:
if method == "GET" && req.PathParameters["id"] != "" {
    return r.getSuretyBond(ctx, req)
}
```

---

## Step 9 — Errors (`internal/domain/errors.go`)

Add a new sentinel error only if the operation introduces a failure mode not yet covered.

```go
// Existing errors cover most cases:
var (
    ErrNotFound     = errors.New("not found")
    ErrConflict     = errors.New("conflict")
    ErrInvalidInput = errors.New("invalid input")
    ErrInternal     = errors.New("internal error")
)
// Add a new one only if truly distinct.
```

---

## Checklist

Use this before considering the endpoint done:

```
Domain
  [ ] Input struct added (if needed)
  [ ] New domain error added (if needed)

Ports
  [ ] SuretyBondRepository interface updated
  [ ] SuretyBondUseCase interface updated

Application
  [ ] Service method implemented with business validation
  [ ] Service unit tests (table-driven, ≥80% coverage on new code)

Adapters — HTTP
  [ ] Request DTO defined with json + example tags (if body needed)
  [ ] Handler method added with full swaggo annotations
  [ ] Handler unit tests (table-driven)

Adapters — Repository
  [ ] DynamoDB method implemented with parameterized queries
  [ ] Repository unit tests (table-driven)

Router
  [ ] Route registered and middleware chain wired

Verification
  [ ] task lint   → passes
  [ ] task test   → passes
  [ ] task coverage → ≥80%
  [ ] task docs:swagger → openapi.json updated with new route
```
