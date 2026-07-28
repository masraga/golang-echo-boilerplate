#!/usr/bin/env bash

set -euo pipefail

WIRE_VERSION="v0.7.0"
OAPI_CODEGEN_VERSION="v2.5.0"
MOCKGEN_VERSION="v0.6.0"
MIGRATE_VERSION="v4.19.0"

TOOLS=(
  "github.com/google/wire/cmd/wire@${WIRE_VERSION}"
  "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}"
  "go.uber.org/mock/mockgen@${MOCKGEN_VERSION}"
)

echo "Installing Go tools..."

for tool in "${TOOLS[@]}"; do
    echo "-> ${tool}"
    go install "${tool}"
done

echo "-> github.com/golang-migrate/migrate/v4/cmd/migrate@${MIGRATE_VERSION}"
CGO_ENABLED=0 go install -tags postgres \
    github.com/golang-migrate/migrate/v4/cmd/migrate@"${MIGRATE_VERSION}"

echo "Done!"