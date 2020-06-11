package da

import "strings"

type request struct {
	Text string `json:"q"`
}

// A Result represents a result of analysis for each clausea phrase.
type Result struct {
	Chunks []Chunk `json:"chunks"`
}

// A Chunk represents information of a clausea phrase.
type Chunk struct {
	Head   int     `json:"head"`
	ID     int     `json:"id"`
	Tokens []Token `json:"tokens"`
}

// String returns a string of chunk.
func (c Chunk) String() string {
	var s strings.Builder
	for _, t := range c.Tokens {
		s.WriteString(t.Surface())
	}
	return s.String()
}

// A Token represents information of a morpheme.
type Token []string

func (t Token) Surface() string {
	if len(t) < 1 {
		return ""
	}

	return t[0]
}

func (t Token) Reading() string {
	if len(t) < 2 {
		return ""
	}

	return t[1]
}

func (t Token) Baseform() string {
	if len(t) < 3 {
		return ""
	}

	return t[2]
}

func (t Token) PartOfSpeech() string {
	if len(t) < 4 {
		return ""
	}

	return t[3]
}

func (t Token) POSDetail() string {
	if len(t) < 5 {
		return ""
	}

	return t[4]
}

func (t Token) ConjugationType() string {
	if len(t) < 6 {
		return ""
	}

	return t[5]
}

func (t Token) ConjugationForm() string {
	if len(t) < 7 {
		return ""
	}

	return t[6]
}
