# Auth OTP Verification

OTP verification validates the current active OTP and enables exactly one successful PIN authentication. The one-time gate is stored as `public.auth.is_otp_valid`.

## API

| Method | Path | Operation |
| --- | --- | --- |
| `POST` | `/api/v1/auth/otp/verify` | `VerifyNewAuthUserOTP` |

The request and response shapes remain unchanged. A successful response means the OTP row is verified, the auth account is permanently marked verified, and the PIN authentication gate is enabled.

## Service Flow

1. Find the active auth account by phone number.
2. Set `is_otp_valid = false` before checking the submitted OTP. This ensures wrong, expired, already-used, or persistence-failed attempts leave the gate closed.
3. Find the OTP by user ID and code, then reject expired or already-verified rows.
4. Begin a transaction.
5. Mark the OTP row verified and set permanent account `is_verified = true`.
6. Set account `is_otp_valid = true`.
7. Commit and return OTP validity, phone number, note, and whether the user still needs an initial PIN.

The handler uses the consolidated `VerifyOtpOutput`; it does not make a second `VerifyUserAccount` service call.

## Lifecycle

- New accounts default to `is_otp_valid = false`.
- Requesting a fresh registration OTP resets existing accounts to false.
- Successful OTP verification sets the gate true.
- Wrong PIN attempts preserve true so the user can retry.
- Successful PIN authentication stores the JWT and resets the gate to false in one transaction.
- The issued JWT remains valid; a later PIN login requires another OTP verification.

## Tests

Coverage includes pre-clear failures, invalid and expired OTP attempts, transaction rollback, successful gate enabling, handler response mapping, and repository updates.

```sh
go test ./internal/service/auth -run 'TestAuthService_VerifyOtp|TestAuthRepository_VerifyOtp|TestAuthRepository_UpdateOtpValidity'
go test ./internal/app/backend/server -run TestServer_VerifyNewAuthUserOTP
```
