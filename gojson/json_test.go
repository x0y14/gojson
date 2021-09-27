package gojson

import "testing"

func TestJson_Tree(t *testing.T) {
	var tests = []struct {
		title string
		json  string
	}{
		{
			"simple object",
			"{\"msg\": \"hello\", \"sub\": {\"age\": 26, \"name\": \"john\"}, \"members\": [\"tanaka\", \"sadako\"]}",
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
