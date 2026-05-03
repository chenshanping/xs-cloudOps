---
description: Create a detailed implementation plan from a spec or requirements, before touching code.
---

# Writing Plans

Write comprehensive implementation plans assuming the engineer has zero codebase context. Document everything: which files to touch, code, testing, docs to check, how to verify. Give the whole plan as bite-sized tasks. DRY. YAGNI. TDD. Frequent commits.

## Steps

1. **Read the spec** — load the design doc from `docs/superpowers/specs/`.

2. **Scope check** — if the spec covers multiple independent subsystems, suggest breaking into separate plans.

3. **Map file structure** — list which files will be created or modified and what each is responsible for. Lock decomposition decisions here.

4. **Write tasks** — each task contains bite-sized steps (2-5 minutes each):
   - Write the failing test
   - Run it to confirm it fails
   - Implement minimal code to make it pass
   - Run tests to confirm they pass
   - Commit

5. **Plan header** — every plan starts with: Goal (one sentence), Architecture (2-3 sentences), Tech Stack.

6. **Save plan** — to `docs/superpowers/plans/YYYY-MM-DD-<feature-name>.md`.

7. **Self-review:**
   - Spec coverage: can you point to a task for each requirement?
   - Placeholder scan: no TBD, TODO, "implement later", "similar to Task N"
   - Type consistency: do names and signatures match across tasks?

8. **Offer execution** — ask user whether to implement inline or keep for later.

## Task Structure

Each task must include:
- **Files:** exact paths to create/modify/test
- **Steps:** with actual code blocks (no "add appropriate error handling")
- **Commands:** exact commands with expected output
- **Commit:** git add + commit message

## No Placeholders

Never write: "TBD", "TODO", "implement later", "fill in details", "similar to Task N", "add appropriate error handling". Every step must contain the actual content needed.
