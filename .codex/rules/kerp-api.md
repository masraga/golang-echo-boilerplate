# KERP API Rules

These rules apply every time Codex works in this project.

## Architecture

Use the handler > service > repository mechanism.

- Handlers live in `internal/app/backend/server`.
- Services live in `internal/service/*`.
- Repositories live in `internal/service/*`.
- Each service package must have one main service type defined in `internal/service/[domain]/service.go`.
- Do not create multiple service instances in one service package. Add new domain behavior as methods on the existing service type.
- Each service package must have one main repository type defined in `internal/service/[domain]/repository.go`.
- Do not create multiple repository instances in one service package. Add new persistence behavior as methods on the existing repository type.
- All service and repository interfaces for a package must live in `internal/service/[domain]/interface.go`.
- Use existing package conventions before adding new structure.

## Technical Docs

Application technical docs live in `.codex/tech-docs`.

- Before modifying application code, read `.codex/tech-docs/README.md`, `.codex/tech-docs/dictionary.md`, and any relevant feature docs under `.codex/tech-docs/features`.
- Use these docs as the feature navigator and dictionary for application behavior, API contracts, important files, and known implementation notes.
- After modifying selected application code, update the relevant technical docs in the same change so the docs continue to match the implementation.
- When creating a new feature, add a feature doc under `.codex/tech-docs/features/[feature-name].md` and link it from `.codex/tech-docs/README.md`.

## File Naming

Use these file names for new backend code:

- Handler implementation: `internal/app/backend/server/impl_[HandlerName].go`
- Handler test: `internal/app/backend/server/impl_[HandlerName]_test.go`
- Service implementation: `internal/service/[domain]/service_[ServiceFunctionName].go`
- Service test: `internal/service/[domain]/service_[ServiceFunctionName]_test.go`
- Repository implementation: `internal/service/[domain]/repo_[RepositoryFunctionName].go`
- Repository test: `internal/service/[domain]/repo_[RepositoryFunctionName]_test.go`
- Main service type and constructor: `internal/service/[domain]/service.go`
- Main repository type and constructor: `internal/service/[domain]/repository.go`
- Service and repository interfaces: `internal/service/[domain]/interface.go`

## Tests

Every test must be table-driven.

Use this common shape:

```go
type args struct {
	// inputs
}

type fields struct {
	// mocked dependencies
}

type test struct {
	name     string
	args     args
	fields   fields
	expected expected
	mock     func(tt *test, ctrl *gomock.Controller)
}

tests := []test{
	// cases
}

for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		// setup
		// execute
		// assert
	})
}
```

Prefer `testutil.RequireResult` for service and repository result checks. Prefer `testutil.RequireHttpResultJson` for handler HTTP response checks.

## Generated Code

- Every time an interface is added or changed, run `make clean init` so generated gomock files and other generated artifacts stay in sync.
- Do this before running tests that depend on gomock.

## Handler Tests

Handler tests must follow the style of `internal/app/backend/server/impl_RegisterPhoneNumber_test.go`.

- Test function name format: `TestServer_[HandlerName]`.
- Example: handler `VerifyNewAuthUserOTP` must use `func TestServer_VerifyNewAuthUserOTP(t *testing.T)`.
- Build requests with `httptest.NewRequest`.
- Set `echo.HeaderContentType` to `echo.MIMEApplicationJSON` for JSON bodies.
- Use Echo context from `e.NewContext(req, rec)`.
- Mock services with gomock.
- Construct the server with `server.NewServer(server.ServerOpts{...})`.
- Assert JSON responses with `testutil.RequireHttpResultJson`.

## Service Tests

Service tests must follow the style of `internal/service/auth/service_CreateOTP_test.go`.

- Test package should match the existing external test package pattern, for example `package auth_test`.
- Test function name format: `Test[ServiceType]_[ServiceFunctionName]`.
- Example: service type `AuthService` and method `CreateOTP` must use `func TestAuthService_CreateOTP(t *testing.T)`.
- Define `type expected = testutil.Result[[OutputType]]`.
- Mock repositories with gomock.
- Build the service through the domain constructor, for example `auth.NewAuthService(auth.AuthServiceOpts{...})`.
- Include required shared dependencies such as `ctxerr.NewCtxErr(ctxerr.CtxErrOpts{})`.
- Assert with `testutil.RequireResult`.

The service type name is defined in each domain's `service.go`.

## Repository Tests

Repository tests must follow the style of `internal/service/auth/repo_FindAuth_test.go`.

- Test package should match the existing external test package pattern, for example `package auth_test`.
- Test function name format: `Test[RepositoryType]_[RepositoryFunctionName]`.
- Example: repository type `AuthRepository` and method `FindAuth` must use `func TestAuthRepository_FindAuth(t *testing.T)`.
- Use `go-sqlmock` for database expectations.
- Use `dbtx.DbTx{Db: dbMock}` when constructing repositories.
- Build repositories through the domain constructor, for example `auth.NewAuthRepository(auth.AuthRepositoryOpts{...})`.
- Include `sqlf.PostgreSQL` and `ctxerr.NewCtxErr(ctxerr.CtxErrOpts{})` when required by the constructor.
- Assert with `testutil.RequireResult`.

## API Definitions

This application uses `oapi-codegen`.

- API source lives in `app/api`.
- Main API source is `app/api/src/main.yaml`.
- Request and response schema files must live in `app/api/src/component/schema`.
- Parameter component files must live in `app/api/src/component/parameter`.
- API path and method contracts must live under `app/api/src/path`.
- Path directories must mirror the endpoint path segments. For example, endpoint `api/v1/user/roles` must be defined at `app/api/src/path/api/v1/user/roles/index.yaml`, with the path item and method contract written in that `index.yaml`.
- Request and response schema file naming format: `[OperationId]Request.yaml` or `[OperationId]Response.yaml`.
- Every request and response field must define an example value.
- Every `operationId` must be unique.
- Every endpoint must be grouped with tags.
- The handler name must match the `operationId`.
- Handler implementation file name format from `operationId`: `internal/app/backend/server/impl_[OperationId].go`.
- Handler test file name format from `operationId`: `internal/app/backend/server/impl_[OperationId]_test.go`.
- Handler functions must implement the generated `ServerInterface` contract for the matching `operationId`.

After creating or changing API definitions, run:

```sh
make clean api.yaml validate-api api-docs
```

When the API contract is ready to implement, create the matching handler implementation and test files in `internal/app/backend/server`, then reinitialize the app:

```sh
make clean init
```

## Migrations

This application uses `golang-migrate`.

- Create migration files with the Make target.
- Migration names must reflect the migration purpose.
- Generated migration files live in `migrations`.
- Every up migration that adds DDL must add comments for every new table field in the same migration.
- For `CREATE TABLE`, every column in the new table must have a matching `COMMENT ON COLUMN schema.table.column IS '...'` statement, including primary key, timestamp, status, and foreign key columns.
- For `ALTER TABLE ... ADD COLUMN`, every added column must have a matching `COMMENT ON COLUMN` statement in the same migration.
- Column comments must describe the business or technical purpose of the field, not only repeat the column name.

Commands:

```sh
make create_migration <migration-name>
make migrate_up db=postgres://masraga@localhost:5433/kerp?sslmode=disable
make migrate_down db=postgres://masraga@localhost:5433/kerp?sslmode=disable
```
