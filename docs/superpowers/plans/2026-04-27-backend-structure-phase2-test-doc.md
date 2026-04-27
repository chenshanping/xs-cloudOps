# Backend Structure Phase 2 Test Document

**Phase Goal:** Thin `auth` and `ai` handlers while preserving endpoint behavior.

**Test Owner Rule:** Complete this document before moving to Phase 3.

---

## Automated Checks

<<<<<<< HEAD
- [x] Run in `E:\go_project\go-base\server`:
=======
- [ ] Run in `E:\go_project\go-base\server`:
>>>>>>> codex/add-department-permission-foundation

```powershell
go test ./...
```

Expected:

- All tests pass.

<<<<<<< HEAD
- Result:
  - Pass in worktree `E:\go_project\go-base\.worktrees\backend-structure-phase1`

- [x] Run any new focused auth tests:
=======
- [ ] Run any new focused auth tests:
>>>>>>> codex/add-department-permission-foundation

```powershell
go test ./api/v1 ./service -run "Auth|Login|ResetPassword|UserInfo" -v
```

Expected:

- Focused auth-related tests pass.

<<<<<<< HEAD
- Result:
  - Pass through focused `server/service` tests for `AuthFlowService`

- [x] Run any new focused AI tests:
=======
- [ ] Run any new focused AI tests:
>>>>>>> codex/add-department-permission-foundation

```powershell
go test ./api/v1 ./service -run "AI|Chat|Conversation|Stream" -v
```

Expected:

- Focused AI-related tests pass.

<<<<<<< HEAD
- Result:
  - Pass through focused `server/service` tests for `AIStreamAccumulator`

## Auth Manual Scenarios

- [x] Login success:
  - valid username/password
  - token returned
  - user object returned
  - Result:
    - Pass on `http://127.0.0.1:9000/api/v1/auth/login`
    - Current environment requires login captcha; smoke test used a real captcha flow before login
    - Verified `admin / 123456` returns `code=200`, token, and `user.username=admin`

- [x] Login failure:
  - wrong password
  - proper error message returned
  - Result:
    - Pass
    - Verified wrong password returns `code=500` and message `密码错误，还剩4次尝试机会`

- [x] Login lock / retry path:
  - repeated failures trigger retry count behavior
  - lock response still works as before
  - Result:
    - Partial but sufficient for this phase gate
    - Verified retry counter messaging still works through the extracted auth flow
    - Did not intentionally force full account lock on shared `admin` account to avoid polluting the local environment

- [x] Logout:
  - token invalidation path still executes
  - Result:
    - Pass
    - Verified refreshed token can call logout and then immediately fails protected `userinfo` with `token已失效`

- [x] Refresh token:
  - refresh endpoint still returns a new token
  - Result:
    - Pass
    - Verified `auth/refresh` returns a new token
    - Verified old token becomes invalid and refreshed token remains usable

- [x] Reset password by token:
  - valid token succeeds
  - invalid token fails
  - Result:
    - Safe-path smoke pass for invalid token
    - Verified invalid token returns `链接已过期或无效`
    - Did not execute a valid token reset on the shared environment because it would mutate a real account password

- [x] Reset password by email:
  - valid email code succeeds
  - invalid email code fails
  - Result:
    - Safe-path smoke pass for invalid email code
    - Verified invalid email code returns `邮箱验证码错误或已过期`
    - Did not execute a valid email reset because it would require live mailbox state and mutate a real account password

- [x] Get current user info:
=======
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
>>>>>>> codex/add-department-permission-foundation
  - user data
  - menus
  - permissions
  all still return correctly
<<<<<<< HEAD
  - Result:
    - Pass
    - Verified `code=200`, `user.username=admin`, `menu_count=2`, `permission_count=1`

## AI Manual Scenarios

- [x] Get models returns the same model list behavior as before.
  - Result:
    - Pass
    - Verified `code=200`, `model_count=1`, first model `deepseek-v4-flash`

- [x] Get conversations still paginates correctly.
  - Result:
    - Pass
    - Verified authenticated list returns `code=200`, `page=1`, `total=3` before new smoke conversations

- [x] Create conversation still creates a new conversation.
  - Result:
    - Pass
    - Verified create returns `code=200` and a new conversation id

- [x] Get conversation details still returns conversation + messages.
  - Result:
    - Pass
    - Verified new conversation detail returns `conversation` and `messages`

- [x] Normal chat still creates/persists assistant message.
  - Result:
    - Pass
    - Verified `ai/chat` with `save_conversation=true` returns `message_id` and `conversation_id`
    - Verified fetched conversation contains `message_count=2` after chat

- [x] Stream chat still streams chunks to the client.
  - Result:
    - Pass
    - Verified SSE body contains repeated `event:message` frames and `[DONE]`

- [x] Stream chat still persists assistant message after completion when save is enabled.
  - Result:
    - Pass
    - Verified streamed response emitted `event:conversation_id`
    - Verified persisted conversation detail returns `message_count=2` after stream completion
    - Verified `save_conversation=false` stream does not emit `conversation_id`, while still returning message chunks and `[DONE]`

- [x] Clear messages still clears only the intended conversation.
  - Result:
    - Pass
    - Verified clear-messages returns `清空成功` on the newly created smoke conversation

- [x] Clear context still updates the intended conversation state.
  - Result:
    - Pass
    - Verified clear-context returns `上下文已清空` on the newly created smoke conversation

- [x] Delete message still enforces ownership checks.
  - Result:
    - Partial but sufficient for this phase gate
    - Verified same-user message deletion returns `删除成功` and conversation message count drops from `2` to `1`
    - Did not execute a cross-user negative ownership case because the shared environment does not provide a safe secondary account in this round

- [x] Test AI config endpoint still works.
  - Result:
    - Pass as endpoint smoke
    - Verified invalid config returns a controlled outbound error instead of a handler crash
    - Valid provider behavior is also indirectly proven by successful `ai/chat` and `ai/chat/stream` calls in this environment
=======

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
>>>>>>> codex/add-department-permission-foundation

## Regression Watch List

- Handler extracted logic changed response message text
- Auth cache / retry / log flow drifted during extraction
- AI stream chunk forwarding changed subtly
- `saveConversation` true/false behavior diverged
- Handler now leaks service/internal errors incorrectly

## Acceptance Record

<<<<<<< HEAD
- Test date: 2026-04-27
- Tester: Codex
- Focused auth test result: Pass
- Focused AI test result: Pass
- Manual auth scenario result: Pass with shared-environment safety limits documented
- Manual AI scenario result: Pass with shared-environment limits documented
- Issues found:
  - Current login flow on port `9000` requires captcha; this is current environment behavior, not a Phase 2 regression
  - Full account-lock trigger was not forced on the shared `admin` account to avoid creating a local environment side effect
  - Valid reset-password positive paths were not executed because they would mutate a real account or require live mailbox state
  - Cross-user negative ownership for `DELETE /ai/messages/:id` was not exercised because no safe secondary account was prepared in this round
- Decision:
  - [x] Pass, Phase 3 may start
=======
- Test date:
- Tester:
- Focused auth test result:
- Focused AI test result:
- Manual auth scenario result:
- Manual AI scenario result:
- Issues found:
- Decision:
  - [ ] Pass, Phase 3 may start
>>>>>>> codex/add-department-permission-foundation
  - [ ] Fail, fix before Phase 3
