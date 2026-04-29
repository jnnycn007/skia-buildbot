# Skia Task Failure Analysis

This document contains instructions for agents to perform efficient and accurate
causal analysis of Skia task failures.

## Workflow

Important:

- Do **not** rely on task creation timestamps to establish chronological order.
  Tasks may be backfilled or retried and therefore the only reliable source is
  the Git commit history.
- Do **not** attempt to manually correlate large JSON dumps of raw database
  objects.

To perform this analysis, employ a hierarchical approach:

### Phase 1: Find Failing Tasks

As the coordinator agent, your job is to identify _which_ tasks are broken.

Retrieve the recent task results, eg.
`sk agent workflow get_task_health_report --limit=50 --revision=main --repository=https://skia.googlesource.com/skia.git`

**Result Format:** This returns a data set containing the following:

- **`commit_graph`:** The authoritative, chronologically ordered list of commits
  (index 0 is newest).
- **`tasks`:** A map of task name to recent runs of that task, each of which
  contains:
  - **`rev`:** The commit hash at which the task ran.
  - **`status`:** The result (`SUCCESS`, `FAILURE`, `MISHAP`).
  - **`blame_size`:** Critical for pinpointing culprits. If a `FAILURE` occurs
    at `rev="hash10"` with `blame_size=5`, the true culprit lies somewhere
    between `hash10` and the 4 commits prior to it. The specialist sub-agent
    must investigate all 5 commits in that window.

### Phase 2: Find Culprit Commit(s)

For each failing task identified in Phase 1, you MUST delegate the causal
analysis to a sub-agent (the "Task Specialist").

1.  Invoke a sub-agent (e.g., `invoke_agent(agent_name="generalist", ...)`).
2.  **Prompt the sub-agent with:**
    - The specific `task name` it is responsible for.
    - The `commit_graph` array from the health report.
    - The specific result array for that task from the health report.
    - **Instruction:** "Analyze the results for this task. Note the `blame_size`
      for the failing result(s) and use the `commit_graph` to find the range of
      potential culprit commits. Determine if this is a hard regression or a
      flaky infrastructure issue. Return a definitive report classifying the
      task as either flaky or persistent failure and identifying the culprit
      range."

### Phase 3: Synthesis

Once all Task Specialists have returned their findings, aggregate their reports.

- If multiple specialists point to the same culprit commit or overlapping commit
  range (e.g., "Commit X broke both Task A and Task B"), the failures may be
  related.
- If two or more tasks started failing within overlapping commit ranges and
  those tasks look similar, you may be able to use the intersection of those
  commit ranges to futher narrow down the culprit or isolate it altogether.
- Present the final holistic analysis to the user.
