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

func (tokenType TokenType) String() string {
	switch tokenType {
	case Null:
		return "Null"
	case String:
		return "String"
	case Number:
		return "Number"
	case True:
		return "True"
	case False:
		return "False"
	case WhiteSpace:
		return "WhiteSpace"
	case Comma:
		return "Comma"
	case Colon:
		return "Colon"
	case LCurlyBracket:
		return "LCurlyBracket"
	case RCurlyBracket:
		return "RCurlyBracket"
	case LSquareBracket:
		return "LSquareBracket"
	case RSquareBracket:
		return "RSquareBracket"
	default:
		return "Unknown"
	}
}
