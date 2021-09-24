package gojson

import "strconv"

type TokenType int

const (
	TUnknown TokenType = iota
	TEof
	TNull
	TString
	TNumber
	TTrue
	TFalse

	TWhiteSpace
	TComma
	TColon
	TLCurlyBracket
	TRCurlyBracket
	TLSquareBracket
	TRSquareBracket
)

type Token struct {
	Type TokenType
	Data []rune

	// StartPos <= Data < EndPos
	StartPos int
	EndPos   int
}

func (tokenType TokenType) String() string {
	switch tokenType {
	case TEof:
		return "TEof"
	case TNull:
		return "TNull"
	case TString:
		return "TString"
	case TNumber:
		return "TNumber"
	case TTrue:
		return "TTrue"
	case TFalse:
		return "TFalse"
	case TWhiteSpace:
		return "TWhiteSpace"
	case TComma:
		return "TComma"
	case TColon:
		return "TColon"
	case TLCurlyBracket:
		return "TLCurlyBracket"
	case TRCurlyBracket:
		return "TRCurlyBracket"
	case TLSquareBracket:
		return "TLSquareBracket"
	case TRSquareBracket:
		return "TRSquareBracket"
	default:
		return "TUnknown"
	}
}

func (t *Token) LoadAsFloat64() float64 {
	if t.Type != TNumber {
		panic(&TokenizerError{
			ErrorType:    IllegalValueLoadingError,
			ErrorMessage: "This Token is not TNumber",
			Letters:      t.Data,
			StartPos:     t.StartPos,
			EndPos:       t.EndPos,
		})
	}
	f, err := strconv.ParseFloat(string(t.Data), 64)
	if err != nil {
		panic(err)
	}
	return f
}

func (t *Token) LoadAsString() string {
	if t.Type != TString {
		panic(&TokenizerError{
			ErrorType:    IllegalValueLoadingError,
			ErrorMessage: "This Token is not TString",
			Letters:      t.Data,
			StartPos:     t.StartPos,
			EndPos:       t.EndPos,
		})
	}
	return string(t.Data)
}

func (t *Token) LoadAsBoolean() bool {
	if t.Type != TTrue && t.Type != TFalse {
		panic(&TokenizerError{
			ErrorType:    IllegalValueLoadingError,
			ErrorMessage: "This Token is neither TTrue nor TFalse",
			Letters:      t.Data,
			StartPos:     t.StartPos,
			EndPos:       t.EndPos,
		})
	}
	switch string(t.Data) {
	case "true":
		return true
	default:
		return false
	}
}

func (t *Token) LoadAsNull() error {
	if t.Type != TNull {
		panic(&TokenizerError{
			ErrorType:    IllegalValueLoadingError,
			ErrorMessage: "This Token is TNull",
			Letters:      t.Data,
			StartPos:     t.StartPos,
			EndPos:       t.EndPos,
		})
	}
	return nil
}
