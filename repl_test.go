package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected []string
	}{
		"simple": {
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := cleanInput(tc.input)

			if len(got) != len(tc.expected) {
				t.Errorf("length mismatch: got %d words, expected %d words", len(got), len(tc.expected))
				return
			}
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("not the same: got %q, expected %q", got, tc.expected)
			}

		})
	}

}
