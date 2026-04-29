# Skia Infra Agents Guide

This document serves as the top-level router for agents operating within Skia
repositories. It contains brief descriptions of workflows and is designed to
keep your primary context window lean.

**CRITICAL:** Do NOT perform complex infrastructure analysis or workflow
orchestration using your default behaviors. You MUST consult any relevant
workflows listed below before attempting these tasks.

## Workflows

Depending on your task, read details for the corresponding workflow, either by
running `sk agent workflow <name>` or by finding it in the `docs/` directory:

### task_failure_analysis

**File:** `docs/task_failure_analysis.md`
**Trigger:** Use this workflow when asked to:

- Summarize recent task failures.
- Investigate "what commit broke task X?"
- Find the root cause of recent regressions.
- Determine if a failing task is a hard regression or flaky infrastructure.

---

_Note: As more complex, agent-specific workflows are developed for this
repository, they should be added as separate markdown files in this directory
and linked here in this document._
