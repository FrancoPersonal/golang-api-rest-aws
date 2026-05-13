## Plan: REST API Cauciones en Go — Arquitectura Hexagonal

**TL;DR:** Proyecto Go 1.26 desde cero con arquitectura hexagonal, API REST de Cauciones desplegada en AWS Lambda + HTTP API Gateway v2, persistencia en DynamoDB, logger de `golang-wappers`, pre-commit con linter + cobertura >80%, deploy con AWS SAM y task runner con Taskfile.

---

### Estructura de directorios

```
golang-api-rest-aws/
├── cmd/lambda/main.go
├── internal/
│   ├── domain/
│   │   ├── caucion.go                          # Entidad + value objects + estados
│   │   ├── errors.go                           # Errores de dominio
│   │   └── ports/
│   │       ├── primary/caucion_usecase.go      # Puerto de entrada (interfaz)
│   │       └── secondary/caucion_repository.go # Puerto de salida (interfaz)
│   ├── application/
│   │   ├── caucion_usecase.go
│   │   └── caucion_usecase_test.go
│   └── adapters/
│       ├── primary/http/
│       │   ├── handler.go
│       │   ├── handler_test.go
│       │   ├── router.go
│       │   ├── request.go
│       │   └── response.go
│       └── secondary/dynamodb/
│           ├── caucion_repository.go
│           └── caucion_repository_test.go
├── pkg/logger/logger.go                        # Wrapper del logger de golang-wappers
├── .github/
│   └── workflows/
│       └── ci.yml                              # GitHub Actions: lint + test en PR/push
├── .gitignore
├── .gitattributes
├── .pre-commit-config.yaml
├── .golangci.yml
├── Taskfile.yml
├── template.yaml                               # AWS SAM: Lambda + HTTP API GW v2 + DynamoDB
├── samconfig.toml                              # SAM deploy config por entorno (dev/prod)
├── README.md                                   # Badge de coverage + instrucciones de uso
└── go.mod
```

---

### Pasos

**Fase 1 — Esqueleto del dominio**
1. `go.mod` — módulo `github.com/FrancoPersonal/golang-api-rest-aws`, Go 1.26, deps: `aws-lambda-go`, `aws-sdk-go-v2`, `golang-wappers`, `uuid`, `aws-lambda-go-api-proxy`, `testify`
2. `pkg/logger/logger.go` — interfaz `Logger` y constructor que envuelve el logger de `golang-wappers`
3. `internal/domain/caucion.go` — entidad `Caucion` con campos: ID, Numero, Tipo, Monto, Moneda, FechaEmision, FechaVencimiento, Estado, Beneficiario, Tomador, CreatedAt, UpdatedAt; tipo `Estado` con constantes y método de transición válida
4. `internal/domain/errors.go` — errores centinela: `ErrNotFound`, `ErrInvalidState`, `ErrInvalidTransition`
5. Puertos: `CaucionUseCase` (primary) y `CaucionRepository` (secondary)

**Fase 2 — Capa de aplicación** *(depende de Fase 1)*
6. `caucion_usecase.go` — implementa los 6 casos de uso: Create, GetByID, List, Update, Delete, ChangeState
7. `caucion_usecase_test.go` — mock de `CaucionRepository` con `testify/mock`, cobertura >80%

**Fase 3 — Adaptadores** *(paralelo, depende de Fase 1)*
8. `dynamodb/caucion_repository.go` — PutItem, GetItem, Scan/Query, UpdateItem, DeleteItem con manejo de errores tipados
9. `dynamodb/caucion_repository_test.go` — usa mock client de `aws-sdk-go-v2`
10. `http/request.go` + `http/response.go` — DTOs con validación básica
11. `http/handler.go` — handlers para POST/GET/GET(list)/PUT/DELETE/PATCH; delega a `CaucionUseCase`
12. `http/handler_test.go` — httptest, mock del use case
13. `http/router.go` — registra rutas en `http.ServeMux`

**Fase 4 — Entry point + Tooling** *(depende de todas las fases)*
14. `cmd/lambda/main.go` — wire manual de dependencias, `httpadapter.NewV2(router)`, `lambda.Start()`
15. `.golangci.yml` — habilita: `errcheck`, `govet`, `staticcheck`, `gosimple`, `unused`, `gofmt`, `goimports`
16. `.pre-commit-config.yaml` — hooks: `golangci-lint run`, `conventional-pre-commit` (Conventional Commits), script de coverage ≥80% + generación del badge de coverage en `README.md` (actualiza el porcentaje en el badge de `shields.io` vía sed o script Go embebido)
17. `Taskfile.yml` — tasks: `lint`, `test`, `coverage`, `build` (cross-compile linux/amd64), `deploy`, `deploy:guided`, `local:api` (sam local start-api), `pre-commit`; la task `coverage` incluye el paso de actualizar el badge en `README.md` y hace `git add README.md` para incluirlo en el commit
18. `template.yaml` — SAM: `AWS::Serverless::Function` (runtime `provided.al2023`), `AWS::Serverless::HttpApi` (HTTP API v2 con CORS), `AWS::DynamoDB::Table` (PAY_PER_REQUEST, `DeletionPolicy: Retain`), env var `CAUCION_TABLE_NAME`
19. `samconfig.toml` — stacks `dev` y `prod` con región, capabilities IAM, resolve_s3
20. `.gitignore` — excluir `bootstrap`, `coverage.out`, `.aws-sam/`, `*.env`
21. `.gitattributes` — forzar LF en todos los archivos de texto; `bootstrap` como binario
22. `.github/workflows/ci.yml` — trigger en PR/push a `main`; jobs: setup-go 1.26 → `golangci-lint-action` → test con race detector → verificar cobertura ≥80% → generar badge con `actions/github-script` o subir a `gist` para exponer el porcentaje en el README
23. `README.md` — badge de coverage dinámico (vía `img.shields.io` apuntando al gist o al artefacto de CI), badge de CI status, instrucciones de setup, `task` commands y endpoint de la API

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
| `cmd/lambda/main.go` | Composición y bootstrap Lambda |
| `template.yaml` | SAM IaC — Lambda, HTTP API GW v2, DynamoDB |
| `samconfig.toml` | Config de deploy SAM por entorno dev/prod |
| `Taskfile.yml` | Task runner — lint, test, build, deploy, local |
| `.pre-commit-config.yaml` | Validación automática pre-commit |
| `.github/workflows/ci.yml` | CI en GitHub Actions |
| `README.md` | Badge de coverage (≥80%) + badge CI + instrucciones |

---

### Verificación

1. `go build ./...` sin errores
2. `task test` pasa sin errores con `-race`
3. `task coverage` muestra ≥80% y no hace exit 1
4. `task lint` pasa sin errores
5. `git commit` con mensaje mal formateado es rechazado por `conventional-pre-commit`
6. `git commit` válido activa hooks y bloquea si linter o coverage fallan
7. `task build` produce binario `bootstrap` para linux/amd64
8. `task local:api` arranca API en `localhost:3000` y responde a todos los endpoints de cauciones
9. Push a rama feature activa `ci.yml` y bloquea merge si lint o coverage fallan
10. `task deploy:guided` despliega stack en AWS y el output `ApiEndpoint` devuelve endpoint funcional
11. El badge de coverage en `README.md` refleja el porcentaje real del último CI exitoso en `main`

---

### Decisiones

- **Inyección de dependencias manual** (no Wire) para mantener simplicidad
- **net/http + `httpadapter/v2`** de `aws-lambda-go-api-proxy` para adaptar `http.Handler` a eventos Lambda HTTP API v2 sin reescribir handlers
- El logger de `golang-wappers` se envuelve detrás de una interfaz local `pkg/logger.Logger` — el dominio no depende de paquetes externos
- Tests usan `testify/mock` para mocks de puertos
- **AWS SAM** como IaC; `template.yaml` define todo el stack sin CDK ni Terraform
- **API Gateway HTTP API v2** (no REST API v1) — menor costo y latencia; soportado por `httpadapter/v2`
- **Taskfile** reemplaza Makefile — mejor soporte cross-platform (Windows/Linux/Mac) y sintaxis YAML
- **Conventional Commits** — formato `type(scope): description`; validado en pre-commit
- **Branch strategy:** `main` estable; features en `feature/*`; no merge sin CI verde
- `DeletionPolicy: Retain` en DynamoDB para proteger datos ante `sam delete`
- **Coverage badge** generado en dos momentos: (1) en **pre-commit** — el hook ejecuta `task coverage`, extrae el porcentaje total, actualiza el badge estático en `README.md` con `sed` y lo agrega al stage con `git add README.md`; (2) en **CI** (`ci.yml`) — misma lógica pero sube el porcentaje a un gist via `actions/github-script` para el badge dinámico de `shields.io`; umbral 80% hard-coded en ambos contextos

---

### Consideraciones

1. **Logger privado:** `golang-wappers` requiere `GOPRIVATE=github.com/FrancoPersonal/*` configurado en el entorno de desarrollo y en CI (GitHub Actions secret `GITHUB_TOKEN` o PAT)
2. **Tabla DynamoDB:** nombre configurable vía variable de entorno `CAUCION_TABLE_NAME` inyectada por SAM
