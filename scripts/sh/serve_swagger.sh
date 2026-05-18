#!/usr/bin/env bash
# scripts/sh/serve_swagger.sh
# Generates a self-contained HTML file with Swagger UI pointing to the local
# OpenAPI spec and opens it in the default browser.
# The spec JSON is embedded directly into the HTML to avoid CORS issues with
# the file:// protocol — no HTTP server required.
set -euo pipefail

DOCS_DIR="docs"
SPEC_FILE="${DOCS_DIR}/openapi.json"
HTML_FILE="${DOCS_DIR}/swagger-ui.html"

if [[ ! -f "${SPEC_FILE}" ]]; then
  echo "✗ ${SPEC_FILE} not found. Run 'task docs:swagger' first."
  exit 1
fi

echo "▶ Building Swagger UI HTML..."

SPEC_JSON=$(cat "${SPEC_FILE}")

cat > "${HTML_FILE}" <<HTMLEOF
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Cauciones API — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; }
    #swagger-ui .topbar { background-color: #1a1a2e; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    const spec = ${SPEC_JSON};

    SwaggerUIBundle({
      spec: spec,
      dom_id: '#swagger-ui',
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.SwaggerUIStandalonePreset
      ],
      layout: 'BaseLayout',
      deepLinking: true,
      displayRequestDuration: true,
      filter: true
    });
  </script>
</body>
</html>
HTMLEOF

echo "✓ Swagger UI generated: ${HTML_FILE}"
echo "▶ Opening in browser..."

# Cross-platform open
if command -v xdg-open &>/dev/null; then
  xdg-open "${HTML_FILE}"
elif command -v open &>/dev/null; then
  open "${HTML_FILE}"
elif command -v start &>/dev/null; then
  start "${HTML_FILE}"
else
  echo "  Could not detect a browser opener. Open manually:"
  echo "  ${HTML_FILE}"
fi
