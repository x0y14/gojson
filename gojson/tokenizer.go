package gojson

import (
	"strings"
	"unicode"
)

func NewTokenizer(text string) *Tokenizer {
	return &Tokenizer{
		Raw:     &text,
		Letters: []rune(text),
		Pos:     0,
	}
}

type Tokenizer struct {
	Raw     *string
	Letters []rune
	Pos     int
}

func (t *Tokenizer) Letter() rune {
	return t.Letters[t.Pos]
}

func (t *Tokenizer) LetterAt(pos int) rune {
	return t.Letters[pos]
}

func (t *Tokenizer) NextLetter() rune {
	return t.Letters[t.Pos+1]
}

func (t *Tokenizer) PrevLetter() rune {
	return t.Letters[t.Pos-1]
}

func (t *Tokenizer) GoNext() {
	t.Pos++
}

func (t *Tokenizer) GoPrev() {
	t.Pos--
}

func (t *Tokenizer) IsEof() bool {
	return t.Pos >= len(t.Letters)
}

func (t *Tokenizer) ConsumeWhiteSpace() Token {
	var data []rune

	startPos := t.Pos
	for len(t.Letters) > t.Pos {
		if unicode.IsSpace(t.Letter()) {
			data = append(data, t.Letter())
			t.GoNext()
		} else {
			break
		}
	}
	endPos := t.Pos

	return Token{
		Type:     TWhiteSpace,
		Data:     data,
		StartPos: startPos,
		EndPos:   endPos,
	}
}

func (t *Tokenizer) ConsumeKeyword() (Token, error) {
	var data []rune

	startPos := t.Pos
	for len(t.Letters) > t.Pos {
		if unicode.IsSpace(t.Letter()) || t.Letter() == ':' || t.Letter() == ',' || t.Letter() == ']' || t.Letter() == '}' {
			break
		} else {
			data = append(data, t.Letter())
			t.GoNext()
		}
	}
	endPos := t.Pos

	var tokenType TokenType
	var err error
	switch string(data) {
	case "true":
		tokenType = TTrue
		err = nil
	case "false":
		tokenType = TFalse
		err = nil
	case "null":
		tokenType = TNull
		err = nil
	default:
		tokenType = TUnknown
		err = &TokenizerError{
			ErrorType:    UndefinedKeywordError,
			ErrorMessage: "undefined keyword",
			Letters:      data,
			StartPos:     startPos,
			EndPos:       endPos,
		}
	}

	return Token{
		Type:     tokenType,
		Data:     data,
		StartPos: startPos,
		EndPos:   endPos,
	}, err
}

func (t *Tokenizer) ConsumeNumber() (Token, error) {
	var data []rune
	startPos := t.Pos

	for len(t.Letters) > t.Pos {
		if unicode.IsDigit(t.Letter()) || t.Letter() == '.' || t.Letter() == '-' || t.Letter() == '+' {
			data = append(data, t.Letter())
			t.GoNext()
		} else {
			break
		}
	}
	endPos := t.Pos

	var err error
	if data[0] == '.' {
		err = &TokenizerError{
			ErrorType:    InvalidDataError,
			ErrorMessage: "Dot must not be at the beginning.",
			Letters:      data,
			StartPos:     startPos,
			EndPos:       endPos,
		}
	} else if len(string(data))-len(strings.ReplaceAll(string(data), ".", "")) > 1 {
		err = &TokenizerError{
			ErrorType:    InvalidDataError,
			ErrorMessage: "There must not be more than one dot.",
			Letters:      data,
			StartPos:     startPos,
			EndPos:       endPos,
		}
	} else if strings.Contains(string(data), "-") && data[0] != '-' {
		err = &TokenizerError{
			ErrorType:    InvalidDataError,
			ErrorMessage: "Minus must be at the beginning.",
			Letters:      data,
			StartPos:     startPos,
			EndPos:       endPos,
		}
	} else if len(string(data))-len(strings.ReplaceAll(string(data), "-", "")) > 1 {
		err = &TokenizerError{
			ErrorType:    InvalidDataError,
			ErrorMessage: "There must not be more than one minus.",
			Letters:      data,
			StartPos:     startPos,
			EndPos:       endPos,
		}
	} else if strings.Contains(string(data), "+") && data[0] != '+' {
		err = &TokenizerError{
			ErrorType:    InvalidDataError,
			ErrorMessage: "Plus must be at the beginning.",
			Letters:      data,
			StartPos:     startPos,
			EndPos:       endPos,
		}
	} else if len(string(data))-len(strings.ReplaceAll(string(data), "+", "")) > 1 {
		err = &TokenizerError{
			ErrorType:    InvalidDataError,
			ErrorMessage: "There must not be more than one plus.",
			Letters:      data,
			StartPos:     startPos,
			EndPos:       endPos,
		}
	}

	return Token{
		Type:     TNumber,
		Data:     data,
		StartPos: startPos,
		EndPos:   endPos,
	}, err
}

func (t *Tokenizer) ConsumeString() (Token, error) {
	var data []rune

	var quotationCount int

	startPos := t.Pos
	for len(t.Letters) > t.Pos {
		if t.Letter() == '"' {
			// escaped quotation
			if t.Pos != 0 && t.PrevLetter() == '\\' {
				data = append(data, t.Letter())
				t.GoNext()
			} else {
				quotationCount++
				t.GoNext()
				if quotationCount >= 2 {
					break
				}
			}
		} else {
			data = append(data, t.Letter())
			t.GoNext()
		}
	}
	endPos := t.Pos

	var err error
	return Token{
		Type:     TString,
		Data:     data,
		StartPos: startPos,
		EndPos:   endPos,
	}, err
}

func (t *Tokenizer) Tokenize() *[]Token {
	var tokens []Token

	for len(t.Letters) > t.Pos {
		if t.IsEof() {
			return &tokens
		}

		if t.Letter() == '"' {
			token, err := t.ConsumeString()
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, token)
			continue
		}

		if unicode.IsDigit(t.Letter()) || t.Letter() == '-' || t.Letter() == '+' {
			token, err := t.ConsumeNumber()
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, token)
			continue
		}

		if unicode.IsLetter(t.Letter()) {
			token, err := t.ConsumeKeyword()
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, token)
			continue
		}

		if unicode.IsSpace(t.Letter()) {
			// ignore whitespace
			// tokens = append(tokens, t.ConsumeWhiteSpace())
			t.ConsumeWhiteSpace()
			continue
		}

		if t.Letter() == ',' {
			tokens = append(tokens, Token{
				Type:     TComma,
				Data:     []rune{t.Letter()},
				StartPos: t.Pos,
				EndPos:   t.Pos + 1,
			})
			t.GoNext()
			continue
		}

		if t.Letter() == ':' {
			tokens = append(tokens, Token{
				Type:     TColon,
				Data:     []rune{t.Letter()},
				StartPos: t.Pos,
				EndPos:   t.Pos + 1,
			})
			t.GoNext()
			continue
		}

		// curly bracket
		if t.Letter() == '{' || t.Letter() == '}' {
			var tokenType TokenType
			if t.Letter() == '{' {
				tokenType = TLCurlyBracket
			} else {
				tokenType = TRCurlyBracket
			}
			tokens = append(tokens, Token{
				Type:     tokenType,
				Data:     []rune{t.Letter()},
				StartPos: t.Pos,
				EndPos:   t.Pos + 1,
			})
			t.GoNext()
			continue
		}

		// square bracket
		if t.Letter() == '[' || t.Letter() == ']' {
			var tokenType TokenType
			if t.Letter() == '[' {
				tokenType = TLSquareBracket
			} else {
				tokenType = TRSquareBracket
			}
			tokens = append(tokens, Token{
				Type:     tokenType,
				Data:     []rune{t.Letter()},
				StartPos: t.Pos,
				EndPos:   t.Pos + 1,
			})
			t.GoNext()
			continue
		}
	}

	tokens = append(tokens, Token{
		Type:     TEof,
		Data:     []rune{},
		StartPos: t.Pos,
		EndPos:   t.Pos + 1,
	})

	return &tokens
}
