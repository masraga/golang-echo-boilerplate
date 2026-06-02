# Auth Roles

## Summary

Auth roles group API contracts and apply that group to one auth user. The feature lives entirely in the auth domain: handlers call `AuthService`, and persistence uses `AuthRepository`.

Assigning a role updates the user's `auth.role_id`, `auth.role_name`, and `auth.created_by`, hard-deletes existing `auth_user_api_contract` rows for that user, then inserts active API contracts from `auth_roles_contract_api`.

## Public API

Management endpoints are protected by JWT and API access-control middleware:

| Method | Path | Operation ID |
| --- | --- | --- |
| `GET` | `/api/v1/auth/roles` | `ListAuthRoles` |
| `POST` | `/api/v1/auth/roles` | `CreateAuthRole` |
| `GET` | `/api/v1/auth/roles/{id}` | `GetAuthRole` |
| `PUT` | `/api/v1/auth/roles/{id}` | `UpdateAuthRole` |
| `DELETE` | `/api/v1/auth/roles/{id}` | `DeleteAuthRole` |
| `GET` | `/api/v1/auth/roles/{roleId}/api-contracts` | `ListAuthRoleContractApis` |
| `POST` | `/api/v1/auth/roles/{roleId}/api-contracts` | `CreateAuthRoleContractApi` |
| `DELETE` | `/api/v1/auth/roles/{roleId}/api-contracts/{id}` | `DeleteAuthRoleContractApi` |
| `PUT` | `/api/v1/auth/users/{userId}/role` | `AssignAuthUserRole` |
| `DELETE` | `/api/v1/auth/users/{userId}/role` | `DeleteAuthUserRole` |

## Implementation Map

| Layer | File | Main Symbols |
| --- | --- | --- |
| Handler | `internal/app/backend/server/impl_AuthRole.go` | Role CRUD handlers |
| Handler | `internal/app/backend/server/impl_AuthRoleContractApi.go` | Role API contract mapping handlers |
| Handler | `internal/app/backend/server/impl_AuthUserRole.go` | User role assignment and removal handlers |
| Service | `internal/service/auth/service_AuthRole.go` | Role CRUD service methods |
| Service | `internal/service/auth/service_AuthRoleContractApi.go` | Role API contract service methods |
| Service | `internal/service/auth/service_AuthUserRole.go` | Role assignment transaction and grant synchronization |
| Repository | `internal/service/auth/repo_AuthRole.go` | `auth_roles` persistence |
| Repository | `internal/service/auth/repo_AuthRoleContractApi.go` | `auth_roles_contract_api` persistence |
| Repository | `internal/service/auth/repo_AuthUserRole.go` | Auth user role field updates |
| Repository | `internal/service/auth/repo_AuthUserApiContract.go` | Hard-delete and role grant insert operations |
| Migration | `migrations/202606031023_add_auth_roles_access_control.up.sql` | Creates role tables, auth columns, and role endpoint contracts |

## Behavior

`auth_roles` rows are soft-deleted through `is_active=false`. Active role names are unique.

`auth_roles_contract_api` rows are soft-deleted through `is_active=false`. A role can have only one active mapping per API contract.

`AssignAuthUserRole` first verifies the role exists, then runs a transaction. The transaction updates the auth user role fields, hard-deletes the user's current API grants, and inserts active API contracts from the assigned role. The response includes the number of grants inserted.

`DeleteAuthUserRole` runs a transaction that clears the auth user's role fields and hard-deletes the user's API grants.

## Test Coverage

| Test File | Coverage |
| --- | --- |
| `internal/service/auth/service_AuthUserRole_test.go` | Role lookup, assignment grant replacement, and role removal grant clearing. |
| `internal/service/auth/repo_AuthUserApiContract_test.go` | Hard delete of user grants and inserting grants from role contracts. |
| `internal/app/backend/server/impl_AuthRole_test.go` | Create role handler mapping. |
| `internal/app/backend/server/impl_AuthRoleContractApi_test.go` | Create role API contract handler mapping. |
| `internal/app/backend/server/impl_AuthUserRole_test.go` | Assign role handler mapping. |

Recommended checks:

```sh
go test ./internal/service/auth -run 'TestAuthService_(AssignAuthUserRole|DeleteAuthUserRole)|TestAuthRepository_(DeleteAuthUserApiContract|InsertAuthUserApiContractsFromRole)'
go test ./internal/app/backend/server -run 'TestServer_(CreateAuthRole|CreateAuthRoleContractApi|AssignAuthUserRole)'
```
