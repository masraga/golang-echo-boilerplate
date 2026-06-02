# Auth Token Validation

## Summary

`RegisterMiddlewares` wires `AuthValidationFilterMiddleware` into Echo so protected routes require a JWT access token before the OpenAPI request validator runs. Clients send the token in the `Authorization` header with the `Bearer <token>` format.

## Public API

Protected endpoints must include:

| Header | Required | Meaning |
| --- | --- | --- |
| `Authorization` | yes | JWT access token returned as `authToken` by `POST /api/v1/auth/validate/pin`, formatted as `Bearer <token>`. |

The current public skip map allows these requests without token validation:

| Method | Path |
| --- | --- |
| `GET` | `/api/v1/ping` |
| `POST` | `/api/v1/ping` |
| `POST` | `/api/v1/auth/register/phone` |
| `POST` | `/api/v1/auth/otp/verify` |
| `POST` | `/api/v1/auth/validate/pin` |

All other registered routes pass through JWT validation. In the current OpenAPI surface, `POST /api/v1/crypto/encrypt` is protected.

## Implementation Map

| Layer | File | Main Symbols |
| --- | --- | --- |
| Middleware registration | `app/backend/api.go` | `RegisterMiddlewares` |
| Auth middleware | `internal/echo/middleware/middleware.go` | `AuthValidationFilterMiddleware` |
| Auth middleware | `internal/echo/middleware/auth_validation_filter.go` | `authValidationFilter` |
| Skip map | `internal/echo/middleware/const.go` | `skipAuthValidationFilterMap` |
| Skip helper | `internal/echo/middleware/util.go` | `skipValidation`, `returnUnauthorized` |
| Service | `internal/service/auth/service_ValidateJwtToken.go` | `AuthService.ValidateJwtToken` |
| Repository | `internal/service/auth/repo_FindAccessToken.go` | `AuthRepository.FindAccessToken` |
| Types | `internal/service/auth/type.go` | `ValidateJwtTokenInput`, `ValidateJwtTokenOutput`, `FindAccessTokenInput`, `FindAccessTokenOutput` |
| Errors | `internal/service/auth/error.go` | `ErrAuthSigInvalid`, `ErrAuthTokenInvalid`, `ErrAuthTokenExpired`, `ErrFindAccessTokenNotFound` |

## Middleware Flow

`RegisterMiddlewares` applies Echo middleware in this order:

1. `RequestID`
2. `RequestLogger`
3. `Recover`
4. `AuthValidationFilterMiddleware`
5. `OapiRequestValidatorWithOptions`

It also configures `middleware.HTTPErrorHandler` as Echo's HTTP error handler before registering the OpenAPI request validator.

For each request, `authValidationFilter`:

1. Reads the route path from `echo.Context.Path()` and method from the request.
2. Skips validation when `skipValidation` matches the path and method in `skipAuthValidationFilterMap`.
3. Reads `Authorization`.
4. Returns `401` with `{"error":"Unauthorized"}` when the header is empty.
5. Extracts the JWT string from `Bearer <token>`.
6. Calls `AuthService.ValidateJwtToken`.
7. Returns `401` with `{"error":"Unauthorized"}` when JWT validation returns an error.
8. Calls the next handler on success.

Current implementation note: non-empty `Authorization` values that do not contain `Bearer ` are split without a length check, so malformed non-Bearer headers can panic and be recovered by Echo instead of returning the middleware's unauthorized response.

## JWT Validation

`AuthService.ValidateJwtToken` parses the token into `ValidateJwtTokenOutput` claims with the configured JWT secret. It requires an HMAC signing method, rejects invalid tokens, and rejects tokens whose `ExpiredAtUtc0` is older than the current Unix millisecond time.

After JWT claim validation, the service calls `AuthRepositoryReader.FindAccessToken` with the raw token string and `claims.UserId`. The repository requires an active row in `public.auth_access_token` where:

| Column | Required value |
| --- | --- |
| `id` | Raw JWT string from the request header. |
| `user_id` | `claims.UserId` parsed from the JWT. |
| `is_active` | `true` |

When no matching active row exists, validation returns `ErrFindAccessTokenNotFound`. The service does not compare `public.auth_access_token.expired_at_utc0` against the current time; JWT claim expiration remains the expiration check.

## Test Coverage

| Test File | Coverage |
| --- | --- |
| `internal/service/auth/service_ValidateJwtToken_test.go` | Valid token with active stored row, missing stored row, expired token, invalid signature. |
| `internal/service/auth/service_ValidateJwtToken_integration_test.go` | JWT creation and validation with real auth service/repository wiring and sqlmock-backed active token lookup. |
| `internal/service/auth/repo_FindAccessToken_test.go` | Active stored access-token lookup failure and success. |

Recommended narrow check:

```sh
go test ./internal/service/auth -run 'TestAuthService_ValidateJwtToken|TestAuthRepository_FindAccessToken'
```

## Change Checklist

- Update `skipAuthValidationFilterMap` and this document together when public/protected route behavior changes.
- Update OpenAPI authorization header documentation for any route that is protected by the middleware.
- Add or update middleware tests when token parsing, unauthorized response behavior, or skip-route matching changes.
