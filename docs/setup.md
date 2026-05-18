# Setup

## Prerequisites

- [Go 1.26+](https://golang.org/dl/)
- [Task](https://taskfile.dev) — `go install github.com/go-task/task/v3/cmd/task@latest`
- [AWS CLI](https://docs.aws.amazon.com/cli/)
- [Serverless Framework](https://www.serverless.com/framework/docs/getting-started)
- [golangci-lint](https://golangci-lint.run/usage/install/) (recommended: latest)
- [pre-commit](https://pre-commit.com/#install)
- Docker (optional, for local workflows)

---

## Install dependencies

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
```

---

## Private module setup

```bash
# Required to access github.com/FrancoPersonal/golang-wappers
export GOPRIVATE=github.com/FrancoPersonal/*
export GONOSUMCHECK=github.com/FrancoPersonal/*
go mod download
```

---

## Troubleshooting (new machine)

```bash
# If pre-commit does not run on commit, reinstall hooks in this clone:
pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg

# Verify hook files exist:
ls -la .git/hooks | grep -E "pre-commit|commit-msg"

# If golangci-lint reports "unsupported version of the configuration":
golangci-lint --version
pre-commit run --all-files
```
