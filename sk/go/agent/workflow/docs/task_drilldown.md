# Skia Task Drilldown

This document contains instructions for agents to dig deeply into failed Skia
infrastructure tasks and enumerate failing tests, extract error messages,
investigate flaky failures, etc.

## Workflow

1. Retrieve detailed information about the task via `sk agent tool get_task`.
2. Retrieve the steps for the task via `sk agent tool get_task_steps`.
3. Find the relevant failed step(s). Note that some step failures may be
   expected, for example a step which tests file existence via some command
   that exits with a non-zero code when it does not exist.
4. If the task is a Task Driver or Recipe, retrieve the logs for the failed
   step via `sk agent tool get_task_driver_step_logs` or
   `sk agent tool get_recipe_step_logs`.
5. Scan the logs and extract a digestible snippet.
   - **WARNING:** Skia tasks (especially `dm` and `nanobench`) sometimes
     produce massive amounts of log spam (e.g., graphics API warnings, compiler
     warnings) that are non-fatal red herrings. Do not get distracted by them.
   - Do **not** assume a task timed out just because it ran for a long time or
     generated a lot of spam.
   - To find the _actual_ cause of failure, always check the **end** of the logs
     first. Specifically, look for a `Failures:` section or explicit test
     failure messages right before the step exits.
   - It may be helpful to pipe the logs directly into a file and then read
     the file in chunks rather than consuming all of the logs directly from the
     tool.
6. Analyze the error(s) and present a report to the user.
