#!/usr/bin/env bash
set -euo pipefail

FLOCI_ENDPOINT=http://localhost:4566

# ---------------------------------------------------------------------------
# Dependencies
# ---------------------------------------------------------------------------
if ! command -v newman &> /dev/null; then
  echo "Error: newman is not installed."
  echo "Install it with: npm install -g newman"
  exit 1
fi

if ! command -v curl &> /dev/null; then
  echo "Error: curl is not installed."
  exit 1
fi

# ---------------------------------------------------------------------------
# Verify Floci is reachable
# ---------------------------------------------------------------------------
if ! curl -sf "${FLOCI_ENDPOINT}/restapis" > /dev/null 2>&1; then
  echo "Error: Floci is not reachable at ${FLOCI_ENDPOINT}."
  echo "Run 'task local:floci' to start the local environment first."
  exit 1
fi

# ---------------------------------------------------------------------------
# Resolve API Gateway endpoint via Floci REST API (no aws CLI required)
# ---------------------------------------------------------------------------
echo "Resolving API Gateway endpoint..."
API_ID=$(curl -s "${FLOCI_ENDPOINT}/restapis" | python3 -c "
import sys, json
items = json.load(sys.stdin).get('item', [])
match = next((a for a in items if a.get('name') == 'cauciones-local'), None)
print(match['id'] if match else '')
")

if [ -z "$API_ID" ]; then
  echo "Error: cauciones-local REST API not found in Floci."
  echo "Run 'task local:floci' to start and configure the local environment first."
  exit 1
fi

BASE_URL="${FLOCI_ENDPOINT}/restapis/${API_ID}/local/_user_request_"
echo "Running integration tests against: ${BASE_URL}"
echo ""

# ---------------------------------------------------------------------------
# Newman
# ---------------------------------------------------------------------------
mkdir -p build

JWT_SECRET="${JWT_SECRET:-dev-secret}"

newman run docs/Cauciones-API-Integration.postman_collection.json \
  --environment docs/floci.postman_environment.json \
  --env-var "base_url=${BASE_URL}" \
  --env-var "jwt_secret=${JWT_SECRET}" \
  --timeout-request 15000 \
  --reporters cli,junit \
  --reporter-junit-export build/integration-results.xml \
  "$@"
