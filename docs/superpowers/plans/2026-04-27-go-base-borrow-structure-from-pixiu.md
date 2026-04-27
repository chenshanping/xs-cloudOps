# go-base Borrow Structure From Pixiu Plan

> **For agentic workers:** Borrow only directory partitioning ideas from `E:\go_project\pixiu`. Do not copy its business abstractions, Kubernetes-oriented modules, or `pkg` layout wholesale into `go-base`.

**Goal:** Define how `go-base` can absorb the useful directory-organization ideas from `pixiu` while preserving the current `server/api/v1` + `server/service` + `server/model` + `server/router/modules` backbone.

**Architecture:** Treat `pixiu` as a reference for boundary clarity, not as a target architecture. The adaptation path is “align names, separate cross-cutting HTTP concerns, split oversized feature files, and keep transport/business/persistence responsibilities easier to locate.”

**Tech Stack:** Go, Gin, GORM

---

## What Pixiu Is Doing Well

From the observed structure in `E:\go_project\pixiu`:

- `cmd/` owns startup and assembly:
  - [cmd/app/server.go](/E:/go_project/pixiu/cmd/app/server.go)
  - [cmd/app/options/options.go](/E:/go_project/pixiu/cmd/app/options/options.go)
- `api/server/` owns HTTP-facing cross-cutting concerns:
  - `middleware`
  - `router`
  - `validator`
  - `httputils`
  - `errors`
  - `httpstatus`
- `api/server/router/<feature>/` is split by business module and maps routes to controller calls.
- `pkg/controller/<feature>/` is also split by business module, mirroring router naming.
- `pkg/db/` and `pkg/db/model/` split persistence access from model definitions.
- `pkg/types/` holds shared request/domain types separately from controller and DB code.
- `pkg/client/`, `pkg/jobmanager/`, `pkg/static/`, `pkg/util/` keep external integration and generic capabilities out of feature files.

This gives `pixiu` three strong properties:

1. You can find one business feature across router/controller/db by the same feature name.
2. HTTP infrastructure is not hidden inside a feature file.
3. Cross-cutting technical code has a home, so feature files stay narrower.

## What go-base Should Learn

`go-base` should borrow the following ideas.

### 1. Mirror feature names across layers

**Pixiu pattern**

- `api/server/router/user`
- `pkg/controller/user`
- `pkg/db/user.go`

**go-base adaptation**

- Keep:
  - `server/router/modules`
  - `server/api/v1`
  - `server/service`
  - `server/model`
- But make file names and ownership more obviously aligned by feature.

**Apply to go-base**

- Align names like `ai`, `auth`, `user`, `role`, `menu`, `dept`, `dict`, `storage`, `file` across router/api/service/request files.
- Reduce mixed naming like:
  - `server/router/modules/sys_api.go`
  - `server/router/modules/sys_config.go`
  while handlers/services are just `api.go` and `config.go`.

**Concrete target**

- If a feature is called `api` in handler/service, router module naming should follow the same feature wording unless there is a strong reason not to.

### 2. Give HTTP cross-cutting code a stable shared location

**Pixiu pattern**

- `api/server/middleware`
- `api/server/validator`
- `api/server/httputils`
- `api/server/errors`

**go-base adaptation**

- Keep `server/middleware` as-is.
- Add or strengthen shared HTTP-support locations rather than leaving helpers buried in feature files.

**Apply to go-base**

- Move shared route registration helper `R(...)` out of [auth.go](/E:/go_project/go-base/server/router/modules/auth.go:43).
- Keep route infrastructure together in `server/router/modules/base.go` or a dedicated shared module file.
- If more HTTP-specific helpers appear later, place them in explicit shared locations under `server/router` or `server/api`, not in random feature files.

### 3. Split request DTOs by feature, not by “one big request.go”

**Pixiu pattern**

- `pkg/types/` acts as a dedicated shared type area.

**go-base adaptation**

- Keep `server/model/request` because it already exists and matches the current repo.
- Do not create a new top-level `types` package just to imitate `pixiu`.

**Apply to go-base**

- Split [request.go](/E:/go_project/go-base/server/model/request/request.go) into:
  - `auth_request.go`
  - `user_request.go`
  - `role_request.go`
  - `dept_request.go`
  - `menu_request.go`
  - `api_request.go`
  - `log_request.go`
- Keep the same `package request`.
- Keep `PageRequest` in a small shared request base file.

This borrows the “types have their own zone” idea without changing `go-base`’s top-level layout.

### 4. Let large business areas own multiple focused files

**Pixiu pattern**

- A heavier area like `pkg/controller/plan/` is split into many files:
  - `plan.go`
  - `plan_task.go`
  - `plan_node.go`
  - `register.go`
  - `render.go`
  - `worker.go`

**go-base adaptation**

- Keep `server/service` as a flat package, but split large features into multiple files within that package.

**Apply to go-base**

- Split `server/service/ai.go` into focused files such as:
  - `ai_conversation.go`
  - `ai_context.go`
  - `ai_files.go`
  - `ai_client.go`
  - `ai_stream.go`
- If needed later, split `server/service/user.go` similarly by responsibility rather than by creating new packages.

This is the highest-value structural lesson from `pixiu`: heavy features should own multiple files instead of one god file.

### 5. Separate external integration helpers from feature orchestration

**Pixiu pattern**

- External or technical dependencies have dedicated homes:
  - `pkg/client`
  - `pkg/static`
  - `pkg/util/*`

**go-base adaptation**

- Do not create a giant new `pkg` directory.
- But do stop mixing external-client behavior, file preprocessing, and business orchestration in one feature service file.

**Apply to go-base**

- In `server/service/ai.go`, separate:
  - AI provider HTTP client logic
  - SSE parsing logic
  - local/remote file preprocessing logic
  - conversation persistence logic
- In other features, if a helper is clearly technical and reused, move it to a precise existing home like `server/utils` or a focused new sub-area under `server/service`.

### 6. Keep bootstrap separate from runtime business code

**Pixiu pattern**

- `cmd/` is clearly the bootstrap layer.

**go-base adaptation**

- `go-base` already has `server/main.go` and `server/initialize/`.
- That is good enough for now; do not migrate to `cmd/` just because `pixiu` does.

**Apply to go-base**

- Keep startup assembly concerns in `main.go` and `initialize/`.
- If bootstrap code grows, prefer moving more setup logic into focused files under `initialize/`, not mixing it into handlers or services.

## What go-base Should NOT Copy

These are the wrong moves for the current repo:

- Do not introduce top-level `pkg/controller`, `pkg/db`, `pkg/types`, or `api/server`.
- Do not rename `server/api/v1` to `controller`.
- Do not add a repository layer just to look more like `pixiu`.
- Do not move all shared code into a new generic `pkg`.
- Do not perform a package explosion for small modules that are still easy to navigate as single files.

`pixiu` is a larger app with a different architectural baseline. `go-base` should borrow its separation habits, not its full package topology.

## Recommended Landing Plan For go-base

### Step 1: Name and file alignment

- Unify feature naming across:
  - `server/router/modules`
  - `server/api/v1`
  - `server/service`
  - `server/model/request`
- Remove surprising feature-independent helpers from feature files.

### Step 2: Shared HTTP support cleanup

- Move route helper `R(...)` into shared router infrastructure.
- Keep HTTP metadata and route registration helpers out of business modules.

### Step 3: Request DTO split

- Break `server/model/request/request.go` by feature.
- Keep import paths unchanged at the package level.

### Step 4: Heavy service split

- Split `server/service/ai.go` first.
- Split `server/service/user.go` second only if the first split proves the pattern useful and low-risk.

### Step 5: Feature-local helper extraction

- For each heavy feature, separate:
  - orchestration
  - persistence-facing logic
  - external integration logic
  - serialization / protocol helpers

This should happen inside the existing package layout, not through top-level package churn.

## Best Reusable Rule To Adopt

The best single rule to import from `pixiu` is:

> A feature may span multiple files and multiple layers, but its name and ownership should be obvious at every layer, and shared infrastructure should never hide inside a random feature file.

For `go-base`, this rule is enough to improve structure substantially without rewriting the backend architecture.

## Execution Priority

If this is turned into code work, execute in this order:

1. Shared router helper extraction
2. `model/request` split
3. `service/ai.go` split
4. `api/v1/auth.go` thinning
5. Naming cleanup for mismatched modules

## Validation Notes

- These recommendations are based only on directory structure, import relationships, and layer mapping observed in `E:\go_project\pixiu`.
- No business logic or implementation behavior from `pixiu` is being copied into `go-base`.
