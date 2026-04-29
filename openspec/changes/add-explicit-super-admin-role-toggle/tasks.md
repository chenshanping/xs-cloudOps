## 1. Spec and schema baseline

- [x] 1.1 Add a new OpenSpec capability for explicit super admin roles and align the proposal/design/spec artifacts.
- [x] 1.2 Confirm this change stays separate from role permission UX changes and dual-track permission selection changes.

## 2. Backend persistence and role APIs

- [x] 2.1 Add a persistent `is_super_admin` field to `sys_role` with a safe MySQL upgrade path.
- [x] 2.2 Extend role create/update requests and role responses to read/write `is_super_admin`.
- [x] 2.3 Keep existing menu/API role bindings intact when toggling super admin on or off.

## 3. Runtime authorization behavior

- [x] 3.1 Remove hidden `admin` name-based Casbin bypass logic from runtime authorization.
- [x] 3.2 Apply explicit super admin role semantics to backend API access.
- [x] 3.3 Apply explicit super admin role semantics to frontend permission-code and visible-menu resolution.
- [x] 3.4 Add or update regression tests proving non-super-admin `admin` roles no longer bypass while explicit super admin roles do bypass.

## 4. Frontend role management

- [x] 4.1 Add an explicit super admin switch to `RoleFormDrawer.vue` and related page state types.
- [x] 4.2 Ensure role edit/reopen correctly replays the saved super admin state.
- [x] 4.3 If needed, show super admin state in the role list for quick identification.

## 5. Bootstrap and compatibility

- [x] 5.1 Define and implement bootstrap/upgrade behavior for historical built-in admin roles.
- [x] 5.2 Verify the migration is idempotent and does not rely on permanent role-code backdoors after upgrade.

## 6. Verification

- [x] 6.1 Run `cd server && go test ./tests -count=1`.
- [x] 6.2 Run `cd server && go test ./...`.
- [x] 6.3 Run `cd web && npm run build`.
- [x] 6.4 Run `cd web && npm run typecheck`, or explicitly report the known existing `vue-tsc` toolchain issue if it still blocks typecheck.
- [ ] 6.5 Manually verify:
  - explicit super admin role can access unassigned menus/buttons/APIs
  - ordinary role with `code=admin` but `is_super_admin=false` cannot bypass
  - toggling super admin off restores ordinary fine-grained behavior
