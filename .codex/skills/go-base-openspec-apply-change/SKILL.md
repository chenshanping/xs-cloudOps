---
name: go-base-openspec-apply-change
description: Implement tasks from an approved OpenSpec change in this go-base repository. Use when ready to start or continue implementation of a proposed change.
---

# OpenSpec Apply Change

Implement tasks from an OpenSpec change proposal.

## Steps

1. **Select change** — if name provided, use it. Otherwise auto-select if only one active change. If ambiguous, list and ask.

2. **Check status:**
   ```bash
   openspec status --change "<name>" --json
   ```

3. **Get apply instructions:**
   ```bash
   openspec instructions apply --change "<name>" --json
   ```
   - If `state: "blocked"`: show message, suggest using openspec-propose to create missing artifacts
   - If `state: "all_done"`: congratulate, suggest `/openspec-archive-change`

4. **Read all context files** from the apply instructions output.

5. **Show progress** — schema, N/M tasks complete, remaining tasks, dynamic instruction.

6. **Implement tasks (loop):**
   - Show which task is being worked on
   - Make code changes
   - Keep changes minimal and focused
   - Mark task complete: `- [ ]` → `- [x]`
   - Continue to next task

7. **Pause if:**
   - Task is unclear → ask
   - Design issue found → suggest updating artifacts
   - Error/blocker → report and wait

8. **On completion** — show tasks completed, overall progress. If all done, suggest `/openspec-archive-change`.

## Guardrails

- Always read context files before starting
- Keep code changes scoped to each task
- Update task checkbox immediately after completing
- Pause on errors, don't guess
