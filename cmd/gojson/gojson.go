package main

import "github.com/x0y14/gojson/gojson"

func main() {
	json := "{\"msg\": \"hello\"}"
	tk := gojson.NewTokenizer(json)
	tk.Tokenize()
}
