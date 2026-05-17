#!/usr/bin/env bash
# scripts/sh/generate_postman.sh
# Generates a Postman Collection v2.1 by parsing docs/openapi.json.
# Requires docs/openapi.json to exist — run 'task docs:swagger' first.
# Output: docs/Cauciones-API.postman_collection.json
set -euo pipefail

SPEC="docs/openapi.json"

if [[ ! -f "${SPEC}" ]]; then
  echo "✗ ${SPEC} not found. Run 'task docs:swagger' first."
  exit 1
fi

echo "▶ Generating Postman collection from ${SPEC}..."
go run ./scripts/go/generate_postman.go

echo "  Import docs/Cauciones-API.postman_collection.json into Postman:"
echo "  1. Open Postman → Import"
echo "  2. Select docs/Cauciones-API.postman_collection.json"
echo "  3. Update variables: base_url, jwt_token"
