## 1. OpenSpec and UX baseline

- [x] 1.1 Add a new OpenSpec capability for grouped role permission selection UX that keeps menu and API permissions separately effective.
- [x] 1.2 Confirm the existing backend role menu/API read-write contracts remain unchanged and sufficient for the grouped drawer.

## 2. Grouped drawer implementation

- [x] 2.1 Refactor the role permission drawer into a same-screen grouped editor with left-side top-level menu navigation and right-side page permission blocks.
- [x] 2.2 Build frontend-only grouping logic that organizes menus, buttons, matched APIs, and uncategorized APIs without introducing automatic cross-authorization.
- [x] 2.3 Keep menu ancestor normalization, role reopen replay, page-level selection helpers, search, and uncategorized API assignment working in the new layout.

## 3. Save and error handling

- [x] 3.1 Keep a single save action that continues to call `assignMenus` and `assignApis` separately.
- [x] 3.2 Distinguish full success, menu-only failure, API-only failure, and double failure without closing the drawer on partial failure.

## 4. Verification

- [ ] 4.1 Run `cd server && go test ./...`.
- [x] 4.2 Run `cd web && npm run build`.
- [x] 4.3 Run `cd web && npm run typecheck`, or explicitly report the known existing `vue-tsc` toolchain issue if it still blocks typecheck.
- [ ] 4.4 Manually verify grouped selection, menu-only, API-only, combined grants, reopen replay, and uncategorized API assignment.
