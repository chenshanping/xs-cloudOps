# Secure Public Config Access Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Stop anonymous callers from reading protected `sys_config` values while keeping login-page/public bootstrap config and authenticated admin config loading functional.

**Architecture:** Keep `POST /api/v1/configs/keys` in place, but split behavior by request mode: anonymous callers only receive server-allowlisted public keys, authenticated callers can still read protected keys after token validation. Split frontend config bootstrap into public and admin key groups so unauthenticated startup stops requesting protected keys.

**Tech Stack:** Go, Gin, GORM, JWT, Redis, Vue 3, Pinia, TypeScript

---

### Task 1: Lock behavior with backend API tests

**Files:**
- Create: `E:\go_project\go-base\server\tests\config_access_test.go`
- Test: `E:\go_project\go-base\server\tests\config_access_test.go`

- [ ] Step 1: Write failing tests for anonymous public-key read, anonymous protected-key denial, and authenticated protected-key read.
- [ ] Step 2: Run `cd server && go test ./tests -run ConfigAccess -v` and confirm at least the protected-key anonymous case fails against current code.

### Task 2: Implement backend config access split

**Files:**
- Modify: `E:\go_project\go-base\server\api\v1\config.go`
- Modify: `E:\go_project\go-base\server\service\configsvc\config.go`

- [ ] Step 1: Add a service-owned public allowlist and helpers to filter anonymous batch requests.
- [ ] Step 2: Add optional token validation inside `ConfigApi.GetConfigsByKeys` and branch between authenticated/full mode and anonymous/allowlist mode.
- [ ] Step 3: Re-run `cd server && go test ./tests -run ConfigAccess -v` until the new tests pass.

### Task 3: Split frontend public and admin config loading

**Files:**
- Modify: `E:\go_project\go-base\web\src\store\config.ts`
- Modify: `E:\go_project\go-base\web\src\main.ts`
- Modify: `E:\go_project\go-base\web\src\App.vue`
- Modify: `E:\go_project\go-base\web\src\router\index.ts`
- Modify: `E:\go_project\go-base\web\src\views\auth\login\index.vue`

- [ ] Step 1: Split config keys into public and admin groups in the store.
- [ ] Step 2: Make default bootstrap load only public keys.
- [ ] Step 3: Add an authenticated admin-config load path and update the existing login/route refresh call sites to use it after authentication.

### Task 4: Verify and close

**Files:**
- Modify: `E:\go_project\go-base\openspec\changes\secure-public-config-access\tasks.md`

- [ ] Step 1: Run `cd server && go test ./...`.
- [ ] Step 2: Run `cd web && npm run build`.
- [ ] Step 3: Run `cd web && npm run typecheck`, or record the known blocker if it still fails for unrelated toolchain reasons.
- [ ] Step 4: Mark completed OpenSpec tasks as done.
