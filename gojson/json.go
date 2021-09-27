package gojson

type Json struct {
	//Children []Element
}

type Object []Token
type Array []Token

type Elements struct {
	es []Elements
}

type Member struct {
	Key   string
	Value interface{}
}
