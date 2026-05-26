# KERP API Technical Docs

This directory is the application technical navigator. Read this file, `dictionary.md`, and any relevant feature docs before modifying application code. Update the relevant docs after changing selected code.

## How to Use

- Start here to find the feature and its owning files.
- Use `dictionary.md` for shared terms, layers, and naming conventions.
- Use `features/*.md` for endpoint behavior, service flow, repository behavior, and tests.
- If code behavior changes, update the feature doc in the same change.

## Feature Navigator

| Feature | Public API | Handler | Service | Repository | Docs |
| --- | --- | --- | --- | --- | --- |
| Auth Validate PIN | `POST /api/v1/auth/validate/pin` | `internal/app/backend/server/impl_AuthValidatePin.go` | `internal/service/auth/service_AuthValidatePin.go` | `internal/service/auth/repo_CreateNewPin.go`, `internal/service/auth/repo_StoreAccessToken.go` | `features/auth-validate-pin.md` |

## Documentation Rules

- Keep docs factual and implementation-specific.
- Prefer links to concrete files and symbols over broad descriptions.
- Record user-visible behavior, error behavior, tests, and known implementation notes.
- Do not document planned behavior as current behavior.
