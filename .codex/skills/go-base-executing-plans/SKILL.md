---
name: go-base-executing-plans
description: Execute an approved implementation plan task by task with verification checkpoints in this go-base repository.
---

# Executing Plans

Load plan, review critically, execute all tasks, report when complete.

## Steps

1. **Load and review plan** — read the plan file from `docs/superpowers/plans/`. If concerns exist, raise them before starting.

2. **Create todo list** — track all tasks with the todo_list tool.

3. **Execute tasks** — for each task:
   - Mark as in_progress
   - Follow each step exactly (plan has bite-sized steps)
   - Run verifications as specified
   - Mark as completed

4. **When blocked** — STOP immediately when:
   - Hit a blocker (missing dependency, test fails, instruction unclear)
   - Plan has critical gaps
   - You don't understand an instruction
   - Verification fails repeatedly
   - Ask for clarification rather than guessing.

5. **On completion** — report:
   - Tasks completed
   - Overall progress
   - Any issues encountered
   - Use `/finishing-branch` to finalize

## Rules

- Review plan critically first
- Follow plan steps exactly
- Don't skip verifications
- Stop when blocked, don't guess
- Never start on main/master without explicit user consent
