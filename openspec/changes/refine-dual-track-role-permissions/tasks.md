## 1. OpenSpec and behavior audit

- [x] 1.1 Add a new OpenSpec capability for admin role permission management that explicitly documents the dual-track menu/API model.
- [x] 1.2 Audit the current backend and frontend permission flows against that capability and record only observable first-order gaps.

## 2. Backend dual-track hardening

- [x] 2.1 Verify role detail loading always returns saved menus and APIs so the permission drawer can reopen with correct selections.
- [x] 2.2 Verify button-only menu grants preserve the required ancestor menu chain for `userinfo.menus` and dynamic menu visibility.
- [x] 2.3 Verify role API assignment immediately refreshes Casbin runtime policies and persisted `casbin_rule` rows.

## 3. Frontend role permission clarity

- [x] 3.1 Update the role permission drawer to keep separate menu and API authorization areas and add concise responsibility hints for each tab.
- [x] 3.2 Keep menu permission codes and API method/path metadata visible enough for administrators to diagnose missing frontend or backend permissions.
- [x] 3.3 Confirm the drawer save flow preserves the current separate menu/API submissions without introducing automatic cross-authorization.

## 4. Verification

- [ ] 4.1 Run `cd server && go test ./...`.
- [x] 4.2 Run `cd web && npm run build`.
- [x] 4.3 Run `cd web && npm run typecheck`, or explicitly report the known existing `vue-tsc` toolchain issue if it still blocks typecheck.
- [ ] 4.4 Manually verify three scenarios: menu-only grants affect frontend visibility only, API-only grants affect backend access only, and combined grants satisfy both sides.
