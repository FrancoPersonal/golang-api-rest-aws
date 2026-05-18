# CI/CD

## Pre-commit hooks

Registered via `.pre-commit-config.yaml`. Run automatically on every `git commit`:

| Hook | Script | Trigger |
|---|---|---|
| `golangci-lint` + lint badge | `scripts/sh/update_lint_badge.sh` | `.go` files changed |
| Coverage check ≥80% + badge | `scripts/sh/check_coverage.sh` | any file changed |
| Swagger generation | `scripts/sh/generate_swagger.sh` | `.go` files changed |
| Conventional Commit | `compilerla/conventional-pre-commit` | commit message |

Install hooks (once per clone):

```bash
pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg
```

Run manually against all files:

```bash
pre-commit run --all-files
```

---

## AWS CodeBuild

Two buildspec files handle the CI/CD pipeline:

| File | Stage | Description |
|---|---|---|
| `buildspec.build.yml` | Build | tidy → lint → test → coverage → build Lambda artifacts |
| `buildspec.deploy.yml` | Deploy | deploy policies → infrastructure → API via Serverless Framework |

### Build pipeline (`buildspec.build.yml`)

1. Install Go toolchain and golangci-lint
2. `task tidy` — verify modules
3. `task lint` — fail on lint errors
4. `task coverage` — enforce ≥80% coverage gate
5. `task build` — produce `build/post_cauciones.zip` and `build/get_cauciones.zip`

Artifacts saved: `build/post_cauciones.zip`, `build/get_cauciones.zip`, `badges/`

### Deploy pipeline (`buildspec.deploy.yml`)

1. Install Node.js and Serverless Framework
2. Deploy IAM policies stack
3. Deploy infrastructure stack
4. Deploy API stack

Stage is controlled via the `STAGE` environment variable (`dev` by default).
JWT_SECRET is injected from AWS Systems Manager Parameter Store at `/cauciones/jwt-secret`.

---

## Commit convention

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
```

Allowed types: `feat`, `fix`, `chore`, `test`, `refactor`, `docs`, `ci`, `build`, `perf`, `revert`

Examples:
```
feat(api): add GET /cauciones/{id} endpoint
fix(auth): handle expired JWT token correctly
chore(ci): update golangci-lint version
```
