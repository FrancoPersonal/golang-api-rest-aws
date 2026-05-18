#!/usr/bin/env bash
# scripts/sh/check_coverage.sh
# Runs tests with coverage, enforces an 80% threshold, generates a local
# SVG badge under badges/coverage.svg, and updates README.md to point to it.
# Used by both the pre-commit hook and Taskfile.
set -euo pipefail

THRESHOLD=80
COVERAGE_FILE="coverage.out"
README="README.md"
BADGE_DIR="badges"
BADGE_FILE="${BADGE_DIR}/coverage.svg"

echo "▶ Running tests with coverage..."
# Exclude lambda bootstrap mains and the packaging script from coverage:
# those are wiring/infrastructure code that cannot be unit-tested without
# a live Lambda runtime or real filesystem tooling.
PKGS=$(go list ./... | grep -v \
  -e 'internal/infrastructure/lambda' \
  -e '^github.com/FrancoPersonal/golang-api-rest-aws/scripts$')
go test -coverprofile="${COVERAGE_FILE}" ${PKGS}

TOTAL=$(go tool cover -func="${COVERAGE_FILE}" | grep "^total:" | awk '{print $3}' | tr -d '%')
TOTAL_INT=${TOTAL%.*}

echo "   Coverage: ${TOTAL}%"

# Enforce threshold
if (( TOTAL_INT < THRESHOLD )); then
  echo "✗ Coverage ${TOTAL}% is below the required ${THRESHOLD}%"
  exit 1
fi

# Determine badge color
if   (( TOTAL_INT >= 80 )); then COLOR="brightgreen"
elif (( TOTAL_INT >= 60 )); then COLOR="yellow"
else                              COLOR="red"
fi

# Map color name to hex value used in the SVG
case "${COLOR}" in
  brightgreen) HEX="#4c1"  ;;
  yellow)      HEX="#dfb317" ;;
  *)           HEX="#e05d44" ;;
esac

# Generate the local SVG badge (flat style, same layout as shields.io)
mkdir -p "${BADGE_DIR}"
LABEL="coverage"
VALUE="${TOTAL_INT}%"
LABEL_W=70
VALUE_W=46
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
    <text x="35" y="15" fill="#010101" fill-opacity=".3">${LABEL}</text>
    <text x="35" y="14">${LABEL}</text>
    <text x="$(( LABEL_W + VALUE_W / 2 ))" y="15" fill="#010101" fill-opacity=".3">${VALUE}</text>
    <text x="$(( LABEL_W + VALUE_W / 2 ))" y="14">${VALUE}</text>
  </g>
</svg>
SVG

# Update README.md to reference the local badge file
sed -i.bak "s|!\[Coverage\]([^)]*)|![Coverage](${BADGE_FILE})|g" "${README}" && rm -f "${README}.bak"

# Stage both artefacts so they are included in the commit
git add "${BADGE_FILE}" "${README}"

echo "✓ Local badge written to ${BADGE_FILE} (${TOTAL_INT}%, ${COLOR})"
