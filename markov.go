package nlm

import "strings"

// MarkovBuilder is a map of strings pointing to counts of letters
// that follow each string. This is used while building a Markov
// chain structure
type MarkovBuilder struct {
	Counts   map[string]map[rune]int
	Fragment []rune
}

const EndOfParagraph = '\u001d'
const EndOfDocument = '\u0003'

func NewMarkovBuilder() *MarkovBuilder {
	return &MarkovBuilder{make(map[string]map[rune]int), []rune{}}
}

// AddText adds text to the Markov Builder, creating a table that will work
// with a chain length of size. It also includes a
// defined endCharacter that indicates the end of the passage.
func (b *MarkovBuilder) AddText(text string, size int, endCharacter rune) {
	letters := make([]rune, len(b.Fragment))
	copy(letters, b.Fragment)
	letters = append(letters, []rune(text)...)

	for i := 0; i <= len(letters)-size; i++ {
		var substring string
		if i+size == len(letters) {
			substring = string(letters[i:])
		} else {
			substring = string(letters[i : i+size])
		}
		v := b.Counts[substring]
		if v == nil {
			v = make(map[rune]int)
			b.Counts[substring] = v
		}

		if i+size >= len(letters) {
			v[endCharacter] += 1
		} else {
			v[letters[i+size]] += 1
		}
	}
	b.Fragment = make([]rune, size)
	copy(b.Fragment, letters[len(letters)-size+1:])
	b.Fragment[size-1] = endCharacter
}

type MarkovSource map[string](*CountWeightedList[rune])

func (b *MarkovBuilder) ConvertToSource() MarkovSource {
	ret := make(map[string](*CountWeightedList[rune]))

	for k, v := range b.Counts {
		ret[k] = NewCountWeightedList(v)
	}

	return ret
}

func (s MarkovSource) GetNextCharacter(v []rune) rune {
	cwl := s[string(v)]

	return cwl.GetRandomItem()
}

func (s MarkovSource) GenerateText(start string, maxLength int) []string {
	search := []rune(start)

	for i := 0; i < maxLength; i++ {
		c := s.GetNextCharacter(search[i:])

		if c == EndOfDocument {
			break
		}
		search = append(search, c)
	}

	subsections := strings.FieldsFunc(string(search), func(c rune) bool {
		return c == EndOfParagraph
	})
	return subsections
}
