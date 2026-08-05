package parser

import (
	"github.com/LightningDev1/golox/internal/ast"
	"github.com/LightningDev1/golox/internal/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() (ast.Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.TOKEN_BANG_EQUAL, scanner.TOKEN_EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.TOKEN_GREATER, scanner.TOKEN_GREATER_EQUAL,
		scanner.TOKEN_LESS, scanner.TOKEN_LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.TOKEN_MINUS, scanner.TOKEN_PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.TOKEN_SLASH, scanner.TOKEN_STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(scanner.TOKEN_BANG, scanner.TOKEN_MINUS) {
		operator := p.previous()

		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.Unary{Operator: operator, Right: right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(scanner.TOKEN_FALSE) {
		return &ast.Literal{Value: false}, nil
	}
	if p.match(scanner.TOKEN_TRUE) {
		return &ast.Literal{Value: true}, nil
	}
	if p.match(scanner.TOKEN_NIL) {
		return &ast.Literal{Value: nil}, nil
	}

	if p.match(scanner.TOKEN_NUMBER, scanner.TOKEN_STRING) {
		return &ast.Literal{Value: p.previous().Literal}, nil
	}

	if p.match(scanner.TOKEN_LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return &ast.Grouping{Expression: expr}, nil
	}

	return nil, NewParseError(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return scanner.Token{}, NewParseError(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.TOKEN_SEMICOLON {
			return
		}

		switch p.peek().Type {
		case scanner.TOKEN_CLASS, scanner.TOKEN_FUN,
			scanner.TOKEN_VAR, scanner.TOKEN_FOR,
			scanner.TOKEN_IF, scanner.TOKEN_WHILE,
			scanner.TOKEN_PRINT, scanner.TOKEN_RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.TOKEN_EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}
