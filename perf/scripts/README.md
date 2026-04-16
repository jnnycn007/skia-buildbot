# Perf Scripts & Tools

This directory contains utility scripts for development, testing, and debugging the Perf
application.

## MultiGraph Monkey Test

A stress-testing tool for the MultiGraph UI (`/m/`) that automates user interactions to find
bugs, crashes, and regressions.

### Files

- **`multigraph_monkey_test.js`**: The core test logic. It's a standalone JavaScript file
  designed to run **inside the browser context**.
- **`run_monkey_test.js`**: A Node.js runner that launches Chrome (via Puppeteer), handles
  Google Login, and injects the monkey test script automatically.

---

### How to Run

#### Option 1: Automated Runner (Recommended)

Use the provided runner to launch a headful Chrome instance and run the test automatically.

```bash
# From the repository root
node perf/scripts/run_monkey_test.js
```

**Configuration:**

- **Headless Mode:** `HEADLESS=true node perf/scripts/run_monkey_test.js` (Default: `false`
  to allow login).
- **Chrome Path:** The script automatically detects Chrome at `/usr/bin/google-chrome`.
  Set `CHROME_BIN` env var to override.

#### Option 2: Manual (DevTools Snippet)

1. Open Chrome DevTools (**F12**).
2. Go to **Sources** > **Snippets**.
3. Create a new snippet and paste the content of `multigraph_monkey_test.js`.
4. Right-click the snippet and select **Run**.

#### Option 3: Tampermonkey (Browser Extension)

1. Install the [Tampermonkey](https://www.tampermonkey.net/) extension.
2. Create a new script and paste the content of `multigraph_monkey_test.js`.
3. Navigate to any MultiGraph page (`https://chrome-perf.corp.goog/m/`).
4. The script will load automatically (check console for "SUPER MONKEY" logs).

---

### What it Tests

The monkey test performs the following actions:

1. **Dynamic Discovery:** Drills down through filter options to find a dataset with ~50 traces.
2. **Graph Rendering:** Verifies graphs are drawn and layout is valid.
3. **URL State:** Checks if URL parameters match the current selection.
4. **Split/Unsplit:** Enables "Split by" mode and verifies graph multiplication.
5. **Trace Manipulation:** Removes/Adds traces and checks graph count.
6. **Sidebar Sync:** Verifies sidebar rows match displayed traces.
7. **Overflow:** Triggers "Too many traces" warning and verifies recovery.
8. **Pagination:** Iterates through all pages of a large dataset.
9. **Subset Recovery:** Restores state after overflow test.
10. **Pagination Stress:** Backtracks to find a high-cardinality field and iterates through pages.
11. **Chrome Internal Scenario:** (Env specific) Tests Speedometer3 split view.
12. **Load All Charts:** Verifies functionality of the "Load All" button for large datasets.
13. **Primary Checkbox:** Verifies splitting works when "Primary" filter is active.
14. **Select All Sync:** Verifies "Select All" checkbox state syncs with manual item selection.
