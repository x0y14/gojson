package gojson

import "fmt"

// EBNF rules?
// Json = Object | Array
// - Object = "{" [Member] "}"
//   - Member = String ":" Element
// - Array = "[" [Element] "]"
// - Element = Object | Array | String | Number | Boolean | Null

func NewParser(tokens *[]Token) *Parser {
	return &Parser{
		Tokens: *tokens,
		Pos:    0,
		Depth:  0,
	}
}

type Parser struct {
	Tokens []Token
	Pos    int
	Depth  int
}

func (p *Parser) Token() Token {
	return p.Tokens[p.Pos]
}

func (p *Parser) NextToken() Token {
	return p.Tokens[p.Pos+1]
}

func (p *Parser) PrevToken() Token {
	return p.Tokens[p.Pos-1]
}

func (p *Parser) GoNext() {
	//if p.Pos > 0 {
	//	fmt.Printf("[%03d] `%v` -> ", p.Pos-1, string(p.PrevToken().Data))
	//} else {
	//	fmt.Printf("SOF -> ")
	//}

	//if len(p.Tokens) > p.Pos {
	//	fmt.Printf("[%03d] `%v` -> ", p.Pos, string(p.Token().Data))
	//}
	p.Pos++
	//if len(p.Tokens) > p.Pos {
	//	fmt.Printf("[%03d] `%v`\n", p.Pos, string(p.Token().Data))
	//} else {
	//	fmt.Printf("EOF\n")
	//}
}

func (p *Parser) GoPrev() {
	p.Pos--
}

func (p *Parser) ParseMember() (*[]Member, error) {
	fmt.Printf("[*] call ParseMember\n")
	var members []Member

	// Member <String : Element>
	// 		  [<String : Element> ","]

	for len(p.Tokens) > p.Pos {
		key := p.Token()
		if key.Type != TString {
			return nil, &ParserError{
				ErrorType:    SyntaxError,
				ErrorMessage: fmt.Sprintf("expected `TString`, but found `%v`", key.Type.String()),
				Tokens:       nil,
				StartPos:     0,
				EndPos:       0,
				AllowTypes:   nil,
				ActualType:   0,
			}
		}
		// go to the next of key
		p.GoNext()

		if colon := p.Token(); colon.Type != TColon {
			return nil, &ParserError{
				ErrorType:    SyntaxError,
				ErrorMessage: fmt.Sprintf("expected `TColon`, but found `%v`", colon.Type.String()),
				Tokens:       nil,
				StartPos:     0,
				EndPos:       0,
				AllowTypes:   nil,
				ActualType:   0,
			}
		}
		// go to the next of colon
		p.GoNext()

		valueToken := p.Token()
		var value interface{}
		var err error

		switch valueToken.Type {
		case TLCurlyBracket:
			value, err = p.ParseObject()
			if err != nil {
				return nil, err
			}
			// At ParseObject,
			// Since GoNext() is called at the end of the function, no operation is required on this side.
			p.GoNext()
		case TString:
			value = valueToken.LoadAsString()
			p.GoNext()
		case TTrue, TFalse:
			value = valueToken.LoadAsBoolean()
			p.GoNext()
		case TNumber:
			value = valueToken.LoadAsFloat64()
			p.GoNext()
		case TNull:
			value = valueToken.LoadAsNull()
			p.GoNext()
			// todo : ParseArray
		}

		members = append(members, Member{
			Key:   key.LoadAsString(),
			Value: value,
		})

		if comma := p.Token(); comma.Type != TComma {
			break
		} else {
			// go to the next of comma
			p.GoNext()
		}
	}

	return &members, nil
}

func (p *Parser) ParseObject() (map[string]interface{}, error) {
	fmt.Printf("[*] call ParseObject\n")

	obj := map[string]interface{}{}

	// {
	if lCurly := p.Token(); lCurly.Type != TLCurlyBracket {
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expected `TLCurlyBracket` but found `%v`", lCurly.Type.String()),
			Tokens:       []Token{lCurly},
			StartPos:     lCurly.StartPos,
			EndPos:       lCurly.EndPos,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}
	// go to the next of lCurly
	p.GoNext()

	stringOrRCurly := p.Token()
	fmt.Printf(">> stringOrRCurly : %v(%v) <<\n", string(stringOrRCurly.Data), stringOrRCurly.Type.String())
	switch stringOrRCurly.Type {
	case TString:
		members, err := p.ParseMember()
		if err != nil {
			return nil, err
		}

		for _, member := range *members {
			obj[member.Key] = member.Value
		}

		//if exceptRCurly := p.Token(); exceptRCurly.Type != TRCurlyBracket {
		//	return nil, &ParserError{
		//				ErrorType:    SyntaxError,
		//				ErrorMessage: fmt.Sprintf("expected `TRCurlyBracket`, but found `%v\n`", exceptRCurly.Type.String()),
		//				Tokens:       nil,
		//				StartPos:     0,
		//				EndPos:       0,
		//				AllowTypes:   nil,
		//				ActualType:   0,
		//			}
		//} else {
		//	break
		//}
	case TRCurlyBracket:
		// empty map
		p.GoNext()
		break
	default:
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expected `TString` or `TRCurlyBracket`, but found `%v\n`", stringOrRCurly.Type.String()),
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}

	fmt.Printf(">> end of ParseObject : `%v` (%v) <<\n", string(p.Token().Data), p.Token().Type.String())

	return obj, nil
}

func (p *Parser) ParseArray() ([]interface{}, error) {
	fmt.Printf("[*] call ParseObject\n")
	return nil, nil
}

func (p *Parser) Parse() (interface{}, error) {
	var obj interface{}
	for len(p.Tokens) > p.Pos {
		switch p.Token().Type {
		case TLCurlyBracket:
			obj, err := p.ParseObject()
			if err != nil {
				return nil, err
			}
			return obj, nil
		case TLSquareBracket:
			obj, err := p.ParseArray()
			if err != nil {
				return nil, err
			}
			return obj, nil
		}

		p.GoNext()
	}
	return obj, nil
}

func IsTypeAllowed(allowing []TokenType, tokenType TokenType) bool {
	for _, tType := range allowing {
		if tType == tokenType {
			return true
		}
	}
	return false
}
