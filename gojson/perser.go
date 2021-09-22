package gojson

func NewParser(tokens *[]Token) *Parser {
	return &Parser{
		Tokens:                     *tokens,
		Pos:                        0,
		unpackedSquareBracketsNest: 0,
		unpackedCurlyBracketsNest:  0,
	}
}

type Parser struct {
	Tokens                     []Token
	Pos                        int
	unpackedSquareBracketsNest int
	unpackedCurlyBracketsNest  int
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

func (p *Parser) Parse() error {
	for len(p.Tokens) > p.Pos {
		switch p.Token().Type {
		case LCurlyBracket:
			p.unpackedCurlyBracketsNest++
		case RCurlyBracket:
			p.unpackedCurlyBracketsNest--
		case LSquareBracket:
			p.unpackedSquareBracketsNest++
		case RSquareBracket:
			p.unpackedSquareBracketsNest--
		}
		p.GoNext()
	}
	if p.unpackedCurlyBracketsNest != 0 {
		return &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: "Unpacked Curly Bracket",
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
		}
	} else if p.unpackedSquareBracketsNest != 0 {
		return &ParserError{
			ErrorType:    SyntaxError,
			ErrorMessage: "Unpacked Square Bracket",
			Tokens:       nil,
			StartPos:     0,
			EndPos:       0,
		}
	}
	return nil
}
