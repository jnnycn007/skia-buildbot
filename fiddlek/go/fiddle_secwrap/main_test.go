package main

import (
	"testing"
)

func TestIsAllowedExec(t *testing.T) {
	testCases := []struct {
		name     string
		binary   string
		args     []string
		envp     []string
		expected bool
	}{
		{
			name:     "Legitimate clang command",
			binary:   "/usr/bin/clang",
			args:     []string{"clang", "-o", "/tmp/fiddle", "main.cpp"},
			envp:     []string{"PATH=/usr/bin", "HOME=/tmp"},
			expected: true,
		},
		{
			name:     "Blocked: LD_PRELOAD",
			binary:   "/usr/bin/clang",
			args:     []string{"clang", "-o", "/tmp/fiddle", "main.cpp"},
			envp:     []string{"LD_PRELOAD=/tmp/malicious.so"},
			expected: false,
		},
		{
			name:     "Blocked: LD_LIBRARY_PATH",
			binary:   "/usr/bin/clang",
			args:     []string{"clang", "-o", "/tmp/fiddle", "main.cpp"},
			envp:     []string{"LD_LIBRARY_PATH=/tmp"},
			expected: false,
		},
		{
			name:     "Blocked: PYTHONHOME",
			binary:   "/usr/bin/clang",
			args:     []string{"clang", "-o", "/tmp/fiddle", "main.cpp"},
			envp:     []string{"PYTHONHOME=/tmp"},
			expected: false,
		},
		{
			name:     "Blocked: Dangerous compiler flag",
			binary:   "/usr/bin/clang",
			args:     []string{"clang", "-fplugin=malicious.so", "-o", "/tmp/fiddle", "main.cpp"},
			envp:     []string{"PATH=/usr/bin"},
			expected: false,
		},
		{
			name:     "Legitimate shell command",
			binary:   "/bin/sh",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle main.cpp"},
			envp:     []string{"PATH=/usr/bin"},
			expected: true,
		},
		{
			name:     "Blocked: Shell command with dangerous env",
			binary:   "/bin/sh",
			args:     []string{"/bin/sh", "-c", "clang -o /tmp/fiddle main.cpp"},
			envp:     []string{"LD_PRELOAD=/tmp/malicious.so"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isAllowedExec(tc.binary, "fiddle", tc.args, tc.envp, true)
			if actual != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
