package gojson

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParser(t *testing.T) {
}

func TestParser_Parse(t *testing.T) {
}

//func TestParser_ParseError(t *testing.T) {
//	var tests = []struct {
//		name string
//		json string
//		want error
//	}{
//		{
//			"unpacked curly brackets",
//			"{\"msg\": \"hello\"",
//			&ParserError{
//				ErrorType:    SyntaxError,
//				ErrorMessage: "Unpacked Curly Bracket",
//				Tokens:       nil,
//				StartPos:     0,
//				EndPos:       0,
//			},
//		},
//		{
//			"unpacked square brackets",
//			"[",
//			&ParserError{
//				ErrorType:    SyntaxError,
//				ErrorMessage: "Unpacked Square Bracket",
//				Tokens:       nil,
//				StartPos:     0,
//				EndPos:       0,
//			},
//		},
//	}
//
//	//for _, tt := range tests {
//	//	tk := NewTokenizer(tt.json)
//		//ps := NewParser(tk.Tokenize())
//		//assert.Equal(t, tt.want, ps.Parse())
//	//}
//}

func TestParser_ParseObject(t *testing.T) {
	var tests = []struct {
		name string
		json string
		key  string
		want interface{}
	}{
		{
			"string",
			"{\"msg\": \"hello\"}",
			"msg",
			"hello",
		},
		{
			"number",
			"{\"age\":20}",
			"age",
			float64(20),
		},
		{
			"boolean",
			"{\"dark_theme\": true}",
			"dark_theme",
			true,
		},
		{
			"null",
			"{\"my_money\": null}",
			"my_money",
			nil,
		},
		{
			"map in map",
			"{\"config\": {\"alarm\": true}}",
			"config",
			map[string]interface{}{"alarm": true},
		},
		{
			"multiple map",
			"{\"my_age\":21, \"sister_age\":18}",
			"my_age",
			float64(21),
		},
		{
			"map and int",
			"{\"info\": {\"alarm\": true}, \"age\": 64}",
			"age",
			float64(64),
		},
		{
			"multi in map",
			"{\"a\":{\"b\": true, \"c\": null}}",
			"a",
			map[string]interface{}{"b": true, "c": nil},
		},
		{
			"multiple map value1",
			"{\"a\":{\"alarm\": true}, \"b\":{\"alarm2\": false}}",
			"b",
			map[string]interface{}{"alarm2": false},
		},
		{
			"multiple map value2",
			"{\"a\":{\"alarm\": true}, \"b\":{\"alarm2\": false}}",
			"a",
			map[string]interface{}{"alarm": true},
		},
		{
			"multiple map value3",
			"{\"a\":{}, \"b\":{}}",
			"b",
			map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		tk := NewTokenizer(tt.json)
		ps := NewParser(tk.Tokenize())
		fmt.Printf("\ntry: %v\n", tt.name)
		obj, err := ps.ParseObject()
		if err != nil {
			t.Fatal(err)
		}
		if res := assert.Equal(t, tt.want, obj[tt.key]); res == true {
			fmt.Printf("    -> success\n")
		} else {
			fmt.Printf("    -> failure\n")
		}
	}
}

//func TestParser_ParseArray(t *testing.T) {
//	var tests = []struct{
//		name string
//		json string
//		index int
//		want interface{}
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		tk := NewTokenizer(tt.json)
//		ps := NewParser(tk.Tokenize())
//		_ = ps.ParseArray()
//		//log.Printf("try: %v\n", tt.name)
//		//assert.Equal(t, tt.want, obj[tt.key])
//	}
//}
