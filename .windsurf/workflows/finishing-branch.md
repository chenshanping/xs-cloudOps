---
description: Finalize development work — verify tests, present merge/PR options, clean up.
---

# Finishing a Development Branch

Verify tests → Present options → Execute choice → Clean up.

## Steps

1. **Verify tests pass:**
   ```bash
   # Backend
   cd server && go test ./...
   # Frontend
   cd web && npm run build
   ```
   If tests fail, STOP. Fix before proceeding.

2. **Present options:**
   1. Merge back to base branch locally
   2. Push and create a Pull Request
   3. Keep the branch as-is (handle later)
   4. Discard this work

3. **Execute choice:**

   **Option 1 — Merge locally:**
   ```bash
   git checkout <base-branch>
   git pull
   git merge <feature-branch>
   # Verify tests again
   git branch -d <feature-branch>
   ```

   **Option 2 — Create PR:**
   ```bash
   git push -u origin <feature-branch>
   # Create PR with summary and test plan
   ```

   **Option 3 — Keep as-is:**
   Report branch name and location.

   **Option 4 — Discard:**
   Confirm first (require typed "discard").
   ```bash
   git checkout <base-branch>
   git branch -D <feature-branch>
   ```

## Rules

- Never proceed with failing tests
- Never merge without verifying tests on result
- Never delete work without confirmation
- Never force-push without explicit request
- Present exactly 4 options, don't add explanation
