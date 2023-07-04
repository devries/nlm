package nlm

import (
	"testing"
)

func TestCountWeightedList(t *testing.T) {
	tests := []map[rune]int{
		{'a': 1},
		{'a': 1, 'b': 1, 'c': 1},
		{'a': 10},
		{'a': 10, 'b': 5, 'c': 5},
		{'a': 1, 'b': 4, 'c': 0},
		{'a': 0, 'b': 0, 'c': 1},
		{'a': 0, 'b': 0, 'c': 5},
		{'a': 0, 'b': 5, 'c': 0},
	}

	for _, test := range tests {
		cwl := NewCountWeightedList(test)
		r := cwl.GetRandomItem()
		if test[r] <= 0 {
			t.Errorf("For %#v: got %d which is not in the list", test, r)
		}
	}
}

func TestEmptyCountWeightedList(t *testing.T) {
	tests := []map[rune]int{
		{},
		{'a': 0},
		{'a': 0, 'b': 0, 'c': 0},
	}

	for _, test := range tests {
		cwl := NewCountWeightedList(test)
		r := cwl.GetRandomItem()
		if r != 0 {
			t.Errorf("For %#v: expected 0, got %d", test, r)
		}
	}
}
