package parser

import (
	"fmt"

	"github.com/LightningDev1/golox/internal/ast"
	"github.com/LightningDev1/golox/internal/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func New(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var statements []ast.Stmt

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) declaration() (stmt ast.Stmt, err error) {
	if p.match(scanner.TOKEN_FUN) {
		stmt, err = p.function("function")
	} else if p.match(scanner.TOKEN_VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) function(kind string) (ast.Stmt, error) {
	name, err := p.consume(scanner.TOKEN_IDENTIFIER,
		fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.TOKEN_LEFT_PAREN,
		fmt.Sprintf("Expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}

	var parameters []scanner.Token
	if !p.check(scanner.TOKEN_RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, NewParseError(p.peek(),
					"Can't have more than 255 parameters.")
			}

			paramName, err := p.consume(scanner.TOKEN_IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, paramName)

			if !p.match(scanner.TOKEN_COMMA) {
				break
			}
		}
	}

	_, err = p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.TOKEN_LEFT_BRACE,
		fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStmt{Name: name, Params: parameters, Body: body}, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(scanner.TOKEN_IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(scanner.TOKEN_EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &ast.VarStmt{Name: name, Initializer: initializer}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(scanner.TOKEN_FOR) {
		return p.forStatement()
	}
	if p.match(scanner.TOKEN_IF) {
		return p.ifStatement()
	}
	if p.match(scanner.TOKEN_PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.TOKEN_RETURN) {
		return p.returnStatement()
	}
	if p.match(scanner.TOKEN_WHILE) {
		return p.whileStatement()
	}
	if p.match(scanner.TOKEN_LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.BlockStmt{Statements: stmts}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Stmt
	if p.match(scanner.TOKEN_SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.TOKEN_VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr
	if !p.check(scanner.TOKEN_SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.TOKEN_SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment ast.Expr
	if !p.check(scanner.TOKEN_RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{body, &ast.ExpressionStmt{Expression: increment}},
		}
	}

	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	if p.match(scanner.TOKEN_ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.TOKEN_SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStmt{Expression: value}, nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()

	var value ast.Expr
	if !p.check(scanner.TOKEN_SEMICOLON) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(scanner.TOKEN_SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &ast.ReturnStmt{Keyword: keyword, Value: value}, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
	_, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{Condition: condition, Body: body}, nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	var statements []ast.Stmt

	for !p.check(scanner.TOKEN_RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	if _, err := p.consume(scanner.TOKEN_RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.TOKEN_SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.TOKEN_EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*ast.Variable); ok {
			name := variable.Name
			return &ast.Assign{Name: name, Value: value}, nil
		}

		return nil, NewParseError(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.TOKEN_OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.TOKEN_AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
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

	return p.call()
}

func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(scanner.TOKEN_LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	var arguments []ast.Expr

	if !p.check(scanner.TOKEN_RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				return nil, NewParseError(p.peek(), "Can't have more than 255 arguments.")
			}

			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if !p.match(scanner.TOKEN_COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.Call{Callee: callee, Paren: paren, Arguments: arguments}, nil
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

	if p.match(scanner.TOKEN_IDENTIFIER) {
		return &ast.Variable{Name: p.previous()}, nil
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
