# Auth Validate PIN

## Summary

`AuthValidatePin` validates an auth user's PIN, stores the current JWT access token, and returns that token. If the user does not have a PIN yet, the same flow creates the first PIN after validating `retypePin`.

## Public API

| Item | Value |
| --- | --- |
| Method | `POST` |
| Path | `/api/v1/auth/validate/pin` |
| Operation ID | `AuthValidatePin` |
| API path file | `app/api/src/path/api/v1/auth/validate/pin/index.yaml` |
| Request schema | `app/api/src/component/schema/AuthValidatePinRequest.yaml` |
| Response schema | `app/api/src/component/schema/AuthValidatePinResponse.yaml` |

Request body:

| Field | Required | Meaning |
| --- | --- | --- |
| `phoneNo` | yes | Encrypted phone number received by the handler. |
| `pin` | yes | PIN to create or validate. |
| `retypePin` | no | Required by the service only when the user does not already have a PIN. |

Response body:

| Field | Meaning |
| --- | --- |
| `userId` | User UUID parsed from the auth service output. |
| `authToken` | JWT token created after a valid PIN flow. |

## Implementation Map

| Layer | File | Main Symbols |
| --- | --- | --- |
| Handler | `internal/app/backend/server/impl_AuthValidatePin.go` | `Server.AuthValidatePin` |
| Service | `internal/service/auth/service_AuthValidatePin.go` | `AuthService.AuthValidatePin`, `AuthService.createNewPin` |
| Repository | `internal/service/auth/repo_CreateNewPin.go` | `AuthRepository.CreateNewPin` |
| Repository | `internal/service/auth/repo_StoreAccessToken.go` | `AuthRepository.StoreAccessToken` |
| Types | `internal/service/auth/type.go` | `AuthValidatePinInput`, `AuthValidatePinOutput`, `CreateNewPinInput`, `CreateNewPinOutput`, `StoreAccessTokenInput`, `StoreAccessTokenOutput` |
| Constants | `internal/service/auth/const.go` | `MinPinLen`, `MaxPinLen`, `JwtTokenExpiredDuration` |
| Errors | `internal/service/auth/error.go` | `ErrValidateRetypePin`, `ErrPinCodeNotMatch`, `ErrPinIsTooLongOrShort`, `ErrCreateNewPin`, `ErrStoreAccessToken` |
| Migration | `migrations/202605300013_add_table_access_token.up.sql` | Creates `public.access_token` |

## Handler Flow

1. Bind `api.AuthValidatePinRequest`.
2. Decrypt `phoneNo` with `CryptoService.Decrypt`.
3. Call `AuthService.AuthValidatePin` with decrypted phone number, `pin`, and optional `retypePin`.
4. Parse `AuthValidatePinOutput.UserId` to UUID.
5. Return `200 OK` with `api.AuthValidatePinResponse`.

Current handler note: decrypt and service errors are returned directly by `Server.AuthValidatePin`; invalid UUID errors are passed through `returnError`, which returns the default unknown internal error response unless mapped.

## Service Flow

1. Load auth account with `AuthRepositoryReader.FindAuth` by phone number.
2. Validate `pin` and optional `retypePin` before starting a transaction when possible.
3. Begin a database transaction through `AuthRepositoryWriter.Begin`.
4. If the account has no stored PIN:
   - require `RetypePinCode`;
   - require `PinCode` to match `RetypePinCode`;
   - require PIN length to be exactly 6 based on `MinPinLen` and `MaxPinLen`;
   - write the new PIN through `AuthRepositoryWriter.CreateNewPin`.
5. If the account has a stored PIN, compare it with the supplied `PinCode`.
6. Create a JWT with `CreateToken` using one computed `ExpiredAtUtc0` value.
7. Store the current JWT through `AuthRepositoryWriter.StoreAccessToken`, which deactivates previous active tokens for the user first.
8. Commit on success or roll back on any error.
9. Return `IsValid`, `UserId`, `Token`, and internal `ExpiredAtUtc0`.

## Repository Behavior

`AuthRepository.CreateNewPin` updates `public.auth`:

| Column | Value |
| --- | --- |
| `pin` | `CreateNewPinInput.PinCode` |
| `updated_at_utc0` | Current Unix milliseconds |

The update targets `id = CreateNewPinInput.UserId`. On SQL execution failure, it wraps the error with `ErrCreateNewPin`.

`AuthRepository.StoreAccessToken` updates `public.access_token`:

| Step | Behavior |
| --- | --- |
| Deactivate old tokens | Set `is_active = false` for active rows with the same `user_id`. |
| Insert current token | Insert `id`, `user_id`, `expired_at_utc0`, and `is_active = true`. |

On SQL execution failure, it wraps the error with `ErrStoreAccessToken`.

The `public.access_token` table is created by `migrations/202605300013_add_table_access_token.up.sql`:

| Column | Meaning |
| --- | --- |
| `id` | JWT access token string. |
| `user_id` | Auth user id. |
| `expired_at_utc0` | Token expiration in Unix milliseconds. |
| `is_active` | Logical delete/current-active indicator. |

## Test Coverage

| Test File | Coverage |
| --- | --- |
| `internal/app/backend/server/impl_AuthValidatePin_test.go` | Successful HTTP response, decrypt failure, service failure, invalid service user id. |
| `internal/service/auth/service_AuthValidatePin_test.go` | Auth not found, missing retype PIN, mismatched new PIN, invalid PIN length, begin transaction failure, create PIN failure, store access token failure, existing PIN mismatch, create new PIN success, existing PIN success. |
| `internal/service/auth/repo_CreateNewPin_test.go` | SQL update failure and success result. |
| `internal/service/auth/repo_StoreAccessToken_test.go` | Old token deactivation failure, current token insert failure, and success result. |

Recommended narrow checks:

```sh
go test ./internal/app/backend/server -run TestServer_AuthValidatePin
go test ./internal/service/auth -run 'TestAuthService_AuthValidatePin|TestAuthRepository_CreateNewPin|TestAuthRepository_StoreAccessToken'
```

## Change Checklist

- Update OpenAPI path/schema files when request or response shape changes.
- Regenerate generated API code when OpenAPI changes.
- Keep handler, service, repository, and tests aligned.
- Keep the access-token migration, `StoreAccessToken` repository, and service transaction behavior aligned.
- Update this document after changing AuthValidatePin behavior, errors, request/response shape, or owning files.
