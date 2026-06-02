# KERP API Dictionary

This dictionary defines shared application terms used across feature docs.

## Layers

| Term | Meaning |
| --- | --- |
| API definition | OpenAPI source under `app/api/src`; generated Go API types live in `generated/api`. |
| Handler | Echo HTTP adapter under `internal/app/backend/server`; binds request bodies, calls services, and formats responses. |
| Service | Domain business logic under `internal/service/[domain]`; validates rules and coordinates repositories and helper services. |
| Repository | Database access under `internal/service/[domain]`; builds SQL with `sqlf` and uses `dbtx` for transaction-aware execution. |
| Test helper | Shared assertion helpers under `internal/testutil`, especially `RequireResult` and `RequireHttpResultJson`. |
| Parser utility | Shared parsing and normalization helpers under `internal/util/parser`. |
| Time utility | Shared timestamp helpers under `internal/util/time`. |

## Auth Domain

| Term | Meaning |
| --- | --- |
| `AuthService` | Auth domain service type defined in `internal/service/auth/service.go`. |
| `AuthRepository` | Auth domain repository type defined in `internal/service/auth/repository.go`. |
| `FindAuth` | Repository lookup for an auth account by phone number or user id. |
| `CreateNewAccount` | Auth registration service method that creates a missing phone account, or reuses an existing account, then issues a fresh OTP. |
| `CreateOTP` | Auth service method that verifies the user, deletes active user OTP rows, generates default OTP data when omitted, and stores the new OTP. |
| `FindAccessToken` | Repository lookup for an active JWT row in `public.auth_access_token` by token id and user id. |
| `CreateNewPin` | Repository operation that writes a PIN to `public.auth.pin`. |
| `StoreAccessToken` | Repository operation that deactivates previous active tokens for a user and stores the current JWT in `public.auth_access_token`. |
| `auth_access_token` | Table that stores JWT access tokens by token string id, auth user id, expiration, and active flag. |
| `auth_otp` | Table that stores OTP codes by auth user id, expiration, verification state, and active flag. |
| `auth_api_contract` | Table that stores active API contracts keyed by OpenAPI `operationId`, endpoint path, and endpoint method. |
| `auth_user_api_contract` | Table that stores active per-user grants to API contracts. Role assignment replaces rows in this table with hard deletes and inserts. |
| `auth_roles` | Table that stores active API access roles with owner, creator, and copied role metadata for assignment to auth users. |
| `auth_roles_contract_api` | Table that stores active role-to-API-contract mappings used to grant user API access during role assignment. |
| API access grant | A direct user-to-API-contract permission checked after JWT validation for protected routes. |
| Endpoint normalization | Shared parser behavior that trims endpoint paths, adds a leading `/` when missing, and lowercases HTTP methods before API access storage or validation. |
| Auth role | A named group of API contracts. Assigning one role to a user updates `auth.role_id` and `auth.role_name`, then syncs `auth_user_api_contract`. |
| Bootstrap access user | Optional `AUTH_ACCESS_BOOTSTRAP_USER_ID` config value used at startup to grant one auth user all active API contracts and at runtime to bypass per-endpoint API access checks after JWT validation. |
| `PinCode` | The PIN value supplied by the user. Current PIN length is controlled by `MinPinLen` and `MaxPinLen` in `internal/service/auth/const.go`. |
| `RetypePinCode` | Optional confirmation PIN. Required when the account does not already have a PIN. |
| JWT token | Auth token returned by `POST /api/v1/auth/validate/pin` and validated by `AuthService.ValidateJwtToken` when protected routes receive `Authorization: Bearer <token>`. |

## Error Terms

| Error | Meaning |
| --- | --- |
| `ErrAuthNotFound` | Auth account lookup did not find an account. |
| `ErrValidateRetypePin` | A new PIN was requested without `retypePin`. |
| `ErrPinCodeNotMatch` | `pin` and `retypePin` differ, or an existing stored PIN differs from the supplied PIN. |
| `ErrPinIsTooLongOrShort` | Supplied PIN length is outside the configured bounds. |
| `ErrCreateNewPin` | Repository failed while writing a new PIN. |
| `ErrStoreAccessToken` | Repository failed while deactivating old tokens or storing the current access token. |
| `ErrFindAccessTokenNotFound` | Active JWT row was not found in `public.auth_access_token` for the token id and user id. |
| `ErrAuthSigInvalid` | JWT parsing failed or the token used an unsupported signing method. |
| `ErrAuthTokenInvalid` | JWT parsed but was not valid. |
| `ErrAuthTokenExpired` | JWT `ExpiredAtUtc0` is older than the current Unix millisecond time. |
| `ErrUserApiContractForbidden` | Access validation did not find an active grant for the user, endpoint path, and method. |
| `ErrFindAuthRoleNotFound` | Role lookup did not find an active `auth_roles` row. |
| `ErrFindAuthRoleContractApiNotFound` | Role API contract mapping lookup or delete did not find an active `auth_roles_contract_api` row. |

## Testing Terms

| Term | Meaning |
| --- | --- |
| Table-driven test | Required test shape for new or changed tests in this project. |
| `gomock` | Mocking tool used for service and handler dependency tests. |
| `go-sqlmock` | SQL expectation tool used for repository tests. |
| `make clean init` | Required regeneration command after adding or changing any interface so gomock and other generated artifacts are recreated before tests. |
