# Cauciones API

[![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)](https://golang.org)
[![Coverage](badges/coverage.svg)](coverage.out)
[![Lint](badges/lint.svg)](https://golangci-lint.run)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

REST API for **Cauciones** built with Go 1.26+, designed for AWS Lambda behind API Gateway, following **hexagonal architecture**.

---

## Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/cauciones` | Create caucion |
| `GET` | `/cauciones` | List cauciones |
| `GET` | `/cauciones/{id}` | Get caucion by ID |
| `PUT` | `/cauciones/{id}` | Update caucion |
| `DELETE` | `/cauciones/{id}` | Delete caucion |
| `PATCH` | `/cauciones/{id}/estado` | Change status |

### State transitions

```
pending ──► active ──► expired
     │              │
     └──────────────┴──► canceled
```

---

## Documentation

| Topic | File |
|---|---|
| Architecture & project structure | [docs/architecture.md](docs/architecture.md) |
| Setup & prerequisites | [docs/setup.md](docs/setup.md) |
| Development & available tasks | [docs/development.md](docs/development.md) |
| Deploy & environment variables | [docs/deploy.md](docs/deploy.md) |
| CI/CD (pre-commit + CodeBuild) | [docs/ci.md](docs/ci.md) |
| OpenAPI specification | [docs/openapi.json](docs/openapi.json) |
| Postman collection | [docs/Cauciones-API.postman_collection.json](docs/Cauciones-API.postman_collection.json) |
| New endpoint guide | [docs/new-endpoint.md](docs/new-endpoint.md) |

---
