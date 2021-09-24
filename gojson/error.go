package gojson

import "fmt"

type ErrorType int

const (
	UnknownError ErrorType = iota

	InvalidDataError
	UndefinedKeywordError
	IllegalValueLoadingError

	SyntaxError
)

func (et ErrorType) String() string {
	switch et {
	case InvalidDataError:
		return "InvalidDataError"
	case UndefinedKeywordError:
		return "UndefinedKeywordError"
	case IllegalValueLoadingError:
		return "IllegalValueLoadingError"
	case SyntaxError:
		return "SyntaxError"
	default:
		return "UnknownError"
	}
}

type TokenizerError struct {
	ErrorType    ErrorType
	ErrorMessage string
	Letters      []rune
	StartPos     int
	EndPos       int
}

func (e *TokenizerError) Error() string {
	return fmt.Sprintf("[t-%v @ %03d-%03d] %v: `%v`", e.ErrorType.String(), e.StartPos, e.EndPos, e.ErrorMessage, string(e.Letters))
}

type ParserError struct {
	ErrorType    ErrorType
	ErrorMessage string
	Tokens       []Token
	StartPos     int
	EndPos       int

	AllowTypes []TokenType
	ActualType TokenType
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("[p-%v @ %03d-%03d] %v", e.ErrorType.String(), e.StartPos, e.EndPos, e.ErrorMessage)
}
