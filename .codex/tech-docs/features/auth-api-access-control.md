# Auth API Access Control

## Summary

Auth API access control validates per-user API grants after JWT token validation. Access is stored in `public.auth_user_api_contract` by auth user id and API contract id; grants can be created directly or synchronized from auth roles.

## Public API

Management endpoints are protected by the same JWT and access-control middleware:

| Method | Path | Operation ID |
| --- | --- | --- |
| `GET` | `/api/v1/auth/api-contracts` | `ListAuthApiContracts` |
| `POST` | `/api/v1/auth/api-contracts` | `CreateAuthApiContract` |
| `GET` | `/api/v1/auth/api-contracts/{id}` | `GetAuthApiContract` |
| `PUT` | `/api/v1/auth/api-contracts/{id}` | `UpdateAuthApiContract` |
| `DELETE` | `/api/v1/auth/api-contracts/{id}` | `DeleteAuthApiContract` |
| `GET` | `/api/v1/auth/user-api-contracts` | `ListAuthUserApiContracts` |
| `POST` | `/api/v1/auth/user-api-contracts` | `CreateAuthUserApiContract` |
| `GET` | `/api/v1/auth/user-api-contracts/{id}` | `GetAuthUserApiContract` |
| `PUT` | `/api/v1/auth/user-api-contracts/{id}` | `UpdateAuthUserApiContract` |
| `DELETE` | `/api/v1/auth/user-api-contracts/{id}` | `DeleteAuthUserApiContract` |

## Implementation Map

| Layer | File | Main Symbols |
| --- | --- | --- |
| Middleware | `internal/echo/middleware/auth_validation_filter.go` | `authValidationFilter` |
| Handler | `internal/app/backend/server/impl_AuthApiContract.go` | API contract CRUD handlers |
| Handler | `internal/app/backend/server/impl_AuthUserApiContract.go` | User API contract CRUD handlers |
| Service | `internal/service/auth` | `AuthService` |
| Repository | `internal/service/auth` | `AuthRepository` |
| Utility | `internal/util/parser/http.go` | Endpoint path and method normalization |
| Utility | `internal/util/time/time.go` | Unix millisecond timestamp creation |
| Migration | `migrations/202606030119_add_auth_api_contract_access_control.up.sql` | Creates and seeds access-control tables |
| Migration | `migrations/202606031023_add_auth_roles_access_control.up.sql` | Adds roles, role contracts, auth role columns, and role endpoint contracts |

## Behavior

`auth_api_contract.endpoint_path` stores Echo route paths, including `:id` for path parameters. `endpoint_method` is stored lowercase. Service inputs normalize endpoint paths and methods before writing or validating.

`ValidateUserApiContract` allows the configured `AUTH_ACCESS_BOOTSTRAP_USER_ID` user without checking `auth_user_api_contract`. For other users, it requires an active `auth_user_api_contract` row joined to an active `auth_api_contract` row for the validated user id, route path, and method. Missing access returns `403 {"error":"Forbidden"}` from middleware.

`DeleteAuthUserApiContract` physically deletes the selected user grant. The repository also hard-deletes all grants for a user when role assignment or role removal replaces the user's access set.

`AUTH_ACCESS_BOOTSTRAP_USER_ID` optionally grants one auth user all active API contracts during startup. The bootstrap insert is idempotent and skips active grants that already exist. The same configured user id also bypasses endpoint access validation at runtime after JWT validation succeeds.

## Test Coverage

| Test File | Coverage |
| --- | --- |
| `internal/echo/middleware/auth_validation_filter_test.go` | JWT-first flow and `403` access denial. |
| `internal/service/auth/service_ValidateUserApiContract_test.go` | Bootstrap-user bypass, service path/method normalization, and forbidden propagation. |
| `internal/service/auth/service_ValidateUserApiContract_integration_test.go` | Concrete service/repository validation for bootstrap bypass and normal grant lookup. |
| `internal/service/auth/repo_ValidateUserApiContract_test.go` | Active grant lookup success and no-grant failure. |
| `internal/service/auth/service_BootstrapUserApiContracts_test.go` | Bootstrap service success and failure. |
| `internal/service/auth/repo_BootstrapUserApiContracts_test.go` | Bootstrap insert success and SQL failure. |
| `internal/service/auth/repo_AuthUserApiContract_test.go` | Hard delete for user grants and role-to-user grant insertion. |
| `internal/app/backend/server/impl_AuthApiContract_test.go` | API contract create handler mapping. |
| `internal/app/backend/server/impl_AuthUserApiContract_test.go` | User API contract create handler mapping. |

Recommended checks:

```sh
go test ./internal/service/auth
go test ./internal/echo/middleware
go test ./internal/app/backend/server -run 'TestServer_CreateAuth(Api|UserApi)Contract'
```
