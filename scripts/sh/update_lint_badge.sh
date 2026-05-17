#!/usr/bin/env bash
# scripts/sh/update_lint_badge.sh
# Runs golangci-lint, generates a local SVG badge under badges/lint.svg,
# and updates README.md to point to the local badge.
set -euo pipefail

README="README.md"
BADGE_DIR="badges"
BADGE_FILE="${BADGE_DIR}/lint.svg"

echo "▶ Running golangci-lint..."
if golangci-lint run ./... 2>/dev/null; then
  STATUS="passing"
  COLOR="brightgreen"
  echo "✓ Lint passed"
else
  STATUS="failing"
  COLOR="red"
  echo "✗ Lint failed"
fi

# Map color name to hex value used in the SVG
case "${COLOR}" in
  brightgreen) HEX="#4c1" ;;
  *)           HEX="#e05d44" ;;
esac

# Generate the local SVG badge (flat style, same layout as shields.io)
mkdir -p "${BADGE_DIR}"
LABEL="lint"
VALUE="${STATUS}"
LABEL_W=45
VALUE_W=64
TOTAL_W=$(( LABEL_W + VALUE_W ))
cat > "${BADGE_FILE}" <<SVG
<svg xmlns="http://www.w3.org/2000/svg" width="${TOTAL_W}" height="20">
  <linearGradient id="s" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <rect rx="3" width="${TOTAL_W}" height="20" fill="#555"/>
  <rect rx="3" x="${LABEL_W}" width="${VALUE_W}" height="20" fill="${HEX}"/>
  <rect rx="3" width="${TOTAL_W}" height="20" fill="url(#s)"/>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="$(( LABEL_W / 2 ))" y="15" fill="#010101" fill-opacity=".3">${LABEL}</text>
    <text x="$(( LABEL_W / 2 ))" y="14">${LABEL}</text>
    <text x="$(( LABEL_W + VALUE_W / 2 ))" y="15" fill="#010101" fill-opacity=".3">${VALUE}</text>
    <text x="$(( LABEL_W + VALUE_W / 2 ))" y="14">${VALUE}</text>
  </g>
</svg>
SVG

# Update README.md to reference the local badge file
sed -i.bak "s|!\[Lint\]([^)]*)|![Lint](${BADGE_FILE})|g" "${README}" && rm -f "${README}.bak"

git add "${BADGE_FILE}" "${README}"

echo "✓ Local lint badge written to ${BADGE_FILE} (${STATUS}, ${COLOR})"

# Exit with failure if lint failed so callers know
[[ "${STATUS}" == "passing" ]]
