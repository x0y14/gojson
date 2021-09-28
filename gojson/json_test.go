package gojson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJson_Tree(t *testing.T) {
	var tests = []struct {
		title string
		json  string
	}{
		{
			"simple object",
			"{" +
				"\"msg\": \"hello\", " +
				"\"sub\": {\"age\": 26, \"name\": \"john\"}, " +
				"\"members\": [\"tanaka\", \"sadako\"], " +
				"\"users\": [" +
				"{\"name\":\"john\", \"age\": 35, \"active\": true, \"skill\": \"jump\"}, " +
				"{\"name\":\"tom\", \"age\": 12, \"active\": false, \"skill\": null}]}",
		},
		{
			"simple object",
			"{" +
				"\"msg\": \"hello\", " +
				"\"sub\": {\"age\": 26, \"name\": \"john\"}, " +
				"\"members\": [\"tanaka\", \"sadako\"], " +
				"\"users\": [" +
				"{\"name\":\"john\", \"age\": 35, \"active\": true, \"skill\": \"jump\"}, " +
				"{\"name\":\"tom\", \"age\": 12, \"active\": false, \"skill\": null}]}",
		},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.json)
		ps := NewParser(tk.Tokenize())
		nd, err := ps.Parse()
		if err != nil {
			t.Fatal(err)
		}
		nd.Tree()
	}
}

func TestJson_Map(t *testing.T) {
	var tests = []struct {
		title string
		json  string
	}{
		{
			"simple object",
			"{" +
				"\"msg\": \"hello\", " +
				"\"sub\": {\"age\": 26, \"name\": \"john\"}, " +
				"\"members\": [\"tanaka\", \"sadako\"], " +
				"\"users\": [" +
				"{\"name\":\"john\", \"age\": 35, \"active\": true, \"skill\": \"jump\"}, " +
				"{\"name\":\"tom\", \"age\": 12, \"active\": false, \"skill\": null}]}",
		},
		{
			"simple object",
			"{" +
				"\"msg\": \"hello\", " +
				"\"sub\": {\"age\": 26, \"name\": \"john\"}, " +
				"\"members\": [\"tanaka\", \"sadako\"], " +
				"\"users\": [" +
				"{\"name\":\"john\", \"age\": 35, \"active\": true, \"skill\": \"jump\"}, " +
				"{\"name\":\"tom\", \"age\": 12, \"active\": false, \"skill\": null}]}",
		},
	}

	for _, tt := range tests {
		tk := NewTokenizer(tt.json)
		ps := NewParser(tk.Tokenize())
		nd, err := ps.Parse()
		if err != nil {
			t.Fatal(err)
		}
		//nd.Tree()
		mp, err := nd.Map()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "hello", mp["msg"])
	}
}

func TestJson_Map2(t *testing.T) {
	json := "{" +
		"\"msg\": \"hello\", " +
		"\"sub\": {\"age\": 26, \"name\": \"john\"}, " +
		"\"members\": [\"tanaka\", \"sadako\"], " +
		"\"users\": [" +
		"{\"name\":\"john\", \"age\": 35, \"active\": true, \"skill\": \"jump\"}, " +
		"{\"name\":\"tom\", \"age\": 12, \"active\": false, \"skill\": null}]" +
		"}"
	tk := NewTokenizer(json)
	ps := NewParser(tk.Tokenize())
	nd, err := ps.Parse()
	if err != nil {
		t.Fatal(err)
	}
	mp, err := nd.Map()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", mp["msg"])
	assert.Equal(t, map[string]interface{}{"age": float64(26), "name": "john"}, mp["sub"])
	assert.Equal(t, []interface{}{"tanaka", "sadako"}, mp["members"])

	members := mp["members"].([]interface{})
	//fmt.Printf("%v", members[0])
	assert.Equal(t, "sadako", members[1])
}
