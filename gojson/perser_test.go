package gojson_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/x0y14/gojson/gojson"
	"testing"
)

func Setup(json string) *gojson.Parser {
	tk := gojson.NewTokenizer(json)
	ps := gojson.NewParser(tk.Tokenize())
	return ps
}

func TestShowPos(t *testing.T) {
	gojson.ShowPos("{\"msg\": \"hello\", \"in\": {\"age\": 20}}")
}

func TestParser_ParseArray(t *testing.T) {
	var tests = []struct {
		title  string
		json   string
		expect interface{}
	}{
		{
			"array only",
			"[\"string\", 123, true, false, null]",
			gojson.NewNode(gojson.NDArray, &[]gojson.Node{
				*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "string", 1, 9)),
				*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TNumber, "123", 11, 14)),
				*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TTrue, "true", 16, 20)),
				*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TFalse, "false", 22, 27)),
				*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TNull, "null", 29, 33)),
			}, "", nil),
		},
		{
			"array in array",
			"[[\"hello\", \"world\"]]",
			gojson.NewNode(gojson.NDArray, &[]gojson.Node{
				*gojson.NewNode(gojson.NDArray, &[]gojson.Node{
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "hello", 2, 9)),
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "world", 11, 18)),
				}, "", nil),
			}, "", nil),
		},
		{
			"array and other",
			"[123, [\"hello\", \"world\"], \"321\"]",
			gojson.NewNode(gojson.NDArray, &[]gojson.Node{
				*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TNumber, "123", 1, 4)),
				*gojson.NewNode(gojson.NDArray, &[]gojson.Node{
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "hello", 7, 14)),
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "world", 16, 23)),
				}, "", nil),
				*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "321", 26, 31)),
			}, "", nil),
		},
	}

	for _, tt := range tests {
		ps := Setup(tt.json)
		actual, err := ps.ParseArray()
		if err != nil {
			t.Fatal(err)
		}
		res := assert.Equal(t, tt.expect, actual)
		fmt.Printf("^^^ %v ^^^ \n  -> success: %v\n", tt.title, res)
	}
}

func TestParser_ParseObject(t *testing.T) {
	var tests = []struct {
		title  string
		json   string
		expect interface{}
	}{
		{
			"simple object",
			"{\"msg\": \"hello\"}",
			gojson.NewNode(gojson.NDObject, &[]gojson.Node{
				*gojson.NewNode(gojson.NDPair, &[]gojson.Node{
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "hello", 8, 15)),
				}, "msg", nil),
			}, "", nil),
		},
		{
			"object multi key",
			"{\"msg\":\"hello\", \"age\": 20}",
			gojson.NewNode(gojson.NDObject, &[]gojson.Node{
				*gojson.NewNode(gojson.NDPair, &[]gojson.Node{
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "hello", 7, 14)),
				}, "msg", nil),
				*gojson.NewNode(gojson.NDPair, &[]gojson.Node{
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TNumber, "20", 23, 25)),
				}, "age", nil),
			}, "", nil),
		},
		{
			"object in object",
			"{\"msg\": \"hello\", \"in\": {\"age\": 20}}",
			gojson.NewNode(gojson.NDObject, &[]gojson.Node{
				*gojson.NewNode(gojson.NDPair, &[]gojson.Node{
					*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TString, "hello", 8, 15)),
				}, "msg", nil),
				*gojson.NewNode(gojson.NDPair, &[]gojson.Node{
					*gojson.NewNode(gojson.NDObject, &[]gojson.Node{
						*gojson.NewNode(gojson.NDPair, &[]gojson.Node{
							*gojson.NewNode(gojson.NDValue, nil, "", gojson.NewToken(gojson.TNumber, "20", 31, 33)),
						}, "age", nil),
					}, "", nil),
				}, "in", nil),
			}, "", nil),
		},
	}
	for _, tt := range tests {
		ps := Setup(tt.json)
		actual, err := ps.ParseObject()
		if err != nil {
			t.Fatal(err)
		}
		res := assert.Equal(t, tt.expect, actual)
		fmt.Printf("^^^ %v ^^^ \n  -> success: %v\n", tt.title, res)
	}
}
