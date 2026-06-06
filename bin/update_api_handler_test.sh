#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
GENERATOR="$SCRIPT_DIR/update_api_handler.sh"
FIXTURE_DIR=$(mktemp -d)

cleanup() {
  rm -rf "$FIXTURE_DIR"
}
trap cleanup EXIT

mkdir -p \
  "$FIXTURE_DIR/generated/api" \
  "$FIXTURE_DIR/internal/app/backend/server"

cat >"$FIXTURE_DIR/go.mod" <<'EOF'
module example.com/new-app

go 1.25
EOF

cat >"$FIXTURE_DIR/generated/api/api.gen.go" <<'EOF'
package api

import "github.com/labstack/echo/v4"

type AuthValidateOtpParams struct{}

type ServerInterface interface {
	Existing(ctx echo.Context) error
	AuthValidateOtp(ctx echo.Context, userId string, params AuthValidateOtpParams) error
}
EOF

cat >"$FIXTURE_DIR/internal/app/backend/server/existing.go" <<'EOF'
package server

func (s *Server) Existing(ctx any) error {
	return nil
}
EOF

first_output=$("$GENERATOR" "$FIXTURE_DIR")
handler="$FIXTURE_DIR/internal/app/backend/server/impl_AuthValidateOtp.go"
handler_test="$FIXTURE_DIR/internal/app/backend/server/impl_AuthValidateOtp_test.go"

[[ -f "$handler" ]]
[[ -f "$handler_test" ]]
grep -Fq '"example.com/new-app/generated/api"' "$handler"
grep -Fq 'func (s *Server) AuthValidateOtp(ctx echo.Context, userId string, params api.AuthValidateOtpParams) error' "$handler"
grep -Fq 'return returnNotImplemented(ctx)' "$handler"
grep -Fq 'func TestServer_AuthValidateOtp(t *testing.T)' "$handler_test"
grep -Fq 'svc.AuthValidateOtp(ctx, "", api.AuthValidateOtpParams{})' "$handler_test"
grep -Fq 'created 1 handler pair(s)' <<<"$first_output"
[[ ! -e "$FIXTURE_DIR/internal/app/backend/server/impl_Existing.go" ]]

handler_checksum=$(cksum "$handler")
second_output=$("$GENERATOR" "$FIXTURE_DIR")
[[ "$(cksum "$handler")" == "$handler_checksum" ]]
grep -Fq 'created 0 handler pair(s)' <<<"$second_output"

echo "update_api_handler.sh tests passed"
