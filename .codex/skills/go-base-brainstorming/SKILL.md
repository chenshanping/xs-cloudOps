---
name: go-base-brainstorming
description: Design exploration for this go-base repository before implementation. Use before creative work such as features, components, workflow changes, or behavior changes.
---

# Brainstorming Ideas Into Designs

Help turn ideas into fully formed designs through collaborative dialogue.

**HARD GATE:** Do NOT write any code or take implementation actions until you have presented a design and the user has approved it.

## Steps

1. **Explore project context** — check relevant files, docs, recent changes.

2. **Ask clarifying questions** — one at a time, understand purpose/constraints/success criteria. Prefer multiple choice when possible.

3. **Propose 2-3 approaches** — with trade-offs and your recommendation. Lead with the recommended option.

4. **Present design** — scale each section to its complexity. Ask after each section whether it looks right. Cover: architecture, components, data flow, error handling, testing.

5. **Write design doc** — save to `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md`.

6. **Spec self-review:**
   - Placeholder scan: any TBD, TODO, incomplete sections? Fix them.
   - Internal consistency: do sections contradict each other?
   - Scope check: focused enough for a single implementation plan?
   - Ambiguity check: could any requirement be interpreted two ways?

7. **User reviews written spec** — ask the user to review before proceeding.

8. **Transition** — invoke `/writing-plans` to create the implementation plan.

## Key Principles

- **One question at a time** — don't overwhelm
- **YAGNI ruthlessly** — remove unnecessary features
- **Explore alternatives** — always propose 2-3 approaches
- **Incremental validation** — present design, get approval, move on
- **Existing codebase** — explore current structure first, follow existing patterns
- **Design for isolation** — smaller units with clear boundaries and well-defined interfaces

## Anti-Pattern

"This is too simple to need a design" — every project goes through this process. The design can be short for simple projects, but you MUST present it and get approval.
