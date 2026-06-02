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

## Auth Domain

| Term | Meaning |
| --- | --- |
| `AuthService` | Auth domain service type defined in `internal/service/auth/service.go`. |
| `AuthRepository` | Auth domain repository type defined in `internal/service/auth/repository.go`. |
| `FindAuth` | Repository lookup for an auth account by phone number or user id. |
| `FindAccessToken` | Repository lookup for an active JWT row in `public.auth_access_token` by token id and user id. |
| `CreateNewPin` | Repository operation that writes a PIN to `public.auth.pin`. |
| `StoreAccessToken` | Repository operation that deactivates previous active tokens for a user and stores the current JWT in `public.auth_access_token`. |
| `auth_access_token` | Table that stores JWT access tokens by token string id, auth user id, expiration, and active flag. |
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

## Testing Terms

| Term | Meaning |
| --- | --- |
| Table-driven test | Required test shape for new or changed tests in this project. |
| `gomock` | Mocking tool used for service and handler dependency tests. |
| `go-sqlmock` | SQL expectation tool used for repository tests. |
| `make clean init` | Required regeneration command after adding or changing any interface so gomock and other generated artifacts are recreated before tests. |
