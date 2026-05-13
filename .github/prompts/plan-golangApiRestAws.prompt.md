## Plan: REST API Cauciones en Go — Arquitectura Hexagonal

**TL;DR:** Proyecto Go 1.26 desde cero con arquitectura hexagonal, API REST de Cauciones desplegada en AWS Lambda + API Gateway, persistencia en DynamoDB, logger de `golang-wappers`, y pre-commit con linter + cobertura >80%.

---

### Estructura de directorios

```
golang-api-rest-aws/
├── cmd/lambda/main.go                          # Bootstrap Lambda
├── internal/
│   ├── domain/
│   │   ├── caucion.go                          # Entidad + value objects + estados
│   │   ├── errors.go                           # Errores de dominio
│   │   └── ports/
│   │       ├── primary/caucion_usecase.go      # Puerto de entrada (interfaz)
│   │       └── secondary/caucion_repository.go # Puerto de salida (interfaz)
│   ├── application/
│   │   ├── caucion_usecase.go                  # Casos de uso (implementa puerto primario)
│   │   └── caucion_usecase_test.go
│   └── adapters/
│       ├── primary/http/
│       │   ├── handler.go                      # Handlers HTTP
│       │   ├── handler_test.go
│       │   ├── router.go                       # Setup de rutas net/http
│       │   ├── request.go                      # DTOs de entrada
│       │   └── response.go                     # DTOs de salida
│       └── secondary/dynamodb/
│           ├── caucion_repository.go           # Adaptador DynamoDB
│           └── caucion_repository_test.go
├── pkg/logger/logger.go                        # Wrapper del logger de golang-wappers
├── .pre-commit-config.yaml
├── .golangci.yml
├── Makefile
└── go.mod
```

---

### Pasos

**Fase 1 — Esqueleto del dominio**
1. `go.mod` — módulo `github.com/FrancoPersonal/golang-api-rest-aws`, Go 1.26, deps: `aws-lambda-go`, `aws-sdk-go-v2`, `golang-wappers`, `uuid`
2. `pkg/logger/logger.go` — interfaz `Logger` y constructor que envuelve el logger de `golang-wappers`
3. `internal/domain/caucion.go` — entidad `Caucion` con campos: ID, Numero, Tipo, Monto, Moneda, FechaEmision, FechaVencimiento, Estado, Beneficiario, Tomador, CreatedAt, UpdatedAt; tipo `Estado` con constantes y método de transición válida
4. `internal/domain/errors.go` — errores centinela: `ErrNotFound`, `ErrInvalidState`, `ErrInvalidTransition`
5. Puertos: `CaucionUseCase` (primary) y `CaucionRepository` (secondary)

**Fase 2 — Capa de aplicación** *(depende de Fase 1)*
6. `caucion_usecase.go` — implementa los 6 casos de uso: Create, GetByID, List, Update, Delete, ChangeState
7. `caucion_usecase_test.go` — mock de `CaucionRepository` con `testify/mock`, cobertura >80%

**Fase 3 — Adaptadores** *(paralelo, depende de Fase 1)*
8. `dynamodb/caucion_repository.go` — PutItem, GetItem, Scan/Query, UpdateItem, DeleteItem con manejo de errores tipados
9. `dynamodb/caucion_repository_test.go` — usa `smithy-go` mock client
10. `http/request.go` + `http/response.go` — DTOs con validación básica
11. `http/handler.go` — handlers para POST/GET/GET(list)/PUT/DELETE/PATCH; delega a `CaucionUseCase`
12. `http/handler_test.go` — httptest, mock del use case
13. `http/router.go` — registra rutas en `http.ServeMux`

**Fase 4 — Entry point + Tooling** *(depende de todas las fases)*
14. `cmd/lambda/main.go` — wire manual de dependencias, `httpadapter.New(router)`, `lambda.Start()`
15. `.golangci.yml` — habilita: `errcheck`, `govet`, `staticcheck`, `gosimple`, `unused`, `gofmt`, `goimports`
16. `.pre-commit-config.yaml` — hooks locales: `golangci-lint run` y script de coverage threshold ≥80%
17. `Makefile` — targets: `lint`, `test`, `coverage`, `build`, `pre-commit`

---

### Archivos clave a crear

| Archivo | Rol |
|---|---|
| `go.mod` | Módulo raíz |
| `internal/domain/ports/primary/caucion_usecase.go` | Contrato de entrada para handlers |
| `internal/domain/ports/secondary/caucion_repository.go` | Contrato de salida para DynamoDB |
| `internal/application/caucion_usecase.go` | Lógica de negocio pura |
| `internal/adapters/secondary/dynamodb/caucion_repository.go` | I/O con DynamoDB |
| `internal/adapters/primary/http/handler.go` | REST → use case |
| `cmd/lambda/main.go` | Composición y bootstrap |
| `.pre-commit-config.yaml` | Validación automática pre-commit |

---

### Verificación

1. `go build ./...` sin errores
2. `go test -coverprofile=coverage.out ./...` y `go tool cover -func=coverage.out` muestra ≥80% total
3. `golangci-lint run` pasa sin errores
4. `git commit` activa el pre-commit y bloquea si linter o coverage fallan
5. Deploy local: `GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda` produce binario Lambda válido

---

### Decisiones

- **Inyección de dependencias manual** (no frameworks como Wire) para mantener simplicidad
- **net/http + `httpadapter`** de `aws-lambda-go-api-proxy` para adaptar el `http.Handler` estándar a eventos Lambda sin reescribir handlers
- El logger de `golang-wappers` se envuelve detrás de una interfaz local `pkg/logger.Logger`, así el dominio no depende de un paquete externo
- Tests usan `testify/mock` para mocks de puertos; el adaptador DynamoDB se testa con mock client de `aws-sdk-go-v2`

---

### Consideraciones

1. **Logger privado:** Como `golang-wappers` es privado, el `go.mod` necesitará configurar `GONOSUMCHECK` y `GOFLAGS=-insecure` o un `replace` directive local hasta confirmar acceso. ¿Tienes un token GOPATH/GOPROXY configurado para ese repo?
2. **Tabla DynamoDB:** ¿El nombre de la tabla y partition key deben ser configurables por variable de entorno (ej. `CAUCION_TABLE_NAME`) o hardcodeados para esta primera versión?
