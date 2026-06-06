# Go Echo Boilerplate

Go Echo Boilerplate is a reusable foundation for creating new Go HTTP
applications. It provides a production-oriented project structure and includes
authentication and push notifications as the default application systems.

Use this repository as a starting point, replace its module and application
identity, then add business domains without rebuilding common backend
infrastructure. It uses Echo for HTTP routing, OpenAPI as the API contract,
PostgreSQL for persistence, Wire for dependency injection, and `oapi-codegen`
for generated request types and server interfaces.

The boilerplate includes:

- Phone-number registration, OTP verification, and PIN authentication
- JWT creation and validation
- Per-user API access contracts
- Role-based API access assignment
- AES encryption helpers
- Firebase Cloud Messaging integration
- OpenAPI request validation

Authentication provides a ready-to-extend user access foundation.
Notifications are provider-independent internally, with Firebase Cloud
Messaging included as the default adapter.

## Architecture

Requests move through the application in this order:

```text
HTTP request
    |
Echo middleware and OpenAPI validation
    |
Handler (internal/app/backend/server)
    |
Service (internal/service/<domain>)
    |
Repository (internal/service/<domain>)
    |
PostgreSQL
```

Handlers deal with HTTP input and output. Services contain business rules.
Repositories contain database access. External providers are isolated behind
interfaces so the internal application does not depend directly on a vendor.

## Project Structure

```text
.
|-- app/
|   |-- api/
|   `-- backend/
|-- bin/
|-- external/
|-- generated/
|-- internal/
|-- migrations/
|-- .env.example
|-- Makefile
`-- go.mod
```

### `app`

`app` contains the executable application entry points and API contract.

- `app/backend` starts the Echo HTTP server, loads configuration, registers
  middleware, and wires application dependencies.
- `app/backend/wire.go` defines the Wire dependency graph.
- `app/api/src` is the source of truth for the OpenAPI definition.
- `app/api/api.yaml` is the bundled OpenAPI document generated from
  `app/api/src`.

Edit files under `app/api/src` when changing an endpoint. Do not edit generated
API bindings directly.

### `bin`

`bin` contains development and maintenance scripts.

The current `bin/format_sql.sh` script formats or checks SQL files with
`pg_format`:

```bash
./bin/format_sql.sh format migrations
./bin/format_sql.sh check migrations
```

`bin/update_api_handler.sh` reads the generated OpenAPI `ServerInterface` and
creates missing handler and test stubs. Run it through:

```bash
make update-api-handler
```

### `external`

`external` contains adapters for third-party systems. These packages implement
interfaces owned by the internal application.

The current `external/fcm` package implements the push notification provider
using Firebase Cloud Messaging. Application code should use
`internal/service/notification` instead of importing the FCM package directly.

Add other vendor integrations here, such as payment gateways, object storage,
email providers, or external business services.

### `internal`

`internal` contains private application code that cannot be imported by other
Go modules.

- `internal/app/backend/server`: Echo handlers implementing the generated
  OpenAPI server interface
- `internal/service/auth`: authentication services and PostgreSQL repositories
- `internal/service/notification`: provider-independent notification logic
- `internal/echo/middleware`: authentication, authorization, error handling,
  and OpenAPI middleware
- `internal/database`: PostgreSQL connection setup
- `internal/dbtx`: database and transaction abstraction
- `internal/crypto`: encryption and decryption service
- `internal/ctxerr`: contextual application error handling
- `internal/testutil`: shared test assertions
- `internal/util`: parsers and other focused utilities

New domain behavior should follow the existing handler -> service -> repository
flow. Keep one main service and one main repository per domain package.

### `migrations`

`migrations` contains versioned PostgreSQL migrations managed by
`golang-migrate`. Each change has an `up.sql` file to apply it and a `down.sql`
file to reverse it.

Create migrations through the Make target:

```bash
make create_migration add_customer_table
```

Use descriptive names. When adding a table or column, include a PostgreSQL
`COMMENT ON COLUMN` statement that explains the purpose of every new field.

### `generated`

`generated` contains code created by project tools:

- `generated/api`: OpenAPI types, embedded specification, and server interfaces
- `generated/app`: application-generation workspace created by the Makefile

This directory is ignored by Git and can be recreated with `make init`. The
compiled backend executable is written to the repository root as `make`.

## Prerequisites

- Go `1.25.8` or a compatible newer version
- PostgreSQL
- GNU Make
- Node.js and npm
- `swagger-cli`
- `oapi-codegen`
- `wire`
- `mockgen`
- `golang-migrate`
- `pg_format` only when using `bin/format_sql.sh`

Install the Go development tools:

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
go install github.com/google/wire/cmd/wire@latest
go install go.uber.org/mock/mockgen@latest
```

Install the OpenAPI bundler:

```bash
npm install --global @apidevtools/swagger-cli
```

Install `golang-migrate` using the package manager for your operating system.
For example, on macOS:

```bash
brew install golang-migrate
```

## Quick Start

1. Create a PostgreSQL database:

   ```bash
   createdb app
   ```

2. Create the local environment file:

   ```bash
   cp .env.example .env
   ```

3. Update `.env`, especially `DATABASE_URL`, `JWT_SECRET`, and `CRYPTO_KEY`.
   `CRYPTO_KEY` must contain exactly 16, 24, or 32 bytes.

4. Generate the API bindings, Wire injector, mocks, log file, and backend
   binary:

   ```bash
   make init
   ```

5. Apply database migrations:

   ```bash
   make migrate_up db='postgres://postgres:postgres@localhost:5433/app?sslmode=disable'
   ```

   The URL can also be taken from `DATABASE_URL`:

   ```bash
   make migrate_up db="$DATABASE_URL"
   ```

6. Start the API:

   ```bash
   go run ./app/backend
   ```

   The default address is `http://localhost:1323`.

7. Check the server:

   ```bash
   curl http://localhost:1323/api/v1/ping
   ```

   The response should contain:

   ```json
   {
     "message": "PONG"
   }
   ```

`make init` also builds a root-level `make` executable. It can be used instead
of `go run`:

```bash
./make
```

## Configuration

Configuration is read from `.env` and then overridden by process environment
variables.

| Variable                         | Purpose                                             |
| -------------------------------- | --------------------------------------------------- |
| `APP_PORT`                       | Echo server port                                    |
| `SHOW_ERR_MODE`                  | Shows internal errors on stdout when enabled        |
| `DATABASE_URL`                   | PostgreSQL connection URL                           |
| `JWT_SECRET`                     | Secret used to sign JWT access tokens               |
| `JWT_EXPIRATION`                 | JWT lifetime in minutes                             |
| `CRYPTO_KEY`                     | AES key with a length of 16, 24, or 32 bytes        |
| `AUTH_ACCESS_BOOTSTRAP_USER_ID`  | Optional user UUID granted all active API contracts |
| `FCM_SERVICE_ACCOUNT_ID`         | Firebase project ID                                 |
| `GOOGLE_APPLICATION_CREDENTIALS` | Absolute path to Firebase credentials               |

Firebase credentials are required when using push notifications. Keep the
service-account JSON file outside this repository.

## Common Commands

```bash
# Recreate all generated files and the backend executable
make clean init

# Bundle and validate the OpenAPI definition
make clean api.yaml validate-api

# Generate missing handler and test stubs for new OpenAPI operations
make update-api-handler

# Test the handler generator with an isolated fixture
make test-update-api-handler

# Build api-docs.html with Redocly
make api-docs

# Run all tests
go test ./...

# Run tests for one package
go test ./internal/service/auth/...

# Apply all pending migrations
make migrate_up db="$DATABASE_URL"

# Roll back all applied migrations after confirmation
make migrate_down db="$DATABASE_URL"

# Show the current migration version
make migrate_version db="$DATABASE_URL"

# Force a migration version after resolving a dirty migration
make migrate_force db="$DATABASE_URL" version=202606052305
```

The `migrate_down` target runs `golang-migrate down`, which asks for
confirmation and rolls back all applied migrations. To roll back only one
migration, use:

```bash
migrate -path ./migrations -database "$DATABASE_URL" down 1
```

## Adding a Feature

1. Define or update the endpoint under `app/api/src`.
2. Bundle and validate the API:

   ```bash
   make api.yaml validate update-api-handler api-docs
   ```

3. Implement the generated handler in `internal/app/backend/server`.
4. Add business logic to the domain service in `internal/service/<domain>`.
5. Add persistence logic to the domain repository when needed.
6. Add table-driven tests for every changed layer.
7. If an interface changed, regenerate the project:

   ```bash
   make clean init
   ```

8. Run the relevant tests:

   ```bash
   go test ./internal/app/backend/server/... ./internal/service/...
   ```

Handler names must match the OpenAPI `operationId`. Follow the existing file
naming conventions, such as `impl_GetPing.go`, `service_CreateOTP.go`, and
`repo_FindAuth.go`.

`make update-api-handler` generates a missing pair for every new operation:

```text
operationId: AuthValidateOtp

internal/app/backend/server/impl_AuthValidateOtp.go
internal/app/backend/server/impl_AuthValidateOtp_test.go
```

The implementation initially returns HTTP 501 Not Implemented, and the
generated table-driven test verifies that response. The command reads exact
method parameters from the generated `api.ServerInterface`. Existing server
methods and existing target files are skipped, so rerunning it does not
overwrite application code.

## Using This Boilerplate for a New App

1. Create a new repository from this boilerplate.
2. Replace the module path `github.com/masraga/golang-echo-boilerplate` in `go.mod` and Go
   imports with the module path of the new application.
3. Change the OpenAPI title in `app/api/src/main.yaml`.
4. Update `.env.example` with suitable local defaults.
5. Keep or adapt the default authentication and notification systems.
6. Add application domains under `internal/service` and expose them through
   handlers under `internal/app/backend/server`.
7. Regenerate the project:

   ```bash
   make clean init
   ```

8. Run the complete test suite:

   ```bash
   go test ./...
   ```

Do not put business rules in handlers. Add each new domain using the existing
handler -> service -> repository architecture, and place third-party adapters
under `external`.

## API Documentation

The bundled API contract is `app/api/api.yaml`. Generate a standalone HTML
reference with:

```bash
make api-docs
```

Open the generated `api-docs.html` in a browser. For implementation-specific
feature notes, see `.codex/tech-docs`.
