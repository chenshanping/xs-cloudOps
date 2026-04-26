## 1. OpenSpec and configuration plumbing

- [x] 1.1 Add a new system config key and default value for controlling whether the admin user profile button is visible, with the default set to hidden.
- [x] 1.2 Expose the new config key through the frontend config store and system config page so administrators can toggle it.

## 2. Backend batch user status implementation

- [x] 2.1 Add a request DTO for batch user status updates and register a new authenticated route for the batch status API.
- [x] 2.2 Implement the user API handler to bind the batch status request, read the current operator identity from JWT context, and return consistent success or failure responses.
- [x] 2.3 Implement the user service batch status update logic with validation for empty IDs, self-disable protection, protected admin protection, transactional status updates, cache clearing, and token invalidation for disabled users.

## 3. Frontend user management updates

- [x] 3.1 Add a frontend API wrapper for the batch user status endpoint.
- [x] 3.2 Update the admin user management page to conditionally render the “身份” action button based on system config, defaulting to hidden when the config is absent.
- [x] 3.3 Add toolbar actions for batch enable and batch disable, including selection checks, confirmation prompts, success feedback, and table refresh behavior.

## 4. Verification

- [x] 4.1 Run backend verification with `cd server && go test ./...`.
- [x] 4.2 Run frontend verification with `cd web && npm run build`.
- [x] 4.3 Run frontend type checking with `cd web && npm run typecheck`, or explicitly report the known existing `vue-tsc` environment issue if it still blocks this check.
