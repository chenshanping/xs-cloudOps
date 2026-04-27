# Backend Structure Phase 3 Test Document

**Phase Goal:** Split oversized services and close remaining boundary leaks without breaking route contracts or business behavior.

**Test Owner Rule:** Complete this document before declaring the 3-phase backend cleanup finished.

---

## Automated Checks

- [ ] Run in `E:\go_project\go-base\server`:

```powershell
go test ./...
```

Expected:

- All tests pass.

- [ ] Run focused AI service tests:

```powershell
go test ./service -run "AI|Conversation|Chat|Stream" -v
```

Expected:

- All AI-related service tests pass after file split.

- [ ] Run focused user service tests if `user.go` was split:

```powershell
go test ./service -run "User|Profile|Password|Status" -v
```

Expected:

- All touched user-related tests pass.

## AI Service Split Verification

- [ ] Confirm `service/ai.go` is no longer the single owner of:
  - conversation CRUD
  - context building
  - file preprocessing
  - provider HTTP client
  - stream parsing

- [ ] Confirm exported AI service entry points still behave identically from the handler perspective.
- [ ] Confirm SSE parsing has a single clear owner after the split.
- [ ] Confirm local file read, remote file read, and image base64 conversion still work.

## Request DTO Boundary Verification

- [ ] For touched `auth`, `ai`, and `user` paths, confirm handlers map HTTP DTOs to service-local inputs instead of passing `request.*` through unchanged where refactoring was planned.
- [ ] Confirm untouched modules were not unnecessarily migrated.

## Metadata Verification

- [ ] Confirm route metadata parsing has one shared rule set.
- [ ] Confirm Swagger generation still works.
- [ ] Confirm API sync still extracts request/response metadata correctly.

Suggested checks:

```powershell
Invoke-WebRequest http://127.0.0.1:8888/swagger/doc.json -UseBasicParsing
```

Expected:

- HTTP 200
- JSON returned

- [ ] If available in the local environment, trigger API sync and verify it completes without parsing errors.

## End-to-End Regression Watch List

- AI conversation title auto-update broken after file split
- Assistant message not persisted after stream completion
- File attachment content preprocessing broken
- Request DTO decoupling changed validation or defaulting behavior
- Swagger and API sync diverged due to duplicated reflection logic

## Acceptance Record

- Test date:
- Tester:
- `go test ./...` result:
- AI split verification result:
- DTO boundary verification result:
- Swagger / metadata verification result:
- Issues found:
- Decision:
  - [ ] Pass, 3-phase cleanup complete
  - [ ] Fail, fix before closing work
