## 1. OpenSpec and baseline alignment

- [x] 1.1 Create the OpenSpec proposal, design, and capability spec for `admin-ai-config-management`.
- [x] 1.2 Confirm the implementation follows existing AI config persistence and startup repair conventions without changing `ai_config` storage shape.

## 2. Backend provider model discovery

- [x] 2.1 Add request and response DTO coverage plus backend tests for provider model fetch URL resolution, success parsing, and secret-safe failure handling.
- [x] 2.2 Add the authenticated AI route and API handler for `POST /api/v1/ai/providers/models/fetch`.
- [x] 2.3 Implement the AI service/client proxy logic to fetch remote OpenAI-compatible model lists with timeout, URL fallback, and sanitized errors.
- [x] 2.4 Extend AI startup repair tests and logic so the new API metadata is granted to admin/system_admin without overwriting customized built-in metadata.

## 3. Frontend AI config page refactor

- [x] 3.1 Add focused frontend tests for AI config view resolution and any new state helpers used by provider/model management.
- [x] 3.2 Refactor the AI config page into a left provider list and right imported-model detail panel while preserving existing dirty-state and unified save behavior.
- [x] 3.3 Add a provider edit Drawer component for create/edit/default-provider flows using the current editing state only.
- [x] 3.4 Add a provider model management Drawer component that fetches remote models from the backend, supports search and selection, and imports by append-with-dedup into local editing state.

## 4. Verification

- [x] 4.1 Run `cd server && go test ./...`.
- [x] 4.2 Run targeted frontend checks for new tests plus `cd web && npm run build`.
- [x] 4.3 Run `cd web && npm run typecheck`, or explicitly report the known `vue-tsc` toolchain issue if it still blocks verification.
