# Skia Task Failure Analysis

This document contains instructions for agents to perform efficient and accurate
causal analysis of Skia task failures.

## General Rules

- Do **not** rely on task creation timestamps to establish chronological order.
  Tasks may be backfilled or retried and therefore the only reliable source is
  the Git commit history.
- Do **not** attempt to manually correlate large JSON dumps of raw database
  objects.
- Do **not** try to use local Git history or search the local checkout unless
  explicitly asked to do so.

To perform this analysis, employ a hierarchical approach:

## Phase 1: Find Failing Tasks

As the coordinator agent, your job is to identify which tasks are broken by
which commits.

Retrieve the recent task results, eg.
`sk agent tool get_task_health_report --limit=35 --revision=main --repo=https://skia.googlesource.com/skia.git`

**Result Format:** This returns a data set containing the following:

- **`Commits`:** The authoritative, chronologically ordered list of commits
  (index 0 is newest).
- **`Task Results`:** A series of results for tasks whose result has changed
  within the given commit range. Each line contains:
  - **`commit`:** The commit hash at which the task ran.
  - **`status`:** The result (`SUCCESS`, `FAILURE`, `MISHAP`).
  - **`id`:** The database ID of the task.

## Phase 2: Find Culprit Commit(s)

To find the culprit commit(s), you MUST delegate the causal analysis for each
task from Phase 1 to a separate sub-agent (the "Task Specialist").

1.  Invoke a sub-agent (e.g., `invoke_agent(agent_name="generalist", ...)`).
2.  **Prompt the sub-agent with:**
    - The specific section for a single task from the health report.
    - **Instruction:** "Determine whether the failures of this task are _flaky_
      or _persistent_. If _persistent_, find the culprit commit or range of
      potential culprit commits. Return a report with your classification, any
      culprit commits, and a rationale for your decision. You MUST base your
      decision on the data presented to you. Do NOT attempt to collect any more
      data by reading files, running commands or MCP server tools, etc."

## Phase 3: Aggregate

Once all Task Specialists have returned their findings, aggregate their reports.

- If multiple specialists point to the same culprit commit or overlapping commit
  range (e.g., "Commit X broke both Task A and Task B"), the failures may be
  related.
- If two or more tasks started failing within overlapping commit ranges and
  those tasks look similar, you may be able to use the intersection of those
  commit ranges to futher narrow down the culprit or isolate it altogether.
- Use the commit subjects found in Phase 1 to make an educated guess about which
  suspect commit is most likely to be the culprit.

## Phase 4: Refinement and Further Analysis

If the results from Phase 3 are definitive, you might be able to stop there.
However, depending on what you were originally asked to do, you may need to
investigate further. Follow the instructions below.

### Persistent Failures

If you found a persistently-failing task, your first priority is to single out
the commit which caused it. If a single culprit has not already been found by
this point, start by retrieving the commit messages for the suspect commits. The
`gerrit_get_commit_message` tool from the pnd MCP server will be your best bet.
If that's not available, you can try using `git log` locally but you may not be
inside of a checkout of the correct repository. If all else fails, try using the
Gitiles HTTP API.

If the culprit is not obvious by correlating the name(s) of the failing
task(s) with the commit message, invoke a sub-agent to run
`sk agent workflow task_drilldown` and follow the instructions it returns.
**IMPORTANT:** Explicitly instruct the sub-agent that this is a log analysis
task and it MUST NOT attempt to read source code or debug the issue. Have the
sub-agent report the findings back to you so that you can aggregate them and
report them to the user.

### Flaky Tasks

These are generally lower priority than recent persistent failures, but they may
warrant investigation based on the original prompt. If so, invoke a sub-agent to
run `sk agent workflow task_drilldown` and follow the instructions it returns.
**IMPORTANT:** Explicitly instruct the sub-agent that this is a log analysis
task and it MUST NOT attempt to read source code or debug the issue.

## Phase 5: Report to the User

This somewhat depends on what was originally requested of you, but generally:

- **Active Persistent Failures:** For any tasks still failing at their most
  recent run, group the failures by root cause and report the task name(s),
  culprit commit (or suspect commit range), any error message you were able
  to extract, and any summary or recommendations you were able to derive.
- **Resolved Persistent Failures:** Leave these out unless the user requested
  to see even the resolved failures. If so, report them in the same way as
  persistent failures.
- **Flaky Failures:** If you investigated these deeply in Phase 4, group them
  by root cause or similar error message and report what you found. Otherwise,
  simply present a list of the flakily-failing task names.
