package gojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParser(t *testing.T) {
}

func TestParser_Parse(t *testing.T) {
}

func TestParser_ParseError(t *testing.T) {
	var tests = []struct {
		json string
		want error
	}{
		{
			"{\"msg\": \"hello\"",
			&ParserError{
				ErrorType:    SyntaxError,
				ErrorMessage: "Unpacked Curly Bracket",
				Tokens:       nil,
				StartPos:     0,
				EndPos:       0,
			},
		},
		{
			"[",
			&ParserError{
				ErrorType:    SyntaxError,
				ErrorMessage: "Unpacked Square Bracket",
				Tokens:       nil,
				StartPos:     0,
				EndPos:       0,
			},
		},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.json)
		ps := NewParser(tk.Tokenize())
		assert.Equal(t, tt.want, ps.Parse())
	}
}
