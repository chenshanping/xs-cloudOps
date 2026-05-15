---
name: go-base-file-reference-guardrails
description: Guardrails for go-base modules that upload, bind, display, or delete business files through sys_file and sys_file_reference. Use when a change touches file IDs, file previews, config images, or reference cleanup.
---

# File Reference Guardrails

Keep uploaded-file ownership consistent across backend, frontend, migration, and runtime reads.

## Checklist

1. **Model / request design**
   - Business tables store only `xxx_file_id` / `xxx_file_ids`.
   - Do not add uploaded URL columns as source-of-truth business data.
   - Response structs may expose transient `xxx_file_url` fields with `gorm:"-"`.

2. **Service write path**
   - Validate referenced `sys_file` rows before binding file IDs.
   - In create / update flows, call `filesvc.Reference.ReplaceRefs(tx, refTable, refID, refs)` in the same transaction as the business write.
   - For multi-file JSON arrays, normalize and deduplicate file IDs before syncing refs.

3. **Delete / clear cleanup**
   - In delete, clear, batch delete, conversation cleanup, or destructive config flows, call `filesvc.Reference.ClearRefs(...)` in the same transaction as the business deletion.
   - Do not rely on `server/service/file/file.go` to scan business tables for reference checks.

4. **Read-time URL resolution**
   - Resolve display URLs from `sys_file.url` when reading data.
   - Batch query file IDs where possible; do not `Preload` FK associations only to expose uploaded URLs.
   - Config display keys such as `sys_logo` should be derived from `*_file_id` keys at read time.

5. **Config image handling**
   - Source-of-truth config keys are `*_file_id`.
   - Keep anonymous/public exposure policy in backend code.
   - Do not persist derived URL keys as side effects during migrations or saves.

6. **Backfill / SQL**
   - If persisted DB state changes, add an idempotent MySQL-safe script under `server/sql/`.
   - Backfill `sys_file_reference` rows from existing business tables or config rows in an idempotent way.
   - Prefer module-scoped rerunnable backfill logic over one-shot data assumptions.

7. **Frontend upload integration**
   - Keep preview URLs in component-local state.
   - Submit only file IDs to backend APIs.
   - If UI exposes an uploaded image/file preview, also complete the preview/open loop in the same change.
