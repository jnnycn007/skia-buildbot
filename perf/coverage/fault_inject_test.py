#!/usr/bin/env python3
import unittest
from unittest.mock import patch, mock_open, MagicMock, call
import argparse
import pathlib
import sys
from typing import Final

# Add the directory containing the module to the Python path.
TOOL_DIR: Final = pathlib.Path(__file__).absolute()
if TOOL_DIR not in sys.path:
  sys.path.insert(0, TOOL_DIR)

from fault_inject import mutate_generated_html, process_puppeteer_test_modules, get_args

class TestMutateGeneratedHTML(unittest.TestCase):

    def test_mutate_generated_html_event_handler(self):
        content = '<button @click=${this.handleClick}>Click Me</button>'
        faults = list(mutate_generated_html(content))
        self.assertEqual(len(faults), 2) # One for @click, one for remove-button

        # Test @click fault
        click_fault = next(f for f in faults if f[0] == 'click-throws-error')
        self.assertEqual(click_fault[0], 'click-throws-error')
        self.assertEqual(click_fault[1], 1)
        expected_line = '<button @click=${() => { throw new Error(\'Injected Fault!\'); }}>Click Me</button>'
        self.assertEqual(click_fault[2], expected_line)

    def test_mutate_generated_html_quoted_event_handler(self):
        content = '<button @click="someFunction()"></button>'
        faults = list(mutate_generated_html(content))
        click_fault = next(f for f in faults if f[0] == 'click-throws-error')
        expected_line = '<button @click=${() => { throw new Error(\'Injected Fault!\'); }}></button>'
        self.assertEqual(click_fault[2], expected_line)

    def test_mutate_generated_html_remove_element(self):
        content = '<div><button>Click Me</button></div>'
        faults = list(mutate_generated_html(content))

        # Test remove-div fault
        div_fault = next(f for f in faults if f[0] == 'remove-div')
        self.assertEqual(div_fault[0], 'remove-div')
        self.assertEqual(div_fault[1], 1)
        self.assertEqual(div_fault[2], '')

        # Test remove-button fault
        button_fault = next(f for f in faults if f[0] == 'remove-button')
        self.assertEqual(button_fault[0], 'remove-button')
        self.assertEqual(button_fault[1], 1)
        self.assertEqual(button_fault[2], '<div></div>')

    def test_mutate_generated_html_no_faults(self):
        content = 'const x = 1;'
        faults = list(mutate_generated_html(content))
        self.assertEqual(len(faults), 0)

    def test_mutate_generated_html_multiline_element(self):
        content = (
            '<div>\n'
            '  <span>Hello</span>\n'
            '</div>'
        )
        faults = list(mutate_generated_html(content))
        div_fault = next(f for f in faults if f[0] == 'remove-div')
        self.assertEqual(div_fault[2], '')

class TestFindPuppeteerTestModules(unittest.TestCase):

    @patch('fault_inject.mutate_generated_html')
    @patch('subprocess.run')
    @patch('os.path.exists')
    @patch('os.walk')
    def test_survivor_is_logged(self, mock_walk, mock_exists, mock_run, mock_mutate_generated_html):
        # Setup mocks
        ts_module_content = '<button @click=${this.handler}>Test</button>'
        faulty_content = '<button @click=${() => { throw new Error(\'Injected Fault!\'); }}>Test</button>'

        mock_walk.return_value = [
            ('perf/modules/my-module', [], ['my-module_puppeteer_test.ts']),
        ]
        mock_exists.return_value = True
        mock_mutate_generated_html.return_value = [('click-throws-error', 1, faulty_content)]

        # Mock subprocess to simulate a test that passes (survivor)
        mock_run.return_value = MagicMock(returncode=0)

        m = mock_open(read_data=ts_module_content)
        with patch('builtins.open', m):
            process_puppeteer_test_modules('log.txt', {'my-module'})

        # Assertions
        mock_mutate_generated_html.assert_called_once_with(ts_module_content)
        mock_run.assert_called_once()
        self.assertIn('bazelisk', mock_run.call_args[0][0])

        # Check that the file was written with the fault and then restored
        handle = m()
        self.assertEqual(handle.write.call_count, 3)
        calls = [call(faulty_content), call().writelines(unittest.mock.ANY), call(ts_module_content)]
        handle.write.assert_has_calls(calls, any_order=False)

        # Check that survivor diff is logged
        log_write_call = handle.write.call_args_list[1]
        self.assertIn('--- SURVIVOR: perf/modules/my-module/my-module.ts', log_write_call[0][0])

    @patch('fault_inject.mutate_generated_html')
    @patch('subprocess.run')
    @patch('os.path.exists')
    @patch('os.walk')
    def test_failing_test_is_not_logged(self, mock_walk, mock_exists, mock_run, mock_mutate_generated_html):
        # Setup mocks
        ts_module_content = '<button>Test</button>'
        faulty_content = ''

        mock_walk.return_value = [
            ('perf/modules/my-module', [], ['my-module_puppeteer_test.ts']),
        ]
        mock_exists.return_value = True
        mock_mutate_generated_html.return_value = [('remove-button', 1, faulty_content)]

        # Mock subprocess to simulate a test that fails (expected)
        mock_run.return_value = MagicMock(returncode=1)

        m = mock_open(read_data=ts_module_content)
        with patch('builtins.open', m):
            process_puppeteer_test_modules('log.txt', {'my-module'})

        # Assertions
        handle = m()
        # open() is called for reading, writing fault, and restoring.
        # No call to write to the log file.
        self.assertEqual(handle.write.call_count, 2)
        calls = [call(faulty_content), call(ts_module_content)]
        handle.write.assert_has_calls(calls, any_order=False)

    @patch('os.walk')
    def test_module_filtering(self, mock_walk):
        mock_walk.return_value = [
            ('perf/modules/module-a', [], ['module-a_puppeteer_test.ts']),
            ('perf/modules/module-b', [], ['module-b_puppeteer_test.ts']),
        ]

        with patch('builtins.open', mock_open(read_data='content')) as m:
            with patch('os.path.exists', return_value=True):
                 with patch('subprocess.run'):
                    process_puppeteer_test_modules('log.txt', {'module-b'})

        # Assert that we only tried to open the files for module-b
        m.assert_called_once_with('perf/modules/module-b/module-b.ts', 'r')


class TestGetArgs(unittest.TestCase):

    @patch('argparse.ArgumentParser.parse_args')
    def test_get_args_parsing(self, mock_parse_args):
        # Test with a single module
        mock_parse_args.return_value = argparse.Namespace(
            log_file='/tmp/log.txt',
            filter_modules='a-module'
        )
        with patch('sys.argv', ['script.py', '--log_file', '/tmp/log.txt', '--filter_modules', 'a-module']):
             args = get_args()
             self.assertEqual(args.log_file, '/tmp/log.txt')
             self.assertEqual(args.filter_modules, 'a-module')

        # Test with multiple modules
        mock_parse_args.return_value = argparse.Namespace(
            log_file='/tmp/log.txt',
            filter_modules='a-module,b-module'
        )
        with patch('sys.argv', ['script.py', '--log_file', '/tmp/log.txt', '--filter_modules', 'a-module,b-module']):
             args = get_args()
             self.assertEqual(args.log_file, '/tmp/log.txt')
             self.assertEqual(args.filter_modules, 'a-module,b-module')

if __name__ == '__main__':
    unittest.main()
