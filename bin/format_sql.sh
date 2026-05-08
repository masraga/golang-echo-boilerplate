#!/usr/bin/env bash

set -euo pipefail

MODE=${1:-format}   # format | check
TARGET_DIR=${2:-.}

# Ensure pg_format exists
if ! command -v pg_format >/dev/null 2>&1; then
  echo "Error: pg_format is not installed."
  exit 1
fi

echo "Mode: $MODE"
echo "Target: $TARGET_DIR"

EXIT_CODE=0

# Find all .sql files recursively
while IFS= read -r -d '' file; do
  if [[ "$MODE" == "format" ]]; then
    echo "Formatting: $file"
    pg_format -i "$file"
  elif [[ "$MODE" == "check" ]]; then
    formatted=$(pg_format "$file")
    current=$(cat "$file")

    if [[ "$formatted" != "$current" ]]; then
      echo "Needs formatting: $file"
      EXIT_CODE=1
    fi
  else
    echo "Invalid mode: $MODE (use 'format' or 'check')"
    exit 1
  fi
done < <(find "$TARGET_DIR" -type f -name "*.sql" -print0)

exit $EXIT_CODE