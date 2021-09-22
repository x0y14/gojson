package gojson

type TokenType int

const (
	Unknown TokenType = iota
	Null
	String
	Number

	True
	False

	WhiteSpace

	Comma
	Colon
	LCurlyBracket
	RCurlyBracket
	LSquareBracket
	RSquareBracket
)

type Token struct {
	Type TokenType
	Data []rune

	// StartPos <= Data < EndPos
	StartPos int
	EndPos   int
}
