## 1. OpenSpec and scope baseline

- [x] 1.1 Add OpenSpec capability deltas for role permission governance and menu API association.
- [x] 1.2 Review the existing `refine-dual-track-role-permissions` change and ensure this new change extends it instead of contradicting it.

## 2. Persistence and backend model groundwork

- [x] 2.1 Add the menu API association persistence model and migration-safe bootstrap/update logic.
- [x] 2.2 Load menu API associations through backend menu services without breaking existing menu tree and user info behavior.
- [x] 2.3 Keep existing `sys_role_menu`, `sys_role_api`, and `casbin_rule` semantics intact while introducing inherited API resolution.

## 3. Menu API association management

- [x] 3.1 Add backend endpoints or request handling for reading and updating menu API associations.
- [x] 3.2 Update menu management UI to let administrators maintain associated APIs for menu and button items.
- [x] 3.3 Ensure directory-type menus do not expose unrelated API association behavior.

## 4. Atomic role permission save

- [x] 4.1 Add a single role permission save endpoint that accepts menu IDs, direct API IDs, and feature data scopes together.
- [x] 4.2 Implement transaction-based replacement of role menus, direct APIs, feature data scopes, final API union, Casbin rules, and cache invalidation.
- [x] 4.3 Preserve rollback behavior so any failure leaves stored role permissions and runtime Casbin state unchanged.

## 5. Role permission UI adaptation

- [x] 5.1 Update the role permission drawer to submit one unified payload instead of separate menu/API/data-scope saves.
- [x] 5.2 Show which APIs are directly granted versus inherited from selected menus.
- [x] 5.3 Keep the current drawer-based workflow compact and aligned with existing admin role page patterns.

## 6. Casbin whitelist tightening

- [x] 6.1 Remove suffix-based whitelist bypasses and replace them with explicit route entries only.
- [ ] 6.2 Verify all intended self-service endpoints remain reachable for authenticated users after the whitelist change.

## 7. Operation log governance

- [x] 7.1 Add centralized request/response field masking for common sensitive values.
- [x] 7.2 Replace uncontrolled per-request goroutine writes with a more controlled async logging strategy.
- [ ] 7.3 Ensure permission-governance endpoints still emit useful but sanitized audit records.

## 8. Verification

- [x] 8.1 Add backend regression tests for menu API association inheritance, atomic role permission saves, rollback on failure, and immediate Casbin updates.
- [x] 8.2 Add backend tests for tightened whitelist behavior.
- [x] 8.3 Add backend tests or focused checks for operation log masking behavior.
- [x] 8.4 Run `cd server && go test ./...`.
- [x] 8.5 Run `cd web && npm run build`.
- [x] 8.6 Run `cd web && npm run typecheck`, or explicitly report the known existing `vue-tsc` toolchain issue if it still blocks typecheck.
- [ ] 8.7 Manually verify: menu selection inherits expected APIs, direct API grants still work without menus, and failed permission saves do not partially apply.
