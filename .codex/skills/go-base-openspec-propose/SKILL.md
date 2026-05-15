---
name: go-base-openspec-propose
description: Create proposal, design, and tasks artifacts for a high-risk or cross-module OpenSpec change in this go-base repository.
---

# OpenSpec Propose

Create a change with all artifacts in one step: proposal.md (what & why), design.md (how), tasks.md (implementation steps).

## Steps

1. **Clarify intent** — if no clear input, ask what the user wants to build. Derive a kebab-case name.

2. **Create change directory:**
   ```bash
   openspec new change "<name>"
   ```

3. **Check status:**
   ```bash
   openspec status --change "<name>" --json
   ```

4. **Create artifacts in dependency order** — for each artifact:
   ```bash
   openspec instructions <artifact-id> --change "<name>" --json
   ```
   - Read completed dependency files for context
   - Use `template` as structure, apply `context` and `rules` as constraints (do NOT copy them into output)
   - Write the artifact file

5. **Continue until all `applyRequires` artifacts are complete.** Re-run status after each artifact.

6. **Show final status:**
   ```bash
   openspec status --change "<name>"
   ```

## Output

- Change name and location
- List of artifacts created with brief descriptions
- "Ready for implementation. Use `/openspec-apply-change` to start."

## Guardrails

- Create ALL artifacts needed for implementation
- Always read dependency artifacts before creating new ones
- If context is unclear, ask the user
- If change already exists, ask whether to continue or create new
- Verify each artifact file exists after writing
