# Auth Phone Registration

Phone registration creates an auth account when the phone number is new and issues an OTP for the account. When the phone number already belongs to an active auth account, the endpoint does not create another `public.auth` row; it replaces the user's active OTP rows and returns the fresh OTP payload.

## API

| Method | Path | Operation |
| --- | --- | --- |
| `POST` | `/api/v1/auth/register/phone` | `RegisterPhoneNumber` |

The request body uses `RegisterPhoneNumberRequest` with required encrypted `phoneNo` and `firebaseId`. The response uses `RegisterPhoneNumberResponse` with the auth account `id` and issued `otpCode`.

`firebaseId` is the Firebase Cloud Messaging registration token for the user's current device. The OpenAPI request validator rejects missing or empty values before the handler.

## Important Files

| Layer | File | Purpose |
| --- | --- | --- |
| Handler | `internal/app/backend/server/impl_RegisterPhoneNumber.go` | Decrypts `phoneNo`, passes the required `firebaseId`, calls `AuthService.CreateNewAccount`, and returns `201`. |
| Service | `internal/service/auth/service_CreateNewAccount.go` | Resolves existing accounts, creates missing accounts or refreshes the validated Firebase ID, and delegates OTP creation to `CreateOTP`. |
| Service | `internal/service/auth/service_CreateOTP.go` | Finds the account, deletes active OTP rows, generates OTP code and expiration when omitted, and stores the OTP. |
| Repository | `internal/service/auth/repo_CreateNewAccount.go` | Inserts new `public.auth` rows with the required registration Firebase ID. |
| Repository | `internal/service/auth/repo_UpdateFirebaseId.go` | Replaces an existing account's Firebase ID and updates its modification timestamp. |
| Repository | `internal/service/auth/repo_CreateOTP.go` | Inserts `public.auth_otp` rows. |
| Repository | `internal/service/auth/repo_DeleteAllUserOTP.go` | Deactivates the user's active OTP rows before inserting a replacement. |

## Behavior

The global OpenAPI middleware validates the required, non-empty Firebase ID before invoking the handler. `AuthService.CreateNewAccount` then calls `FindAuth` by `phoneNo`.

If the account exists, the service uses the existing user id, skips `CreateNewAccount`, and replaces the stored Firebase ID. The update and OTP replacement use the same transaction.

If the account is missing, the service generates a UUID, inserts the account with the Firebase ID, then creates the OTP inside the same transaction.

Both paths call the same private registration OTP helper, which delegates code generation, expiration defaults, prior OTP deletion, and OTP insertion to `AuthService.CreateOTP`. This keeps OTP behavior shared between new and existing users.

## Tests

| File | Coverage |
| --- | --- |
| `internal/service/auth/service_CreateNewAccount_test.go` | Existing-user token update, begin failure, rollback behavior, and new-account success. |
| `internal/service/auth/service_CreateOTP_test.go` | OTP user lookup, prior OTP deletion, OTP insert failure, and success. |
| `internal/service/auth/repo_CreateNewAccount_test.go` | Account insert failure and returned account id. |
| `internal/service/auth/repo_UpdateFirebaseId_test.go` | Existing-account Firebase ID update success and failure. |
| `internal/app/backend/server/impl_RegisterPhoneNumber_test.go` | Handler success with the required Firebase ID for new and existing accounts. |

## Useful Commands

```sh
go test ./internal/service/auth -run 'TestAuthService_CreateNewAccount|TestAuthService_CreateOTP|TestAuthRepository_CreateNewAccount|TestAuthRepository_UpdateFirebaseId'
go test ./internal/app/backend/server -run TestServer_RegisterPhoneNumber
```
