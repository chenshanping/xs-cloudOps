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
- Local project skills: `.codex/skills/`
- Superpowers design/plan docs: `docs/superpowers/`

## Current Persistent Decisions

- This repository uses `server` and `web` as the stable backend/frontend directory names.
- Do not reintroduce backend/frontend code generator features or generator-driven workflow. The project has already moved away from that mode.
- For admin create/edit/other non-trivial popup interactions, default to `Drawer` unless the user explicitly asks for `Modal` or another interaction.
- Non-trivial popup/drawer content should be extracted into local `components/` instead of being kept inline in a large page file.
- Theme/layout work should follow the current Ant Design Vue + Pinia layout preference approach already present in `web/src/layouts` and `web/src/store/ui.ts`.

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

## Frontend Rules

- Follow existing layout, table, form, and permission patterns before introducing new UI structure.
- Do not expose clickable UI affordances without a real handler or visible feedback.
- Reuse shared components when they already fit: uploads, previews, tables, icons, permission helpers.
- Keep list/detail CRUD pages compact; do not add decorative wrappers the user did not request.

## Backend Rules

- Follow existing flat module placement under `server/api/v1`, `server/service`, `server/model`, and `server/router/modules`.
- Keep business logic in services, not handlers.
- Reuse existing response helpers, auth flow, cache invalidation, and permission refresh patterns.

## OpenSpec Conventions In This Repo

- Project config lives at `openspec/config.yaml`.
- Current behavior source of truth belongs in `openspec/specs/`.
- Planned work belongs in `openspec/changes/<change-name>/`.
- Proposal/design/spec/tasks should be reviewed before implementation when the task is non-trivial.

## Superpowers Conventions In This Repo

- Local Superpowers skills live in `.codex/skills/`.
- Brainstorm specs should be written to `docs/superpowers/specs/`.
- Implementation plans should be written to `docs/superpowers/plans/`.
- Prefer `subagent-driven-development` or `executing-plans` only after design and plan are approved.

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
