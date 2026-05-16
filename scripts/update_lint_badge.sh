#!/usr/bin/env bash
# scripts/update_lint_badge.sh
# Runs golangci-lint and updates the lint badge in README.md.
set -euo pipefail

README="README.md"

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

# Update badge URL in README.md
BADGE_URL="https://img.shields.io/badge/lint-${STATUS}-${COLOR}"
sed -i.bak "s|https://img\.shields\.io/badge/lint-[^)]*|${BADGE_URL}|g" "${README}" && rm -f "${README}.bak"

git add "${README}"

echo "✓ Lint badge updated: ${STATUS} (${COLOR})"

# Exit with failure if lint failed so callers know
[[ "${STATUS}" == "passing" ]]
