# Backend Structure Phase 1 Test Document

**Phase Goal:** Structural alignment only. No intended behavior change.

**Test Owner Rule:** Complete this document before moving to Phase 2.

---

## Automated Checks

- [ ] Run in `E:\go_project\go-base\server`:

```powershell
go test ./...
```

Expected:

- All packages pass as they did before Phase 1.

- [ ] Run:

```powershell
go test ./router/... ./model/request/...
```

Expected:

- No package-level compile or import errors after file moves and renames.

## Manual Structure Checks

- [ ] Confirm `R(...)` no longer lives in `server/router/modules/auth.go`.
- [ ] Confirm route module files contain only feature route registration logic.
- [ ] Confirm request DTOs are split by feature but still use `package request`.
- [ ] Confirm no request type name changed externally.
- [ ] Confirm no new top-level backend architecture layer was introduced.

## Functional Smoke Checks

- [ ] Start the backend locally.
- [ ] Confirm server boot succeeds with no import/runtime panic caused by moved files.
- [ ] Open Swagger JSON endpoint and confirm it still returns data:

Suggested check:

```powershell
Invoke-WebRequest http://127.0.0.1:8888/swagger/doc.json -UseBasicParsing
```

Expected:

- HTTP 200
- Response body is JSON

- [ ] Confirm auth route registration still exists.
- [ ] Confirm AI route registration still exists.
- [ ] Confirm API sync route still exists.

## Regression Watch List

- Missing request types due to wrong file split
- Route helper moved incorrectly and some routes not registered
- Swagger route metadata broken by naming cleanup
- Imports still referencing deleted file-local helpers

## Acceptance Record

- Test date:
- Tester:
- `go test ./...` result:
- Swagger endpoint result:
- Route registration check result:
- Issues found:
- Decision:
  - [ ] Pass, Phase 2 may start
  - [ ] Fail, fix before Phase 2
