package gojson

import (
	"fmt"
)

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
	p.Pos++
}

func (p *Parser) GoPrev() {
	p.Pos--
}

func (p *Parser) IsValidToken() bool {
	return p.Token().Type != TEof
}

func (p *Parser) Parse() (*Node, error) {
	var nd *Node
	var err error
	for p.IsValidToken() {
		if p.Token().Type == TLCurlyBracket {
			nd, err = p.ParseObject()
			if err != nil {
				return nil, err
			}
		} else if p.Token().Type == TLSquareBracket {
			nd, err = p.ParseArray()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, &ParserError{
				ErrorType:    SyntaxError,
				ErrorMessage: fmt.Sprintf("expected `[` or `{`, but found `%v`", string(p.Token().Data)),
				Tokens:       nil,
				StartPos:     0,
				EndPos:       0,
				AllowTypes:   nil,
				ActualType:   0,
			}
		}
	}
	return nd, nil
}

func (p *Parser) ParseObject() (*Node, error) {
	if p.Token().Type != TLCurlyBracket {
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expected `{`, but found `%v`", string(p.Token().Data)),
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}
	// consume '{'
	p.GoNext()

	members, err := p.ParseMember()
	if err != nil {
		return nil, err
	}

	if p.Token().Type != TRCurlyBracket {
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expected `}`, but found `%v`", string(p.Token().Data)),
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}
	// consume '}'
	p.GoNext()

	// { } eof
	//      ^
	// now we are in eof

	//if p.Token().Type != TEof {
	//	return nil, &ParserError{
	//		ErrorType:    SyntaxError,
	//		ErrorMessage: fmt.Sprintf("expected `EOF`, but found `%v`", string(p.Token().Data)),
	//		Tokens:       nil,
	//		StartPos:     0,
	//		EndPos:       0,
	//		AllowTypes:   nil,
	//		ActualType:   0,
	//	}
	//}

	obj := NewNode(NDObject, members, "", nil)

	return obj, nil
}

func (p *Parser) ParseMember() (*[]Node, error) {
	var member []Node

ValueLoop:
	for p.IsValidToken() {
		token := p.Token()
		switch token.Type {
		case TString:
			pair, err := p.ParsePair()
			if err != nil {
				return nil, nil
			}
			member = append(member, *pair)
		case TComma:
			p.GoNext()
			continue
		default:
			break ValueLoop
		}
	}

	return &member, nil
}

func (p *Parser) ParsePair() (*Node, error) {
	tkKey := p.Token()
	if tkKey.Type != TString {
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expected `TString`, but found `%v`", string(p.Token().Data)),
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}
	p.GoNext()

	if tkColon := p.Token(); tkColon.Type != TColon {
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expected `:`, but found `%v`", string(p.Token().Data)),
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}
	// consume ":"
	p.GoNext()

	tkVal, err := p.ParseValue()
	if err != nil {
		return nil, err
	}

	pair := NewNode(NDPair, &[]Node{*tkVal}, string(tkKey.Data), nil)

	return pair, nil
}

func (p *Parser) ParseArray() (*Node, error) {
	token := p.Token()
	if token.Type != TLSquareBracket {
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expect `[`, but found %v", string(token.Data)),
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}
	// consume '['
	p.GoNext()

	el, err := p.ParseElement()
	if err != nil {
		return nil, err
	}

	token = p.Token()
	if token.Type != TRSquareBracket {
		return nil, &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: fmt.Sprintf("expect `]`, but found %v", string(token.Data)),
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
			AllowTypes:   nil,
			ActualType:   0,
		}
	}
	// consume ']'
	p.GoNext()

	arr := NewNode(NDArray, el, "", nil)

	return arr, nil
}

func (p *Parser) ParseElement() (*[]Node, error) {
	var children []Node

ValueLoop:
	for p.IsValidToken() {
		fmt.Printf("[ParseElement @ %v-%v] `%v` (%v)\n", p.Token().StartPos, p.Token().EndPos, string(p.Token().Data), p.Token().Type.String())
		switch p.Token().Type {
		case TString, TNumber, TTrue, TFalse, TNull, TLCurlyBracket, TLSquareBracket:
			nd, err := p.ParseValue()
			if err != nil {
				return nil, err
			}
			children = append(children, *nd)
		case TComma:
			p.GoNext()
			continue
		default:
			break ValueLoop
		}
	}

	return &children, nil
}

func (p *Parser) ParseValue() (*Node, error) {
	token := p.Token()
	var nd *Node

	switch token.Type {
	// their kinds can access the value directly
	case TString, TNumber, TTrue, TFalse, TNull:
		nd = NewNode(NDValue, nil, "", &token)
		p.GoNext()
	case TLCurlyBracket:
		child, err := p.ParseObject()
		if err != nil {
			return nil, err
		}
		//nd = NewNode(NDObject, &[]Node{*child}, nil)
		nd = child
	case TLSquareBracket:
		child, err := p.ParseArray()
		if err != nil {
			return nil, err
		}
		nd = child
	}
	return nd, nil
}
