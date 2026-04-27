# Backend Structure 3-Phase Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reorganize the `go-base` backend in three controlled phases, with a dedicated test document and acceptance gate for each phase before moving forward.

**Architecture:** Preserve the current `server/api/v1` + `server/service` + `server/model` + `server/router/modules` backbone. Phase 1 is mechanical structure cleanup, Phase 2 thins the heaviest handlers without changing external behavior, and Phase 3 splits oversized services and closes the remaining boundary leaks. Each phase must stay releasable on its own.

**Tech Stack:** Go, Gin, GORM, Redis, Casbin

---

## Phase Overview

### Phase 1: Structural Alignment Without Behavior Change

**Goal**

- Improve navigability and file ownership without changing runtime behavior.

**Scope**

- Move the shared route helper `R(...)` out of `server/router/modules/auth.go`.
- Split `server/model/request/request.go` into feature files while keeping `package request`.
- Normalize obviously inconsistent feature file names where the rename is mechanical and low-risk.

**Files**

- Modify:
  - `server/router/modules/base.go`
  - `server/router/modules/auth.go`
  - `server/router/modules/sys_api.go`
  - `server/router/modules/sys_config.go`
  - `server/model/request/request.go`
- Create:
  - `server/router/modules/register.go` if shared helper is not folded into `base.go`
  - `server/model/request/base_request.go`
  - `server/model/request/auth_request.go`
  - `server/model/request/user_request.go`
  - `server/model/request/role_request.go`
  - `server/model/request/dept_request.go`
  - `server/model/request/menu_request.go`
  - `server/model/request/api_request.go`
  - `server/model/request/log_request.go`

- [ ] Read the existing docs first:
  - `docs/superpowers/plans/2026-04-27-backend-structure-cleanup-plan.md`
  - `docs/superpowers/plans/2026-04-27-go-base-borrow-structure-from-pixiu.md`

- [ ] Move `R(...)` into shared router infrastructure and leave feature route files feature-only.

- [ ] Split `request.go` by feature while keeping all existing exported type names unchanged.

- [ ] Run the Phase 1 test document:
  - `docs/superpowers/plans/2026-04-27-backend-structure-phase1-test-doc.md`

- [ ] Do not proceed to Phase 2 unless all Phase 1 checks pass.

### Phase 2: Thin Handlers And Lock Behavior With Tests

**Goal**

- Move orchestration-heavy logic out of the heaviest handlers while preserving endpoint behavior.

**Scope**

- Thin `server/api/v1/auth.go`.
- Thin `server/api/v1/ai.go`.
- Add focused tests around auth and AI transport behavior before moving logic.

**Files**

- Modify:
  - `server/api/v1/auth.go`
  - `server/api/v1/ai.go`
  - `server/service/cache.go`
  - `server/service/captcha.go`
  - `server/service/log.go`
  - `server/service/user.go`
  - `server/service/ai.go` or newly split AI service files if Phase 3 extraction starts earlier for safety
- Create:
  - `server/service/auth_flow.go`
  - `server/service/auth_session.go`
  - `server/service/ai_stream_adapter.go` if needed to isolate transport-facing orchestration
  - `server/api/v1/auth_test.go` if handler-level tests are added
  - `server/api/v1/ai_test.go` if handler-level tests are added

- [ ] Add or update tests around the current auth behavior before extracting orchestration.

- [ ] Extract these auth flows out of `api/v1/auth.go`:
  - login orchestration
  - reset-password token lookup / validation
  - current-user aggregate loading

- [ ] Add or update tests around the current AI handler behavior before moving stream logic.

- [ ] Reduce `api/v1/ai.go` to request binding, response wiring, and SSE response writing only.

- [ ] Run the Phase 2 test document:
  - `docs/superpowers/plans/2026-04-27-backend-structure-phase2-test-doc.md`

- [ ] Do not proceed to Phase 3 unless all Phase 2 checks pass.

### Phase 3: Split Oversized Services And Close Remaining Boundary Leaks

**Goal**

- Turn the largest backend files into smaller, focused units while preserving public behavior and route contracts.

**Scope**

- Split `server/service/ai.go` by responsibility.
- Split `server/service/user.go` only if Phase 2 proves the extraction pattern safe.
- Start reducing direct `request.*` DTO coupling in service signatures for the highest-churn paths.
- Consolidate duplicated route metadata parsing only after behavior is locked.

**Files**

- Modify:
  - `server/service/ai.go`
  - `server/service/user.go`
  - `server/router/registry/registry.go`
  - `server/swagger/swagger.go`
  - `server/api/v1/api.go`
- Create:
  - `server/service/ai_conversation.go`
  - `server/service/ai_context.go`
  - `server/service/ai_files.go`
  - `server/service/ai_client.go`
  - `server/service/ai_stream.go`
  - `server/service/ai_types.go` if chunk or request shapes need a stable shared home
  - `server/service/user_query.go` or similarly focused files if `user.go` is split
  - `server/router/registry/fields.go` if shared metadata parsing is extracted

- [ ] Split `service/ai.go` first, keeping exported entry points stable.

- [ ] Only split `service/user.go` if the Phase 2 and early Phase 3 test coverage is sufficient to catch regressions.

- [ ] Replace service signatures that take `request.*` structs only in touched hot paths:
  - `auth`
  - `ai`
  - `user`

- [ ] Consolidate route metadata field parsing so Swagger and API sync consume the same parsing rules.

- [ ] Run the Phase 3 test document:
  - `docs/superpowers/plans/2026-04-27-backend-structure-phase3-test-doc.md`

- [ ] Stop after Phase 3 and review whether further cleanup is still worth the risk.

## Phase Gates

The following rules are mandatory:

- Do not overlap phases in one unchecked batch.
- Each phase must end with:
  - code complete
  - targeted tests complete
  - `go test ./...` complete in `server/`
  - a human-readable test record filled from that phase’s test doc
- If one phase exposes unrelated structure issues, record them separately and do not silently expand scope.

## Test Document Mapping

- Phase 1:
  - `docs/superpowers/plans/2026-04-27-backend-structure-phase1-test-doc.md`
- Phase 2:
  - `docs/superpowers/plans/2026-04-27-backend-structure-phase2-test-doc.md`
- Phase 3:
  - `docs/superpowers/plans/2026-04-27-backend-structure-phase3-test-doc.md`

## Acceptance Standard

Phase completion requires all three conditions:

1. Automated verification passes.
2. Manual scenario checklist is completed.
3. No route, auth, AI, cache, or Swagger regression is found in the phase-specific test document.

## Execution Notes

- Phase 1 is the safest starting point and should not intentionally change runtime behavior.
- Phase 2 is the first phase that can easily break production behavior, so tests must be written before extraction.
- Phase 3 has the highest internal churn and must be entered only after the first two phases are stable.
