package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "bulbasaur   charmander squirtle",
			expected: []string{"bulbasaur", "charmander", "squirtle"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Length mismatch: got %d, expected %d for input '%s'", len(actual), len(c.expected), c.input)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Test failed: Actual word: %s and Expected word: %s do not match", word, expectedWord)
			}

		}
	}
}
