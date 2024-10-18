package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestSplitFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple newline split",
			input:    "line1\nline2\nline3\n",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Carriage return split",
			input:    "line1\rline2\rline3\r",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Mixed newline and carriage return",
			input:    "line1\r\nline2\rline3\n",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "No newline at EOF",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			scanner.Split(splitFunc)

			var result []string
			for scanner.Scan() {
				result = append(result, scanner.Text())
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d lines, got %d", len(tt.expected), len(result))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("expected line %d to be %q, got %q", i, tt.expected[i], result[i])
				}
			}
		})
	}
}
