# KERP API Rules

These rules apply every time Codex works in this project.

## Architecture

Use the handler > service > repository mechanism.

- Handlers live in `internal/app/backend/server`.
- Services live in `internal/service/*`.
- Repositories live in `internal/service/*`.
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
- Request and response schema files live in `app/api/src/component/schema`.
- Schema file naming format: `[OperationId]Request.yaml` or `[OperationId]Response.yaml`.
- Every `operationId` must be unique.
- Every endpoint must be grouped with tags.
- The handler name must match the `operationId`.
- Handler file name format from `operationId`: `impl_[OperationId].go`.

After creating or changing API definitions, run:

```sh
make clean api.yaml validate-api api-docs
```

After the API handler is generated, reinitialize the app:

```sh
make clean init
```

## Migrations

This application uses `golang-migrate`.

- Create migration files with the Make target.
- Migration names must reflect the migration purpose.
- Generated migration files live in `migrations`.

Commands:

```sh
make create_migration <migration-name>
make migrate_up db=postgres://masraga@localhost:5433/kerp?sslmode=disable
make migrate_down db=postgres://masraga@localhost:5433/kerp?sslmode=disable
```
