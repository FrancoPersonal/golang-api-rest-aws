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
- All documentation and code must be in English.
- Use English for README content, comments, logs, error messages, commit messages, variable names, and identifiers.

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

All APIs must be developed in Golang 1.26 or higher.

### Architecture

Use hexagonal architecture:

- domain
- application
- adapters
- infrastructure

Separate:
- handlers
- services
- repositories
- dto
- clients

### AWS

The API must run on AWS Lambda behind API Gateway.

Use:
- aws-lambda-go
- API Gateway Proxy Request/Response
- context.Context

### Rules

- Clean and decoupled code
- Explicit error handling
- Structured logs
- If a reusable library is needed for other projects, it must be created in a private repository and used as a dependency (e.g., github.com/FrancoPersonal/golang-wappers)
- Interfaces for ports
- Unit tests
- DTOs separated from the domain
- Do not use global variables
- Dependency injection
- All documentation and code must be written in English.
- All variable names in code must be written in English.

### Database

If persistence is required:
- use repository pattern
- use DynamoDB as database
- use parameterized queries

### Security

- OAuth2/JWT
- Scope validation
- Do not hardcode secrets
- Use Secrets Manager for credentials

### Serverless

Generate:
- serverless.yml
- serverless.policies.yml
- serverless.infrastructure.yml
- minimum required IAM
- stage-based variables
- optimized packaging
- Lambda binaries must be generated under the build/ directory
- each Lambda must generate its own separated binary artifact
- IAM policies and infrastructure must be split into independent serverless files
- policy and infrastructure files must be independently deployable

### Testing

Generate:
- mocks
- table tests
- coverage

### Expected format

Each endpoint must include:
- handler
- request dto
- response dto
- service
- repository
- tests
