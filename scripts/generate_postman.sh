#!/usr/bin/env bash
# scripts/generate_postman.sh
# Generates Postman collection JSON for the Cauciones API.
# Output: docs/Cauciones-API.postman_collection.json
set -euo pipefail

DOCS_DIR="docs"
POSTMAN_FILE="${DOCS_DIR}/Cauciones-API.postman_collection.json"

echo "▶ Generating Postman collection..."

mkdir -p "${DOCS_DIR}"

cat > "${POSTMAN_FILE}" <<'EOF'
{
  "info": {
    "_postman_id": "cauciones-api-collection",
    "name": "Cauciones API",
    "description": "REST API for managing surety bonds (cauciones) with AWS Lambda backend",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
    "_exporter_id": "franco-personal"
  },
  "item": [
    {
      "name": "Cauciones",
      "item": [
        {
          "name": "Create Surety Bond",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Content-Type",
                "value": "application/json"
              },
              {
                "key": "Authorization",
                "value": "Bearer {{jwt_token}}"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"numero\": \"CB-2024-001\",\n  \"tipo\": \"performance\",\n  \"monto\": 50000,\n  \"moneda\": \"USD\",\n  \"beneficiario\": \"ACME Corp\",\n  \"tomador\": \"John Doe\",\n  \"fecha_emision\": \"2024-01-15T10:00:00Z\",\n  \"fecha_vencimiento\": \"2025-01-15T10:00:00Z\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/cauciones",
              "host": ["{{base_url}}"],
              "path": ["cauciones"]
            },
            "description": "Create a new surety bond"
          },
          "response": []
        },
        {
          "name": "List Surety Bonds",
          "request": {
            "method": "GET",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{jwt_token}}"
              }
            ],
            "url": {
              "raw": "{{base_url}}/cauciones",
              "host": ["{{base_url}}"],
              "path": ["cauciones"]
            },
            "description": "Get list of all surety bonds"
          },
          "response": []
        },
        {
          "name": "Get Surety Bond",
          "request": {
            "method": "GET",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{jwt_token}}"
              }
            ],
            "url": {
              "raw": "{{base_url}}/cauciones/{{caucion_id}}",
              "host": ["{{base_url}}"],
              "path": ["cauciones", "{{caucion_id}}"]
            },
            "description": "Get surety bond by ID"
          },
          "response": []
        },
        {
          "name": "Update Surety Bond",
          "request": {
            "method": "PUT",
            "header": [
              {
                "key": "Content-Type",
                "value": "application/json"
              },
              {
                "key": "Authorization",
                "value": "Bearer {{jwt_token}}"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"monto\": 60000,\n  \"beneficiario\": \"Updated Corp\",\n  \"tomador\": \"Jane Doe\",\n  \"fecha_vencimiento\": \"2025-06-15T10:00:00Z\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/cauciones/{{caucion_id}}",
              "host": ["{{base_url}}"],
              "path": ["cauciones", "{{caucion_id}}"]
            },
            "description": "Update surety bond details"
          },
          "response": []
        },
        {
          "name": "Change Surety Bond Status",
          "request": {
            "method": "PATCH",
            "header": [
              {
                "key": "Content-Type",
                "value": "application/json"
              },
              {
                "key": "Authorization",
                "value": "Bearer {{jwt_token}}"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"estado\": \"active\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/cauciones/{{caucion_id}}/estado",
              "host": ["{{base_url}}"],
              "path": ["cauciones", "{{caucion_id}}", "estado"]
            },
            "description": "Change surety bond status (pending → active → expired or canceled)"
          },
          "response": []
        },
        {
          "name": "Delete Surety Bond",
          "request": {
            "method": "DELETE",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{jwt_token}}"
              }
            ],
            "url": {
              "raw": "{{base_url}}/cauciones/{{caucion_id}}",
              "host": ["{{base_url}}"],
              "path": ["cauciones", "{{caucion_id}}"]
            },
            "description": "Delete surety bond"
          },
          "response": []
        }
      ]
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:3000",
      "type": "string"
    },
    {
      "key": "jwt_token",
      "value": "your-jwt-token-here",
      "type": "string"
    },
    {
      "key": "caucion_id",
      "value": "550e8400-e29b-41d4-a716-446655440000",
      "type": "string"
    }
  ]
}
EOF

echo "✓ Postman collection generated: ${POSTMAN_FILE}"
echo "  Import this file into Postman:"
echo "  1. Open Postman"
echo "  2. Click 'Import'"
echo "  3. Select ${POSTMAN_FILE}"
echo "  4. Update variables: base_url, jwt_token, caucion_id"
