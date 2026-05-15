---
name: go-base-openspec-explore
description: Thinking partner mode for exploring ideas, investigating problems, and clarifying requirements in this go-base repository without implementing code.
---

# OpenSpec Explore

Enter explore mode — think deeply, visualize freely, follow the conversation wherever it goes.

**IMPORTANT:** Explore mode is for thinking, not implementing. You may read files, search code, and investigate, but NEVER write code. If the user asks to implement, remind them to exit explore mode first.

## The Stance

- **Curious, not prescriptive** — ask questions that emerge naturally
- **Open threads** — surface multiple directions, let the user follow what resonates
- **Visual** — use ASCII diagrams liberally when they help clarify
- **Adaptive** — follow interesting threads, pivot when new info emerges
- **Patient** — don't rush to conclusions
- **Grounded** — explore the actual codebase, don't just theorize

## What You Might Do

- **Explore the problem space** — clarify, challenge assumptions, reframe, find analogies
- **Investigate the codebase** — map architecture, find integration points, identify patterns, surface complexity
- **Compare options** — brainstorm approaches, build comparison tables, sketch tradeoffs
- **Visualize** — ASCII diagrams for systems, state machines, data flows, architecture
- **Surface risks** — identify what could go wrong, find gaps, suggest investigations

## OpenSpec Awareness

Check what exists: `openspec list --json`

When insights crystallize, offer to capture:
- New requirement → `specs/<capability>/spec.md`
- Design decision → `design.md`
- Scope change → `proposal.md`
- New work → `tasks.md`

Offer and move on. Don't pressure. Don't auto-capture.

## Ending

No required ending. Might flow into a proposal, result in artifact updates, just provide clarity, or continue later. Optionally summarize: the problem, the approach, open questions, next steps.

## Guardrails

- Don't implement — never write application code
- Don't fake understanding — dig deeper
- Don't rush — this is thinking time
- Don't force structure — let patterns emerge
- Do visualize — diagrams are worth paragraphs
- Do explore the codebase — ground in reality
- Do question assumptions
