# Test Coverage Tools

## Stryker

[Stryker](https://stryker-mutator.io/) is a third-party tool utilized for
mutating TypeScript code to evaluate the effectiveness of unit tests. The core
premise is that robust tests should fail when the underlying code is modified; a
passing test suite in this context may highlight insufficient coverage.

- Pros

  - Available at no cost.
  - Provides straightforward configuration for selecting target tests and
    execution methods

- Cons

  - Optimized primarily for unit-level testing.
  - No mutation in the generated HTML output.

### Operational Guide

As a Node.js-based utility, Stryker can be installed by following the official
[documentation](https://stryker-mutator.io/docs/General/example/). To set up npm
on a cloudtop environment, execute the following within the `skia/buildbot/perf`
directory:

```bash
sudo apt install npm
```

Once installed, modify `skia/buildbot/perf/stryker.config.json` to define
mutation targets and test execution parameters. For Karma or Puppeteer
integration, the `testRunner` should be set to command. Below is a configuration
example targeting the `triage-menu-sk.ts` module:

```json
{
  "$schema": "./node_modules/@stryker-mutator/core/schema/stryker-schema.json",
  "packageManager": "npm",
  "reporters": ["html", "clear-text", "progress"],
  "mutate": ["modules/triage-menu-sk/triage-menu-sk.ts"],
  "testRunner": "command",
  "commandRunner": {
    "command": "bazelisk test --config=mayberemote --test_output=all
      --nocache_test_results //perf/modules/triage-menu-sk/..."
  },
  "coverageAnalysis": "off",
  "plugins": []
}
```

To initiate the mutation run, use this command in the `skia/buildbot/perf/`
folder:

```shell
npx stryker run
```

Execution time varies based on the module size and valid mutation count. The
final report, including any "surviving" mutants (changes that didn't trigger a
test failure), will be available in the console output and a generated HTML file
located at the path specified at the end of the run.

[Here](https://paste.googleplex.com/5422945469071360) is an example of a
survival report. Analyzing these results is essential for identifying specific
scenarios that require additional unit or integration tests.

## Mutate Generated HTML

While the Stryker tool effectively mutates TypeScript code, it lacks the ability
to mutate generated HTML. This limitation is significant for identifying gaps in
integration tests, as many CUJs depend on web button elements and event
handlers. To address this, a new tool was developed to detect HTML code
generation within TypeScript components and apply the following mutations:

- Locate events such as ["@click", "@change", "@input", "@submit"] and replace
  their handlers with throw new Error().
- Identify and remove critical elements that impact page behavior, such as
  ['div', 'button', 'md-switch', 'pivot-table-sk'], from the HTML generation
  script.

The goal is to ensure that while the TypeScript remains functional, the
execution tests fail as a result of the corrupted HTML.

### Operational Procedures

The `perf/coverage/fault_inject.py` script facilitates HTML mutation and
initiates Puppeteer integration tests. The tool requires two primary arguments:

- `--filter_modules=`: Mandatory. Specify the target perf module name (e.g.,
  `triage-menu-sk`). The tool automatically prepends the `perf/modules/` path.
  While multiple modules can be comma-separated, evaluating them individually
  is recommended due to the duration of test cycles.
- `--log_file=`: Mandatory. While a summary is displayed in the console, the
  comprehensive analysis report is saved to this designated file.

The following example demonstrates the tool running for the "triage-menu-sk"
module. You can view the sample output
[here](https://paste.googleplex.com/6223122207473664).

```bash
python3 perf/coverage/fault_inject.py \
  --log_file=/tmp/fault-inject.txt \
  --filter_modules=triage-menu-sk
```
