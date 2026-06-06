#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR=${1:-.}
GENERATED_API_PATH="$ROOT_DIR/generated/api/api.gen.go"
SERVER_DIR="$ROOT_DIR/internal/app/backend/server"
GO_MOD_PATH="$ROOT_DIR/go.mod"

trim() {
  local value=$1
  value="${value#"${value%%[![:space:]]*}"}"
  value="${value%"${value##*[![:space:]]}"}"
  printf '%s' "$value"
}

is_builtin_type() {
  case "$1" in
    bool|byte|complex64|complex128|error|float32|float64|int|int8|int16|int32|int64|rune|string|uint|uint8|uint16|uint32|uint64|uintptr|any)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

render_type() {
  local type_name
  type_name=$(trim "$1")

  case "$type_name" in
    echo.Context)
      printf '%s' "$type_name"
      ;;
    \**)
      printf '*%s' "$(render_type "${type_name#\*}")"
      ;;
    \[\]*)
      printf '[]%s' "$(render_type "${type_name#\[\]}")"
      ;;
    map\[*|interface\{\})
      printf '%s' "$type_name"
      ;;
    *.*)
      printf '%s' "$type_name"
      ;;
    *)
      if is_builtin_type "$type_name"; then
        printf '%s' "$type_name"
      else
        printf 'api.%s' "$type_name"
      fi
      ;;
  esac
}

zero_value() {
  local type_name
  type_name=$(render_type "$1")

  case "$type_name" in
    string)
      printf '""'
      ;;
    bool)
      printf 'false'
      ;;
    error|any|interface\{\}|\**|\[\]*|map\[*)
      printf 'nil'
      ;;
    byte|complex64|complex128|float32|float64|int|int8|int16|int32|int64|rune|uint|uint8|uint16|uint32|uint64|uintptr)
      printf '0'
      ;;
    api.*)
      printf '%s{}' "$type_name"
      ;;
    *)
      printf '%s{}' "$type_name"
      ;;
  esac
}

require_file() {
  if [[ ! -f "$1" ]]; then
    echo "required file not found: $1" >&2
    exit 1
  fi
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "required command not found: $1" >&2
    exit 1
  fi
}

server_method_exists() {
  local operation_id=$1
  local file

  while IFS= read -r file; do
    if grep -Eq "^[[:space:]]*func[[:space:]]+\\([^)]*\\*Server\\)[[:space:]]+${operation_id}\\(" "$file"; then
      return 0
    fi
  done < <(find "$SERVER_DIR" -maxdepth 1 -type f -name '*.go' ! -name '*_test.go' -print)

  return 1
}

parse_parameters() {
  local raw_parameters=$1
  local parameter
  local name
  local type_name
  local rendered_type

  PARAMETERS=""
  TEST_ARGUMENTS=""
  CONTEXT_NAME=""
  USES_API_TYPE=0

  while IFS= read -r parameter; do
    parameter=$(trim "$parameter")
    [[ -z "$parameter" ]] && continue

    name=${parameter%%[[:space:]]*}
    type_name=$(trim "${parameter#"$name"}")
    if [[ -z "$name" || -z "$type_name" ]]; then
      echo "unsupported generated parameter: $parameter" >&2
      return 1
    fi

    rendered_type=$(render_type "$type_name")
    if [[ -n "$PARAMETERS" ]]; then
      PARAMETERS="$PARAMETERS, "
      TEST_ARGUMENTS="$TEST_ARGUMENTS, "
    fi
    PARAMETERS="${PARAMETERS}${name} ${rendered_type}"

    if [[ "$rendered_type" == "echo.Context" ]]; then
      CONTEXT_NAME=$name
      TEST_ARGUMENTS="${TEST_ARGUMENTS}ctx"
    else
      TEST_ARGUMENTS="${TEST_ARGUMENTS}$(zero_value "$type_name")"
    fi

    if [[ "$rendered_type" == *api.* ]]; then
      USES_API_TYPE=1
    fi
  done < <(printf '%s\n' "$raw_parameters" | awk -v RS=, '{ print }')

  if [[ -z "$CONTEXT_NAME" ]]; then
    echo "generated handler has no echo.Context parameter" >&2
    return 1
  fi
}

write_handler() {
  local output_path=$1
  local operation_id=$2

  {
    echo "package server"
    echo
    echo "import ("
    echo $'\t'"\"github.com/labstack/echo/v4\""
    if [[ "$USES_API_TYPE" -eq 1 ]]; then
      echo $'\t'"\"${MODULE_PATH}/generated/api\""
    fi
    echo ")"
    echo
    echo "func (s *Server) ${operation_id}(${PARAMETERS}) error {"
    echo $'\t'"return returnNotImplemented(${CONTEXT_NAME})"
    echo "}"
  } >"$output_path"
}

write_handler_test() {
  local output_path=$1
  local operation_id=$2

  {
    echo "package server_test"
    echo
    echo "import ("
    echo $'\t'"\"net/http\""
    echo $'\t'"\"net/http/httptest\""
    echo $'\t'"\"testing\""
    echo
    echo $'\t'"\"github.com/labstack/echo/v4\""
    if [[ "$USES_API_TYPE" -eq 1 ]]; then
      echo $'\t'"\"${MODULE_PATH}/generated/api\""
    fi
    echo $'\t'"\"${MODULE_PATH}/internal/app/backend/server\""
    echo ")"
    echo
    echo "func TestServer_${operation_id}(t *testing.T) {"
    echo $'\t'"type test struct {"
    echo $'\t\t'"name         string"
    echo $'\t\t'"expectedCode int"
    echo $'\t'"}"
    echo
    echo $'\t'"tests := []test{"
    echo $'\t\t'{
    echo $'\t\t\t'"name:         \"not implemented\","
    echo $'\t\t\t'"expectedCode: http.StatusNotImplemented,"
    printf '\t\t},\n'
    echo $'\t'"}"
    echo
    echo $'\t'"for _, tt := range tests {"
    echo $'\t\t'"t.Run(tt.name, func(t *testing.T) {"
    echo $'\t\t\t'"e := echo.New()"
    echo $'\t\t\t'"req := httptest.NewRequest(http.MethodGet, \"/\", nil)"
    echo $'\t\t\t'"rec := httptest.NewRecorder()"
    echo $'\t\t\t'"ctx := e.NewContext(req, rec)"
    echo $'\t\t\t'"svc := server.NewServer(server.ServerOpts{})"
    echo
    echo $'\t\t\t'"if err := svc.${operation_id}(${TEST_ARGUMENTS}); err != nil {"
    echo $'\t\t\t\t'"t.Fatal(err)"
    echo $'\t\t\t'"}"
    echo $'\t\t\t'"if rec.Code != tt.expectedCode {"
    echo $'\t\t\t\t'"t.Fatalf(\"expected status %d, got %d\", tt.expectedCode, rec.Code)"
    echo $'\t\t\t'"}"
    echo $'\t\t'"})"
    echo $'\t'"}"
    echo "}"
  } >"$output_path"
}

require_command awk
require_command find
require_command gofmt
require_command grep
require_file "$GO_MOD_PATH"
require_file "$GENERATED_API_PATH"

if [[ ! -d "$SERVER_DIR" ]]; then
  echo "server directory not found: $SERVER_DIR" >&2
  exit 1
fi

MODULE_PATH=$(awk '$1 == "module" { print $2; exit }' "$GO_MOD_PATH")
if [[ -z "$MODULE_PATH" ]]; then
  echo "module path not found in $GO_MOD_PATH" >&2
  exit 1
fi

SIGNATURES=$(awk '
  /^type ServerInterface interface \{/ {
    in_interface = 1
    next
  }
  in_interface && /^}/ {
    exit
  }
  in_interface {
    line = $0
    sub(/^[[:space:]]+/, "", line)
    if (line ~ /^[[:alpha:]_][[:alnum:]_]*\(.*\) error$/) {
      print line
    }
  }
' "$GENERATED_API_PATH")

if [[ -z "$SIGNATURES" ]]; then
  echo "ServerInterface methods not found in $GENERATED_API_PATH" >&2
  exit 1
fi

created=0
while IFS= read -r signature; do
  operation_id=${signature%%(*}
  raw_parameters=${signature#*(}
  raw_parameters=${raw_parameters%) error}

  if server_method_exists "$operation_id"; then
    echo "skip $operation_id: handler already exists"
    continue
  fi

  handler_path="$SERVER_DIR/impl_${operation_id}.go"
  test_path="$SERVER_DIR/impl_${operation_id}_test.go"
  if [[ -e "$handler_path" || -e "$test_path" ]]; then
    echo "skip $operation_id: target file already exists"
    continue
  fi

  parse_parameters "$raw_parameters"

  handler_tmp=$(mktemp "$SERVER_DIR/.impl_${operation_id}.XXXXXX")
  test_tmp=$(mktemp "$SERVER_DIR/.impl_${operation_id}_test.XXXXXX")
  cleanup_tmp() {
    rm -f "$handler_tmp" "$test_tmp"
  }
  trap cleanup_tmp EXIT

  write_handler "$handler_tmp" "$operation_id"
  write_handler_test "$test_tmp" "$operation_id"
  gofmt -w "$handler_tmp" "$test_tmp"

  mv "$handler_tmp" "$handler_path"
  mv "$test_tmp" "$test_path"
  trap - EXIT

  created=$((created + 1))
  echo "create $handler_path"
  echo "create $test_path"
done <<<"$SIGNATURES"

echo "created $created handler pair(s)"
