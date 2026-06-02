# Auth Phone Registration

Phone registration creates an auth account when the phone number is new and issues an OTP for the account. When the phone number already belongs to an active auth account, the endpoint does not create another `public.auth` row; it replaces the user's active OTP rows and returns the fresh OTP payload.

## API

| Method | Path | Operation |
| --- | --- | --- |
| `POST` | `/api/v1/auth/register/phone` | `RegisterPhoneNumber` |

The request body uses `RegisterPhoneNumberRequest` with encrypted `phoneNo`. The response uses `RegisterPhoneNumberResponse` with the auth account `id` and issued `otpCode`.

## Important Files

| Layer | File | Purpose |
| --- | --- | --- |
| Handler | `internal/app/backend/server/impl_RegisterPhoneNumber.go` | Decrypts `phoneNo`, calls `AuthService.CreateNewAccount`, and returns `201`. |
| Service | `internal/service/auth/service_CreateNewAccount.go` | Resolves existing account by phone number, creates a missing account, and delegates OTP creation to `CreateOTP`. |
| Service | `internal/service/auth/service_CreateOTP.go` | Finds the account, deletes active OTP rows, generates OTP code and expiration when omitted, and stores the OTP. |
| Repository | `internal/service/auth/repo_CreateNewAccount.go` | Inserts new `public.auth` rows. |
| Repository | `internal/service/auth/repo_CreateOTP.go` | Inserts `public.auth_otp` rows. |
| Repository | `internal/service/auth/repo_DeleteAllUserOTP.go` | Deactivates the user's active OTP rows before inserting a replacement. |

## Behavior

`AuthService.CreateNewAccount` first calls `FindAuth` by `phoneNo`.

If the account exists, the service uses the existing user id and skips `CreateNewAccount`. It still opens a transaction so OTP deletion and creation are handled together.

If the account is missing, the service generates a UUID, inserts the account, then creates the OTP inside the same transaction.

Both paths call the same private registration OTP helper, which delegates code generation, expiration defaults, prior OTP deletion, and OTP insertion to `AuthService.CreateOTP`. This keeps OTP behavior shared between new and existing users.

## Tests

| File | Coverage |
| --- | --- |
| `internal/service/auth/service_CreateNewAccount_test.go` | Existing-user OTP creation without account insert, begin failure, OTP rollback, and new-account success. |
| `internal/service/auth/service_CreateOTP_test.go` | OTP user lookup, prior OTP deletion, OTP insert failure, and success. |
| `internal/service/auth/repo_CreateNewAccount_test.go` | Account insert failure and returned account id. |
| `internal/app/backend/server/impl_RegisterPhoneNumber_test.go` | Handler success for new and existing registration calls. |

## Useful Commands

```sh
go test ./internal/service/auth -run 'TestAuthService_CreateNewAccount|TestAuthService_CreateOTP|TestAuthRepository_CreateNewAccount'
go test ./internal/app/backend/server -run TestServer_RegisterPhoneNumber
```
