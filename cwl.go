package nlm

import "math/rand"

// CountWeightedList holds a weighted list of runes and allows
// selection of random runes weighted by count.
type CountWeightedList[T comparable] struct {
	Items             []T   `json:"items"`
	CumulativeWeights []int `json:"weights"`
	Total             int   `json:"total"`
}

// NewCountWeightedList returns a pointer to a CountWeightedList structure
// from which you can select random elements weighted by count.
func NewCountWeightedList[T comparable](counts map[T]int) *CountWeightedList[T] {
	n := len(counts)

	cwl := CountWeightedList[T]{make([]T, n), make([]int, n), 0}

	sum := 0
	i := 0
	for k, v := range counts {
		sum += v
		cwl.Items[i] = k
		cwl.CumulativeWeights[i] = sum
		i++
	}
	cwl.Total = sum

	return &cwl
}

// GetRandomItem returns a random rune from the weighted list.
func (cwl *CountWeightedList[T]) GetRandomItem() T {
	if cwl.Total == 0 {
		var ret T
		return ret
	}
	i := rand.Int() % cwl.Total
	idx := binarySearch(cwl.CumulativeWeights, i)

	return cwl.Items[idx]
}

func binarySearch(a []int, n int) int {
	l := len(a)
	si := l / 2
	p := si

	for {
		switch a[p] > n {
		case true:
			switch si > 1 {
			case true:
				si = si / 2
				p -= si
			case false:
				for i := p; i > 0; i-- {
					if a[i-1] <= n {
						return i
					}
				}
				return 0
			}
		case false:
			switch si > 1 {
			case true:
				si = si / 2
				p += si
			case false:
				for i := p + 1; i < len(a); i++ {
					if a[i] > n {
						return i
					}
				}
				return len(a)
			}
		}
	}
}
