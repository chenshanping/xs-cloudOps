---
description: Standardize CRUD work for the go-base workspace. Use when creating/refactoring admin modules, DB tables, Go backend APIs, or Vue frontend pages.
---

# Backend CRUD Frontend

Discover the live module shape first, then extend it.

## Core Rule

- Do NOT assume Spring Boot paths or XTMS conventions
- Do NOT assume each module has its own directory — many are flat files
- Prefer the current neighboring module over memory

## Read First

Inspect these anchors before editing:
- `server/api/v1/user.go`, `server/service/user.go`, `server/model/sys_user.go`
- `server/router/modules/user.go`, `server/model/request/request.go`
- `web/src/api/user.ts`, `web/src/views/admin/system/user/index.vue`
- `web/src/components/ProTable.vue`, `web/src/utils/permission.ts`
- `web/src/types/index.ts`

## Backend Conventions

- API handlers in `server/api/v1/<module>.go`
- Services in `server/service/<module>.go`
- Models in `server/model/<module>.go`
- Routes in `server/router/modules/<module>.go`
- Response helpers: `response.BadRequest`, `response.Fail`, `response.OkWithData`, `response.OkWithPage`
- No database foreign keys

## Frontend Conventions

- Use `ProTable` for standard search + table pages
- Use `useTableColumns(...)` for permission-aware action columns
- Use `v-permission` on buttons, filter dropdown items in JS
- Default create/edit flows to Drawer (not Modal)
- Reuse shared components: `AvatarUpload`, `ImageUpload`, `FileUpload`, `FilePreview`
- Support dark mode: use `useUiStore().isDark` and semantic CSS variables

## Verification

// turbo
- Backend: `go build ./...` (in `server/`)
// turbo
- Frontend: `npm run build` (in `web/`)

## Self-Check

- Did I inspect a real neighboring module first?
- Did I follow current flat-file backend structure?
- Did I reuse ProTable, permission helpers, shared components?
- Did I keep interactions usable (no placeholder-only buttons)?
- Did I verify dark mode for touched surfaces?
