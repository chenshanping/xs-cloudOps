## 1. OpenSpec and contracts

- [x] 1.1 Confirm and keep the proposal, design, and delta specs aligned with the approved “permission foundation” direction.
- [x] 1.2 Define backend/frontend request and response contracts for departments, user department assignment, and role data scope fields.

## 2. Backend department foundation

- [x] 2.1 Add department, role-data-scope, and user-department models plus request DTOs, including validation and association definitions.
- [x] 2.2 Implement department service, API handlers, and router registration for tree query, detail, create, update, and delete.
- [x] 2.3 Implement shared data-scope resolution and query filtering helpers, then enforce them across user management list and protected user operations.
- [x] 2.4 Extend role and user services/APIs to save and return `data_scope`, custom department selections, and `dept_id`.
- [x] 2.5 Add compatibility handling for historical users that remained bound to a department that later became a parent, allowing unchanged legacy bindings during edit while still blocking new parent-department bindings.

## 3. Persistence and built-in bootstrap

- [x] 3.1 Update DB initialization and built-in data repair to migrate new models, seed the root department, and backfill compatible defaults for existing users and roles.
- [x] 3.2 Add safe incremental SQL upgrade scripts for new tables, new columns, root department backfill, and department menu/API/config seed data.

## 4. Frontend department and permission UI

- [x] 4.1 Add department API/types and build the department management page with tree display plus Drawer-based create/edit form component.
- [x] 4.2 Update the user management page and form to display/select departments and send `dept_id`.
- [x] 4.3 Add user-edit compatibility messaging and validation for historical parent-department bindings, while keeping create and reassignment limited to leaf departments.
- [x] 4.4 Update the role management page and form to edit `data_scope`, and when needed, choose custom departments.
- [x] 4.5 Respect the department module display config when exposing department navigation, without removing underlying data handling.

## 5. Verification

- [x] 5.1 Add or update backend tests for department hierarchy validation, delete restrictions, and user data-scope enforcement.
- [x] 5.2 Run backend verification with `cd server && go test ./...`.
- [x] 5.3 Run frontend verification with `cd web && npm run build`.
- [x] 5.4 Run frontend type checking with `cd web && npm run typecheck`, or explicitly report the known `vue-tsc` environment issue if it still blocks this check.
