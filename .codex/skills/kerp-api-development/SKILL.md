---
name: kerp-api-development
description: Follow KERP API project conventions for Go backend changes. Use when creating or editing handlers in internal/app/backend/server, services or repositories in internal/service/*, table-driven tests, OpenAPI definitions under app/api, generated API handlers, or migrations.
---

# KERP API Development

Before changing this repository, read `.codex/rules/kerp-api.md` and apply it as the project-specific source of truth.

## Core Workflow

1. Keep backend code in the existing handler > service > repository flow.
2. Use the repository's file naming rules for implementation and test files.
3. Write every new or changed test as a table-driven test.
4. Match test function names to the concrete type and method under test.
5. For API definition changes, update `app/api/src/main.yaml` and referenced schema/path files, then run the required Make targets.
6. For migrations, create files through the Make target, use descriptive migration names, and add purpose-driven comments for every field introduced by DDL.
7. Prefer existing examples in the same package as the implementation source for imports, setup, mocks, assertions, and error handling.

## Required Checks

Run the narrowest relevant Go tests after backend changes. For API definition changes, run:

```sh
make clean api.yaml validate-api api-docs
make clean init
```

For migrations, create and apply or roll back through:

```sh
make create_migration <migration-name>
make migrate_up db=$kerpdb
make migrate_down db=$kerpdb
```
