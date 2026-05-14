# Role Page Layout Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Align the admin role management page with the approved mockup while preserving the existing backend API contracts and permission logic.

**Architecture:** Keep the existing `ProTable` and permission drawers, then refine layout, local search, action grouping, and drawer shells in place. Shared drawer context and responsive width handling will be implemented in frontend-only helpers/components.

**Tech Stack:** Vue 3, TypeScript, Ant Design Vue, existing local role page hooks/components

---

### Task 1: Lock Role List Presentation Scope

**Files:**
- Modify: `web/src/views/admin/system/role/index.vue`
- Modify: `web/src/views/admin/system/role/hooks/useRolePage.ts`

- [ ] Add local search state and filtered role list support in the role page hook.
- [ ] Replace the plain panel header with the approved title/subtitle/search/action layout.
- [ ] Keep summary cards and action semantics intact while tightening spacing and visual hierarchy.

### Task 2: Unify Drawer Shells

**Files:**
- Create: `web/src/views/admin/system/role/components/useResponsiveDrawerWidth.ts`
- Modify: `web/src/views/admin/system/role/components/RolePermissionMenuDrawer.vue`
- Modify: `web/src/views/admin/system/role/components/RolePermissionApiDrawer.vue`
- Modify: `web/src/views/admin/system/role/components/RoleDataScopeDrawer.vue`
- Modify: `web/src/views/admin/system/role/components/RoleUsersDrawer.vue`
- Modify: `web/src/views/admin/system/role/components/rolePermissionDrawerShared.css`

- [ ] Add one shared responsive width helper for right drawers.
- [ ] Replace string drawer titles with unified title slots and subtitle context.
- [ ] Apply `maskClosable=false`, `destroyOnClose`, and responsive widths.

### Task 3: Verify Wiring

**Files:**
- Modify: `web/src/views/admin/system/role/index.vue`
- Modify: `web/src/views/admin/system/role/hooks/useRolePage.ts`

- [ ] Read back edited files and confirm template bindings match exported hook state.
- [ ] Search for new symbols and confirm there are no missing imports or stale template references.
- [ ] Hand off a concise click-through checklist for the running dev server.
