---
name: go-base-systematic-debugging
description: Systematic debugging methodology for this go-base repository. Use for any bug, test failure, or unexpected behavior before proposing fixes.
---

# Systematic Debugging

**Core principle:** ALWAYS find root cause before attempting fixes. NO FIXES WITHOUT ROOT CAUSE INVESTIGATION FIRST.

## Phase 1: Root Cause Investigation

BEFORE attempting ANY fix:

1. **Read error messages carefully** — don't skip errors or warnings. Read stack traces completely. Note line numbers, file paths, error codes.

2. **Reproduce consistently** — can you trigger it reliably? What are exact steps? If not reproducible, gather more data, don't guess.

3. **Check recent changes** — git diff, recent commits, new dependencies, config changes, environmental differences.

4. **Gather evidence in multi-component systems** — for EACH component boundary: log what enters, what exits, verify config propagation. Run once to find WHERE it breaks.

5. **Trace data flow** — where does bad value originate? What called this with bad value? Trace up until you find the source. Fix at source, not symptom.

## Phase 2: Pattern Analysis

1. **Find working examples** — locate similar working code in same codebase
2. **Compare against references** — read reference implementations COMPLETELY
3. **Identify differences** — list every difference between working and broken
4. **Understand dependencies** — what components, settings, config, environment does this need?

## Phase 3: Hypothesis and Testing

1. **Form single hypothesis** — "I think X is the root cause because Y"
2. **Test minimally** — smallest possible change, one variable at a time
3. **Verify** — worked? → Phase 4. Didn't? → new hypothesis. DON'T stack fixes.

## Phase 4: Implementation

1. **Create failing test case** — simplest reproduction, automated if possible
2. **Implement single fix** — address root cause, ONE change, no "while I'm here" improvements
3. **Verify fix** — test passes? No other tests broken?
4. **If 3+ fixes failed** — STOP. Question the architecture. Discuss with user before more attempts.

## Red Flags — STOP and Return to Phase 1

- "Quick fix for now, investigate later"
- "Just try changing X"
- "Add multiple changes, run tests"
- "It's probably X, let me fix that"
- Proposing solutions before tracing data flow
- Each fix reveals new problem in different place
