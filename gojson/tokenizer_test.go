package gojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTokenizer(t *testing.T) {
	text := "{\"msg\": \"hello\"}"
	tk := NewTokenizer(text)

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
		assert.Equal(t, tt.want, tk.TokenAt(tt.pos))
	}
}
