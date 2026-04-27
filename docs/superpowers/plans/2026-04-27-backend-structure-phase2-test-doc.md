# Backend Structure Phase 2 Test Document

**Phase Goal:** Thin `auth` and `ai` handlers while preserving endpoint behavior.

**Test Owner Rule:** Complete this document before moving to Phase 3.

---

## Automated Checks

- [ ] Run in `E:\go_project\go-base\server`:

```powershell
go test ./...
```

Expected:

- All tests pass.

- [ ] Run any new focused auth tests:

```powershell
go test ./api/v1 ./service -run "Auth|Login|ResetPassword|UserInfo" -v
```

Expected:

- Focused auth-related tests pass.

- [ ] Run any new focused AI tests:

```powershell
go test ./api/v1 ./service -run "AI|Chat|Conversation|Stream" -v
```

Expected:

- Focused AI-related tests pass.

## Auth Manual Scenarios

- [ ] Login success:
  - valid username/password
  - token returned
  - user object returned

- [ ] Login failure:
  - wrong password
  - proper error message returned

- [ ] Login lock / retry path:
  - repeated failures trigger retry count behavior
  - lock response still works as before

- [ ] Logout:
  - token invalidation path still executes

- [ ] Refresh token:
  - refresh endpoint still returns a new token

- [ ] Reset password by token:
  - valid token succeeds
  - invalid token fails

- [ ] Reset password by email:
  - valid email code succeeds
  - invalid email code fails

- [ ] Get current user info:
  - user data
  - menus
  - permissions
  all still return correctly

## AI Manual Scenarios

- [ ] Get models returns the same model list behavior as before.
- [ ] Get conversations still paginates correctly.
- [ ] Create conversation still creates a new conversation.
- [ ] Get conversation details still returns conversation + messages.
- [ ] Normal chat still creates/persists assistant message.
- [ ] Stream chat still streams chunks to the client.
- [ ] Stream chat still persists assistant message after completion when save is enabled.
- [ ] Clear messages still clears only the intended conversation.
- [ ] Clear context still updates the intended conversation state.
- [ ] Delete message still enforces ownership checks.
- [ ] Test AI config endpoint still works.

## Regression Watch List

- Handler extracted logic changed response message text
- Auth cache / retry / log flow drifted during extraction
- AI stream chunk forwarding changed subtly
- `saveConversation` true/false behavior diverged
- Handler now leaks service/internal errors incorrectly

## Acceptance Record

- Test date:
- Tester:
- Focused auth test result:
- Focused AI test result:
- Manual auth scenario result:
- Manual AI scenario result:
- Issues found:
- Decision:
  - [ ] Pass, Phase 3 may start
  - [ ] Fail, fix before Phase 3
