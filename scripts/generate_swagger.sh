#!/usr/bin/env bash
# scripts/generate_swagger.sh
# Parses swaggo annotations from Go source files and generates OpenAPI 2.0 spec.
# Output: docs/swagger.json and docs/swagger.yaml
#
# Annotations live in:
#   internal/adapters/http/swagger_docs.go   → global API metadata
#   internal/adapters/http/handlers/         → @Router, @Param, @Success, @Failure
#   internal/adapters/http/dto/              → request/response schemas
#
# Re-run this script (or `task docs:swagger`) whenever handlers or DTOs change.
set -euo pipefail

DOCS_DIR="docs"
GENERAL_INFO="internal/adapters/http/swagger_docs.go"

echo "▶ Checking for swag CLI..."
if ! command -v swag &>/dev/null; then
  echo "  swag not found — installing via go install..."
  go install github.com/swaggo/swag/cmd/swag@latest
fi

echo "▶ Generating OpenAPI specification from Go annotations..."
swag init \
  --generalInfo "${GENERAL_INFO}" \
  --dir "." \
  --output "${DOCS_DIR}" \
  --outputTypes json,yaml \
  --parseInternal

# Rename to conventional filename
mv -f "${DOCS_DIR}/swagger.json" "${DOCS_DIR}/openapi.json"
mv -f "${DOCS_DIR}/swagger.yaml" "${DOCS_DIR}/openapi.yaml"
# docs.go is swag's runtime file — not needed without a live HTTP server
rm -f "${DOCS_DIR}/docs.go"

echo "✓ OpenAPI specification generated:"
echo "  ${DOCS_DIR}/openapi.json"
echo "  ${DOCS_DIR}/openapi.yaml"

# Keep the legacy heredoc below this comment removed — spec is now code-driven.
exit 0
