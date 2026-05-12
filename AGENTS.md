# go-base Agent Rules

## Mandatory

- Use UTF-8 for every file edit.
- Before any project-related implementation, refactor, config change, or command with side effects, read this `AGENTS.md`.
- If a deeper directory later contains another `AGENTS.md`, the deeper file overrides this one for that subtree.

## Priority

1. User explicit instructions
2. This `AGENTS.md`
3. Confirmed decisions already made in the current conversation
4. Existing live code patterns in this repository
5. Default model habits

## Project Map

- Backend: `server/`
- Frontend: `web/`
- Specs and change management: `openspec/`
- Local project workflows: `.windsurf/workflows/`
- Superpowers design/plan docs: `docs/superpowers/`

## Current Persistent Decisions

- This repository uses `server` and `web` as the stable backend/frontend directory names.
- Do not reintroduce backend/frontend code generator features or generator-driven workflow. The project has already moved away from that mode.
- For admin create/edit/other non-trivial popup interactions, default to `Drawer` unless the user explicitly asks for `Modal` or another interaction.
- Non-trivial popup/drawer content should be extracted into local `components/` instead of being kept inline in a large page file.
- Theme/layout work should follow the current Ant Design Vue + Pinia layout preference approach already present in `web/src/layouts` and `web/src/store/ui.ts`.
- For system config features that define security boundaries such as anonymous-readable config keys, auth bypass, exposure white/blacklists, secret visibility, or other high-risk access controls, default to backend-coded policy. Do not add admin UI or DB-configurable toggles for these controls unless the user explicitly asks for that operational model.

## OpenSpec Decision Rule

Decide whether to use `openspec-propose` based on the user's actual requirement, not by default.

Default behavior:

- If the requirement is clear and low risk, implement directly.
- If the requirement is medium-sized but the scope and neighboring patterns are clear, do a short local analysis and then implement directly.
- Only use `openspec-propose` when the change is high risk, cross-module, or changes important system behavior.

Implement directly without `openspec-propose` for:

- Copy changes
- Style or layout adjustments
- Small form field changes
- Small CRUD changes with clear scope
- Small frontend-only or backend-only adjustments
- Tiny config or typo fixes with no behavior change

Usually analyze briefly and then implement directly for:

- Medium feature changes with clear boundaries
- Work that has adjacent implementations to copy from
- Changes limited to one or two modules with controllable risk

Use the full OpenSpec workflow for:

- New business capability
- API behavior change
- Permission, audit, auth, data consistency, or workflow logic changes
- Cross-frontend/backend/SQL permission linkage
- Cross-module refactors
- Any task that needs reviewable specs, resumability, and traceability
- Any request whose scope is still unclear and must be locked before coding

If the user explicitly says `直接做`, prefer direct implementation unless the change is clearly high risk.

## Workflow Auto-Trigger Rules

When the user's request matches one of these patterns, **automatically read the corresponding workflow file and follow its Steps** without the user needing to type the slash command:

| Trigger condition | Workflow to read and follow |
|---|---|
| Creating/modifying backend API + frontend page, CRUD module, new admin page, adding buttons/actions, menu/permission work | `.windsurf/workflows/backend-crud-frontend.md` |
| Writing or modifying SQL upgrade scripts under `server/sql/` | `.windsurf/workflows/sql-upgrade-guardrails.md` |

- Read the workflow file with the `read_file` tool, then follow its `## Steps` sequentially.
- If the user says "直接做" or the task is trivially small (one-line fix, typo, style tweak), skip the workflow and implement directly.
- When in doubt, follow the workflow — it is cheaper to follow steps than to miss permissions or bootstrap.

## Required Workflow For Non-Trivial Work

1. Read the relevant code and neighboring modules first. Do not invent structure from memory.
2. First decide whether this request really needs OpenSpec by using the rule above. Do not enter `openspec-propose` automatically.
3. If the request is unclear or has design tradeoffs, start with Superpowers exploration and planning:
   - `brainstorming`
   - `writing-plans`
4. If the work is high risk or changes core behavior, API contract, permission model, audit semantics, data consistency, or architecture, create or update an OpenSpec change before coding:
   - `openspec-explore`
   - `openspec-propose`
   - `openspec-apply-change`
   - `openspec-archive-change`
5. Otherwise, implement directly in small steps and keep the scope tight.
6. Prefer an isolated branch or worktree for risky work.
7. Run real verification commands before claiming completion.
8. Keep code, tasks, and specs aligned whenever OpenSpec is used. Do not archive if they diverge.

## Search-First Rule

For architecture, workflow, dependency, bug, or integration tasks:

- Search first.
- Prefer multi-round search when risk or ambiguity is non-trivial.
- Use the available search tooling before guessing.

## Project Commands

### Backend

- Install deps: `go mod tidy`
- Test: `go test ./...`

Run these in: `server/`

### Frontend

- Install deps: `npm install`
- Dev: `npm run dev`
- Build: `npm run build`
- Type check: `npm run typecheck`

Run these in: `web/`

Known caveat:

- `web` currently has an existing `vue-tsc` environment issue in some sessions. If `npm run typecheck` fails with the known toolchain string replacement error, report it explicitly and still run the strongest unaffected verification available, usually `npm run build` plus targeted checks.

### Frontend Verification — Token Budget Rule (MANDATORY)

**Do NOT run `npm run build` or `npm run typecheck` for verification purposes** in agent sessions. These commands take 1–2 minutes and dump hundreds of lines of asset listings into the chat, which wastes the user's token budget for negligible signal.

Default frontend verification is **manual user click-through on the running dev server**. The agent's job is to:

1. Write code that follows existing patterns (template, script setup, imports).
2. Read the file back after editing to confirm structure.
3. Run targeted static checks ONLY when needed:
   - `grep_search` for the modified symbol/path to confirm no broken references.
   - `find_by_name` to confirm new files exist where expected.
4. Hand off to the user with a concise list of what to click and what to expect.

The agent may run `npm run build` / `npm run typecheck` ONLY when:

- The user explicitly asks for it (e.g. "跑一下 build" / "verify with typecheck").
- A backend route/handler signature changed and silently broke a TS type — and only after the user agrees.
- Final CI-style verification at the end of a feature branch, with the user's explicit approval.

Backend tests (`go test ./...`, `go build ./...`) are NOT covered by this rule — they are fast and high-signal, keep running them.

### Git Operations — Owned by User (MANDATORY)

**The agent does NOT run `git add` / `git commit` / `git push` / `git reset` / `git checkout` / `git diff` / `git log` / `git status` in normal feature work.** Git operations are owned by the user.

Why: long multi-line commit messages with full root-cause analysis routinely cost 600–1000 tokens each; PowerShell encoding pitfalls (BOM via `Out-File -Encoding utf8`, GBK mojibake on Chinese literals, single-quote escape errors) cause retry loops that waste another 500–1500 tokens with zero deliverable value. Over a multi-resource feature branch this adds up to ~15k tokens of pure overhead.

The agent's job is to:

1. Write / edit the code files (this is the actual deliverable).
2. Run verification that is fast and high-signal:
   - `go test ./service/k8s -run TestXxxService -count=1` for the resource just changed.
   - `go build ./...` to catch handler / router / DTO type errors.
   - Targeted `grep_search` / `find_by_name` for cross-file reference integrity.
3. Hand off to the user with a concise summary of changed files + acceptance click list. The user decides commit boundaries, writes the commit message they want, and pushes.

The agent may run git ONLY when:

- The user explicitly asks (e.g. "帮我提交一下" / "push" / "commit this").
- The user asks for branch / history inspection that requires `git log` / `git diff`.

When the agent does run git on request:

- Prefer one-line `git commit -m "title"` over multi-paragraph `-F message_file`.
- Body, if needed, in 2–3 short ASCII bullets only; no Chinese in commit body to avoid PowerShell encoding loops.
- Do NOT chain `git status` / `git log` after every commit; only inspect when scope is unclear.

This rule supersedes any prior implicit "commit after each logical change" behavior.

## Frontend Rules

- Follow existing layout, table, form, and permission patterns before introducing new UI structure.
- Do not expose clickable UI affordances without a real handler or visible feedback.
- Reuse shared components when they already fit: uploads, previews, tables, icons, permission helpers.
- Keep list/detail CRUD pages compact; do not add decorative wrappers the user did not request.

## Backend Rules

- Follow existing flat module placement under `server/api/v1`, `server/service`, `server/model`, and `server/router/modules`.
- Keep business logic in services, not handlers.
- Reuse existing response helpers, auth flow, cache invalidation, and permission refresh patterns.
- Uploaded file references must be centralized in `sys_file_reference`. Store uploaded file IDs in business data, resolve URLs from `sys_file` at read time, and do not add module-specific reference scans to `server/service/file/file.go`.
- Treat security-sensitive system config rules as backend-owned policy by default. If a config item changes anonymous exposure, auth boundaries, secret visibility, or other high-risk behavior, prefer code-defined allow/deny lists and reviewed deployment changes over admin-page runtime configuration.
- **No database foreign keys.** Do not use GORM association tags that generate FK constraints (e.g. `gorm:"foreignKey:..."`). Keep referential integrity in application code. If a model needs to reference another table, store the ID column only and query the related record manually in the service layer. This avoids AutoMigrate FK failures when nullable/zero-value IDs exist.
- **Soft delete + unique index must be explicitly resolved.** Any model embedding `BaseModel` (or otherwise containing `gorm.DeletedAt`) that also declares any `uniqueIndex` / `uniqueIndex:xxx` field MUST pick a documented strategy. See `## GORM Soft Delete + Unique Index Rule` below. The default soft delete + single-column unique index combination silently breaks "delete then recreate same-name" workflows and must not be left implicit.

## GORM Soft Delete + Unique Index Rule

### Why This Rule Exists

`BaseModel` embeds `gorm.DeletedAt`, so by default `db.Delete(&row)` performs a **soft delete** (sets `deleted_at = now()`, the row stays). A field marked with `uniqueIndex` becomes a **single-column** unique index that does NOT include `deleted_at`. Therefore a soft-deleted row continues to occupy the unique key, and any subsequent `INSERT` with the same value fails with `Duplicate entry`.

Concrete failure mode actually hit in this repo (K8s cluster module): user registers cluster `prod`, deletes it, tries to register `prod` again, gets `集群名已存在` — because the original row is still physically present with `deleted_at` set.

This pattern exists across many existing models (`SysUser.Username`, `SysRole.Code`, `SysDictType.Type`, `SysConfig.Key`, `CmdbHostMetric.HostID`, `CmdbHostAgent.HostID`, `AIProviderConfig.Name`, `K8sCluster.Name`, `K8sClusterCredential.ClusterID`, …). They all share this latent bug; only the K8s ones have been fixed so far.

### When This Rule Applies

A model triggers this rule when BOTH conditions are true:

1. It embeds `BaseModel` or otherwise carries a `gorm.DeletedAt` field.
2. It declares any field tagged with `uniqueIndex` or `uniqueIndex:<name>`.

If both are true, the author MUST pick one of the strategies below and state the choice in the PR/commit description.

### Allowed Strategies

**Strategy A — Physical delete (default, prefer this).**

- Service-layer `Delete` MUST use `db.Unscoped().Delete(...)` for both the main model and any soft-deletable related rows uniqueIndexed by foreign id (e.g. `K8sClusterCredential.ClusterID`).
- Use this when: the entity has no audit/recovery requirement, OR delete must be semantically permanent (e.g. credential wiping, cache rows, ephemeral metadata).
- Tradeoff: loses ability to "restore" a deleted row.

**Strategy B — Composite unique index `(field, deleted_at)`.**

- Change the tag to `uniqueIndex:uk_xxx;index:uk_xxx,priority:1` plus a second tagged `uniqueIndex:uk_xxx;index:uk_xxx,priority:2` on `deleted_at`, OR add a raw composite index in the SQL upgrade script and remove the single-column tag.
- Use this when: audit trail of past entities matters AND you want to keep soft delete behavior.
- MySQL allows multiple NULLs in a composite unique index, so multiple historic deleted rows coexist while only one live row (`deleted_at IS NULL`) is allowed.
- Requires a matching SQL upgrade script (`server/sql/`) — index changes affect existing installations.

**Strategy C — Revival on conflict in Create path.**

- `Create` first looks up any soft-deleted row with the same unique value; if found, it physically removes the residue (or restores it) before inserting.
- Use this rarely — it complicates Create logic. Prefer A or B.

### Mandatory Test Coverage

Every model that triggers this rule MUST include a service-layer test like:

```
TestXxxService_Delete_AllowsSameUniqueValueRecreate
```

The test must: create a row → delete it → create another row with the same unique field value → assert the second create succeeds. This single test is the regression gate; without it, future PRs can silently regress to the broken default.

### Decision Record

When adding a new model that triggers this rule, state the strategy in commit message or PR description, e.g.:

```
soft-delete-strategy: A (physical delete)
```

When auditing or fixing an existing model, the same record applies.

### Existing Known-Affected Models (Audit Backlog)

The following models currently use the broken default (single-column unique index + soft delete) but have NOT yet been fixed. They should be fixed lazily as bugs are reported, or proactively in a dedicated cleanup pass:

- `SysUser.Username`
- `SysRole.Code`
- `SysDictType.Type`
- `SysConfig.Key`
- `CmdbHostMetric.HostID`
- `CmdbHostAgent.HostID`
- `AIProviderConfig.Name`

(`K8sCluster.Name` and `K8sClusterCredential.ClusterID` are fixed via Strategy A.)

Do not silently fix unrelated models when a user request is scoped to one feature; note them in PR review and let the user decide.

## Built-In Bootstrap Rules

- Treat startup repair and built-in data bootstrap under `server/initialize/` as `fill missing data only` by default, not `sync defaults back into existing rows`.
- Do not overwrite user-editable menu/config/API metadata on restart unless the user explicitly asks for a forced reset or a versioned migration:
  - menu fields such as `name`, `icon`, `sort`, `hidden`, `path`, `component`
  - config display values
  - other admin-maintained presentation metadata
- For built-in menu/config/API repair, prefer create-only or missing-field-only patterns such as `FirstOrCreate + Attrs`, `OnConflict DoNothing`, or explicit null/missing checks.
- Do not use `Assign(...)` in startup repair paths unless the requirement is explicitly to push defaults into existing records.
- Any change to startup bootstrap, seed repair, or built-in menu/API/config补齐 logic must include a regression test proving that customized existing data is not overwritten on restart.

## SQL Upgrade Rules

### When An Upgrade Script Is Mandatory

Any change that affects the **persisted state of existing installations** must ship a matching incremental SQL script under `server/sql/`, even when the change is made through Go code (GORM `AutoMigrate`, `initialize/` bootstrap, seed helpers). This includes but is not limited to:

- Adding, renaming, or removing tables, columns, indexes, or unique constraints
- Changing column types, defaults, nullability, or character sets
- Changing built-in menu rows: `path`, `name`, `icon`, `sort`, `hidden`, `component`, `parent_id`, `permission`, `status`
- Changing built-in API rows, permission rows, role-menu / role-api bindings
- Changing built-in config rows: default value, display label, exposure flag, scope
- Changing seed data that existing deployments have already received
- Any data migration or backfill required by a behavior change

Rule of thumb: **if the `initialize/` code uses a create-only pattern (`FirstOrCreate + Attrs`, `OnConflict DoNothing`, missing-field checks)**, then any change to the intended values for existing rows will NOT take effect on existing installs — an upgrade script is required to reconcile the drift.

Upgrade scripts must:

- Be idempotent and safe to rerun
- Only update rows that actually need updating (use `WHERE` guards on the old value)
- Be placed in `server/sql/` with a descriptive filename, e.g. `update_<feature>_<change>.sql`
- Be referenced in the PR/commit description or task notes so operators know to run it

### SQL Script Authoring Rules

- Any change under `server/sql/` must use the `/sql-upgrade-guardrails` workflow before writing or modifying the script.
- Treat this repository as `Oracle MySQL` by default, not MariaDB. Do not assume MySQL supports MariaDB DDL syntax.
- **No foreign key constraints in DDL.** Do not add `FOREIGN KEY` or `REFERENCES` clauses in `CREATE TABLE` or `ALTER TABLE` statements. Use plain columns with indexes for cross-table references. Enforce referential integrity in application code only.
- Before editing an incremental SQL script, inspect the baseline snapshot `go-base.sql` and the nearest related upgrade scripts.
- Do not use unsupported MySQL incremental DDL patterns such as:
  - `ALTER TABLE ... ADD COLUMN IF NOT EXISTS ...`
  - `ALTER TABLE ... ADD INDEX IF NOT EXISTS ...`
  - other unverified `IF [NOT] EXISTS` forms inside `ALTER TABLE`
- For additive DDL that may run on mixed states, use idempotent guards based on `information_schema`, dynamic SQL, or the repository's migration mechanism.
- Seed data, permission rows, menu rows, API rows, and config rows in upgrade scripts must be duplicate-safe.
- Keep incremental SQL limited to the current feature. Do not mix unrelated schema cleanup into the same script.
- Do not rewrite `go-base.sql` for normal feature delivery unless the user explicitly asks for a baseline refresh.
- If SQL compatibility is uncertain, search first and verify the exact MySQL syntax before editing.

## OpenSpec Conventions In This Repo

- Project config lives at `openspec/config.yaml`.
- Current behavior source of truth belongs in `openspec/specs/`.
- Planned work belongs in `openspec/changes/<change-name>/`.
- Proposal/design/spec/tasks should be reviewed before implementation when the task is non-trivial.

## Superpowers Conventions In This Repo

- Workflow definitions live in `.windsurf/workflows/` and are invoked with slash commands (e.g. `/brainstorming`, `/writing-plans`, `/executing-plans`, `/openspec-propose`).
- Brainstorm specs should be written to `docs/superpowers/specs/`.
- Implementation plans should be written to `docs/superpowers/plans/`.
- Prefer `/executing-plans` only after design and plan are approved.

## Scope Discipline

- Do not expand a local task into a system-wide refactor without user approval.
- Do not add dependencies, abstractions, or polish work just because they seem nicer.
- If you find adjacent issues, note them separately instead of silently fixing everything.

## Completion Gate

Before declaring work complete, verify:

- Relevant specs or plans were created/updated when required
- Code changes match the agreed scope
- Verification commands were actually run
- Any known blocker or existing unrelated failure is reported clearly
- **If the change affects persisted DB state of existing installs** (schema, built-in menu/API/config/permission rows, seed data), confirm a matching SQL upgrade script exists under `server/sql/` and is reported to the user. Changes made only through Go `initialize/` code without an SQL script are considered incomplete for existing installations.
- If `server/sql/` changed, report whether the upgrade script was verified for MySQL syntax and whether rerun idempotence was checked
- If `server/initialize/` changed, report whether startup rerun behavior was verified and whether customized built-in data survives restart.
