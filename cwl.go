package nlm

import "math/rand"

// CountWeightedList holds a weighted list of runes and allows
// selection of random runes weighted by count.
type CountWeightedList struct {
	letters           []rune
	cumulativeWeights []int
	total             int
}

// NewCountWeightedList returns a pointer to a CountWeightedList structure
// from which you can select random elements weighted by count.
func NewCountWeightedList(counts map[rune]int) *CountWeightedList {
	n := len(counts)

	cwl := CountWeightedList{make([]rune, n), make([]int, n), 0}

	sum := 0
	i := 0
	for k, v := range counts {
		sum += v
		cwl.letters[i] = k
		cwl.cumulativeWeights[i] = sum
		i++
	}
	cwl.total = sum

	return &cwl
}

// GetRandomItem returns a random rune from the weighted list.
func (cwl *CountWeightedList) GetRandomItem() rune {
	if cwl.total == 0 {
		return 0
	}
	i := rand.Int() % cwl.total
	idx := binarySearch(cwl.cumulativeWeights, i)

	return cwl.letters[idx]
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
