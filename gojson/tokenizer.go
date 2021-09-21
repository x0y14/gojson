package gojson

func NewTokenizer(text string) *Tokenizer {
	return &Tokenizer{
		Original: text,
		Runes:    []rune(text),
		Pos:      0,
	}
}

type Tokenizer struct {
	Original string
	Runes    []rune
	Pos      int
}

func (t *Tokenizer) Token() rune {
	return t.Runes[t.Pos]
}

func (t *Tokenizer) TokenAt(pos int) rune {
	return t.Runes[pos]
}

func (t *Tokenizer) Next() {
	t.Pos++
}

func (t *Tokenizer) Prev() {
	t.Pos--
}

func (t *Tokenizer) ShowNext() rune {
	return t.Runes[t.Pos+1]
}

func (t *Tokenizer) ShowPrev() rune {
	return t.Runes[t.Pos-1]
}
