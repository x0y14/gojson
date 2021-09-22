package gojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func SetUpSimpleJson() *Tokenizer {
	text := "{\"msg\": \"hello\"}"
	tk := NewTokenizer(text)
	return tk
}

func TestNewTokenizer(t *testing.T) {
	tk := SetUpSimpleJson()

	var tests = []struct {
		pos  int
		want rune
	}{
		{0, '{'},
		{1, '"'},
		{2, 'm'},
		{3, 's'},
		{4, 'g'},
		{5, '"'},
		{6, ':'},
		{7, ' '},
		{8, '"'},
		{9, 'h'},
		{10, 'e'},
		{11, 'l'},
		{12, 'l'},
		{13, 'o'},
		{14, '"'},
		{15, '}'},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tk.LetterAt(tt.pos))
	}
}

func TestTokenizer_GoNext(t *testing.T) {
	tk := SetUpSimpleJson()
	assert.Equal(t, '{', tk.Letter())
	tk.GoNext()
	assert.Equal(t, '"', tk.Letter())
}

func TestTokenizer_GoPrev(t *testing.T) {
	tk := SetUpSimpleJson()
	tk.Pos = 1
	assert.Equal(t, '"', tk.Letter())
	tk.GoPrev()
	assert.Equal(t, '{', tk.Letter())
}

func TestTokenizer_ConsumeWhiteSpace(t *testing.T) {
	text := "{ \t\n }"
	tk := NewTokenizer(text)
	assert.Equal(t, '{', tk.Letter())
	tk.GoNext()
	assert.Equal(t, ' ', tk.Letter())
	tk.ConsumeWhiteSpace()
	assert.Equal(t, '}', tk.Letter())
}

func TestTokenizer_Tokenize(t *testing.T) {
	var tests = []struct {
		json string
		want []Token
	}{
		{"[true]", []Token{
			{
				Type:     LSquareBracket,
				Data:     []rune{'['},
				StartPos: 0,
				EndPos:   1,
			},
			{
				Type:     True,
				Data:     []rune("true"),
				StartPos: 1,
				EndPos:   5,
			},
			{
				Type:     RSquareBracket,
				Data:     []rune{']'},
				StartPos: 5,
				EndPos:   6,
			},
		}},
		{
			"\"hello\"", []Token{
				{
					Type:     String,
					Data:     []rune("hello"),
					StartPos: 0,
					EndPos:   7,
				},
			},
		},
		{
			"{\"msg\": \"hello\"}", []Token{
				{
					Type:     LCurlyBracket,
					Data:     []rune("{"),
					StartPos: 0,
					EndPos:   1,
				},
				{
					Type:     String,
					Data:     []rune("msg"),
					StartPos: 1,
					EndPos:   6,
				},
				{
					Type:     Colon,
					Data:     []rune(":"),
					StartPos: 6,
					EndPos:   7,
				},
				{
					Type:     String,
					Data:     []rune("hello"),
					StartPos: 8,
					EndPos:   15,
				},
				{
					Type:     RCurlyBracket,
					Data:     []rune("}"),
					StartPos: 15,
					EndPos:   16,
				},
			},
		},
	}

	for _, tt := range tests {
		tokens := NewTokenizer(tt.json).Tokenize()
		assert.Equal(t, tt.want, *tokens)
	}
}
