package gojson

// EBNF rules?
// Json = Object | Array
// - Object = "{" [Member] "}"
//   - Member = String ":" Element
// - Array = "[" [Element] "]"
// - Element = Object | Array | String | Number | Boolean | Null

type Json struct {
	//Children []Element
}

//type ElementType int
//const (
//	EUnknown ElementType = iota
//	EObject
//	EArray
//	EString
//	ENumber
//	EBoolean
//	ENull
//)
//
//type Element struct {
//	Type ElementType
//	Tokens []Token
//}

type Object []Token
type Array []Token

type Elements struct {
	es []Elements
}

type Member struct {
	Key   string
	Value interface{}
}
