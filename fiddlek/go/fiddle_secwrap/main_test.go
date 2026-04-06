package main

import (
	"testing"
)

func TestIsAllowedShellCmd(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "Legitimate clang command",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle main.cpp"},
			expected: true,
		},
		{
			name:     "Legitimate clang++ command (absolute path)",
			args:     []string{"/bin/sh", "-c", "/usr/bin/clang++ -o /tmp/fiddle main.cpp"},
			expected: true,
		},
		{
			name:     "Legitimate ar command",
			args:     []string{"/bin/sh", "-c", "ar rcs /tmp/lib.a main.o"},
			expected: true,
		},
		{
			name:     "Legitimate ninja command",
			args:     []string{"/bin/sh", "-c", "ninja -C out/Static fiddle"},
			expected: true,
		},
		{
			name:     "Blocked: Wrong number of arguments",
			args:     []string{"/bin/sh", "-c"},
			expected: false,
		},
		{
			name:     "Blocked: Wrong flag",
			args:     []string{"/bin/sh", "-x", "clang -o /tmp/fiddle main.cpp"},
			expected: false,
		},
		{
			name:     "Blocked: Disallowed binary (ls)",
			args:     []string{"/bin/sh", "-c", "ls -l"},
			expected: false,
		},
		{
			name:     "Blocked: Disallowed binary (sh recursively)",
			args:     []string{"/bin/sh", "-c", "/bin/sh -c clang"},
			expected: false,
		},
		{
			name:     "Blocked: Shell injection (semicolon)",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle main.cpp; rm -rf /"},
			expected: false,
		},
		{
			name:     "Blocked: Shell injection (pipe)",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle main.cpp | cat"},
			expected: false,
		},
		{
			name:     "Blocked: Shell injection (redirection)",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle main.cpp > /etc/passwd"},
			expected: false,
		},
		{
			name:     "Blocked: Command substitution (backticks)",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle `whoami`.cpp"},
			expected: false,
		},
		{
			name:     "Blocked: Command substitution (dollar-paren)",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle $(whoami).cpp"},
			expected: false,
		},
		{
			name:     "Blocked: Variable expansion",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle $HOME/main.cpp"},
			expected: false,
		},
		{
			name:     "Blocked: Dangerous flag (-fplugin)",
			args:     []string{"/bin/sh", "-c", "clang -fplugin=malicious.so -o /tmp/fiddle main.cpp"},
			expected: false,
		},
		{
			name:     "Blocked: Dangerous flag (-Xclang)",
			args:     []string{"/bin/sh", "-c", "clang -Xclang -load -Xclang malicious.so -o /tmp/fiddle main.cpp"},
			expected: false,
		},
		{
			name:     "Blocked: Dangerous flag (-specs)",
			args:     []string{"/bin/sh", "-c", "clang -specs=malicious.specs -o /tmp/fiddle main.cpp"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isAllowedShellCmd(tc.args)
			if actual != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
