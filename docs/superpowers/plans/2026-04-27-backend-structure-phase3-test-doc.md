# Backend Structure Phase 3 Test Document

**Phase Goal:** Split oversized services and close remaining boundary leaks without breaking route contracts or business behavior.

**Test Owner Rule:** Complete this document before declaring the 3-phase backend cleanup finished.

---

## Automated Checks

- [x] Run in `E:\go_project\go-base\server`:

```powershell
go test ./...
```

Expected:

- All tests pass.

- Result:
  - Pass in worktree `E:\go_project\go-base\.worktrees\backend-structure-phase1\server`

- [x] Run focused AI service tests:

```powershell
go test ./service -run "AI|Conversation|Chat|Stream" -v
```

Expected:

- All AI-related service tests pass after file split.

- Result:
  - Pass
  - Verified AI stream accumulator tests still pass
  - Added and passed file-chain tests for:
    - local text file read
    - remote file read
    - local image to base64 conversion

- [x] Run focused user service tests if `user.go` was split:

```powershell
go test ./service -run "User|Profile|Password|Status" -v
```

Expected:

- All touched user-related tests pass.

- Result:
  - Not needed in this phase
  - `server/service/user.go` was intentionally left untouched to keep Phase 3 scope controlled

## AI Service Split Verification

- [x] Confirm `service/ai.go` is no longer the single owner of:
  - conversation CRUD
  - context building
  - file preprocessing
  - provider HTTP client
  - stream parsing
  - Result:
    - `service/ai.go` now only keeps `AIService` declaration and global singleton
    - Conversation orchestration moved to `service/ai_conversation.go`
    - Context/message assembly moved to `service/ai_context.go`
    - Local/remote file handling moved to `service/ai_files.go`
    - Provider HTTP client and config test path moved to `service/ai_client.go`
    - Shared AI types and service-local inputs moved to `service/ai_types.go`

- [x] Confirm exported AI service entry points still behave identically from the handler perspective.
  - Result:
    - Pass through compile verification and `go test ./...`
    - `api/v1/ai.go` still calls the same exported entry points on `service.AI`
    - Existing Phase 2 auth/AI smoke results remain the behavioral baseline; Phase 3 changed file ownership and DTO boundaries, not route contracts

- [x] Confirm SSE parsing has a single clear owner after the split.
  - Result:
    - Pass
    - Removed the unused `ParseSSEStream` copy from the old monolithic `service/ai.go`
    - `service/ai_stream_accumulator.go` remains the single SSE line/chunk accumulator owner used by the handler

- [x] Confirm local file read, remote file read, and image base64 conversion still work.
  - Result:
    - Pass
    - Verified by new focused tests in `server/service/ai_files_test.go`

## Request DTO Boundary Verification

- [x] For touched `auth`, `ai`, and `user` paths, confirm handlers map HTTP DTOs to service-local inputs instead of passing `request.*` through unchanged where refactoring was planned.
  - Result:
    - Pass for touched Phase 3 paths
    - `api/v1/auth.go` now maps `request.LoginRequest` into `service.AuthLoginInput`
    - `api/v1/ai.go` now maps request DTOs into:
      - `service.ConversationListInput`
      - `service.CreateConversationInput`
      - `service.AIChatInput`
    - `service/user.go` was not touched in this phase, so no user DTO migration was performed

- [x] Confirm untouched modules were not unnecessarily migrated.
  - Result:
    - Pass
    - Phase 3 stayed inside `auth`, `ai`, `router/registry`, and `swagger`

## Metadata Verification

- [x] Confirm route metadata parsing has one shared rule set.
  - Result:
    - Pass
    - Added `router/registry/fields.go`
    - `registry.ParseStructFields` and `swagger.Generator` now both consume the same shared field parsing rules
    - Verified by new tests in:
      - `server/router/registry/fields_test.go`
      - `server/swagger/swagger_test.go`

- [x] Confirm Swagger generation still works.
  - Result:
    - Pass
    - `go test ./swagger -v` passes
    - Live smoke on current `9000` listener returns `HTTP 200`, title `Go Base Server API`, basePath `/api/v1`

- [x] Confirm API sync still extracts request/response metadata correctly.
  - Result:
    - Pass by code-path verification
    - `api/v1/api.go` still uses `registry.ParseStructFields`
    - Shared field parser tests cover embedded fields, `form/json` tag fallback, and `binding:\"required\"`
    - No dedicated live sync trigger was executed in this round to avoid mutating shared API metadata

Suggested checks:

```powershell
Invoke-WebRequest http://127.0.0.1:8888/swagger/doc.json -UseBasicParsing
```

Expected:

- HTTP 200
- JSON returned

- [x] If available in the local environment, trigger API sync and verify it completes without parsing errors.
  - Result:
    - Skipped intentionally in this round
    - Current environment appears to be a shared local admin instance; did not mutate API metadata just to prove parser wiring
    - Parser behavior is covered by new registry/swagger tests and compile verification

## End-to-End Regression Watch List

- AI conversation title auto-update broken after file split
- Assistant message not persisted after stream completion
- File attachment content preprocessing broken
- Request DTO decoupling changed validation or defaulting behavior
- Swagger and API sync diverged due to duplicated reflection logic

## Acceptance Record

- Test date: 2026-04-27
- Tester: Codex
- `go test ./...` result: Pass
- AI split verification result: Pass
- DTO boundary verification result: Pass
- Swagger / metadata verification result: Pass
- Issues found:
  - The existing listener on port `9000` is a pre-running GoLand-built process. It was not force-restarted in this round, so live `9000` smoke confirms interface availability, not that the running binary was replaced by this exact worktree build.
  - `service/user.go` was intentionally not split in this phase to keep scope aligned with the minimum safe Phase 3 slice.
  - Live API sync was not triggered to avoid mutating shared API metadata in the local environment.
- Decision:
  - [x] Pass, 3-phase cleanup complete
  - [ ] Fail, fix before closing work
