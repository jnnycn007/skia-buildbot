#!/usr/bin/env python3
"""This tool identifies gaps in test coverage by utilizing fault injection.

The tool currently identifies gaps in test coverage for components within
`perf/modules` that include Puppeteer tests.

Document: go/perf-module-mutation

How to use:
skia/buildbot$ python3 perf/coverage/fault_inject.py \
  --log_file=/tmp/fault-inject.txt \
  --filter_modules=triage-menu-sk,new-bug-dialog-sk

The output log contains a `git diff` of the injected faults that the tests
failed to detect. Please analyze this report to identify the specific scenarios
needed to update the Puppeteer tests.

Note: Presubmit checks do not include `fault_inject_test.py`. Please ensure you
run the unit tests locally after making changes to this file.
skia/buildbot$ python3 perf/coverage/fault_inject_test.py
"""
import argparse
import os
import re
import subprocess
import difflib

# Fault to inject in event handlers to throw an error.
THROW_ERROR_FAULT = "=${() => { throw new Error('Injected Fault!'); }}"
BUSY_WAIT_FAULT = "=${() => { const start = Date.now(); while (Date.now() - start < 10000) {} }}"

def mutate_generated_html(ts_module_content):
  """
  A generator that yields modified versions of the ts_module_content,
  each with a single fault injected, along with a description of the fault.
  """
  lines = ts_module_content.split('\n')
  event_handlers = ["@click", "@change", "@input", "@submit"]
  for i, line in enumerate(lines):
    for handler in event_handlers:
      handler_name = handler[1:]
      # Fault: Replace event handlers to throw an error.
      if f'{handler}=' in line:
        # This regex finds handler="..." or handler='...' or handler=${...} and
        # replaces the content with an expression that throws an error.
        faulty_line = re.sub(
            handler + r"""=(?:(['"]).*?\1|\$\{[^}]+\})""",
            handler + THROW_ERROR_FAULT,
            line
        )
        if faulty_line != line:
          modified_lines = lines[:i] + [faulty_line] + lines[i+1:]
          yield (f'{handler_name}-throws-error', i + 1, '\n'.join(modified_lines))

  # Fault: Remove various HTML elements, including multi-line ones.
  # We do this separately as it's not a line-by-line transformation.
  elements_to_remove = [
      'div', 'button', 'span', 'ul', 'li', 'md-icon-button', 'md-dialog',
      'form', 'md-switch', 'pivot-table-sk', 'pivot-query-sk', 'toast-sk'
  ]
  for tag in elements_to_remove:
    # The regex looks for <tag ...> ... </tag>
    pattern = re.compile(rf'<{tag}.*?</{tag}>', re.DOTALL)
    # We must use a list to collect matches before iterating to avoid issues
    # with finditer on a string that is being modified.
    matches = list(pattern.finditer(ts_module_content))
    for match in matches:
      modified_content = ts_module_content[:match.start()] + ts_module_content[match.end():]
      # For logging, find the line number of the start of the match.
      line_num = ts_module_content.count('\n', 0, match.start()) + 1
      yield (f'remove-{tag}', line_num, modified_content)

def log_survivor(log_file, module_ts_file, fault_type, line_num, original_content, modified_content):
  """Logs a test that survived fault injection."""
  print(f'  [SURVIVED] Tests passed with fault: {fault_type} at line {line_num}')
  diff = difflib.unified_diff(
      original_content.splitlines(keepends=True),
      modified_content.splitlines(keepends=True),
      fromfile=f'a/{module_ts_file}', tofile=f'b/{module_ts_file}')
  with open(log_file, 'a') as log:
    log.write(
        f'--- SURVIVOR: {module_ts_file} | FAULT: {fault_type} at line {line_num} ---\n'
    )
    log.writelines(diff)


def process_puppeteer_test_modules(log_file, filter_modules):
  """
  Finds Puppeteer test files within the `perf/modules` directory. If a module
  name matches the filter_modules argument, the script detects supported
  mutation patterns, applies a mutation, and executes the Puppeteer test. If
  the test passes, it is logged as a potential candidate for improved
  integration test coverage.
  """
  modules_dir = 'perf/modules'
  for root, _, files in os.walk(modules_dir):
    for filename in files:
      if filename.endswith('_puppeteer_test.ts'):
        module_base_name = filename.removesuffix('_puppeteer_test.ts')
        if module_base_name not in filter_modules:
          continue

        module_ts_file = os.path.join(root, f'{module_base_name}.ts')
        if not os.path.exists(module_ts_file):
          print(f'{module_ts_file} is missing!')
          continue

        print(f'--- Starting fault injection for {module_ts_file} ---')
        total = 0
        servived = 0

        with open(module_ts_file, 'r') as f:
          ts_module_content = f.read()

        for fault_type, line_num, modified_content in mutate_generated_html(ts_module_content):
          total += 1
          print(f'  Injecting fault: {fault_type} at line {line_num}')
          with open(module_ts_file, 'w') as f:
            f.write(modified_content)

          try:
            test_target = f'//{root}:{module_base_name}_puppeteer_test'
            command = [
                'bazelisk', 'test',
                '--config=mayberemote',
                '--test_output=all',
                '--nocache_test_results',
                test_target
            ]
            # We expect tests to fail. If they don't, it's a "survivor".
            result = subprocess.run(command, check=False, capture_output=True, text=True)

            if result.returncode == 0 and fault_type.endswith('-throws-error'):
              print(f'  [INFO] Survived error throw, trying busy-wait for {fault_type} at line {line_num}')
              # The test survived an error being thrown. Let's try a busy-wait
              # loop to see if it's truly a survivor. This can catch cases
              # where the test doesn't wait for async operations.
              busy_wait_content = modified_content.replace(
                  THROW_ERROR_FAULT,
                  BUSY_WAIT_FAULT
              )
              with open(module_ts_file, 'w') as f:
                f.write(busy_wait_content)

              result_busy_wait = subprocess.run(command, check=False, capture_output=True, text=True)

              if result_busy_wait.returncode == 0:
                servived += 1
                log_survivor(log_file, module_ts_file, fault_type, line_num,
                             ts_module_content, modified_content)
            elif result.returncode == 0:
                servived += 1
                log_survivor(log_file, module_ts_file, fault_type, line_num,
                             ts_module_content, modified_content)

          finally:
            # Restore original file content
            with open(module_ts_file, 'w') as f:
              f.write(ts_module_content)

        score = 100 * (total - servived) / total if total > 0 else 100
        print(f'{module_ts_file}: {total}, {servived}, {score:.2f}%')

def get_args():
  """Parses and returns command-line arguments."""
  parser = argparse.ArgumentParser(
      description='Inject faults into TypeScript modules and run Puppeteer tests.')
  parser.add_argument(
      '--log_file',
      required=True,
      help='Path to the log file for fault injection survivors.')
  parser.add_argument(
      '--filter_modules',
      type=lambda arg: arg.split(','),
      required=True,
      help='A comma-separated list of modules to run fault injection on.')
  return parser.parse_args()

if __name__ == '__main__':
  args = get_args()
  filter_modules_set = set(args.filter_modules) if args.filter_modules else None
  process_puppeteer_test_modules(args.log_file, filter_modules_set)
