# Backend Structure Cleanup Plan

> **For agentic workers:** Keep the existing `server/api/v1` + `server/service` + `server/model` + `server/router/modules` backbone. Do not rewrite this into Clean Architecture, DDD, or `internal/domain/repository` in one pass.

**Goal:** Audit the current backend structure, identify the concrete sources of “目录和逻辑混乱”, and define an incremental cleanup plan that another engineer can execute without changing system behavior.

**Architecture:** Preserve the current flat layered layout. Improve boundaries inside the existing layers by thinning handlers, splitting oversized service files by responsibility, extracting shared infrastructure helpers out of feature files, and reducing direct coupling between service code and HTTP DTOs.

**Tech Stack:** Go, Gin, GORM, Redis, Casbin

---

## Confirmed Constraints

- Keep `server/` as the backend root.
- Keep the current top-level layers: `api/v1`, `service`, `model`, `router/modules`.
- Do not introduce a new repository layer, code generator flow, or a full dependency-injection rewrite in this round.
- Prefer medium cleanup inside existing layers over structural reinvention.

## Findings

### 1. Shared router infrastructure is stored inside a feature file

**Evidence**

- `server/router/modules/auth.go:43-56` defines the shared `R(...)` helper that every route module uses.
- `server/router/modules/base.go:5-23` already contains shared module infrastructure (`RouterModule`, `RegisterModule`, `GetAllModules`).

**Why this is messy**

- `auth.go` is no longer only “认证模块”; it is also a package-wide infrastructure file.
- Route registration behavior now depends on engineers knowing that the shared helper lives in an unrelated feature file.

**Required cleanup**

- Move `R(...)` into `server/router/modules/base.go` or a new shared file such as `server/router/modules/register.go`.
- Keep feature files limited to feature routes.

### 2. `api/v1/auth.go` is doing orchestration and infrastructure work, not only HTTP transport

**Evidence**

- `server/api/v1/auth.go:30-88` contains login lock checks, captcha policy branching, retry counting, lock messaging, and login log assembly.
- `server/api/v1/auth.go:214-244` reads Redis directly and manually parses the reset token’s user id.
- `server/api/v1/auth.go:324-368` contains cache hit/fallback logic plus user/menu/permission aggregation.

**Why this is messy**

- Handler code is mixing HTTP concerns with business process orchestration, cache policy, Redis access, and log object assembly.
- Auth behavior is scattered across `service.User`, `service.Captcha`, `service.Log`, `service.Cache`, `global.Redis`, and `utils`, with the handler acting as the coordinator.

**Required cleanup**

- Keep handlers in `api/v1`, but reduce them to `bind -> call service orchestration -> respond`.
- Add a focused auth-oriented service entry point to own login flow, reset-token flow, and current-user aggregate loading.

### 3. AI streaming transport logic is split across handler and service layers

**Evidence**

- `server/api/v1/ai.go:115-195` handles SSE headers, scanner loop, chunk parsing, event formatting, and final persistence trigger.
- `server/service/ai.go:99-114` defines stream chunk DTOs.
- `server/service/ai.go:788-813` also contains `ParseSSEStream(...)`, another SSE parsing path.

**Why this is messy**

- The transport protocol boundary is unclear. Some stream concerns live in the handler, some in the service.
- Similar parsing logic exists in two places, which increases drift risk when the AI provider response shape changes.

**Required cleanup**

- Choose one owner for stream parsing and chunk transformation.
- Recommended: keep SSE response writing in `api/v1`, but move provider stream parsing and normalized chunk emission into split service files.

### 4. `service/ai.go` is a god file with multiple unrelated responsibilities

**Evidence**

- `server/service/ai.go:116-199` handles conversation listing/creation.
- `server/service/ai.go:218-402` handles chat flow, persistence, and title updates.
- `server/service/ai.go:468-579` handles context assembly and message/file composition.
- `server/service/ai.go:581-677` handles local file reads, HTTP file fetch, base64 conversion.
- `server/service/ai.go:679-876` handles external AI HTTP calls, stream setup, and config testing.
- File size is currently the largest service file in the repo (`server/service/ai.go`, about 25 KB by file size).

**Why this is messy**

- One file currently owns conversation CRUD, persistence, file preprocessing, provider client behavior, and stream parsing.
- The file is large enough that safe local changes require reading too many unrelated concerns.

**Required cleanup**

- Split inside `server/service` without changing the package layout:
  - `ai_conversation.go`
  - `ai_context.go`
  - `ai_files.go`
  - `ai_client.go`
  - `ai_stream.go`
- Keep the exported service surface stable during the split to reduce blast radius.

### 5. Service layer is tightly coupled to HTTP request DTOs

**Evidence**

- Multiple service files import `server/model/request`, including:
  - `server/service/ai.go`
  - `server/service/api.go`
  - `server/service/dept.go`
  - `server/service/dict.go`
  - `server/service/log.go`
  - `server/service/menu.go`
  - `server/service/role.go`
  - `server/service/user.go`
- Concrete examples:
  - `server/service/ai.go:138` uses `*request.ConversationListRequest`
  - `server/service/user.go:130` uses `*request.UserListRequest`
  - `server/service/api.go:15` uses `*request.ApiListRequest`

**Why this is messy**

- The service layer is not consuming business-oriented inputs; it is consuming Gin binding DTOs directly.
- This makes transport concerns leak downward and makes non-HTTP reuse harder.

**Required cleanup**

- Do not rewrite every service signature at once.
- Start with the highest-churn modules (`auth`, `ai`, `user`) and introduce service-local input structs or feature-specific query structs.

### 6. Global singleton dependencies are the default boundary, which hides real module dependencies

**Evidence**

- `server/global/global.go:13-20` exposes `Config`, `Log`, `DB`, `Redis`, and `Enforcer` as global singletons.
- Service files reference globals heavily. Current usage counts from a code scan:
  - `server/service/ai.go`: 30 direct `global.*` references
  - `server/service/user.go`: 27
  - `server/service/dict.go`: 23
  - `server/service/cache.go`: 21
  - `server/service/role.go`: 21

**Why this is messy**

- Real dependencies are implicit rather than visible at construction boundaries.
- File splitting is still possible, but unit isolation and future refactors become harder because everything can reach everything through globals.

**Required cleanup**

- Do not attempt a full DI migration now.
- For newly extracted helpers, prefer narrow structs or functions that accept the dependencies they need.
- Keep the old global-backed service entry points as compatibility adapters during migration.

### 7. Request DTOs are concentrated in one large dumping-ground file

**Evidence**

- `server/model/request/request.go` contains auth, user, role, dept, menu, api, profile, and log request types in a single file.
- The file is already over 10 KB and spans unrelated features from `LoginRequest` through `SlowLogListRequest`.
- Only a small subset has been split out so far (`server/model/request/ai_request.go`, `dict_request.go`), which means the split strategy is inconsistent.

**Why this is messy**

- Navigation cost is high.
- Feature boundaries in the request layer do not match the rest of the backend layout.
- New request DTOs are likely to keep accumulating in the same file.

**Required cleanup**

- Keep `package request`, but split files by feature:
  - `auth_request.go`
  - `user_request.go`
  - `role_request.go`
  - `menu_request.go`
  - `api_request.go`
  - `log_request.go`

### 8. Route metadata is a cross-cutting concern, but its parsing logic is duplicated

**Evidence**

- `server/router/modules/*.go` declare route metadata through `R(...)`.
- `server/router/registry/registry.go:99-192` parses request struct fields for route metadata.
- `server/swagger/swagger.go:120-230` consumes registry metadata to generate Swagger.
- `server/swagger/swagger.go:260-374` implements another reflection-based struct-field parsing path.
- `server/api/v1/api.go:118-155` also reads registry metadata and serializes request/response definitions into database records.

**Why this is messy**

- Registry metadata is now consumed by routing, Swagger generation, and API sync, but the reflection/parsing logic is split between packages.
- Tag parsing behavior can drift between `registry` and `swagger`.

**Required cleanup**

- Keep `registry` as the metadata source of truth.
- Extract shared struct-tag parsing into one reusable helper used by both registry and Swagger generation.
- Keep API sync as a consumer, not an alternate parser owner.

### 9. Naming is inconsistent across layers for the same feature

**Evidence**

- Router modules use `server/router/modules/sys_api.go` and `server/router/modules/sys_config.go`.
- Handlers/services use `server/api/v1/api.go`, `server/service/api.go`, `server/api/v1/config.go`, `server/service/config.go`.
- Models use `server/model/sys_api.go`, `server/model/sys_config.go`.

**Why this is messy**

- One feature is referenced as `api`, `sys_api`, `config`, and `sys_config` depending on the layer.
- This increases grep friction and makes cross-layer ownership less obvious.

**Required cleanup**

- Pick one naming convention per feature and apply it consistently at the file level where safe.
- Recommended: keep model names aligned with table/domain naming, but align file names across `api/service/router/modules` for easier navigation.

## Cleanup Principles

1. Preserve the top-level layout; clean inside it.
2. Thin handlers first; do not move business logic upward.
3. Split by responsibility, not by abstract pattern worship.
4. Prefer file-level extraction before package-level reorganization.
5. Avoid big-bang renames that change behavior and structure at the same time.
6. Introduce explicit boundaries only where they pay for themselves immediately.

## Minimum Executable Cleanup Plan

### Phase 1: Low-risk mechanical cleanup

- Move shared route helper `R(...)` out of `server/router/modules/auth.go`.
- Split `server/model/request/request.go` into feature files while keeping `package request`.
- Normalize file naming between `router/modules`, `api/v1`, and `service` where the rename is mechanical and low-risk.

**Expected outcome**

- Navigation improves immediately.
- No business behavior changes.

### Phase 2: Thin the heaviest handlers

- Extract auth flow orchestration out of `server/api/v1/auth.go`:
  - login flow
  - reset-token flow
  - current-user aggregate loading
- Extract AI stream parsing/orchestration so `server/api/v1/ai.go` keeps only SSE response writing and request/response plumbing.

**Expected outcome**

- `api/v1` becomes a transport layer again.
- Auth and AI behavior become easier to test and reason about.

### Phase 3: Split oversized services inside the existing `service` package

- Split `server/service/ai.go` by responsibility first.
- After AI, split any next oversized file only if it still blocks change safety (`user.go` is the next candidate).
- Keep exported service entry points stable during the split to reduce call-site churn.

**Expected outcome**

- Smaller review units.
- Lower conflict risk.
- Fewer “touch one line, read 800 lines” changes.

### Phase 4: Reduce request-DTO coupling in service signatures

- For `auth`, `ai`, and `user`, replace direct `request.*` inputs with service-local command/query structs or smaller explicit parameter groups.
- Keep handler-to-service mapping local to `api/v1`.

**Expected outcome**

- Service layer stops depending directly on HTTP-bound DTO packages.
- Future non-HTTP reuse and tests become easier.

### Phase 5: Contain cross-cutting metadata concerns

- Consolidate struct-tag reflection helpers shared by `registry` and `swagger`.
- Keep `api.SyncApis` as a consumer of route metadata, not a second metadata definition path.

**Expected outcome**

- One parsing rule set for API metadata.
- Lower drift risk between Swagger and API sync output.

## Do First

- Move the shared `R(...)` helper.
- Split `request.go`.
- Thin `auth.go`.
- Split `service/ai.go` by file responsibility, not by package rewrite.

## Do Not Touch Yet

- Do not replace the whole backend with `internal/` + repository + usecase layers.
- Do not remove all globals in one pass.
- Do not redesign SQL schema, permission model, or startup bootstrap as part of this cleanup.
- Do not mix this work with frontend, deployment, or generator reintroduction.

## Risks And Regression Points

- Auth regressions:
  - login lock / retry count
  - captcha validation branches
  - login success/failure log persistence
  - reset-password token validation
  - token refresh / logout invalidation
- AI regressions:
  - `saveConversation` true/false paths
  - SSE chunk forwarding
  - assistant message persistence after stream completion
  - file attachment handling and local-file base64 conversion
  - conversation title auto-update and clear-context behavior
- Metadata regressions:
  - Swagger route generation
  - API sync output consistency
  - auth/public route registration integrity

## Verification Expectations For The Future Refactor

- Run `go test ./...` in `server/` after each phase.
- Add or update focused tests before touching auth orchestration or AI stream behavior.
- For any route metadata refactor, verify:
  - registered routes still load
  - `/swagger/doc.json` remains structurally correct
  - API sync still captures request/response metadata
