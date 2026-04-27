# Backend Structure Phase 1 Test Document

**Phase Goal:** Structural alignment only. No intended behavior change.

**Test Owner Rule:** Complete this document before moving to Phase 2.

---

## Automated Checks

- [x] Run in `E:\go_project\go-base\server`:

```powershell
go test ./...
```

Expected:

- All packages pass as they did before Phase 1.

- Result:
  - Pass in worktree `E:\go_project\go-base\.worktrees\backend-structure-phase1`

- [x] Run:

```powershell
go test ./router/... ./model/request/...
```

Expected:

- No package-level compile or import errors after file moves and renames.

- Result:
  - Pass in worktree `E:\go_project\go-base\.worktrees\backend-structure-phase1`

## Manual Structure Checks

- [x] Confirm `R(...)` no longer lives in `server/router/modules/auth.go`.
- [x] Confirm route module files contain only feature route registration logic.
- [x] Confirm request DTOs are split by feature but still use `package request`.
- [x] Confirm no request type name changed externally.
- [x] Confirm no new top-level backend architecture layer was introduced.

## Functional Smoke Checks

- [x] Start the backend locally.
- [x] Confirm server boot succeeds with no import/runtime panic caused by moved files.
- [x] Open Swagger JSON endpoint and confirm it still returns data:

Suggested check:

```powershell
Invoke-WebRequest http://127.0.0.1:8888/swagger/doc.json -UseBasicParsing
```

Expected:

- HTTP 200
- Response body is JSON

- Current result:
  - `http://127.0.0.1:9000/swagger/doc.json` returns HTTP 200
  - Response contains Swagger JSON with title `Go Base Server API`

- [x] Confirm auth route registration still exists.
- [x] Confirm AI route registration still exists.
- [x] Confirm API sync route still exists.

## Regression Watch List

- Missing request types due to wrong file split
- Route helper moved incorrectly and some routes not registered
- Swagger route metadata broken by naming cleanup
- Imports still referencing deleted file-local helpers

## Acceptance Record

- Test date: 2026-04-27
- Tester: Codex
- `go test ./...` result: Pass
- Swagger endpoint result: Pass on `http://127.0.0.1:9000/swagger/doc.json`
- Route registration check result: Pass (`/auth/login`, `/ai/models`, `/apis/sync` confirmed in route registration output)
- Issues found:
  - None for Phase 1 after switching smoke verification to port `9000`
- Decision:
  - [x] Pass, Phase 2 may start
  - [ ] Fail, fix before Phase 2
