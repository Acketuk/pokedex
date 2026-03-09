package helpers

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    HeLlo worlD   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: make([]string, 0),
		},
		{
			input:    "🍎 or 🫐",
			expected: []string{"🍎", "or", "🫐"},
		},
	}

	for i, current := range cases {
		got := CleanInput(current.input)

		if len(got) != len(current.expected) {
			// test failed. length's are not equal
			t.Errorf("⛔ test %d: got: %#v, expected: %#v\nlength got: %d, length expected: %d", i+1, got, current.expected, len(got), len(current.expected))

			for j := range got {
				word := got[j]
				expectedWord := current.expected[j]

				if word != expectedWord {
					t.Errorf("⛔ test %d: %s != %s\n", i, word, expectedWord)
				}
			}
		}
	}
}
