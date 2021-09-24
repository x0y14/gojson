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
				Type:     TLSquareBracket,
				Data:     []rune{'['},
				StartPos: 0,
				EndPos:   1,
			},
			{
				Type:     TTrue,
				Data:     []rune("true"),
				StartPos: 1,
				EndPos:   5,
			},
			{
				Type:     TRSquareBracket,
				Data:     []rune{']'},
				StartPos: 5,
				EndPos:   6,
			},
			{
				Type:     TEof,
				Data:     []rune{},
				StartPos: 6,
				EndPos:   7,
			},
		}},
		{
			"\"hello\"", []Token{
				{
					Type:     TString,
					Data:     []rune("hello"),
					StartPos: 0,
					EndPos:   7,
				},
				{
					Type:     TEof,
					Data:     []rune{},
					StartPos: 7,
					EndPos:   8,
				},
			},
		},
		{
			"{\"msg\": \"hello\"}", []Token{
				{
					Type:     TLCurlyBracket,
					Data:     []rune("{"),
					StartPos: 0,
					EndPos:   1,
				},
				{
					Type:     TString,
					Data:     []rune("msg"),
					StartPos: 1,
					EndPos:   6,
				},
				{
					Type:     TColon,
					Data:     []rune(":"),
					StartPos: 6,
					EndPos:   7,
				},
				{
					Type:     TString,
					Data:     []rune("hello"),
					StartPos: 8,
					EndPos:   15,
				},
				{
					Type:     TRCurlyBracket,
					Data:     []rune("}"),
					StartPos: 15,
					EndPos:   16,
				},
				{
					Type:     TEof,
					Data:     []rune{},
					StartPos: 16,
					EndPos:   17,
				},
			},
		},
	}

	for _, tt := range tests {
		tokens := NewTokenizer(tt.json).Tokenize()
		assert.Equal(t, tt.want, *tokens)
	}
}
