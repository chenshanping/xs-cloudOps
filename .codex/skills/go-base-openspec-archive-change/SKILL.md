---
name: go-base-openspec-archive-change
description: Archive a completed OpenSpec change in this go-base repository after implementation and artifact checks are done.
---

# OpenSpec Archive Change

Archive a completed change after implementation.

## Steps

1. **Select change** — if not provided, run `openspec list --json` and ask user to select.

2. **Check artifact completion:**
   ```bash
   openspec status --change "<name>" --json
   ```
   If incomplete artifacts: warn and confirm before proceeding.

3. **Check task completion** — read tasks.md, count incomplete `- [ ]` vs complete `- [x]`. Warn if incomplete.

4. **Check delta specs** — if `openspec/changes/<name>/specs/` has delta specs, compare with main specs, show summary, offer sync.

5. **Archive:**
   ```bash
   mkdir -p openspec/changes/archive
   mv openspec/changes/<name> openspec/changes/archive/YYYY-MM-DD-<name>
   ```

6. **Summary** — change name, schema, archive location, spec sync status, any warnings.

## Guardrails

- Always prompt for selection if not provided
- Don't block on warnings, just inform and confirm
- Preserve .openspec.yaml when archiving
- Get confirmation before discarding incomplete work
