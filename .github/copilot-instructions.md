# Project Instructions

This repository uses Go and AWS Lambda with hexagonal architecture.

## Rule Priority (No Ambiguity)

- Highest priority: everything under `## Copilot Instructions`.
- The sections above `## Copilot Instructions` are baseline guidance and apply only when they do not conflict with it.
- If two rules conflict, follow `## Copilot Instructions` without exceptions.
- If there is still ambiguity after applying this priority, stop and ask the user before implementing changes.

## Core Principles

- Think before coding: state assumptions and ask when requirements are unclear.
- Keep solutions simple: implement only what was requested.
- Make surgical changes: modify only files and lines required by the task.
- Use goal-driven execution: define a verifiable success check for each change.

## Architecture Boundaries

- Domain logic belongs in `internal/domain`.
- Application orchestration belongs in `internal/application`.
- HTTP adapters belong in `internal/adapters/http`.
- Repository adapters (DynamoDB/SQL Server) belong in `internal/adapters/repositories`.
- Lambda bootstrap and wiring belong in `internal/infrastructure/lambda`.

Do not move business rules into adapters.

## Coding Conventions

- Match existing style and naming patterns in the touched file.
- Do not refactor unrelated code.
- Keep public interfaces stable unless a task explicitly asks for API changes.
- If your change introduces dead code, remove only what your change made unused.

## Testing and Verification

Before considering a task done, run relevant checks:

- `task test` for unit tests
- `task lint` for lint checks
- `task coverage` when behavior changes or new paths are added

If a command cannot run in the environment, state that clearly in the final report.

## API and State Rules

- Preserve current endpoint contracts unless explicitly asked to change them.
- Respect caucion state transitions documented in `README.md`.
- Validate input at the HTTP boundary; keep business validation in use cases.

## Scope Control

- No speculative abstractions.
- No new dependencies unless required.
- No broad formatting-only edits.
- Keep diffs minimal and focused on the request.

## Copilot Instructions

Todas las APIs deben desarrollarse en Golang 1.26 o superior.

### Arquitectura

Usar arquitectura hexagonal:

- domain
- application
- adapters
- infrastructure

Separar:
- handlers
- services
- repositories
- dto
- clients

### AWS

La API debe ejecutarse en AWS Lambda detras de API Gateway.

Usar:
- aws-lambda-go
- API Gateway Proxy Request/Response
- context.Context

### Reglas

- Codigo limpio y desacoplado
- Manejo explicito de errores
- Logs estructurados
- de necesitarse alguna biblioteca que pueda ser reutilizada en otros proyectos debe ser creada en un repositorio privado y usada como dependencia (ej: github.com/FrancoPersonal/golang-wappers)
- Interfaces para puertos
- Unit tests
- DTOs separados del dominio
- No usar variables globales
- Inyeccion de dependencias

### Base de datos

Si hay persistencia:
- usar repository pattern
- DynamoDB como base de datos
- queries parametrizadas

### Seguridad

- OAuth2/JWT
- Validacion de scopes
- No hardcodear secretos
- Secrets Manager para credenciales

### Serverless

Generar:
- serverless.yml
- serverless.policies.yml
- serverless.infrastructure.yml
- IAM minimo necesario
- variables por stage
- empaquetado optimizado
- los compilados de Lambda deben generarse en la carpeta build/
- cada Lambda debe generar su propio compilado separado (un artefacto independiente por funcion)
- las politicas IAM y la infraestructura deben estar separadas en archivos serverless independientes
- los archivos de politicas e infraestructura deben poder desplegarse de forma independiente

### Testing

Generar:
- mocks
- table tests
- coverage

### Formato esperado

Cada endpoint debe incluir:
- handler
- request dto
- response dto
- service
- repository
- tests
