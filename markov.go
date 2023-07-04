package nlm

// MarkovBuilder is a map of strings pointing to counts of letters
// that follow each string. This is used while building a Markov
// chain structure
type MarkovBuilder map[string]map[rune]int

const EndOfParagraph = '\u001d'
const EndOfDocument = '\u0003'

func NewMarkovBuilder() MarkovBuilder {
	return make(map[string]map[rune]int)
}

// AddText adds text to the Markov Builder, creating a table that will work
// with a chain length of size. It also includes a
// defined endCharacter that indicates the end of the passage.
func (b MarkovBuilder) AddText(text string, size int, endCharacter rune) {
	letters := []rune(text)

	for i := 0; i < len(letters)-size; i++ {
		substring := string(letters[i : i+size])
		v := b[substring]
		if v == nil {
			v = make(map[rune]int)
			b[substring] = v
		}

		if i+size >= len(letters) {
			v[endCharacter] += 1
		} else {
			v[letters[i+size]] += 1
		}
	}
}

type MarkovSource map[string](*CountWeightedList[rune])

func (b MarkovBuilder) ConvertToSource() MarkovSource {
	ret := make(map[string](*CountWeightedList[rune]))

	for k, v := range b {
		ret[k] = NewCountWeightedList(v)
	}

	return ret
}

func (s MarkovSource) GetNextCharacter(v []rune) rune {
	cwl := s[string(v)]

	return cwl.GetRandomItem()
}
