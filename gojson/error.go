package gojson

import "fmt"

type ErrorType int

const (
	UnknownError ErrorType = iota
	InvalidDataError
	UndefinedError
)

type TKError struct {
	ErrorType    ErrorType
	ErrorMessage string
	Letters      []rune
	StartPos     int
	EndPos       int
}

func (e *TKError) Error() string {
	return fmt.Sprintf("[e-%v @ %03d-%03d] %v: `%v`", e.ErrorType, e.StartPos, e.EndPos, e.ErrorMessage, string(e.Letters))
}
