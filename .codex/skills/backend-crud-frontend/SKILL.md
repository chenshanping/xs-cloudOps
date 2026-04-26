---
name: backend-crud-frontend
description: Standardize CRUD work for the current go-base workspace. Use when creating or refactoring admin modules, database tables or seed SQL, Go backend APIs/services/models/requests/routes, or Vue frontend APIs/pages/types in the current Gin + Gorm + Vue3 + Ant Design Vue project.
---

# Backend CRUD Frontend

Use this skill to keep CRUD work aligned with the current `go-base` workspace instead of copying conventions from `XTMS`, Spring Boot projects, or older repos.

## Core Rule

Discover the live module shape first, then extend it.

- Do not assume Spring Boot paths such as `admin.system`, `admin.business`, or `vue/apps/web-antd`.
- Do not assume each backend module has its own directory; many current modules are flat files.
- Do not assume each request struct lives in a dedicated file; many request structs live in `server/model/request/request.go`.
- Prefer the current neighboring module over memory.

## Read First

Inspect these anchors before editing:

- `server/api/v1/user.go`
- `server/api/v1/role.go`
- `server/api/v1/menu.go`
- `server/service/user.go`
- `server/service/role.go`
- `server/service/menu.go`
- `server/model/sys_user.go`
- `server/model/sys_role.go`
- `server/model/sys_menu.go`
- `server/model/request/request.go`
- `server/router/modules/user.go`
- `server/router/modules/role.go`
- `server/router/modules/menu.go`
- `web/src/api/user.ts`
- `web/src/api/role.ts`
- `web/src/api/menu.ts`
- `web/src/views/system/user/index.vue`
- `web/src/views/system/role/index.vue`
- `web/src/views/system/menu/index.vue`
- `web/src/components/ProTable.vue`
- `web/src/components/AvatarUpload.vue`
- `web/src/components/ImageUpload.vue`
- `web/src/components/FileUpload.vue`
- `web/src/components/FilePreview.vue`
- `web/src/utils/permission.ts`
- `web/src/directives/permission.ts`
- `web/src/store/user.ts`
- `web/src/types/index.ts`

Load these references when relevant:

- `references/list-page-pattern.md`

## Dynamic Discovery Rules

When a target file does not exist, search the nearest real module and follow that pattern.

- For standard admin CRUD, inspect `user`, `role`, `menu`, `dict`, and `storage` first.
- Keep backend modules flat unless the current repo already uses a subdirectory for that module.
- Keep frontend APIs flat under `web/src/api/*.ts` unless the neighboring module already uses another pattern.
- For frontend pages, default to `web/src/views/system/<module>/`.
- For frontend types, choose between `web/src/types/index.ts` and `web/src/types/<module>.ts` by following the nearest live example.

Do not cite or depend on external docs that are not present in the current workspace.

## Current Project Architecture

### Backend

This workspace is a single Gin application under `server`.

- API handlers live in `server/api/v1`
- Business logic lives in `server/service`
- Gorm models live in `server/model`
- Request DTOs live in `server/model/request`
- Response helpers live in `server/model/response`
- Route registration lives in `server/router/modules`
- Shared state and infrastructure live in `server/global`

Current backend patterns:

- Parse path ids with `strconv.ParseUint`
- Bind query params with `ShouldBindQuery`
- Bind JSON bodies with `ShouldBindJSON`
- Return through `response.BadRequest`, `response.Fail`, `response.OkWithData`, `response.OkWithMessage`, or `response.OkWithPage`
- Use `global.DB` and Gorm transactions in service code
- Keep route metadata registration aligned with `R(...)`, `registry.WithAuth()`, and `registry.WithRequest(...)`

### Frontend

The frontend is a single Vite Vue 3 application under `web`.

- Admin APIs live in `web/src/api/*.ts`
- Admin pages live in `web/src/views/system/*`
- Shared components live in `web/src/components/*`
- Permission helpers live in `web/src/utils/permission.ts`
- Permission directive lives in `web/src/directives/permission.ts`
- Session state lives in `web/src/store/user.ts`
- Shared types live in `web/src/types/index.ts`

Current frontend patterns:

- Use `ProTable` for standard search + toolbar + table pages
- Use `useTableColumns(...)` to hide the action column when the user has no row-action permission
- Use `v-permission` on buttons and row actions
- Keep CRUD pages compact; do not add decorative hero headers by default
- For create, edit, and similar action surfaces, default to a drawer-based flow unless the user explicitly asks for a modal or another pattern
- Reuse `AvatarUpload`, `ImageUpload`, `FileUpload`, and `FilePreview` when the module needs them

## Backend Conventions

### File Placement

Default backend targets for a CRUD-style module:

- `server/api/v1/<module>.go`
- `server/service/<module>.go`
- `server/model/<module>.go`
- `server/router/modules/<module>.go`
- `server/model/request/request.go` or a neighboring `*_request.go` file

Do not introduce a new backend package tree unless the repo already uses that pattern for the same area.

### Layering

- API layer handles path/query/body parsing and top-level response shaping
- Service layer handles validation, transactions, associations, soft-delete safety, cache clearing, and token invalidation
- Model layer holds Gorm schema and helper methods such as file URL filling
- Route modules register endpoints and request metadata

### File and Relation Handling

- Prefer storing file ids such as `avatar_file_id`, not only raw URLs
- Reuse existing `Fill*URL` helpers on models when the module exposes files or images
- Keep association writes in the service layer transaction
- Follow the existing soft-delete unique-field protection pattern when a unique field exists

### Permission and Auth Side Effects

When a mutation changes roles, menus, APIs, or current-user access behavior:

- Inspect `server/service/role.go`
- Inspect `server/service/user.go`
- Reuse existing cache-clearing and token invalidation flows
- Do not invent a separate permission refresh mechanism if the existing cache or token path already handles it

## Frontend Conventions

### Page Container Pattern

For a standard admin page:

- Keep `index.vue` as the orchestration container
- Use `ProTable` for the main table layout
- Keep search, toolbar, and table in one compact page
- Put create, edit, and other non-trivial popup or drawer content into local `components/`
- If the user does not specify the interaction form, default create and edit flows to drawers
- Extract larger forms, drawers, and preview blocks into local `components/` when the page becomes large

### Table and Permission Pattern

- Define base columns first
- Use `useTableColumns(baseColumns, actionColumn, rowActionPermissions)`
- Put toolbar permissions on toolbar buttons with `v-permission`
- Put row-action permissions on the row buttons themselves with `v-permission`
- Do not let add-only permissions create an empty action column

See `references/list-page-pattern.md` for the concrete pattern.

### API and Type Pattern

- Keep frontend request wrappers in `web/src/api/<module>.ts`
- Reuse `ApiResponse<T>` and `PageResponse<T>` from `web/src/types/index.ts`
- Add module-local types under `web/src/types/<module>.ts` only when the neighboring module or generator pattern already does that

### Frontend Closure Default

When adding or exposing an interaction:

- Do not leave a visible button without a handler
- Do not show files or images without preview or open behavior
- Do not stop at data display if the expected local CRUD loop is still obviously incomplete
- Do not keep create or edit form markup inline in `index.vue` when it should be a reusable drawer or popup component

## Workflow

1. Identify whether the task is extending an existing module or creating a new one.
2. Inspect the neighboring backend files and frontend files first.
3. Keep backend placement, route registration, request structs, and service logic aligned with current live code.
4. Keep frontend page, API, permission, and shared-component usage aligned with current live code.
5. Complete the smallest usable CRUD loop instead of leaving half-finished UI or route placeholders.
6. Run targeted verification.

## Verification

After substantial edits, prefer these checks when available:

- backend: `go test ./...`
- frontend typecheck: `npm run typecheck`
- frontend build when the change is broad: `npm run build`

Run narrower commands when the workspace has known unrelated failures and the change only affects one side.

## Hard Rules

- Use UTF-8.
- Follow current workspace files instead of `XTMS` or Spring Boot conventions.
- Prefer live files over remembered patterns.
- Do not create new nested backend package trees by habit.
- Do not move admin CRUD pages out of `web/src/views/system` unless the user explicitly asks for a different product structure.
- Do not bypass `ProTable`, `useTableColumns`, `v-permission`, or shared upload/preview components without a concrete gap.
- Do not leave an empty action column.
- Do not leave clickable frontend affordances without a real handler.
- Do not default to modal-based create or edit flows unless the user explicitly asks for a modal.

## Self-Check

Before finishing, verify:

- Did I inspect a real neighboring module first?
- Did I follow the current flat-file backend structure?
- Did I keep backend and frontend placement consistent with the live repo?
- Did I reuse `ProTable`, permission helpers, and shared upload/preview components where applicable?
- Did I avoid introducing Spring Boot or `XTMS` path assumptions?
- Did I keep newly exposed interactions usable instead of placeholder-only?
