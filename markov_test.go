package nlm

import (
	"testing"
)

func TestSequences(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"abcde", "abcdefghijklmnopqrstuvwxyz"},
		{"lmnop", "lmnopqrstuvwxyz"},
	}

	mb := NewMarkovBuilder()
	mb.AddText("abcdefghijklmnopqrstuvwxyz", 5, EndOfDocument)

	ms := mb.ConvertToSource()

	for _, test := range tests {
		res := ms.GenerateText(test.input, 50)
		if res[0] != test.output {
			t.Errorf("For %s: expected %s but got %s", test.input, test.output, res[0])
		}
	}
}
