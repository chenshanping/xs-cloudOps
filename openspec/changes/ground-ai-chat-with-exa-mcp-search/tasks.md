## 1. OpenSpec And Backend Search Client

- [x] 1.1 Add failing backend tests for Exa MCP initialize, tools/list parsing, and `tools/call` response handling.
- [x] 1.2 Implement a minimal Exa MCP HTTP client under `server/service/` with request timeout, UTF-8 JSON-RPC payloads, and safe error handling.
- [x] 1.3 Implement search-result normalization that extracts a small set of titles, URLs, publish times, and highlights for grounding and visible sources.

## 2. Chat Grounding Integration

- [x] 2.1 Add failing backend tests proving `enable_search=true` executes Exa search before model invocation and injects grounding context into the model messages.
- [x] 2.2 Implement shared grounding flow for `Chat` and `ChatStream` so both paths reuse the same Exa search preparation logic.
- [x] 2.3 Append a compact Markdown source section to successful search-grounded assistant replies so saved conversations retain visible source links.
- [x] 2.4 Add backend tests covering “insufficient evidence must not guess” prompt guidance and Exa failure/error fallback behavior.

## 3. Frontend And Verification

- [ ] 3.1 Verify the existing AI chat UI renders appended source links correctly for live replies and reopened conversations.
- [ ] 3.2 Run `cd server && go test ./...` and fix any regressions introduced by the Exa grounding change.
- [x] 3.3 Run `cd web && npm run build` and `cd web && npm run typecheck`; if `vue-tsc` still hits the known toolchain issue, record it explicitly instead of treating it as a new regression.
