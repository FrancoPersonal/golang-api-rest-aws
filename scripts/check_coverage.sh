#!/usr/bin/env bash
# scripts/check_coverage.sh
# Runs tests with coverage, enforces an 80% threshold, and updates the badge in README.md.
# Used by both the pre-commit hook and Taskfile.
set -euo pipefail

THRESHOLD=80
COVERAGE_FILE="coverage.out"
README="README.md"

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

# Update badge URL in README.md
BADGE_URL="https://img.shields.io/badge/coverage-${TOTAL_INT}%25-${COLOR}"
sed -i.bak "s|https://img\.shields\.io/badge/coverage-[^)]*|${BADGE_URL}|g" "${README}" && rm -f "${README}.bak"

# Stage the updated README so it is included in the commit
git add "${README}"

echo "✓ Coverage badge updated to ${TOTAL_INT}% (${COLOR})"
